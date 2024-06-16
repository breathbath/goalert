package notificationchannel

import (
	"strings"

	"github.com/breathbath/goalert/validation/validate"
	"github.com/google/uuid"
)

type Channel struct {
	ID    string
	Name  string
	Type  Type
	Value string
}

func (c Channel) Normalize() (*Channel, error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	err := validate.Many(
		validate.UUID("ID", c.ID),
		validate.Text("Name", c.Name, 1, 255),
		validate.OneOf("Type", c.Type, TypeSlackChan, TypeWebhook, TypeSlackUG),
	)

	switch c.Type {
	case TypeSlackUG:
		grp, ch, _ := strings.Cut(c.Value, ":")
		err = validate.Many(err,
			validate.RequiredText("Value.GroupID", grp, 1, 32),
			validate.RequiredText("Value.ChannelID", ch, 1, 32),
		)
	case TypeSlackChan:
		err = validate.Many(err, validate.RequiredText("Value", c.Value, 1, 32))
	case TypeWebhook:
		err = validate.Many(err, validate.URL("Value", c.Value))
	}

	return &c, err
}
