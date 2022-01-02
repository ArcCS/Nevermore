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
	interval    int

	effect    func()
	effectOff func()
}

//NewHook Use -1 for unlimited executions, create a new hook.
func NewHook(executions int, length string, interval int, effect func(), effectOff func()) *Hook {
	lengthTime, _ := strconv.Atoi(length)
	parseLength := time.Duration(lengthTime) * time.Second
	return &Hook{ executions,
		time.Now(),
		parseLength,
		time.Now(),
		interval,
		effect,
		effectOff}
}

func (s *Hook) RunHook(){
	s.effect()
	s.lastTrigger = time.Now()
	if s.remainingExecutions >= 1 {
		s.remainingExecutions -= 1
	}
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
	return float64(s.interval) - time.Now().Sub(s.lastTrigger).Minutes()
}
