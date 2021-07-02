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
var PeekCD = 4
var StealCD = 4
var HideChance = 15
var SneakChance = 15
var StealChance = 15
var StealChancePerLevel = 5
var BackStabChance = 15
var BackStabChancePerLevel = 5
var SnipeChance = 5
var HideChancePerPoint = 3
var SneakChancePerPoint = 2
var StealChancePerPoint = 2
var BackStabChancePerPoint = 2
var SnipeChancePerPoint = 2
var SnipeChancePerLevel = 5
var SnipeFumbleChance = 20

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

// Mob Things
var MobBlock = 70
var MobBlockPerLevel = 10
var MobFollow = 60
var MobFollowPerLevel = 5

// Str Mods
var StrCarryMod = 10 // Per Point
var BaseCarryWeight = 40
var StrDamageMod = .01 // Per Point

// Con Mods
var ConArmorMod = .01
var ConBonusHealth = 1
var ConBonusHealthDiv = 5
var ConHealRegenMod = .8
var ConMonkArmor = 2 // 2 Armor Extra Per Con

// Dex Mods
var DexDodgeMod = .0025 //Chance to dodge
var DexGlobalMod = .05  // Seconds to subtract from global ticker

// Int Mods
var IntOffensiveMod = .01
var IntManaPool = 2             // Number of points of mana to add
var IntManaPoolDiv = 5          // Number to divide by
var IntSpellEffectDuration = 30 // Seconds to add
var IntBroadDaily = 1
var IntEvalDaily = 1
var IntEvalDailyDiv = 3
var BaseEvals = 3

// Piety Mods
var PieRegenMod = .8 // Regen Mana per tick
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
	return (tier * Classes[AvailableClasses[class]].Health) + (tier * ((con / ConBonusHealthDiv) * ConBonusHealth))
}

func CalcStamina(tier int, con int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Stamina) + (tier * ((con / ConBonusHealthDiv) * ConBonusHealth))
}

func CalcMana(tier int, intel int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Mana) + (tier * ((intel / IntManaPoolDiv) * IntManaPool))
}

func CalcHaste(tier int) int {
	if tier < 10 {
		return 2
	}else if tier >=10 && tier < 15 {
		return 3
	}else if tier > 15 {
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

func RollLethal(skill int) bool {
	if skill > 0 {
		dRoll := utils.Roll(1000000, 1, 0)
		if dRoll <= LethalDamage[skill] {
			return true
		}
	}
	return false
}
