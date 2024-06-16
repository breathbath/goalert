package app

import (
	"crypto/tls"
	"time"

	"github.com/breathbath/goalert/config"
	"github.com/breathbath/goalert/expflag"
	"github.com/breathbath/goalert/keyring"
	"github.com/breathbath/goalert/swo"
	"github.com/breathbath/goalert/util/log"
)

type Config struct {
	Logger *log.Logger

	ExpFlags expflag.FlagSet

	ListenAddr  string
	Verbose     bool
	JSON        bool
	LogRequests bool
	APIOnly     bool
	LogEngine   bool

	PublicURL string

	TLSListenAddr string
	TLSConfig     *tls.Config

	SysAPIListenAddr string
	SysAPICertFile   string
	SysAPIKeyFile    string
	SysAPICAFile     string

	SMTPListenAddr        string
	SMTPListenAddrTLS     string
	SMTPMaxRecipients     int
	TLSConfigSMTP         *tls.Config
	SMTPAdditionalDomains string

	EmailIntegrationDomain string

	HTTPPrefix string

	DBMaxOpen int
	DBMaxIdle int

	MaxReqBodyBytes   int64
	MaxReqHeaderBytes int

	DisableHTTPSRedirect bool

	TwilioBaseURL     string
	SlackBaseURL      string
	PinpointAWSRegion string

	DBURL     string
	DBURLNext string

	StatusAddr string

	EngineCycleTime time.Duration

	EncryptionKeys keyring.Keys

	RegionName string

	StubNotifiers bool

	UIDir string

	// InitialConfig will be pushed into the config store
	// if specified before the engine is started.
	InitialConfig *config.Config

	// SWO should be set to operate in switchover mode.
	SWO *swo.Manager
}
