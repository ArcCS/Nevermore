package objects

import (
	"strconv"
	"time"
)

type Effect struct {
	startTime time.Time
	length    time.Duration

	lastTrigger time.Time
	interval    int
	triggers    int
	effect      func(triggers int)
	effectOff   func()
	magnitude   int
}

func (e *Effect) AlterTime(duration float64) {
	e.length = time.Duration(duration * float64(time.Second))
}

func (e *Effect) ExtendDuration(duration float64) {
	calc := e.length - (e.length - (time.Now().Sub(e.startTime)))
	e.length = time.Duration(duration)*time.Second - time.Duration(calc.Seconds())
}

func NewEffect(length string, interval int, magnitude int, effect func(triggers int), effectOff func()) *Effect {
	lengthTime, _ := strconv.Atoi(length)
	parseLength := time.Duration(lengthTime) * time.Second
	return &Effect{time.Now(),
		parseLength,
		time.Now(),
		interval,
		0,
		effect,
		effectOff,
		magnitude}
}

func (e *Effect) RunEffect() {
	e.effect(e.triggers)
	e.triggers += 1
	e.lastTrigger = time.Now()
}

func (e *Effect) Reset(t string) {
	lengthTime, _ := strconv.Atoi(t)
	parseLength := time.Duration(lengthTime) * time.Second
	e.startTime = time.Now()
	e.length = parseLength
	e.triggers = 0
}

func (e *Effect) TimeRemaining() float64 {
	calc := e.length - (time.Now().Sub(e.startTime))
	return calc.Seconds()
}

func (e *Effect) LastTriggerInterval() int {
	lTrigger := time.Now().Sub(e.lastTrigger)
	calc := e.interval - int(lTrigger.Seconds())
	return calc
}

// Function to return only the modifiable properties
func (e *Effect) ReturnEffectProps() map[string]interface{} {
	serialList := map[string]interface{}{
		"timeRemaining": e.TimeRemaining(),
		"magnitude":     e.magnitude,
	}
	return serialList
}
