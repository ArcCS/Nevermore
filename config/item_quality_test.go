package config

import "testing"

func TestWeaponQualityBands(t *testing.T) {
	cases := []struct {
		name       string
		roll       WeaponRoll
		twoHanded  bool
		alwaysCrit bool
		special    bool
		want       int
	}{
		// Tier 10 one-handed: real cap 59, target 80, cuts at avg 46.02/47.79/51.33.
		{"t10 under good", WeaponRoll{1, 19, 36}, false, false, false, QualityUnremarkable},
		{"t10 good", WeaponRoll{1, 18, 37}, false, false, false, QualityGood},
		{"t10 top of good", WeaponRoll{1, 17, 38}, false, false, false, QualityGood},
		{"t10 excellent", WeaponRoll{1, 19, 38}, false, false, false, QualityExcellent},
		{"t10 top of excellent", WeaponRoll{1, 15, 43}, false, false, false, QualityExcellent},
		{"t10 exceptional", WeaponRoll{1, 15, 44}, false, false, false, QualityExceptional},
		// The flag bonus moves a weapon up once it crosses a cut, and only then.
		{"t10 flag lifts to good", WeaponRoll{1, 19, 36}, false, false, true, QualityGood},
		{"t10 flag lifts to exceptional", WeaponRoll{1, 15, 43}, false, false, true, QualityExceptional},
		// Tier 5 target 78: real cap 34, cuts at avg 25.84/26.86/28.9.
		{"t5 under good", WeaponRoll{1, 11, 19}, false, false, false, QualityUnremarkable},
		{"t5 exceptional", WeaponRoll{1, 9, 25}, false, false, false, QualityExceptional},
		// Tier 15 target 82: real cap 84, cuts at avg 67.2/69.72/74.76.
		{"t15 under good", WeaponRoll{1, 30, 50}, false, false, false, QualityUnremarkable},
		{"t15 good", WeaponRoll{1, 30, 52}, false, false, false, QualityGood},
		// Tier 20 target 84: real cap 109, cuts at avg 89.38/92.65/99.19.
		{"t20 under good", WeaponRoll{1, 40, 65}, false, false, false, QualityUnremarkable},
		{"t20 excellent", WeaponRoll{1, 30, 78}, false, false, false, QualityExcellent},
		// Tier 10 two-handed: real cap 69, cuts at avg 53.82/55.89/60.03.
		{"2h t10 good", WeaponRoll{1, 22, 43}, true, false, false, QualityGood},
		{"2h t10 exceptional", WeaponRoll{1, 17, 52}, true, false, false, QualityExceptional},
		{"always crit", WeaponRoll{2, 4, 34}, false, true, false, QualityExceptional},
		{"past every cap", WeaponRoll{100, 100, 100}, false, false, false, QualityExceptional},
	}
	for _, c := range cases {
		if got := WeaponQuality(c.roll, c.twoHanded, c.alwaysCrit, c.special); got != c.want {
			t.Errorf("%s: %dd%d+%d got %s, want %s", c.name, c.roll.Dice, c.roll.Sides, c.roll.Plus, QualityNames[got], QualityNames[c.want])
		}
	}
}

func TestWeaponTierForMatchesCanWield(t *testing.T) {
	for _, class := range []int{FIGHTER, BARBARIAN} {
		for _, twoHanded := range []bool{false, true} {
			for max := 1; max <= 200; max++ {
				tier := WeaponTierFor(max, twoHanded, class)
				if tier == 0 || tier > LevelCap {
					if CanWield(LevelCap, class, max, twoHanded) {
						t.Errorf("class %d 2h %v max %d: no tier line, but wieldable at the level cap", class, twoHanded, max)
					}
					continue
				}
				if !CanWield(tier, class, max, twoHanded) {
					t.Errorf("class %d 2h %v max %d: reads tier %d but can't wield", class, twoHanded, max, tier)
				}
				if tier > 1 && CanWield(tier-1, class, max, twoHanded) {
					t.Errorf("class %d 2h %v max %d: reads tier %d but wieldable a tier sooner", class, twoHanded, max, tier)
				}
			}
		}
	}
}

func TestArmorQuality(t *testing.T) {
	cases := []struct {
		name  string
		aType int
		val   int
		want  int
	}{
		{"legs at cap", 20, 60, QualityExceptional},
		{"legs 1 under", 20, 59, QualityExcellent},
		{"legs 3 under", 20, 57, QualityExcellent},
		{"legs 4 under", 20, 56, QualityGood},
		{"legs standard", 20, 55, QualityGood},
		{"legs under standard", 20, 54, QualityUnremarkable},
		{"body at cap", 5, 90, QualityExceptional},
		{"body 1 under", 5, 89, QualityExcellent},
		{"body 5 under", 5, 85, QualityExcellent},
		{"body 6 under", 5, 84, QualityGood},
		{"body standard", 5, 80, QualityGood},
		{"body under standard", 5, 79, QualityUnremarkable},
		{"shield 5 under", 23, 35, QualityExcellent},
		{"ring at cap", 24, 8, QualityExceptional},
		{"ring 1 under", 24, 7, QualityGood},
		{"ring 2 under", 24, 6, QualityExceptional}, // 6 is the band below's cap
		{"no armor", 20, 0, QualityUnremarkable},
		{"past every cap", 5, 131, QualityExceptional},
	}
	for _, c := range cases {
		if got := ArmorQuality(c.aType, c.val); got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, QualityNames[got], QualityNames[c.want])
		}
	}
}

// Standard drops in prod (2026-11-04) sit at cap - step for every slot and band.
func TestArmorStandardDropsRateGood(t *testing.T) {
	for _, band := range armorBandTiers {
		for _, aType := range ArmorTypes {
			step, ok := armorStandardDrop[aType]
			if !ok {
				step = armorStandardDropDefault
			}
			val := MaxArmor[band][aType] - step
			if got := ArmorQuality(aType, val); got != QualityGood {
				t.Errorf("type %d band %d armor %d: got %s", aType, band, val, QualityNames[got])
			}
		}
	}
}

func TestArmorTierMatchesCheckArmor(t *testing.T) {
	for _, aType := range ArmorTypes {
		for val := 0; val <= 140; val++ {
			tier := ArmorTier(aType, val)
			if tier == 0 {
				if CheckArmor(aType, LevelCap, val) {
					t.Errorf("type %d armor %d: no tier, but wearable at the level cap", aType, val)
				}
				continue
			}
			if !CheckArmor(aType, tier, val) {
				t.Errorf("type %d armor %d: reads tier %d but can't wear", aType, val, tier)
			}
			if tier > 1 && CheckArmor(aType, tier-1, val) {
				t.Errorf("type %d armor %d: reads tier %d but wearable a tier sooner", aType, val, tier)
			}
		}
	}
}

func TestChargeDescription(t *testing.T) {
	nearly := "This appears to be nearly used up."
	heavily := "This appears to be heavily used."
	lightly := "This appears to be lightly used."
	plump := "This appears plump and full."
	swelling := "This is swelling unnaturally."
	burst := "This looks like it's starting to burst."
	cases := []struct {
		uses, base int
		want       string
	}{
		{0, 100, nearly},
		{15, 100, nearly},
		{16, 100, heavily},
		{50, 100, heavily},
		{51, 100, lightly},
		{99, 100, lightly},
		{100, 100, plump},
		{150, 100, plump},
		{151, 100, swelling},
		{1000, 100, swelling},
		{1001, 100, burst},
		{5, 0, ""},
	}
	for _, c := range cases {
		if got := ChargeDescription(c.uses, c.base, ChargeTypeDefault); got != c.want {
			t.Errorf("%d of %d: got %q, want %q", c.uses, c.base, got, c.want)
		}
	}
}

func TestChargeWordSets(t *testing.T) {
	cases := []struct {
		chargeType, uses int
		want             string
	}{
		{ChargeTypeLiquid, 10, "This container is nearly empty."},
		{ChargeTypeLiquid, 100, "This container is full."},
		{ChargeTypeFood, 60, "A few bites are missing."},
		{ChargeTypeFood, 2000, "This looks like it's starting to burst."},
		{ChargeTypePlural, 10, "These are almost gone."},
		{ChargeTypePlural, 500, "There are an unusual number of these."},
		{ChargeTypePlural, 2000, "These are overflowing."},
		{99, 100, "This appears plump and full."},
	}
	for _, c := range cases {
		if got := ChargeDescription(c.uses, 100, c.chargeType); got != c.want {
			t.Errorf("type %d at %d: got %q, want %q", c.chargeType, c.uses, got, c.want)
		}
	}
}

func TestParseChargeType(t *testing.T) {
	cases := []struct {
		value string
		want  int
		ok    bool
	}{
		{"default", ChargeTypeDefault, true},
		{"LIQUID", ChargeTypeLiquid, true},
		{"food", ChargeTypeFood, true},
		{"3", ChargeTypePlural, true},
		{"potion", ChargeTypeDefault, false},
		{"4", ChargeTypeDefault, false},
	}
	for _, c := range cases {
		if got, ok := ParseChargeType(c.value); got != c.want || ok != c.ok {
			t.Errorf("%q: got %d %v, want %d %v", c.value, got, ok, c.want, c.ok)
		}
	}
}
