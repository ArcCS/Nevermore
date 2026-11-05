package config

import "github.com/ArcCS/Nevermore/text"

// Qualitative item ratings shown when a player looks at weapons and armor.
const (
	QualityUnremarkable = iota
	QualityGood
	QualityExcellent
	QualityExceptional
)

var QualityNames = map[int]string{
	QualityUnremarkable: "unremarkable",
	QualityGood:         "good quality",
	QualityExcellent:    "excellent quality",
	QualityExceptional:  "exceptional quality",
}

var QualityColors = map[int]string{
	QualityUnremarkable: text.White,
	QualityGood:         text.LightBlue,
	QualityExcellent:    text.LightPurple,
	QualityExceptional:  text.LightOrange,
}

// Where a weapon's rating changes, in points of the real cap off the tier's
// WeaponRollTargetPercent.
const WeaponQualityGoodOffset = -2
const WeaponQualityExcellentOffset = 1
const WeaponQualityExceptionalOffset = 3

// WeaponQualityFlagBonus is the points of the real cap a weapon with a special
// flag rates above its dice.
const WeaponQualityFlagBonus = 2

// WeaponQuality rates a weapon's base roll against its tier's target.  Always
// crit weapons have deliberately low dice, so they rate on what they do, and a
// roll past every cap has no tier to rate against.
func WeaponQuality(roll WeaponRoll, twoHanded bool, alwaysCrit bool, specialFlag bool) int {
	tier := WeaponTier(roll.Max(), twoHanded)
	if alwaysCrit || tier == 0 {
		return QualityExceptional
	}
	realCap := RealWeaponCap(tier, twoHanded)
	target := WeaponRollTargetPercent(tier)
	// In hundredths of a point, so half point averages stay whole.
	score := roll.Avg2() * 50
	if specialFlag {
		score += realCap * WeaponQualityFlagBonus
	}
	switch {
	case score >= realCap*(target+WeaponQualityExceptionalOffset):
		return QualityExceptional
	case score >= realCap*(target+WeaponQualityExcellentOffset):
		return QualityExcellent
	case score >= realCap*(target+WeaponQualityGoodOffset):
		return QualityGood
	}
	return QualityUnremarkable
}

// WeaponTierFor is the lowest tier a class can wield a weapon at, or 0 if its
// max roll is past every cap.  Fighters read one tier ahead, see WeaponDamageCap.
func WeaponTierFor(max int, twoHanded bool, class int) int {
	tier := WeaponTier(max, twoHanded)
	if class == FIGHTER && tier > 1 {
		tier -= 1
	}
	return tier
}

// armorBandTiers lists MaxArmor's rows in order.
var armorBandTiers = []int{4, 9, 14, 19, 25}

// armorStandardDrop is how far under its band's cap a slot's standard drop
// sits, and armorExcellentReach how close to the cap a piece must be to rate
// above standard.  Slots not listed use the defaults.
const armorStandardDropDefault = 5
const armorExcellentReachDefault = 3

var armorStandardDrop = map[int]int{5: 10, 23: 10, 24: 1}
var armorExcellentReach = map[int]int{5: 5, 23: 5, 24: 0}

// ArmorQuality rates a piece's base armor against the lowest band that can wear it.
func ArmorQuality(aType int, val int) int {
	for _, band := range armorBandTiers {
		cap := MaxArmor[band][aType]
		if val > cap {
			continue
		}
		standard, ok := armorStandardDrop[aType]
		if !ok {
			standard = armorStandardDropDefault
		}
		reach, ok := armorExcellentReach[aType]
		if !ok {
			reach = armorExcellentReachDefault
		}
		switch under := cap - val; {
		case under == 0:
			return QualityExceptional
		case under <= reach:
			return QualityExcellent
		case under <= standard:
			return QualityGood
		}
		return QualityUnremarkable
	}
	return QualityExceptional
}

// ArmorTier is the lowest tier that can wear a piece's base armor, or 0 if it
// is past every cap.
func ArmorTier(aType int, val int) int {
	for tier := 1; tier <= LevelCap; tier++ {
		if CheckArmor(aType, tier, val) {
			return tier
		}
	}
	return 0
}

// ChargeDescription tells how used a device or wand looks from its remaining
// uses against its template's, or "" if the template has none.  Melding can
// take an item well past its template.
func ChargeDescription(uses int, base int) string {
	if base <= 0 {
		return ""
	}
	switch percent := uses * 100 / base; {
	case percent <= 15:
		return "This item appears to be nearly used up."
	case percent <= 50:
		return "This item appears to be heavily used."
	case percent <= 99:
		return "This item appears to be lightly used."
	case percent <= 150:
		return "This item appears plump and full."
	case percent <= 1000:
		return "This item is swelling unnaturally."
	}
	return "This item looks like it's starting to burst."
}
