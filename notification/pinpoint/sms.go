package pinpoint

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
	"github.com/breathbath/goalert/config"
	"github.com/breathbath/goalert/notification"
	"github.com/breathbath/goalert/util/log"
	"github.com/pkg/errors"
)

// SMS implements a notification.Sender for AWS Pinpoint SMS.
type SMS struct {
	c *Config
	r notification.Receiver

	limit *replyLimiter
}

var (
	_ notification.ReceiverSetter = &SMS{}
	_ notification.Sender         = &SMS{}
	_ notification.StatusChecker  = &SMS{}
)

// NewSMS performs operations like validating essential parameters, registering the AWS Pinpoint client and db
// and adding routes for successful and unsuccessful message delivery to AWS Pinpoint
func NewSMS(ctx context.Context, c *Config) (*SMS, error) {
	s := &SMS{
		c: c,

		limit: newReplyLimiter(),
	}

	return s, nil
}

// SetReceiver sets the notification.Receiver for incoming messages and status updates.
func (s *SMS) SetReceiver(r notification.Receiver) { s.r = r }

// Status provides the current status of a message.
func (s *SMS) Status(ctx context.Context, externalID string) (*notification.Status, error) {
	return &notification.Status{
		State: notification.StateUnknown, //as of today AWS Pinpoint doesn't have an API for checking the status of a message see https://docs.aws.amazon.com/pinpoint/latest/apireference_smsvoicev2/API_Operations.html
	}, nil
}

// Send implements the notification.Sender interface.
func (s *SMS) Send(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)
	if !cfg.PinPoint.Enable {
		return nil, errors.New("PinPoint provider is disabled")
	}
	if msg.Destination().Type != notification.DestTypeSMS {
		return nil, errors.Errorf("unsupported destination type %s; expected SMS", msg.Destination().Type)
	}
	destNumber := msg.Destination().Value
	if destNumber == cfg.PinPoint.OriginationIdentity {
		return nil, errors.New("refusing to send outgoing SMS to PhoneNumber, PhoneNumberId, PhoneNumberArn, SenderId, SenderIdArn, PoolId, or PoolArn.")
	}

	ctx = log.WithFields(ctx, log.Fields{
		"Phone": destNumber,
		"Type":  "PinpointSMS",
	})

	makeSMSCode := func(alertID int, serviceID string) int {
		return 0
	}

	var message string
	var err error
	switch t := msg.(type) {
	case notification.AlertStatus:
		message, err = renderAlertStatusMessage(cfg.ApplicationName(), t)
	case notification.AlertBundle:
		var link string
		if canContainURL(ctx, destNumber) {
			link = cfg.CallbackURL(fmt.Sprintf("/services/%s/alerts", t.ServiceID))
		}

		message, err = renderAlertBundleMessage(cfg.ApplicationName(), t, link, makeSMSCode(0, t.ServiceID))
	case notification.Alert:
		var link string
		if canContainURL(ctx, destNumber) {
			link = cfg.CallbackURL(fmt.Sprintf("/alerts/%d", t.AlertID))
		}

		message, err = renderAlertMessage(cfg.ApplicationName(), t, link, makeSMSCode(t.AlertID, ""))
	case notification.Test:
		message = fmt.Sprintf("%s: Test message.", cfg.ApplicationName())
	case notification.Verification:
		message = fmt.Sprintf("%s: Verification code: %d", cfg.ApplicationName(), t.Code)
	default:
		return nil, errors.Errorf("unhandled message type %T", t)
	}
	if err != nil {
		return nil, errors.Wrap(err, "render message")
	}

	cl, err := NewClient(ctx, cfg, s.c)
	if err != nil {
		return nil, errors.Wrap(err, "wrong settings")
	}

	params := &pinpointsmsvoicev2.SendTextMessageInput{
		DestinationPhoneNumber:       &destNumber,
		ConfigurationSetName:         &cfg.PinPoint.ConfigurationSetName,
		Context:                      cfg.PinPoint.Context,
		DestinationCountryParameters: cfg.PinPoint.DestinationCountryParameters,
		Keyword:                      &cfg.PinPoint.Keyword,
		MaxPrice:                     &cfg.PinPoint.MaxPrice,
		MessageBody:                  &message,
		MessageType:                  types.MessageTypeTransactional,
		OriginationIdentity:          &cfg.PinPoint.OriginationIdentity,
		ProtectConfigurationId:       &cfg.PinPoint.ProtectConfigurationId,
		TimeToLive:                   &cfg.PinPoint.TimeToLive,
	}

	// Actually send notification to end user & receive Message Status
	resp, err := cl.SendTextMessage(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "send message")
	}

	// If the message was sent successfully, reset reply limits.
	s.limit.Reset(destNumber)

	return &notification.SentMessage{
		ExternalID: *resp.MessageId,
		State:      notification.StateSent,
	}, nil
}
