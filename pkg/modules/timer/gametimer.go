package timer

import (
	"fmt"
	"sync"
)

// GameTimer provides thread-safe timer for client time tracking.
//
// It maintains three related but distinct time values:
//   - clientTime: Current logical time that advances with real clock
//   - lastReturnedTime: Last value returned to ensure monotonicity
//   - lastPhysicsTime: Time for physics updates (33ms increments)
//
// The timer ensures that:
//   - Time never goes backward (monotonic)
//   - Physics time advances in fixed 33ms steps (30 FPS)
//   - Multiple goroutines can safely access time values
type GameTimer struct {
	// clientTime is the current logical time in milliseconds.
	// It advances with real time but can be set by server for sync.
	clientTime int64

	// lastTimestamp is the real system time (ms) when clientTime was last updated.
	// Used to calculate elapsed time between calls.
	lastTimestamp int64

	// lastReturnedTime is the last value returned to caller.
	// Ensures monotonicity by preventing backward jumps.
	lastReturnedTime int64

	// lastPhysicsTime is the time of last physics update.
	// Always a multiple of 33ms and never less than lastReturnedTime.
	lastPhysicsTime int64

	timeProvider TimeProvider
	mu           sync.RWMutex
}

func NewGameTimerWithProvider(tp TimeProvider) *GameTimer {
	return &GameTimer{
		clientTime:       0,
		lastReturnedTime: 0,
		lastPhysicsTime:  0,
		lastTimestamp:    tp.Now(),
		timeProvider:     tp,
	}
}

// NewGameTimer initializes GameTimer with RealTimeProvider
func NewGameTimer() *GameTimer {
	tp := &RealTimeProvider{}
	return &GameTimer{
		clientTime:       0,
		lastReturnedTime: 0,
		lastPhysicsTime:  0,
		lastTimestamp:    tp.Now(),
		timeProvider:     tp,
	}
}

// updateClientTime updates client (logical) time based on elapsed time.
// It is not thread-safe.
func (t *GameTimer) updateClientTime() {
	currentTimestamp := t.timeProvider.Now()
	elapsedTime := currentTimestamp - t.lastTimestamp

	t.clientTime += elapsedTime
	t.lastTimestamp = currentTimestamp

	if t.clientTime < t.lastReturnedTime {
		t.clientTime = t.lastReturnedTime
	}
}

// ClientTime gets current client time with thread-safe updates
func (t *GameTimer) ClientTime() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.updateClientTime()
	t.lastReturnedTime = t.clientTime
	return t.clientTime
}

// SetClientTime sets client time (thread-safe)
func (t *GameTimer) SetClientTime(value int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if value < t.lastReturnedTime {
		return
	}

	t.clientTime = value
	t.lastTimestamp = t.timeProvider.Now()
	t.lastReturnedTime = value

	// Also update physics time if needed
	if value >= t.lastPhysicsTime {
		t.lastPhysicsTime = value
	}
}

// PhysicsTime gets time with consistent 33ms intervals between calls
func (t *GameTimer) PhysicsTime() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.updateClientTime()
	timeValue := t.clientTime

	// Calculate next physics time based on 33ms increments from last physics time
	var result int64
	if t.lastPhysicsTime == 0 {
		// First physics update - just use current time
		result = timeValue
	} else {
		// Subsequent updates - maintain exact 33ms intervals
		timeDiff := timeValue - t.lastPhysicsTime
		steps := max((timeDiff+16)/33, 1)
		result = t.lastPhysicsTime + (steps * 33)
	}

	if result < t.lastReturnedTime {
		for result < t.lastReturnedTime {
			result += 33
		}
	}

	t.lastPhysicsTime = result
	t.lastReturnedTime = result
	return result
}

// SetPhysicsTime allows manual setting of physics time with validation
func (t *GameTimer) SetPhysicsTime(value int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Ensure value is a multiple of 33ms from the last physics time
	remainder := (value - t.lastPhysicsTime) % 33
	if remainder != 0 {
		return fmt.Errorf("physics time must be a multiple of 33ms: %d != 0", remainder)
	}

	if value < t.lastReturnedTime {
		return fmt.Errorf("physics time cannot go backward: %d < %d", value, t.lastReturnedTime)
	}

	t.lastPhysicsTime = value
	t.lastReturnedTime = value
	return nil
}

// Reset resets the timer to initial state
func (t *GameTimer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.clientTime = 0
	t.lastTimestamp = t.timeProvider.Now()
	t.lastReturnedTime = 0
	t.lastPhysicsTime = 0
}
