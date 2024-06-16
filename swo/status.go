package swo

import (
	"github.com/breathbath/goalert/swo/swogrp"
	"github.com/google/uuid"
)

// Status represents the current status of the switchover process.
type Status struct {
	swogrp.Status

	MainDBID      uuid.UUID
	NextDBID      uuid.UUID
	MainDBVersion string
	NextDBVersion string
}
