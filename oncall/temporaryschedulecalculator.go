package oncall

import (
	"sort"

	"github.com/breathbath/goalert/schedule"
)

// TemporaryScheduleCalculator will calculate active state and active users for a set of TemporarySchedules.
type TemporaryScheduleCalculator struct {
	*TimeIterator

	act *ActiveCalculator
	usr *UserCalculator
}

// NewTemporaryScheduleCalculator will create a new TemporaryScheduleCalculator bound to the TimeIterator.
func (t *TimeIterator) NewTemporaryScheduleCalculator(tempScheds []schedule.TemporarySchedule) *TemporaryScheduleCalculator {
	ts := &TemporaryScheduleCalculator{
		TimeIterator: t,
		act:          t.NewActiveCalculator(),
		usr:          t.NewUserCalculator(),
	}

	sort.Slice(tempScheds, func(i, j int) bool { return tempScheds[i].Start.Before(tempScheds[j].Start) })
	var allShifts []schedule.FixedShift

	for _, temp := range tempScheds {
		ts.act.SetSpan(temp.Start, temp.End)
		allShifts = append(allShifts, temp.Shifts...)
	}
	sort.Slice(allShifts, func(i, j int) bool { return allShifts[i].Start.Before(allShifts[j].Start) })

	for _, s := range allShifts {
		ts.usr.SetSpan(s.Start, s.End, s.UserID)
	}
	ts.act.Init()
	ts.usr.Init()

	return ts
}

// Active will return true if a TemporarySchedule is currently active.
func (fg *TemporaryScheduleCalculator) Active() bool { return fg.act.Active() }

// ActiveUsers will return the current set of ActiveUsers. It is only valid if `Active()` is true.
//
// It is only valid if `Active()` is true and until the following Next() call. It should not be modified.
func (fg *TemporaryScheduleCalculator) ActiveUsers() []string { return fg.usr.ActiveUsers() }
