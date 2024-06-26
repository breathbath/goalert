package app

import (
	"context"

	"github.com/breathbath/goalert/expflag"
	"github.com/breathbath/goalert/util/log"
)

// Context returns a new context with the App's configuration for
// experimental flags and logger.
//
// It should be used for calls from other packages to ensure that
// the correct configuration is used.
func (app *App) Context(ctx context.Context) context.Context {
	ctx = expflag.Context(ctx, app.cfg.ExpFlags)
	ctx = log.WithLogger(ctx, app.cfg.Logger)

	if app.ConfigStore != nil {
		ctx = app.ConfigStore.Config().Context(ctx)
	}

	return ctx
}
