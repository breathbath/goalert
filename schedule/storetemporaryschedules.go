package schedule

import (
	"context"
	"database/sql"
	"time"

	"github.com/breathbath/goalert/permission"
	"github.com/breathbath/goalert/util/sqlutil"
	"github.com/breathbath/goalert/validation"
	"github.com/breathbath/goalert/validation/validate"
	"github.com/google/uuid"
)

// FixedShiftsPerTemporaryScheduleLimit is the maximum number of shifts that can be configured for a single TemporarySchedule at a time.
const FixedShiftsPerTemporaryScheduleLimit = 150

// TemporarySchedules will return the current set for the provided scheduleID.
func (store *Store) TemporarySchedules(ctx context.Context, tx *sql.Tx, scheduleID uuid.UUID) ([]TemporarySchedule, error) {
	err := permission.LimitCheckAny(ctx, permission.User)
	if err != nil {
		return nil, err
	}

	data, err := store.scheduleData(ctx, tx, scheduleID)
	if err != nil {
		return nil, err
	}

	check, err := store.usr.UserExists(ctx)
	if err != nil {
		return nil, err
	}

	// omit shifts for non-existent users
	for i, tmp := range data.V1.TemporarySchedules {
		shifts := tmp.Shifts[:0]
		for _, shift := range tmp.Shifts {
			if !check.UserExistsString(shift.UserID) {
				continue
			}
			shifts = append(shifts, shift)
		}
		tmp.Shifts = shifts
		data.V1.TemporarySchedules[i] = tmp
	}

	data.V1.TemporarySchedules = MergeTemporarySchedules(data.V1.TemporarySchedules)

	return data.V1.TemporarySchedules, nil
}

func isDataPkeyConflict(err error) bool {
	dbErr := sqlutil.MapError(err)
	if dbErr == nil {
		return false
	}
	return dbErr.ConstraintName == "schedule_data_pkey"
}

func validateFuture(fieldName string, t time.Time) error {
	if time.Until(t) > 5*time.Minute {
		return nil
	}
	return validation.NewFieldError(fieldName, "must be at least 5 min the future")
}

// SetTemporarySchedule will cause the schedule to use only, and exactly, the provided set of shifts between the provided start and end times.
func (store *Store) SetTemporarySchedule(ctx context.Context, tx *sql.Tx, scheduleID uuid.UUID, temp TemporarySchedule) error {
	err := permission.LimitCheckAny(ctx, permission.User)
	if err != nil {
		return err
	}

	check, err := store.usr.UserExists(ctx)
	if err != nil {
		return err
	}

	newTemp, err := temp.Normalize(check)
	if err != nil {
		return err
	}

	return store.updateScheduleData(ctx, tx, scheduleID, func(data *Data) error {
		data.V1.TemporarySchedules = setFixedShifts(data.V1.TemporarySchedules, *newTemp)
		return nil
	})
}

// SetClearTemporarySchedules works like SetTemporarySchedule after clearing out any existing TemporarySchedules between clearStart and clearEnd.
func (store *Store) SetClearTemporarySchedule(ctx context.Context, tx *sql.Tx, scheduleID uuid.UUID, temp TemporarySchedule, clearStart, clearEnd time.Time) error {
	err := permission.LimitCheckAny(ctx, permission.User)
	if err != nil {
		return err
	}

	check, err := store.usr.UserExists(ctx)
	if err != nil {
		return err
	}

	newTemp, err := temp.Normalize(check)
	if err != nil {
		return err
	}

	err = validateTimeRange("Clear", clearStart, clearEnd)
	if err != nil {
		return err
	}

	now := time.Now()
	if clearStart.Before(now) {
		clearStart = now
	}

	return store.updateScheduleData(ctx, tx, scheduleID, func(data *Data) error {
		data.V1.TemporarySchedules = deleteFixedShifts(data.V1.TemporarySchedules, clearStart, clearEnd)
		data.V1.TemporarySchedules = setFixedShifts(data.V1.TemporarySchedules, *newTemp)
		return nil
	})
}

// ClearTemporarySchedules will clear out (or split, if needed) any defined TemporarySchedules that exist between the start and end time.
func (store *Store) ClearTemporarySchedules(ctx context.Context, tx *sql.Tx, scheduleID uuid.UUID, start, end time.Time) error {
	err := permission.LimitCheckAny(ctx, permission.User)
	if err != nil {
		return err
	}

	err = validate.Many(
		validateFuture("End", end),
		validateTimeRange("", start, end),
	)
	if err != nil {
		return err
	}
	if time.Since(start) > 0 {
		start = time.Now()
	}

	return store.updateScheduleData(ctx, tx, scheduleID, func(data *Data) error {
		data.V1.TemporarySchedules = deleteFixedShifts(data.V1.TemporarySchedules, start, end)
		return nil
	})
}
