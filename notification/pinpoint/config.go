package pinpoint

import (
	"database/sql"
	"github.com/breathbath/goalert/user/contactmethod"
	"net/http"
)

// Config contains the details needed to interact with AWS Pinpoint for SMS
type Config struct {

	// Client is an optional net/http client to use, if nil the global default is used.
	Client *http.Client

	// BaseURL can be used to override the base AWS API URL.
	BaseURL string

	// CMStore is used for storing and fetching metadata (like carrier information).
	CMStore *contactmethod.Store

	// DB is used for storing DB connection data (needed for carrier metadata dbtx).
	DB *sql.DB
}
