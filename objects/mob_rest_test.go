package objects

import (
	"testing"
	"time"

	"github.com/ArcCS/Nevermore/config"
)

// A permanent mob heals for the minutes its room sat empty and re-arms the
// HP triggers it is back above.
func TestMobRest(t *testing.T) {
	m := &Mob{Stam: Meter{Max: 1000, Current: 100}, CurrentTarget: "Leaver"}
	m.StoppedAt = time.Now().Add(-3 * time.Minute)
	st := m.scriptState()
	for _, threshold := range []int{75, 50, 25, 10} {
		st.hpFired[threshold] = true
	}
	st.windup = &mobWindup{}

	m.Rest()

	want := 100 + 3*1000*config.PermMobRegenPctPerMinute/100
	if m.Stam.Current < want-1 || m.Stam.Current > want+1 {
		t.Errorf("after 3 minutes HP is %d, want about %d", m.Stam.Current, want)
	}
	if !st.hpFired[75] || !st.hpFired[50] || st.hpFired[25] || st.hpFired[10] {
		t.Errorf("hpFired = %v, want 75 and 50 still fired, 25 and 10 re-armed", st.hpFired)
	}
	if st.windup != nil {
		t.Error("a windup survived the room emptying")
	}
	if m.CurrentTarget != "" {
		t.Errorf("still targeting %q after the room emptied", m.CurrentTarget)
	}

	m.StoppedAt = time.Now().Add(-time.Hour)
	m.Rest()
	if m.Stam.Current != m.Stam.Max {
		t.Errorf("after an hour HP is %d, want full %d", m.Stam.Current, m.Stam.Max)
	}
}
