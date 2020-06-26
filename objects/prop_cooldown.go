package objects

import "time"

type Cooldown struct {
 timeStart time.Time
 end   time.Time
}

func NewCooldown(seconds int) *Cooldown {
 return &Cooldown{time.Now(), time.Now().Add(time.Duration(seconds) * time.Second)}
}

func (s *Cooldown) Reset(t int) {
 s.timeStart = time.Now()
 s.end = time.Now().Add(time.Duration(t) * time.Second)
}

func (s *Cooldown) TimeRemaining() float64 {
 return s.end.Sub(time.Now()).Seconds()
}