package config

import "github.com/ArcCS/Nevermore/utils"

var CombatModifiers = map[string]int{
	// Attack Modifiers
	"critical": 5,
	"double":   2,

	// Bash
	"thunk":    100,
	"crushing": 10,
	"thwomp":   2,

	// Amount of damage per strength point
	"berserk": 5,

	// Sneaky Types
	"backstab": 5,
	"snipe":    4,
}

var CombatCooldown = 8
var UnequipCooldown = 2

var RoomClearTimer = 3            // Seconds
var RoomEffectInvocation = 18     // Seconds
var RoomDefaultEncounterSpeed = 8 // Seconds

var IntMajorPenalty = 7
var PieMajorPenalty = 5

var MobAugmentPerCharacter = 3

var FreeDeathTier = 4

var SpecialAbilityTier = 7
var MinorAbilityTier = 5

var MobVital = 3
var MobCritical = 4
var MobDouble = 10
var MobFollowVital = 35

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
var StealCD = 8
var HideChance = 20
var SneakChance = 20
var SneakBonus = 10
var StealChance = 20
var StealChancePerSkillLevel = 4
var BackStabChance = 20
var BackStabChancePerLevel = 5
var BackStabChancePerSkillLevel = 3
var SnipeChance = 15
var HideChancePerPoint = 3
var SneakChancePerPoint = 1
var SneakChancePerTier = 1
var StealChancePerPoint = 1
var BackStabChancePerPoint = 1
var SnipeChancePerPoint = 1
var SnipeChancePerLevel = 5
var SnipeFumbleChance = 20
var MobStealRevengeVitalChance = 5
var MobBSRevengeVitalChance = 10
var VitalStrikeScale = 2

// Monk
var TodMax = 5
var TodScaleDown = 10
var MonkArmorPerLevel = 15
var TodTimer = 600
var TodCost = 10
var VitalChance = 15
var MeditateTime = 600

// Paladin/Cleri
var TurnMax = 50
var TurnScaleDown = 10
var DisintegrateChance = 5
var TurnTimer = 60
var SlamTimer = 30
var ShieldDamage = 3
var ShieldStun = .4

// Bard Stuff
var ScalePerPiety = 1
var DurationPerCon = 10

// Mob Stuns:
var ParryStuns = 2
var CircleStuns = 1
var CircleTimer = 16
var HamTimer = 24
var BashStuns = 16
var BashTimer = 45

// Mob Things
var MobBlock = 25
var MobBlockPerLevel = 5
var MobFollow = 25
var MobFollowPerLevel = 2
var MobTakeChance = 20 // Percent

// Str Mods
var StrCarryMod = 10 // Per Point
var BaseCarryWeight = 40
var StatDamageMod = .01 // Per Point

// Con Mods
var ReduceSickCon = 1
var SickConBonus = 2
var ConBonusHealthDiv = 5
var ConHealRegenMod = .10
var ConMonkArmor = 2 // 2 Armor Extra Per Con
var ConFallDamageMod = 1
var ConArmorMod = .01

// Dex Mods
var HitPerDex = 1
var MissPerDex = 1
var DexDodgeMod = .0025 //Chance to dodge
var DexFallDamageMod = 1

var FallDamage = .20

// Int Mods
var IntResistMagicBase = 10
var IntResistMagicPerPoint = 1
var IntManaPool = 2             // Number of points of mana to add
var IntManaPoolDiv = 5          // Number to divide by
var IntSpellEffectDuration = 30 // Seconds to add
var IntBroad = 1                // Number of broadcasts per int point
var IntEvalDivInt = 3           //Divide int by this number to get eval
var BaseEvals = 1
var BaseBroads = 5
var FizzleSave = 50

// Piety Mods
var PieRegenMod = .4 // Regen Mana per tick
var PieHealMod = .7  // Per point

// Armor Values
var ArmorReduction = .007
var ArmorReductionPoints = 10

var MobArmorReduction = .5

//Party
var ExperienceReduction = map[int]float64{
	1: .9,
	2: .7,
	3: .6,
	4: .5,
	5: .45,
}

func MaxWeight(str int) int {
	return BaseCarryWeight + (str * StrCarryMod)
}

func CalcHealth(tier int, con int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Health) + int(float64(tier)*(float64(con)/float64(ConBonusHealthDiv)))
}

func CalcStamina(tier int, con int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Stamina) + int(float64(tier)*(float64(con)/float64(ConBonusHealthDiv)))
}

func CalcMana(tier int, intel int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Mana) + int(float64(tier)*(float64(intel)/float64(IntManaPoolDiv))*float64(IntManaPool))
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
	1: {125, 600, 1200, 2400},
	2: {250, 1000, 2000, 4000},
	3: {500, 2000, 4000, 8000},
	4: {750, 3000, 6000, 12000},
	5: {1000, 4000, 8000, 16000},
	6: {1250, 5000, 10000, 20000},
	7: {1500, 6000, 12000, 24000},
	8: {1875, 7500, 15000, 30000},
	9: {3000, 12000, 24000, 48000},
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
