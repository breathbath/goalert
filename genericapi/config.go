package genericapi

import (
	"github.com/breathbath/goalert/alert"
	"github.com/breathbath/goalert/heartbeat"
	"github.com/breathbath/goalert/integrationkey"
	"github.com/breathbath/goalert/user"
)

// Config contains the values needed to implement the generic API handler.
type Config struct {
	AlertStore          *alert.Store
	IntegrationKeyStore *integrationkey.Store
	HeartbeatStore      *heartbeat.Store
	UserStore           *user.Store
}
