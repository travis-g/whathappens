package whathappens

import (
	"errors"
	"time"
)

// ErrNotImplemented is an error returned when planned functionality is not yet
// implemented.
var ErrNotImplemented = errors.New("not yet implemented")

// ElapsedSince returns the time since a start time in a floating-point number
// of milliseconds.
func ElapsedSince(start time.Time) float32 {
	now := time.Now()
	return float32(now.Sub(start).Nanoseconds()) / float32(time.Millisecond)
}
