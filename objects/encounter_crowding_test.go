package objects

import (
	"testing"
	"time"
)

func TestEncounterCrowding(t *testing.T) {
	cases := []struct {
		chars, mobs int
		wantDouble  bool
		wantPct     int
	}{
		{1, 0, true, 100},
		{1, 1, true, 100},
		{3, 3, true, 100},
		{1, 2, false, 88},
		{1, 3, false, 75},
		{1, 4, false, 63},
		{1, 5, false, 50},
		{1, 9, false, 50},
		{4, 5, false, 88},
		{4, 8, false, 50},
	}
	for _, c := range cases {
		gotDouble, gotPct := EncounterCrowding(c.chars, c.mobs)
		if gotDouble != c.wantDouble || gotPct != c.wantPct {
			t.Errorf("EncounterCrowding(%d chars, %d mobs) = (%v, %d), want (%v, %d)",
				c.chars, c.mobs, gotDouble, gotPct, c.wantDouble, c.wantPct)
		}
	}
}

func TestDelayFirstEncounterComesDueAfterDelay(t *testing.T) {
	r := &Room{EncounterSpeed: 10, EncounterJigger: 4}
	for _, delay := range []int{0, 5, 9} {
		r.DelayFirstEncounter(delay)
		due := r.LastEncounterTime.Add(time.Duration(r.EncounterSpeed+r.EncounterJigger) * time.Second)
		if left := time.Until(due); left > time.Duration(delay)*time.Second || left < time.Duration(delay)*time.Second-time.Second {
			t.Errorf("delay %d: next roll due in %v", delay, left)
		}
	}
}
