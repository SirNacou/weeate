package clock

import "time"

type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
}

type RealClock struct {
}

func NewRealClock() *RealClock {
	return &RealClock{}
}

func (rc *RealClock) Now() time.Time {
	return time.Now()
}

func (rc *RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (rc *RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
