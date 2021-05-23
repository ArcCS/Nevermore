package config

import "github.com/ArcCS/Nevermore/utils"

var CombatModifiers = map[string]int{
	// Attack Modifiers
	"critical": 10,
	"double":   2,

	// Bash
	"thunk":    100,
	"crushing": 10,
	"thwomp":   2,

	// Amount of damage per strength point
	"berserk": 5,

	// Sneaky Types
	"backstab": 8,
	"snipe":    8,
}

//Thief & Ranger
var HideChance = 15
var SneakChance = 15
var StealChance = 15
var BackStabChance = 15
var SnipeChance = 5
var HideChancePerPoint = 3
var SneakChancePerPoint = 2
var StealChancePerPoint = 2
var BackStabChancePerPoint = 2
var SnipeChancePerPoint = 1

var MobBSRevengeVitalChance = 10
var VitalStrikeScale = 2

// Monk
var TodMax = .5
var TodFailDamage = .5
var TodScaleDown = .1

// Paladin
var TurnMax = .5
var TurnFailDamage = .5
var TurnScaleDown = .1

var CombatCooldown = 8

// Mob Stuns:
var ParryStuns = 2
var CircleStuns = 1
var CircleTimer = 30
var BashStuns = 15
var BashTimer = 180
var ThumpRoll = 10
var ThwompRoll = 50
var CrushingRoll = 500
var ThunkRoll = 1000

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
