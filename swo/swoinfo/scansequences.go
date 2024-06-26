package swoinfo

import (
	"context"
	"sort"

	"github.com/breathbath/goalert/swo/swodb"
	"github.com/jackc/pgx/v5"
)

// ScanSequences will return the names of all sequences in the database, ordered
// by name.
func ScanSequences(ctx context.Context, conn *pgx.Conn) ([]string, error) {
	names, err := swodb.New(conn).SequenceNames(ctx)
	if err != nil {
		return nil, err
	}

	sort.Strings(names)
	return names, nil
}
