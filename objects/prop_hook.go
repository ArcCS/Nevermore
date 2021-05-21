package objects

import (
	"strconv"
	"time"
)

type Hook struct {
	remainingExecutions int
	startTime time.Time
	length    time.Duration

	lastTrigger time.Time
	interval    time.Duration

	effect    func()
	effectOff func()
}

func NewHook(executions int, length string, interval string, effect func(), effectOff func()) *Hook {
	lengthTime, _ := strconv.Atoi(length)
	parseLength := time.Duration(lengthTime) * time.Second
	parseInterval, _ := time.ParseDuration(interval)
	return &Hook{ executions,
		time.Now(),
		parseLength,
		time.Now(),
		parseInterval,
		effect,
		effectOff}
}

func (s *Hook) RunHook(){
	s.effect()
	s.lastTrigger = time.Now()
}

func (s *Hook) Reset(t time.Duration) {
	s.startTime = time.Now()
	s.length = t
}

func (s *Hook) TimeRemaining() float64 {
	calc := s.length - (time.Now().Sub(s.startTime))
	return calc.Minutes()
}

func (s *Hook) LastTriggerInterval() float64 {
	calc := s.interval - (time.Now().Sub(s.lastTrigger))
	return calc.Minutes()
}
