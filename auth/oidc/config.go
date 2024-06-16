package oidc

import (
	"github.com/breathbath/goalert/auth/nonce"
	"github.com/breathbath/goalert/keyring"
)

// Config provides necessary parameters for OIDC authentication.
type Config struct {
	Keyring    keyring.Keyring
	NonceStore *nonce.Store
}
