package app

import (
	"context"
	"github.com/breathbath/goalert/notification/pinpoint"

	"github.com/breathbath/goalert/notification"
	"github.com/pkg/errors"
)

func (app *App) initPinpoint(ctx context.Context) error {
	app.pinpointConfig = &pinpoint.Config{
		BaseURL: app.cfg.PinpointBaseURL,
		CMStore: app.ContactMethodStore,
		DB:      app.db,
	}

	var err error
	app.pinpointSMS, err = pinpoint.NewSMS(ctx, app.pinpointConfig)
	if err != nil {
		return errors.Wrap(err, "init PinpointSMS")
	}
	app.notificationManager.RegisterSender(notification.DestTypeSMS, "Pinpoint-SMS", app.pinpointSMS)

	return nil
}
