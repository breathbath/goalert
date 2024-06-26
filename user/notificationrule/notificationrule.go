package notificationrule

import (
	"github.com/breathbath/goalert/validation/validate"
)

type NotificationRule struct {
	ID              string `json:"id"`
	UserID          string `json:"-"`
	DelayMinutes    int    `json:"delay"`
	ContactMethodID string `json:"contact_method_id"`
}

func validateDelay(d int) error {
	return validate.Range("DelayMinutes", d, 0, 9000)
}

func (n NotificationRule) Normalize(update bool) (*NotificationRule, error) {
	err := validateDelay(n.DelayMinutes)

	if !update {
		err = validate.Many(
			err,
			validate.UUID("ContactMethodID", n.ContactMethodID),
			validate.UUID("UserID", n.UserID),
		)
	}
	if err != nil {
		return nil, err
	}

	return &n, nil
}
