package rotation

import (
	"time"

	"github.com/breathbath/goalert/validation/validate"
)

type State struct {
	RotationID    string
	ParticipantID string
	Position      int
	ShiftStart    time.Time
}

func (s State) Normalize() (*State, error) {
	err := validate.Many(
		validate.UUID("ParticipantID", s.ParticipantID),
		validate.Range("Position", s.Position, 0, 9000),
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
