package config

import "github.com/ArcCS/Nevermore/utils"

var CombatModifiers = map[string]int {
	// Attack Modifiers
	"critical": 10,
	"double":	2,

	// Bash
	"thunk": 	100,
	"crushing": 10,
	"thwomp":	2,

	// Sneaky Types
	"backstab":	8,
	"snipe":	8,
}

// TODO: Capture these attacks in their respective commands
/*
	Thump		Triple Stun duration	Bash
	Berserk	6-7 damage per point of Strength added to base damage	5 Bonus Strength	Berserk
	Touch of Death		50% Halve Hit Points / 50% Fatal	Touch of Death
	Turn		50% Halve Hit Points / 50% Fatal	Turn
 */

var CombatCooldown = 8


// Mob Stuns:
var ParryStuns = 2
var CircleStuns = 1
var CircleTimer = 30

// Double Damage is out of 100
var Parry = []int{
	0,
	2,
	3,
	5,
	6,
	8,
	10,
	14,
	15,
	18,
	20,
}

func RollParry(skill int) bool {
	if skill > 0 {
		dRoll := utils.Roll(100, 1, 0)
		if dRoll <= Parry[skill] {
			return true
		}
	}
	return false
}

// Double Damage is out of 100
var DoubleDamage = []int{
	0,
	1,
	2,
	4,
	6,
	8,
	10,
	12,
	15,
	20,
	25,
}

func RollDouble(skill int) bool {
	if skill > 0 {
		dRoll := utils.Roll(100, 1, 0)
		if dRoll <= DoubleDamage[skill] {
			return true
		}
	}
	return false
}

// Criticals are out of 1000
var CriticalDamage = []int{
	0,
	1,
	2,
	3,
	4,
	5,
	6,
	7,
	8,
	10,
	12,
}

func RollCritical(skill int) bool {
	if skill > 0 {
		dRoll := utils.Roll(1000, 1, 0)
		if dRoll <= CriticalDamage[skill] {
			return true
		}
	}
	return false
}

// Lethals are 1000000 chance rolls.
var LethalDamage = []int{
	0,
	125,
	250,
	500,
	750,
	1000,
	1250,
	1500,
	1875,
	2500,
	3125,
}

func RollLethal(skill int) bool {
	if skill > 0 {
		dRoll := utils.Roll(1000000, 1, 0)
		if dRoll <= LethalDamage[skill] {
			return true
		}
	}
	return false
}
