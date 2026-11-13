package config

import (
	"math"
	"testing"
)

func TestSongDuration(t *testing.T) {
	if got := SongDuration(1); got != 182 {
		t.Errorf("tier 1 song = %d, want 182", got)
	}
	if got := SongDuration(LevelCap); got != 230 {
		t.Errorf("tier %d song = %d, want 230", LevelCap, got)
	}
	// A song never outlasts the cooldown or the mix timer, so a bard can't
	// sing forever or mix twice in one song.
	for tier := 1; tier <= LevelCap+1; tier++ {
		if d := SongDuration(tier); d >= SongCooldown || d >= SongMixCooldown {
			t.Errorf("tier %d song lasts %ds, not under the %ds cooldown / %ds mix timer", tier, d, SongCooldown, SongMixCooldown)
		}
	}
}

func TestSongTickSeconds(t *testing.T) {
	for tier, want := range map[int]int{1: 8, 14: 8, 15: 6, 25: 6} {
		if got := SongTickSeconds(tier); got != want {
			t.Errorf("tier %d tick = %ds, want %ds", tier, got, want)
		}
	}
}

func TestSongStopRefundSeconds(t *testing.T) {
	cases := []struct {
		remaining, length, want float64
	}{
		{206, 206, 270}, // stopped at once: 90% of the cooldown back
		{103, 206, 135}, // halfway
		{0, 206, 0},     // song already over
		{-5, 206, 0},    // effect overdue for removal
		{500, 206, 270}, // clamped to the song length
		{100, 0, 0},     // no recorded length
	}
	for _, c := range cases {
		if got := SongStopRefundSeconds(c.remaining, c.length); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("refund(%v of %v) = %v, want %v", c.remaining, c.length, got, c.want)
		}
	}
}

// Stopping early, over and over, never keeps a song up more of the time than
// singing every song to its end.
func TestSongStopNeverBeatsFullSongUptime(t *testing.T) {
	for tier := 1; tier <= LevelCap; tier++ {
		length := float64(SongDuration(tier))
		full := length / float64(SongCooldown)
		for sung := 1.0; sung <= length; sung++ {
			readyAt := float64(SongCooldown) - SongStopRefundSeconds(length-sung, length)
			if uptime := sung / readyAt; uptime > full+1e-9 {
				t.Fatalf("tier %d stop at %vs: uptime %.3f beats a full song's %.3f", tier, sung, uptime, full)
			}
		}
	}
}

func TestDraensTaleMana(t *testing.T) {
	cases := []struct{ roll, pie, want int }{
		{55, 18, 15}, // 5.5 + 9 = 14.5, rounds up
		{31, 18, 12}, // 3.1 + 9
		{0, 18, 9},
		{103, 0, 10},
		{-50, 0, 0}, // a negative adjustment never drains mana
	}
	for _, c := range cases {
		if got := DraensTaleMana(c.roll, c.pie); got != c.want {
			t.Errorf("DraensTaleMana(%d, %d) = %d, want %d", c.roll, c.pie, got, c.want)
		}
	}
}

func TestCelebrationHeal(t *testing.T) {
	cases := []struct{ roll, players, want int }{
		{100, 1, 85},
		{100, 2, 54}, // 85 / 2^(2/3) = 53.5
		{100, 5, 29}, // 85 / 5^(2/3) = 29.1
		{100, 6, 28}, // divisor capped at 3
		{100, 10, 28},
		{100, 0, 85}, // nobody counted reads as the bard alone
		{2, 6, 1},    // small rolls round instead of truncating to 0
	}
	for _, c := range cases {
		if got := CelebrationHeal(c.roll, c.players); got != c.want {
			t.Errorf("CelebrationHeal(%d, %d) = %d, want %d", c.roll, c.players, got, c.want)
		}
	}
}

func TestBardSingCastMod(t *testing.T) {
	if got := BardSingCastMod(BARD, 1, true); got != .8 {
		t.Errorf("tier 1 bard singing = %v, want 0.8", got)
	}
	if got := BardSingCastMod(BARD, BardSingCastFreeTier-1, true); got != .8 {
		t.Errorf("tier %d bard singing = %v, want 0.8", BardSingCastFreeTier-1, got)
	}
	if got := BardSingCastMod(BARD, BardSingCastFreeTier, true); got != 1 {
		t.Errorf("tier %d bard singing = %v, want 1", BardSingCastFreeTier, got)
	}
	if got := BardSingCastMod(BARD, 1, false); got != 1 {
		t.Errorf("bard not singing = %v, want 1", got)
	}
	if got := BardSingCastMod(CLERIC, 1, true); got != 1 {
		t.Errorf("non-bard = %v, want 1", got)
	}
}
