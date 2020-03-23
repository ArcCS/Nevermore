package objects

import (
	"time"
)

// Todo: Effects are really timers with sub-tickers in some cases.
// Use should include putting a timer for a spell that expires and removes any toggles
// Or it expires and removes anything ticking.
// Ticking may have environmental impact, capture here

type Effect struct{
	timer *time.Timer
	timeEnd time.Time
	// Does it do something?
	repetition int64
	effectFunc string
	callback string
}

func NewEffect(t time.Duration, rep int64, effect string, callback string) *Effect {
	return &Effect{time.NewTimer(t), time.Now().Add(t), rep, effect, callback}
}

func (s *Effect) Reset(t time.Duration) {
	s.timer.Reset(t)
	s.timeEnd = time.Now().Add(t)
}

func (s *Effect) Stop() {
	s.timer.Stop()
}

func (s *Effect) TimeRemaining() time.Duration {
	return s.timeEnd.Sub(time.Now())
}