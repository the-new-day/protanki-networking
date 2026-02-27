package timer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockTimeProvider struct {
	currentTime int64
	mu          sync.Mutex
}

func (m *MockTimeProvider) Now() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentTime
}

func (m *MockTimeProvider) Advance(ms int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTime += ms
}

func TestNewTankiTimer_InitializesCorrectly(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	assert.Equal(t, int64(0), timer.clientTime)
	assert.Equal(t, int64(1000), timer.lastTimestamp)
	assert.Equal(t, int64(0), timer.lastReturnedTime)
	assert.Equal(t, int64(0), timer.lastPhysicsTime)
}

func TestClientTime_IncreasesWithRealTime(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	first := timer.ClientTime()
	assert.Equal(t, int64(0), first)

	mock.Advance(500)
	second := timer.ClientTime()
	assert.Equal(t, int64(500), second)

	mock.Advance(300)
	third := timer.ClientTime()
	assert.Equal(t, int64(800), third)
}

func TestClientTime_NeverGoesBackward(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	values := []int64{}
	for i := 0; i < 10; i++ {
		mock.Advance(50)
		values = append(values, timer.ClientTime())
	}

	for i := 1; i < len(values); i++ {
		assert.True(t, values[i] > values[i-1])
	}
}

func TestClientTime_SetToFutureValue(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	timer.SetClientTime(5000)
	assert.Equal(t, int64(5000), timer.clientTime)
	assert.Equal(t, int64(1000), timer.lastTimestamp)

	mock.Advance(200)
	_ = timer.ClientTime()
	assert.Equal(t, int64(1200), timer.lastTimestamp)
}

func TestClientTime_IgnoresPastValue(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	timer.SetClientTime(5000)
	timer.SetClientTime(4000)
	assert.Equal(t, int64(5000), timer.clientTime)
}

func TestPhysicsTime_FirstCallReturnsClientTime(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	mock.Advance(150)
	result := timer.PhysicsTime()
	assert.Equal(t, int64(150), result)
}

func TestPhysicsTime_SecondCallReturnsCalculatedTime(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)
	mock.Advance(1) // imitate real timer (1ms)

	timer.PhysicsTime()

	mock.Advance(40)

	second := timer.PhysicsTime()
	assert.Equal(t, int64(34), second)
}

func TestPhysicsTime_StepsBy33ms(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)
	mock.Advance(1)

	first := timer.PhysicsTime()
	assert.Equal(t, int64(1), first)

	second := timer.PhysicsTime()
	assert.Equal(t, int64(34), second)

	third := timer.PhysicsTime()
	assert.Equal(t, int64(67), third)

	fourth := timer.PhysicsTime()
	assert.Equal(t, int64(100), fourth)
}

func TestPhysicsTime_AccumulatesMultipleSteps(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)
	mock.Advance(1)

	first := timer.PhysicsTime()
	assert.Equal(t, int64(1), first)

	mock.Advance(100)
	second := timer.PhysicsTime()
	assert.Equal(t, int64(100), second)
}

func TestPhysicsTime_HandlesLargeClientTimeJump(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)
	mock.Advance(1)

	first := timer.PhysicsTime()
	assert.Equal(t, int64(1), first)

	mock.currentTime = 2000
	timer.SetClientTime(2000)

	second := timer.PhysicsTime()
	assert.Equal(t, int64(2033), second)
}

func TestSetPhysicsTime_ValidMultipleSucceeds(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	base := timer.PhysicsTime()

	err := timer.SetPhysicsTime(base + 33)
	assert.NoError(t, err)
	assert.Equal(t, base+33, timer.lastPhysicsTime)

	err = timer.SetPhysicsTime(base + 66)
	assert.NoError(t, err)
	assert.Equal(t, base+66, timer.lastPhysicsTime)
}

func TestSetPhysicsTime_InvalidMultipleFails(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	base := timer.PhysicsTime()
	err := timer.SetPhysicsTime(base + 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple of 33ms")
}

func TestSetPhysicsTime_BackwardValueFails(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	base := timer.PhysicsTime()
	err := timer.SetPhysicsTime(base - 33)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot go backward")
}

func TestReset_ClearsAllCounters(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)

	mock.Advance(500)
	timer.ClientTime()
	timer.PhysicsTime()

	timer.Reset()

	assert.Equal(t, int64(0), timer.clientTime)
	assert.Equal(t, int64(1500), timer.lastTimestamp)
	assert.Equal(t, int64(0), timer.lastReturnedTime)
	assert.Equal(t, int64(0), timer.lastPhysicsTime)

	mock.Advance(100)
	newTime := timer.ClientTime()
	assert.Equal(t, int64(100), newTime)
}

func TestConcurrentAccess_NoDeadlocks(t *testing.T) {
	mock := &MockTimeProvider{currentTime: 1000}
	timer := NewTankiTimerWithProvider(mock)
	done := make(chan bool)

	for range 10 {
		go func() {
			for range 100 {
				timer.ClientTime()
				timer.PhysicsTime()
				mock.Advance(1)
			}
			done <- true
		}()
	}

	for range 10 {
		<-done
	}
}
