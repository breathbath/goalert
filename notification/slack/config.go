package slack

import (
	"github.com/breathbath/goalert/user"
)

// Config contains values used for the Slack notification sender.
type Config struct {
	BaseURL   string
	UserStore *user.Store
}
