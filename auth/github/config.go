package github

import (
	"github.com/breathbath/goalert/auth/nonce"
	"github.com/breathbath/goalert/keyring"
)

// Config is used to configure the GitHub OAuth2 provider. If none of Organization, Teams, or Users are
// specified as criteria, any valid user will be accepted.
type Config struct {
	Keyring    keyring.Keyring
	NonceStore *nonce.Store
}
