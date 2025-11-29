package clock

import (
	"sync"
	"time"
)

// MockClock is a controllable clock implementation for testing
type MockClock struct {
	mu       sync.RWMutex
	now      time.Time
	channels []chan time.Time
}

func NewMockClock(initialTime time.Time) *MockClock {
	return &MockClock{
		now:      initialTime,
		channels: make([]chan time.Time, 0),
	}
}

func (mc *MockClock) Now() time.Time {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.now
}

func (mc *MockClock) Sleep(d time.Duration) {
	mc.Advance(d)
}

func (mc *MockClock) After(d time.Duration) <-chan time.Time {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	ch := make(chan time.Time, 1)
	mc.channels = append(mc.channels, ch)

	// Spawn goroutine to wait for the duration and send the time
	go func(duration time.Duration, targetTime time.Time) {
		for {
			mc.mu.RLock()
			currentTime := mc.now
			mc.mu.RUnlock()

			if !currentTime.Before(targetTime) {
				ch <- currentTime
				close(ch)
				return
			}
			time.Sleep(10 * time.Millisecond) // Poll interval
		}
	}(d, mc.now.Add(d))

	return ch
}

// Advance moves the mock clock forward by the given duration
// and triggers any waiting After channels that should fire
func (mc *MockClock) Advance(d time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.now = mc.now.Add(d)
}

// Set sets the mock clock to a specific time
func (mc *MockClock) Set(t time.Time) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.now = t
}
