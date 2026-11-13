package objects

import (
	"strings"
	"testing"
	"time"

	"github.com/ArcCS/Nevermore/config"
)

func TestEverySongHasAnEffectAndATune(t *testing.T) {
	for name, s := range Songs {
		if _, _, ok := SongEffectFor(name); !ok {
			t.Errorf("song %q has no effect (effect key %q)", name, s["effect"])
		}
		if s["tune"] == "" {
			t.Errorf("song %q has no tune", name)
		}
	}
}

// songRoom puts the given characters in a fresh room and returns its id.
func songRoom(t *testing.T, id int, chars ...*Character) *Room {
	r := &Room{RoomId: id, Flags: map[string]bool{}, Chars: &CharInventory{ParentId: id, Contents: chars}, Mobs: &MobInventory{ParentId: id}}
	Rooms[id] = r
	t.Cleanup(func() { delete(Rooms, id) })
	for _, c := range chars {
		c.ParentId = id
	}
	return r
}

// songBard is a bard holding an instrument that always rolls 55.
func songBard(name string, pie int) *Character {
	b := &Character{
		Object:    Object{Name: name},
		Flags:     map[string]bool{},
		Modifiers: map[string]int{},
		Timers:    map[string]time.Time{},
		Effects:   map[string]*Effect{},
		Equipment: &Equipment{Off: &Item{ItemType: 16, NumDice: 1, SidesDice: 1, PlusDice: 54, Flags: map[string]bool{}}},
		Class:     config.BARD,
		Tier:      10,
	}
	b.Pie.Current = pie
	b.Mana.Max = 1000
	return b
}

func songListener(name string) *Character {
	c := &Character{Object: Object{Name: name}, Flags: map[string]bool{}, Modifiers: map[string]int{}, Effects: map[string]*Effect{}}
	c.Mana.Max = 1000
	return c
}

func TestDraensTaleGivesEveryoneTheSameMana(t *testing.T) {
	bard := songBard("Lala", 18)
	a, b := songListener("Ann"), songListener("Bob")
	songRoom(t, -7101, bard, a, b)
	bard.CurrentSong = "draens-tale"

	bard.songPulse()

	// round(0.1*55 + 18/2) = round(14.5) = 15
	for _, c := range []*Character{bard, a, b} {
		if c.Mana.Current != 15 {
			t.Errorf("%s got %d mana, want 15", c.Name, c.Mana.Current)
		}
	}
}

func TestMixSongPulsesTheNewSongAtOnce(t *testing.T) {
	bard := songBard("Lala", 18)
	ally := songListener("Ann")
	songRoom(t, -7102, bard, ally)
	bard.CurrentSong = "curious-canticle"

	bard.MixSong("draens-tale")

	if bard.CurrentSong != "draens-tale" {
		t.Fatalf("current song = %q, want draens-tale", bard.CurrentSong)
	}
	if ally.Mana.Current != 15 {
		t.Errorf("mix should pulse at once: ally mana = %d, want 15", ally.Mana.Current)
	}
}

func TestRoomListShowsTheTuneAndFollowsAMix(t *testing.T) {
	bard := songBard("Lala", 18)
	bard.Flags["singing"] = true
	bard.CurrentSong = "celebration-night"
	observer := songListener("Ann")
	inv := &CharInventory{Contents: []*Character{bard, observer}}

	if got := strings.Join(inv.List(observer), ","); got != "Lala (singing a joyous tune)" {
		t.Errorf("room list = %q", got)
	}
	bard.Pose = "sitting"
	if got := strings.Join(inv.ListPoses(observer), ","); got != "Lala is sitting and singing a joyous tune." {
		t.Errorf("pose line = %q", got)
	}

	bard.CurrentSong = "banshees-lament"
	if got := strings.Join(inv.ListPoses(observer), ","); got != "Lala is sitting and singing a shrill dirge." {
		t.Errorf("pose line after mix = %q", got)
	}
}

func TestSongDescriptorFallsBack(t *testing.T) {
	c := &Character{}
	if got := c.SongDescriptor(); got != "a song" {
		t.Errorf("silent descriptor = %q, want \"a song\"", got)
	}
}

func TestRefundTimer(t *testing.T) {
	c := &Character{Timers: map[string]time.Time{}}
	c.Timers["combat_sing"] = time.Now().Add(300 * time.Second)

	c.RefundTimer("combat_sing", 270)
	if left := c.TimerRemaining("combat_sing"); left < 29 || left > 31 {
		t.Errorf("after refunding 270 of 300s, %v left, want ~30", left)
	}

	c.RefundTimer("combat_sing", 1000)
	if left := c.TimerRemaining("combat_sing"); left > 0.1 {
		t.Errorf("a refund past now should stop at now, %v left", left)
	}

	// No timer (staff never get one) stays unset.
	c.RefundTimer("combat_sing_mix", 50)
	if _, ok := c.Timers["combat_sing_mix"]; ok {
		t.Error("refund created a timer that was never set")
	}
}

func TestRunAwayNoLongerStopsFollowing(t *testing.T) {
	const fromId, toId = -7103, -7104
	target := songListener("Ann")
	from := songRoom(t, fromId, target)
	to := songRoom(t, toId)
	to.Flags["active"] = true

	followed := false
	for range 50 {
		m := &Mob{Object: Object{Name: "goblin"}, Flags: map[string]bool{"follows": true, "run_away": true},
			CurrentTarget: "Ann", Stam: Meter{Max: 10, Current: 10}}
		from.Mobs.Contents = []*Mob{m}
		m.ParentId = fromId
		if m.FollowChar(target, from, to) {
			followed = true
			break
		}
	}
	if !followed {
		t.Error("a run_away mob should still be able to follow")
	}

	m := &Mob{Object: Object{Name: "goblin"}, Flags: map[string]bool{"follows": true, "curious_canticle": true},
		CurrentTarget: "Ann", Stam: Meter{Max: 10, Current: 10}, ParentId: fromId}
	from.Mobs.Contents = []*Mob{m}
	if m.FollowChar(target, from, to) {
		t.Error("a curious_canticle mob should not follow")
	}
}

func TestLevelCapIgnoresInstrumentCap(t *testing.T) {
	lute := &Item{Object: Object{Name: "lute"}, ItemType: 16, SidesDice: 100, NumDice: 2, Flags: map[string]bool{}}
	below := &Character{Class: config.BARD, Tier: config.LevelCap - 1}
	if ok, _ := below.CanEquip(lute); ok {
		t.Error("over-cap instrument should be refused below LevelCap")
	}
	capped := &Character{Class: config.BARD, Tier: config.LevelCap}
	if ok, msg := capped.CanEquip(lute); !ok {
		t.Errorf("over-cap instrument should be allowed at LevelCap: %s", msg)
	}
}
