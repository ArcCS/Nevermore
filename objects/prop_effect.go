package objects

import (
	"strconv"
	"time"
)

type Effect struct {
	startTime time.Time
	length    time.Duration

	lastTrigger time.Time
	interval    time.Duration

	effect    func()
	effectOff func()
}

func (s *Effect) AlterTime(duration float64) {
	s.length = time.Duration(duration)*time.Minute
}

func NewEffect(length string, interval string, effect func(), effectOff func()) *Effect {
	lengthTime, _ := strconv.Atoi(length)
	parseLength := time.Duration(lengthTime) * time.Second
	parseInterval, _ := time.ParseDuration(interval)
	return &Effect{time.Now(),
		parseLength,
		time.Now(),
		parseInterval,
		effect,
		effectOff}
}

func (s *Effect) RunEffect(){
	s.effect()
	s.lastTrigger = time.Now()
}

func (s *Effect) Reset(t time.Duration) {
	s.startTime = time.Now()
	s.length = t
}

func (s *Effect) TimeRemaining() float64 {
	calc := s.length - (time.Now().Sub(s.startTime))
	return calc.Minutes()
}

func (s *Effect) LastTriggerInterval() float64 {
	calc := s.interval - (time.Now().Sub(s.lastTrigger))
	return calc.Minutes()
}
