package engine

import (
	"time"

	"github.com/breathbath/goalert/alert"
	"github.com/breathbath/goalert/alert/alertlog"
	"github.com/breathbath/goalert/auth/authlink"
	"github.com/breathbath/goalert/config"
	"github.com/breathbath/goalert/keyring"
	"github.com/breathbath/goalert/notification"
	"github.com/breathbath/goalert/notification/slack"
	"github.com/breathbath/goalert/notificationchannel"
	"github.com/breathbath/goalert/oncall"
	"github.com/breathbath/goalert/schedule"
	"github.com/breathbath/goalert/user"
	"github.com/breathbath/goalert/user/contactmethod"
)

// Config contains parameters for controlling how the Engine operates.
type Config struct {
	AlertLogStore       *alertlog.Store
	AlertStore          *alert.Store
	ContactMethodStore  *contactmethod.Store
	NotificationManager *notification.Manager
	UserStore           *user.Store
	NotificationStore   *notification.Store
	NCStore             *notificationchannel.Store
	OnCallStore         *oncall.Store
	ScheduleStore       *schedule.Store
	AuthLinkStore       *authlink.Store
	SlackStore          *slack.ChannelSender

	ConfigSource config.Source

	Keys keyring.Keys

	MaxMessages int

	DisableCycle bool
	LogCycles    bool

	CycleTime time.Duration
}
