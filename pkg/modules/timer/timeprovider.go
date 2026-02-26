package timer

import "time"

// TimeProvider is an interface used by timers for getting current time.
type TimeProvider interface {
	// Returns milliseconds
	Now() int64
}

// RealTimeProvider implements TimeProvider interface using built-in time package.
type RealTimeProvider struct{}

// Returns current time in milliseconds using time.Now().
func (tp *RealTimeProvider) Now() int64 {
	return time.Now().UnixMilli()
}
