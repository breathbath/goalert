package graphqlapp

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/breathbath/goalert/calsub"
	"github.com/breathbath/goalert/config"
	"github.com/breathbath/goalert/graphql2"
	"github.com/breathbath/goalert/permission"
	"github.com/breathbath/goalert/schedule"
)

type UserCalendarSubscription App

func (a *App) UserCalendarSubscription() graphql2.UserCalendarSubscriptionResolver {
	return (*UserCalendarSubscription)(a)
}

func (a *UserCalendarSubscription) ReminderMinutes(ctx context.Context, obj *calsub.Subscription) ([]int, error) {
	return obj.Config.ReminderMinutes, nil
}
func (a *UserCalendarSubscription) FullSchedule(ctx context.Context, obj *calsub.Subscription) (bool, error) {
	return obj.Config.FullSchedule, nil
}

func (a *UserCalendarSubscription) Schedule(ctx context.Context, obj *calsub.Subscription) (*schedule.Schedule, error) {
	return a.ScheduleStore.FindOne(ctx, obj.ScheduleID)
}

func (a *UserCalendarSubscription) URL(ctx context.Context, obj *calsub.Subscription) (*string, error) {
	tok := obj.Token()
	if tok == "" {
		return nil, nil
	}

	v := make(url.Values)
	v.Set("token", tok)

	cfg := config.FromContext(ctx)
	callback := cfg.CallbackURL("/api/v2/calendar", v)
	return &callback, nil
}

func (q *Query) UserCalendarSubscription(ctx context.Context, id string) (*calsub.Subscription, error) {
	return q.CalSubStore.FindOne(ctx, id)
}

// todo: return UserCalendarSubscription with generated url once endpoint has been created
func (m *Mutation) CreateUserCalendarSubscription(ctx context.Context, input graphql2.CreateUserCalendarSubscriptionInput) (cs *calsub.Subscription, err error) {
	cs = &calsub.Subscription{
		Name:       input.Name,
		ScheduleID: input.ScheduleID,
		UserID:     permission.UserID(ctx),
	}
	if input.Disabled != nil {
		cs.Disabled = *input.Disabled
	}
	cs.Config.ReminderMinutes = input.ReminderMinutes
	if input.FullSchedule != nil {
		cs.Config.FullSchedule = *input.FullSchedule
	}
	err = withContextTx(ctx, m.DB, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		cs, err = m.CalSubStore.CreateTx(ctx, tx, cs)
		return err
	})

	return cs, err
}

func (m *Mutation) UpdateUserCalendarSubscription(ctx context.Context, input graphql2.UpdateUserCalendarSubscriptionInput) (bool, error) {
	err := withContextTx(ctx, m.DB, func(ctx context.Context, tx *sql.Tx) error {
		cs, err := m.CalSubStore.FindOneForUpdate(ctx, tx, input.ID)
		if err != nil {
			return err
		}

		if input.Name != nil {
			cs.Name = *input.Name
		}
		if input.Disabled != nil {
			cs.Disabled = *input.Disabled
		}
		if input.ReminderMinutes != nil {
			cs.Config.ReminderMinutes = input.ReminderMinutes
		}
		if input.FullSchedule != nil {
			cs.Config.FullSchedule = *input.FullSchedule
		}

		return m.CalSubStore.UpdateTx(ctx, tx, cs)
	})

	return err == nil, err
}
