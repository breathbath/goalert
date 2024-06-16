package app

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/breathbath/goalert/notification/pinpoint"
	"net"
	"net/http"

	"github.com/breathbath/goalert/alert"
	"github.com/breathbath/goalert/alert/alertlog"
	"github.com/breathbath/goalert/alert/alertmetrics"
	"github.com/breathbath/goalert/apikey"
	"github.com/breathbath/goalert/app/lifecycle"
	"github.com/breathbath/goalert/auth"
	"github.com/breathbath/goalert/auth/authlink"
	"github.com/breathbath/goalert/auth/basic"
	"github.com/breathbath/goalert/auth/nonce"
	"github.com/breathbath/goalert/calsub"
	"github.com/breathbath/goalert/config"
	"github.com/breathbath/goalert/engine"
	"github.com/breathbath/goalert/escalation"
	"github.com/breathbath/goalert/graphql2/graphqlapp"
	"github.com/breathbath/goalert/heartbeat"
	"github.com/breathbath/goalert/integrationkey"
	"github.com/breathbath/goalert/integrationkey/uik"
	"github.com/breathbath/goalert/keyring"
	"github.com/breathbath/goalert/label"
	"github.com/breathbath/goalert/limit"
	"github.com/breathbath/goalert/notice"
	"github.com/breathbath/goalert/notification"
	"github.com/breathbath/goalert/notification/slack"
	"github.com/breathbath/goalert/notification/twilio"
	"github.com/breathbath/goalert/notificationchannel"
	"github.com/breathbath/goalert/oncall"
	"github.com/breathbath/goalert/override"
	"github.com/breathbath/goalert/permission"
	"github.com/breathbath/goalert/schedule"
	"github.com/breathbath/goalert/schedule/rotation"
	"github.com/breathbath/goalert/schedule/rule"
	"github.com/breathbath/goalert/service"
	"github.com/breathbath/goalert/smtpsrv"
	"github.com/breathbath/goalert/timezone"
	"github.com/breathbath/goalert/user"
	"github.com/breathbath/goalert/user/contactmethod"
	"github.com/breathbath/goalert/user/favorite"
	"github.com/breathbath/goalert/user/notificationrule"
	"github.com/breathbath/goalert/util/log"
	"github.com/breathbath/goalert/util/sqlutil"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// App represents an instance of the GoAlert application.
type App struct {
	cfg Config

	mgr *lifecycle.Manager

	db     *sql.DB
	l      net.Listener
	events *sqlutil.Listener

	doneCh chan struct{}

	sysAPIL   net.Listener
	sysAPISrv *grpc.Server
	hSrv      *health.Server

	srv        *http.Server
	smtpsrv    *smtpsrv.Server
	smtpsrvL   net.Listener
	startupErr error

	notificationManager *notification.Manager
	Engine              *engine.Engine
	graphql2            *graphqlapp.App
	AuthHandler         *auth.Handler

	twilioSMS    *twilio.SMS
	twilioVoice  *twilio.Voice
	twilioConfig *twilio.Config

	pinpointSMS    *pinpoint.SMS
	pinpointConfig *pinpoint.Config

	slackChan *slack.ChannelSender

	ConfigStore *config.Store

	AlertStore        *alert.Store
	AlertLogStore     *alertlog.Store
	AlertMetricsStore *alertmetrics.Store

	AuthBasicStore        *basic.Store
	UserStore             *user.Store
	ContactMethodStore    *contactmethod.Store
	NotificationRuleStore *notificationrule.Store
	FavoriteStore         *favorite.Store

	ServiceStore        *service.Store
	EscalationStore     *escalation.Store
	IntegrationKeyStore *integrationkey.Store
	UIKHandler          *uik.Handler
	ScheduleRuleStore   *rule.Store
	NotificationStore   *notification.Store
	ScheduleStore       *schedule.Store
	RotationStore       *rotation.Store

	CalSubStore    *calsub.Store
	OverrideStore  *override.Store
	LimitStore     *limit.Store
	HeartbeatStore *heartbeat.Store

	OAuthKeyring    keyring.Keyring
	SessionKeyring  keyring.Keyring
	APIKeyring      keyring.Keyring
	AuthLinkKeyring keyring.Keyring

	NonceStore    *nonce.Store
	LabelStore    *label.Store
	OnCallStore   *oncall.Store
	NCStore       *notificationchannel.Store
	TimeZoneStore *timezone.Store
	NoticeStore   *notice.Store
	AuthLinkStore *authlink.Store
	APIKeyStore   *apikey.Store
}

// NewApp constructs a new App and binds the listening socket.
func NewApp(c Config, db *sql.DB) (*App, error) {
	var err error
	permission.SudoContext(context.Background(), func(ctx context.Context) {
		// Should not be possible for the app to ever see `use_next_db` unless misconfigured.
		//
		// In switchover mode, the connector wrapper will check this and provide the app with
		// a connection to the next DB instead, if this was set.
		//
		// This is a sanity check to ensure that the app is not accidentally using the previous DB
		// after a switchover.
		err = db.QueryRowContext(ctx, `select true from switchover_state where current_state = 'use_next_db'`).Scan(new(bool))
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		if err != nil {
			return
		}

		err = fmt.Errorf("refusing to connect to stale database (switchover_state table has use_next_db set)")
	})
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", c.ListenAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "bind address %s", c.ListenAddr)
	}

	if c.TLSListenAddr != "" {
		l2, err := tls.Listen("tcp", c.TLSListenAddr, c.TLSConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "listen %s", c.TLSListenAddr)
		}
		l = newMultiListener(c.Logger, l, l2)
	}

	c.Logger.AddErrorMapper(func(ctx context.Context, err error) context.Context {
		if e := sqlutil.MapError(err); e != nil && e.Detail != "" {
			ctx = log.WithField(ctx, "SQLErrDetails", e.Detail)
		}

		return ctx
	})

	app := &App{
		l:      l,
		db:     db,
		cfg:    c,
		doneCh: make(chan struct{}),
	}

	if c.StatusAddr != "" {
		err = listenStatus(c.StatusAddr, app.doneCh)
		if err != nil {
			return nil, errors.Wrap(err, "start status listener")
		}
	}

	app.db.SetMaxIdleConns(c.DBMaxIdle)
	app.db.SetMaxOpenConns(c.DBMaxOpen)

	app.mgr = lifecycle.NewManager(app._Run, app._Shutdown)
	err = app.mgr.SetStartupFunc(app.startup)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// WaitForStartup will wait until the startup sequence is completed or the context is expired.
func (a *App) WaitForStartup(ctx context.Context) error {
	return a.mgr.WaitForStartup(a.Context(ctx))
}

// DB returns the sql.DB instance used by the application.
func (a *App) DB() *sql.DB { return a.db }

// URL returns the non-TLS listener URL of the application.
func (a *App) URL() string {
	return "http://" + a.l.Addr().String()
}

func (a *App) SMTPAddr() string {
	if a.smtpsrvL == nil {
		return ""
	}

	return a.smtpsrvL.Addr().String()
}
