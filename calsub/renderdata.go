package calsub

import (
	"time"

	"github.com/breathbath/goalert/oncall"
	"github.com/google/uuid"
)

type renderData struct {
	ApplicationName string
	ScheduleID      uuid.UUID
	ScheduleName    string
	Shifts          []oncall.Shift
	ReminderMinutes []int
	Version         string
	GeneratedAt     time.Time
	FullSchedule    bool
	UserNames       map[string]string
}
