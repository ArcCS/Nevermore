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

var MobVital = 3
var MobCritical = 10
var MobDouble = 25
var MobFollowVital = 25

var BindCost = 75000
var RenameCost = 150000

var MissPerLevel = 8 // This is a percentage
var SearchPerInt = 3 // This is a percentage

var SurgeExtraDamage = .15
var SurgeDamageBonus = .20     // Percentage added when using surge
var InertialDamageIgnore = .20 // Percentage ignored when using inertial barrier
var ReflectDamagePerInt = .02  // Percentage of damage reflected per int point
var ReflectDamageFromMob = .15 // Percentage of damage reflected from mob

var DodgeDamagePerDex = .01
var FullDodgeChancePerDex = .01

// Thief & Ranger
var PeekCD = 8
var PeekFailCD = 32
var StealCD = 8
var HideChance = 15
var SneakChance = 15
var SneakBonus = 10
var StealChance = 15
var StealChancePerLevel = 5
var BackStabChance = 15
var BackStabChancePerLevel = 5
var SnipeChance = 5
var HideChancePerPoint = 3
var SneakChancePerPoint = 1
var SneakChancePerTier = 1
var StealChancePerPoint = 2
var BackStabChancePerPoint = 2
var SnipeChancePerPoint = 2
var SnipeChancePerLevel = 5
var SnipeFumbleChance = 20
var MobStealRevengeVitalChance = 5
var MobBSRevengeVitalChance = 10
var VitalStrikeScale = 2

// Monk
var TodMax = 5
var TodFailDamage = 50
var TodScaleDown = 10
var MonkArmorPerLevel = 15
var TodTimer = 600
var TodCost = 10
var VitalChance = 15
var MonkDexPerDice = .25
var MeditateTime = 600

// Paladin/Cleri
var TurnMax = 50
var TurnFailDamage = 50
var TurnScaleDown = 10
var DisintegrateChance = 5
var TurnTimer = 60
var SlamTimer = 30
var ShieldDamage = 3
var ShieldStun = 1
var CombatCooldown = 8

// Bard Stuff
var ScalePerPiety = 1
var DurationPerCon = 10

// Mob Stuns:
var ParryStuns = 2
var CircleStuns = 1
var CircleTimer = 30
var BashStuns = 16
var BashTimer = 45

// Mob Things
var MobBlock = 25
var MobBlockPerLevel = 5
var MobFollow = 25
var MobFollowPerLevel = 2
var MobTakeChance = 10 // Percent

// Str Mods
var StrCarryMod = 10 // Per Point
var BaseCarryWeight = 40
var StrDamageMod = .03 // Per Point

// Con Mods
var ConArmorMod = .01
var ConBonusHealth = 1
var ConBonusHealthDiv = 5
var ConHealRegenMod = .05
var ConMonkArmor = 2 // 2 Armor Extra Per Con

// Dex Mods
var DexDodgeMod = .0025 //Chance to dodge
var DexGlobalMod = .05  // Seconds to subtract from global ticker

// Int Mods
var IntOffensiveMod = .01
var IntManaPool = 2             // Number of points of mana to add
var IntManaPoolDiv = 5          // Number to divide by
var IntSpellEffectDuration = 30 // Seconds to add
var IntBroad = 1                // Number of broadcasts per int point
var IntEvalDivInt = 3           //Divide int by this number to get eval
var BaseEvals = 1
var BaseBroads = 5
var IntMinCast = 5
var IntNoFizzle = 10
var FizzleSave = 45

// Piety Mods
var PieRegenMod = .1 // Regen Mana per tick
var PieHealMod = .3  // Per point

// Armor Values
var ArmorReduction = .01
var ArmorReductionPoints = 10

var MobArmorReduction = .03
var MobArmorReductionPoints = 10

func MaxWeight(str int) int {
	return BaseCarryWeight + (str * StrCarryMod)
}

func CalcHealth(tier int, con int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Health) + int(float64(tier)*float64(con/ConBonusHealthDiv)*float64(ConBonusHealth))
}

func CalcStamina(tier int, con int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Stamina) + int(float64(tier)*float64(con/ConBonusHealthDiv)*float64(ConBonusHealth))
}

func CalcMana(tier int, intel int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Mana) + int(float64(tier)*float64(intel/IntManaPoolDiv)*float64(IntManaPool))
}

func CalcHaste(tier int) int {
	if tier < 10 {
		return 2
	} else if tier >= 10 && tier < 15 {
		return 3
	} else if tier > 15 {
		return 4
	}
	return 0
}

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
		if dRoll <= Parry[skill-1] {
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

// BashChances Skill = Thunk, Crushing, Thwomp, Thump
var BashChances = map[int][]int{
	0: {0, 0, 0, 0},
	1: {125, 0, 0, 0},
	2: {250, 0, 0, 0},
	3: {500, 0, 0, 0},
	4: {750, 0, 0, 0},
	5: {1000, 0, 0, 0},
	6: {1250, 0, 0, 0},
	7: {1500, 0, 0, 0},
	8: {1875, 0, 0, 0},
	9: {3000, 0, 0, 0},
}

func RollBash(skill int) (damModifier int, stunModifier int, output string) {
	/*
		var ThumpRoll = 10
		var ThwompRoll = 50
		var CrushingRoll = 500
		var ThunkRoll = 1000

	*/
	damModifier = 1
	stunModifier = 1
	bashRoll := utils.Roll(1000000, 1, 0)
	if bashRoll <= BashChances[skill][0] { // Thunk
		damModifier = CombatModifiers["thunk"]
		output = "Thunk!!"
	} else if bashRoll <= BashChances[skill][1] { // Crushing
		damModifier = CombatModifiers["crushing"]
		output = "Craaackk!!"
	} else if bashRoll <= BashChances[skill][2] { // Thwomp
		damModifier = CombatModifiers["thwomp"]
		output = "Thwomp!!"
	} else if bashRoll <= BashChances[skill][3] { // Thump
		stunModifier = 3
		output = "Thump!!"
	}
	return
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

func BreatheDamage(level int) int {
	switch {
	case level < 5:
		return 8
	case level < 10:
		return 20
	case level < 15:
		return 40
	case level < 20:
		return 90
	case level < 25:
		return 125
	default:
		return 8
	}
}
