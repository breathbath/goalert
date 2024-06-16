package pinpoint

import (
	"database/sql"
	"github.com/breathbath/goalert/user/contactmethod"
)

// Config contains the details needed to interact with AWS Pinpoint for SMS
type Config struct {

	// Region choose the AWS Region from which you will be sending messages.
	Region string

	// CMStore is used for storing and fetching metadata (like carrier information).
	CMStore *contactmethod.Store

	// DB is used for storing DB connection data (needed for carrier metadata dbtx).
	DB *sql.DB
}
