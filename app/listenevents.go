package app

import (
	"context"

	"github.com/breathbath/goalert/permission"
	"github.com/breathbath/goalert/util/log"
	"github.com/breathbath/goalert/util/sqlutil"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

func (app *App) listenEvents(ctx context.Context) (<-chan struct{}, error) {
	l, err := sqlutil.NewListener(ctx, app.cfg.Logger, app.db, "/goalert/config-refresh")
	if err != nil {
		return nil, err
	}
	app.events = l
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-l.Errors():
				log.Log(ctx, errors.Wrap(err, "listen events"))
			}
		}
	}()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		for {
			var n *pgconn.Notification
			select {
			case n = <-l.Notifications():
			case <-ctx.Done():
				return
			}
			if n == nil {
				return
			}

			log.Debugf(log.WithFields(ctx, log.Fields{
				"Channel": n.Channel,
				"PID":     n.PID,
				"Payload": n.Payload,
			}), "NOTIFY")

			switch n.Channel {
			case "/goalert/config-refresh":
				permission.SudoContext(ctx, func(ctx context.Context) {
					log.Log(ctx, app.ConfigStore.Reload(ctx))
				})
			}
		}
	}()

	return doneCh, nil
}
