package auth

import (
	"github.com/breathbath/goalert/apikey"
	"github.com/breathbath/goalert/calsub"
	"github.com/breathbath/goalert/integrationkey"
	"github.com/breathbath/goalert/keyring"
	"github.com/breathbath/goalert/user"
)

// HandlerConfig provides configuration for the auth handler.
type HandlerConfig struct {
	UserStore      *user.Store
	SessionKeyring keyring.Keyring
	APIKeyring     keyring.Keyring
	IntKeyStore    *integrationkey.Store
	CalSubStore    *calsub.Store
	APIKeyStore    *apikey.Store
}
