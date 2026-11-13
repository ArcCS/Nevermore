package config

import (
	"math"
	"testing"
)

func TestCalcManaMonkScalesOnPiety(t *testing.T) {
	if got := CalcMana(1, 30, 10, MONK); got != 7 {
		t.Errorf("tier 1 pie 10 = %d, want 7", got)
	}
	if got := CalcMana(25, 5, 25, MONK); got != 287 {
		t.Errorf("tier 25 pie 25 = %d, want 287", got)
	}
	// Intelligence must not matter for monks.
	if CalcMana(10, 5, 15, MONK) != CalcMana(10, 40, 15, MONK) {
		t.Error("monk chi changed with intelligence")
	}
	// Other classes keep the old formula and ignore piety.
	want := (10 * Classes["mage"].Mana) + int(10*20*ManaPerStatPerTier)
	if got := CalcMana(10, 20, 99, MAGE); got != want {
		t.Errorf("mage mana = %d, want %d", got, want)
	}
}

func TestTodKillChance(t *testing.T) {
	cases := []struct {
		tier       int
		c, h, want float64
	}{
		{10, 1, 1, 0},
		{10, 0, 0.5, 0},
		{10, 1, 0.9, 0.02},
		{10, 1, 0.7, 0.18},
		{10, 1, 0.5, 0.5},
		{10, 1, 0.4, 0.72},
		{10, 1, 0.35, 0.845},
		{10, 1, 0.34, TodMaxChance},
		{10, 1, 0, TodMaxChance},
		// Chi is linear: half the commitment, half the chance, and no
		// threshold where the mob's health outruns the chi.
		{10, 0.5, 0.5, 0.25},
		{10, 0.5, 0.9, 0.01},
		// The coefficient grows with tier: 2.5 at 20, 2.75 at 25.
		{20, 1, 0.5, 0.625},
		{25, 1, 0.5, 0.6875},
		{25, 1, 0.44, TodMaxChance},
	}
	for _, tc := range cases {
		if got := TodKillChance(tc.tier, tc.c, tc.h); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("TodKillChance(%d, %v, %v) = %v, want %v", tc.tier, tc.c, tc.h, got, tc.want)
		}
	}
	if TodChanceScaleFor(5) != TodChanceScale || TodChanceScaleFor(10) != TodChanceScale {
		t.Error("scale grew below the touch tier")
	}
	if math.Abs(TodChanceScaleFor(25)-2.75) > 1e-9 {
		t.Errorf("scale at 25 = %v, want 2.75", TodChanceScaleFor(25))
	}
}

func TestTodReferenceGrows(t *testing.T) {
	prev := TodReference(1)
	for tier := 2; tier <= 25; tier++ {
		if r := TodReference(tier); r <= prev {
			t.Fatalf("TodReference(%d) = %d not above %d", tier, r, prev)
		} else {
			prev = r
		}
	}
	// R sits below even a low-piety monk's full bar, so they keep a little
	// back after a full-chi touch, and the surplus never shrinks as they
	// level; a typical monk keeps a larger surplus for flurry.
	prevSurplus := 0
	for tier := MonkTodTier; tier <= 25; tier++ {
		surplus := MonkMaxChi(tier, 10) - TodReference(tier)
		if surplus <= 0 || surplus < prevSurplus {
			t.Errorf("tier %d: low-piety surplus over R = %d (previous %d)", tier, surplus, prevSurplus)
		}
		prevSurplus = surplus
		if MonkMaxChi(tier, 12+tier/4) <= TodReference(tier) {
			t.Errorf("tier %d: typical monk has no chi surplus over R", tier)
		}
	}
}

func TestTodCooldown(t *testing.T) {
	cases := map[int]int{10: 300, 15: 300, 16: 276, 19: 204, 20: 180, 25: 180}
	for tier, want := range cases {
		if got := TodCooldown(tier); got != want {
			t.Errorf("TodCooldown(%d) = %d, want %d", tier, got, want)
		}
	}
}

func TestMonkFlurryFollowsFighterTable(t *testing.T) {
	adept := MultiAttackMultipliers[FlurryMinSkill]
	for skill := 0; skill < FlurryMinSkill; skill++ {
		got := MonkFlurryFor(skill)
		if len(got) != len(adept) || got[0] != adept[0] || got[1] != adept[1] {
			t.Errorf("skill %d flurry = %v, want the Adept row %v", skill, got, adept)
		}
	}
	for skill := FlurryMinSkill; skill <= 10; skill++ {
		want := MultiAttackMultipliers[skill]
		got := MonkFlurryFor(skill)
		if len(got) != len(want) {
			t.Errorf("skill %d flurry = %v, want %v", skill, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("skill %d flurry = %v, want %v", skill, got, want)
				break
			}
		}
	}
	if got := MonkFlurryFor(10); len(got) != 5 {
		t.Errorf("Grandmaster flurry = %v, want all five swings", got)
	}
	if len(MonkFlurryFor(8)) <= len(MonkFlurryFor(5)) {
		t.Error("Specialist flurry should out-swing Adept")
	}
}

func TestFlurryChiCostScalesWithSkill(t *testing.T) {
	cases := map[int]int{0: 6, 6: 6, 7: 8, 8: 8, 9: 8, 10: 8}
	for skill, want := range cases {
		if got := FlurryChiCost(skill); got != want {
			t.Errorf("FlurryChiCost(%d) = %d, want %d", skill, got, want)
		}
	}
}

func TestMonkDodgeChance(t *testing.T) {
	base := MonkDodgeChance(0, 0, false)
	if base != MonkDodgeBase {
		t.Errorf("base dodge = %d, want %d", base, MonkDodgeBase)
	}
	if MonkDodgeChance(20, 5, false) <= base {
		t.Error("dex and skill did not raise dodge")
	}
	plain := MonkDodgeChance(20, 5, false)
	if MonkDodgeChance(20, 5, true) != plain+FeintDodgeBonus {
		t.Error("feint bonus not applied")
	}
	if MonkDodgeChance(45, 10, true) != MonkDodgeCap {
		t.Errorf("dodge not capped: %d", MonkDodgeChance(45, 10, true))
	}
}

func TestChiHelpers(t *testing.T) {
	if ChiDecayAmount(5) != 1 {
		t.Errorf("decay on a tiny pool = %d, want 1", ChiDecayAmount(5))
	}
	if ChiDecayAmount(100) != 100*ChiDecayPercent/100 {
		t.Errorf("decay on 100 = %d", ChiDecayAmount(100))
	}
	// The piety share is floored at 1, so low piety still earns 4.
	if ChiPerHitFor(0) != 4 || ChiPerHitFor(4) != 4 {
		t.Errorf("chi per hit at low piety = %d / %d, want 4", ChiPerHitFor(0), ChiPerHitFor(4))
	}
	if ChiPerHitFor(15) != 6 || ChiPerHitFor(25) != 8 {
		t.Errorf("chi per hit at pie 15 / 25 = %d / %d, want 6 / 8", ChiPerHitFor(15), ChiPerHitFor(25))
	}
	if ChiGainCap(5) != 1 {
		t.Errorf("gain cap on a tiny pool = %d, want 1", ChiGainCap(5))
	}
	if ChiGainCap(100) != 100*ChiGainCapPercent/100 {
		t.Errorf("gain cap on 100 = %d", ChiGainCap(100))
	}
}

func TestMeditateHelpers(t *testing.T) {
	if MeditateRestorePercent(1, 12) != MeditateBasePercent+12/MeditatePieDiv+1 {
		t.Errorf("restore percent at tier 1 pie 12 = %d", MeditateRestorePercent(1, 12))
	}
	if MeditateDuration(25) != MeditateBaseDuration+25*MeditateDurationPerTier {
		t.Errorf("duration at tier 25 = %d", MeditateDuration(25))
	}
}

func TestMeditateCharge(t *testing.T) {
	cases := map[float64]float64{-5: 0, 0: 0, 30: 0.1, 150: 0.5, 300: 1, 900: 1}
	for elapsed, want := range cases {
		if got := MeditateCharge(elapsed); got != want {
			t.Errorf("MeditateCharge(%v) = %v, want %v", elapsed, got, want)
		}
	}
	// Tier 10's full trance is 110s: 11s at the minimum charge, 55s at half.
	if got := MeditateTranceFor(10, 0.1); got != 11 {
		t.Errorf("tier 10 trance at 10%% = %d, want 11", got)
	}
	if got := MeditateTranceFor(10, 0.5); got != 55 {
		t.Errorf("tier 10 trance at half = %d, want 55", got)
	}
	if got := MeditateTranceFor(10, 0); got != 1 {
		t.Errorf("trance never drops below a second, got %d", got)
	}
}

func TestMonkReachesGrandmaster(t *testing.T) {
	if WeaponLevel(SkillExpLevels[10], MONK, HandSkill) != 10 {
		t.Error("monk capped below weapon level 10")
	}
	if WeaponExpTitle(SkillExpLevels[10], MONK, HandSkill) != WeaponTitles[10] {
		t.Error("monk denied the grandmaster title")
	}
	if WeaponExpNext(SkillExpLevels[9], MONK, HandSkill) != SkillExpLevels[10] {
		t.Error("monk next-level lookup stops at 9")
	}
	if WeaponLevel(SkillExpLevels[10], THIEF, 0) != 9 {
		t.Error("thief unexpectedly reaches level 10")
	}
}

func TestMonkUnarmedRange(t *testing.T) {
	// Tier 9, str 20: base 27, str share 12, roll 3d9, less the lighter 4 -> 38..62.
	if lo, hi := MonkUnarmedRange(9, 20); lo != 38 || hi != 62 {
		t.Errorf("tier 9 str 20 range = %d..%d, want 38..62", lo, hi)
	}
	// Tier 10, str 20: base 30, str share 14, roll 3d10, cut still 4 -> 43..70.
	if lo, hi := MonkUnarmedRange(10, 20); lo != 43 || hi != 70 {
		t.Errorf("tier 10 str 20 range = %d..%d, want 43..70", lo, hi)
	}
	// Tier 14, str 20: base 40, str share 18, roll 3d13, cut tapered to 2 -> 59..95.
	if lo, hi := MonkUnarmedRange(14, 20); lo != 59 || hi != 95 {
		t.Errorf("tier 14 str 20 range = %d..%d, want 59..95", lo, hi)
	}
	// Tier 21, str 20: base 60, str share 27, roll 3d20, no cut -> 90..147.
	if lo, hi := MonkUnarmedRange(21, 20); lo != 90 || hi != 147 {
		t.Errorf("tier 21 str 20 range = %d..%d, want 90..147", lo, hi)
	}
	// Tier 20, str 30: base 55, str share 37, roll 3d18, no cut -> 95..146.
	if lo, hi := MonkUnarmedRange(20, 30); lo != 95 || hi != 146 {
		t.Errorf("tier 20 str 30 range = %d..%d, want 95..146", lo, hi)
	}
	for tier, want := range map[int]int{9: 4, 11: 4, 12: 3, 16: 1, 17: 1, 18: 0, 25: 0} {
		if got := MonkDamageCut(tier); got != want {
			t.Errorf("tier %d cut = %d, want %d", tier, got, want)
		}
	}
	// The flat cut never drives a hit below the floor.
	if lo, _ := MonkUnarmedRange(1, 0); lo < MonkDamageFloor {
		t.Errorf("tier 1 str 0 floor = %d, want at least %d", lo, MonkDamageFloor)
	}
	prevHi := 0
	for tier := 1; tier <= 25; tier++ {
		_, hi := MonkUnarmedRange(tier, 20)
		if hi <= prevHi {
			t.Fatalf("tier %d max %d not above tier %d", tier, hi, tier-1)
		}
		prevHi = hi
	}
}

func TestMonkNaturalArmor(t *testing.T) {
	// Tier 20, con 20: base 300 + 40 = 340, plus 10% from con = 374.
	if got := MonkNaturalArmor(20, 20); got != 374 {
		t.Errorf("tier 20 con 20 armor = %d, want 374", got)
	}
	// Con 0 leaves only the tier term.
	if got := MonkNaturalArmor(10, 0); got != 10*MonkArmorPerLevel {
		t.Errorf("tier 10 con 0 armor = %d, want %d", got, 10*MonkArmorPerLevel)
	}
	if MonkNaturalArmor(10, 30) <= MonkNaturalArmor(10, 10) {
		t.Error("more con should mean more armor")
	}
}

func TestMonkIronBody(t *testing.T) {
	if MonkVitalReduction(14) != 0 || MonkCriticalReduction(14) != 0 {
		t.Error("iron body applied below its tier")
	}
	cases := []struct {
		tier            int
		vital, critical float64
	}{
		{15, 0.15, 0.225},
		{25, 0.25, 0.375},
	}
	for _, tc := range cases {
		if got := MonkVitalReduction(tc.tier); math.Abs(got-tc.vital) > 1e-9 {
			t.Errorf("vital reduction at %d = %v, want %v", tc.tier, got, tc.vital)
		}
		if got := MonkCriticalReduction(tc.tier); math.Abs(got-tc.critical) > 1e-9 {
			t.Errorf("critical reduction at %d = %v, want %v", tc.tier, got, tc.critical)
		}
	}
}

func TestRangerWeaponRules(t *testing.T) {
	if WeaponLevel(SkillExpLevels[10], RANGER, MissileSkill) != 10 {
		t.Error("ranger denied grandmaster with a bow")
	}
	if WeaponLevel(SkillExpLevels[10], RANGER, 0) != 9 {
		t.Error("ranger reached grandmaster with a sword")
	}
	if WeaponAdvancementFor(RANGER, MissileSkill) != Classes["ranger"].WeaponAdvancement {
		t.Error("ranger missile advancement should be the class value")
	}
	if WeaponAdvancementFor(RANGER, 0) != RangerMeleeAdvancement {
		t.Error("ranger melee advancement should be the reduced rate")
	}
	if WeaponAdvancementFor(FIGHTER, 0) != Classes["fighter"].WeaponAdvancement {
		t.Error("non-ranger advancement changed")
	}
}

func TestMeditateSaveChance(t *testing.T) {
	// Tier 10 vs a level 10 mob with 20 con: base plus half con.
	if got := MeditateSaveChance(10, 10, 20); got != 60 {
		t.Errorf("t10 vs l10 con20 = %d, want 60", got)
	}
	// Five tiers over the mob and five past the gate: +10 and +5.
	if got := MeditateSaveChance(15, 10, 20); got != 75 {
		t.Errorf("t15 vs l10 con20 = %d, want 75", got)
	}
	// Outleveled: the diff cuts both ways.
	if got := MeditateSaveChance(10, 15, 20); got != 50 {
		t.Errorf("t10 vs l15 con20 = %d, want 50", got)
	}
	if got := MeditateSaveChance(25, 10, 30); got != MeditateSaveCap {
		t.Errorf("chance not capped: %d", got)
	}
	if got := MeditateSaveChance(10, 40, 0); got != 0 {
		t.Errorf("chance not floored: %d", got)
	}
}

func TestChiDecayFloorIsTierCappedAtTen(t *testing.T) {
	cases := map[int]int{1: 1, 5: 5, 10: 10, 15: 10, 25: 10}
	for tier, want := range cases {
		if got := ChiDecayFloor(tier); got != want {
			t.Errorf("ChiDecayFloor(%d) = %d, want %d", tier, got, want)
		}
	}
}

func TestTodMissPenaltyFor(t *testing.T) {
	// 30 less tier/2 less pie/3: 25 at tier 10 with no piety, 13 at tier 15 with pie 32.
	cases := []struct{ tier, pie, want int }{{10, 0, 25}, {10, 12, 21}, {15, 32, 13}, {20, 15, 15}}
	for _, c := range cases {
		if got := TodMissPenaltyFor(c.tier, c.pie); got != c.want {
			t.Errorf("tier %d pie %d penalty = %d, want %d", c.tier, c.pie, got, c.want)
		}
	}
}

func TestMonkMagicResistPercent(t *testing.T) {
	cases := map[int]int{0: 10, 15: 13, 20: 14, 44: 18, 45: 19, 60: 19}
	for con, want := range cases {
		if got := MonkMagicResistPercent(con); got != want {
			t.Errorf("MonkMagicResistPercent(%d) = %d, want %d", con, got, want)
		}
	}
}

func TestChiWard(t *testing.T) {
	cases := []struct{ damage, percent, chi, absorbed, cost int }{
		{20, 14, 50, 2, 0},   // 1-4 turned aside is free
		{40, 14, 50, 5, 1},   // 5 rounds up to a chi
		{100, 14, 50, 14, 1}, // 14 rounds down
		{150, 19, 50, 28, 3}, // 28 rounds up
		{150, 19, 1, 14, 1},  // one chi covers at most 14
		{150, 19, 0, 0, 0},   // no chi, no ward
		{0, 19, 50, 0, 0},
	}
	for _, c := range cases {
		absorbed, cost := ChiWard(c.damage, c.percent, c.chi)
		if absorbed != c.absorbed || cost != c.cost {
			t.Errorf("ChiWard(%d, %d, %d) = %d, %d; want %d, %d", c.damage, c.percent, c.chi, absorbed, cost, c.absorbed, c.cost)
		}
	}
}
