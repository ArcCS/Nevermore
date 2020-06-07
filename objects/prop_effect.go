package objects

import (
	"time"
)

type Effect struct{
	startTime time.Time
	length time.Duration

	lastTrigger time.Time
	interval time.Duration

	effect string
	effectOff string
}

func NewEffect(t time.Duration, length string, interval string,  effect string, effectOff string) *Effect {
	parseLength,_ := time.ParseDuration(length)
	parseInterval, _ := time.ParseDuration(interval)
	return &Effect{time.Now(),
		parseLength,
		time.Now(),
		parseInterval,
		effect,
		effectOff }
}

func (s *Effect) Reset(t time.Duration) {
	s.startTime = time.Now()
	s.length = t
}

func (s *Effect) TimeRemaining() float64 {
	calc := s.length - (time.Now().Sub(s.startTime))
	return calc.Minutes()
}
