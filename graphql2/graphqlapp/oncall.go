package graphqlapp

import (
	context "context"
	"github.com/breathbath/goalert/graphql2"
	"github.com/breathbath/goalert/oncall"
	"github.com/breathbath/goalert/user"
)

type OnCallShift App

func (a *App) OnCallShift() graphql2.OnCallShiftResolver { return (*OnCallShift)(a) }

func (oc *OnCallShift) User(ctx context.Context, raw *oncall.Shift) (*user.User, error) {
	return (*App)(oc).FindOneUser(ctx, raw.UserID)
}
