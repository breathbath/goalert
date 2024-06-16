package graphql2

import (
	"github.com/breathbath/goalert/assignment"
	"github.com/breathbath/goalert/schedule"
)

type OnCallNotificationRuleInput struct {
	schedule.OnCallNotificationRule
	Target assignment.RawTarget
}
