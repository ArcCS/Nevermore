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

var (
	MultiLower = .1
	MultiUpper = .2

	ProximityChance = 80
	ProximityStep   = 10

	BerserkCooldown = 60 * 5
	CombatCooldown  = 8
	UnequipCooldown = 2

	RoomClearTimer            = 3  // Seconds
	RoomEffectInvocation      = 18 // Seconds
	RoomDefaultEncounterSpeed = 10 // Seconds
	RoomMaxJigger             = 4
	RoomEncNoDoubles          = 6
	RoomEncNoTriples          = 5

	BaseDevicePiety = 8.0

	IntMajorPenalty = 7
	PieMajorPenalty = 5

	MobAugmentPerCharacter = 3

	FreeDeathTier = 4

	SpecialAbilityTier = 7
	MinorAbilityTier   = 5

	MobVital       = 3
	MobCritical    = 4
	MobDouble      = 10
	MobFollowVital = 35
	MobFollMult    = 3

	BindCost   = 75000
	RenameCost = 150000

	MissPerLevel = 8 // This is a percentage
	SearchPerInt = 3 // This is a percentage

	SurgeExtraDamage     = .15
	SurgeDamageBonus     = .20 // Percentage added when using surge
	InertialDamageIgnore = .20 // Percentage ignored when using inertial barrier
	ReflectDamagePerInt  = .02 // Percentage of damage reflected per int point
	ReflectDamageFromMob = .15 // Percentage of damage reflected from mob

	DodgeDamagePerDex     = .01
	FullDodgeChancePerDex = .01

	PeekCD                      = 8
	StealCD                     = 8
	HideChance                  = 20
	SneakChance                 = 20
	SneakBonus                  = 10
	StealChance                 = 20
	StealChancePerSkillLevel    = 4
	BackStabChance              = 20
	BackStabChancePerSkillLevel = 3
	SnipeChance                 = 15
	HideChancePerPoint          = 3
	SneakChancePerPoint         = 1
	SneakChancePerTier          = 1
	StealChancePerPoint         = 1
	BackStabChancePerPoint      = 1
	SnipeChancePerPoint         = 1
	SnipeChancePerLevel         = 5
	SnipeFumbleChance           = 20
	MobStealRevengeVitalChance  = 15
	MobBSRevengeVitalChance     = 25
	VitalStrikeScale            = 2
	BackstabCooldown            = 30
	TrackCooldown               = 16
	TrackChance                 = 20
	TrackChancePerLevel         = 5
	TrackChancePerPoint         = 1

	TodMax            = 5
	TodScaleDown      = 10
	MonkArmorPerLevel = 15
	TodTimer          = 600
	TodCost           = 10
	VitalChance       = 15
	MeditateTime      = 600

	TurnMax            = 50
	TurnScaleDown      = 10
	DisintegrateChance = 5
	TurnTimer          = 60
	SlamTimer          = 30
	ShieldDamage       = 3
	ShieldStun         = .4

	ScalePerPiety  = 1
	DurationPerCon = 10

	ParryStuns  = 2
	CircleStuns = 1
	CircleTimer = 16
	HamTimer    = 24
	BashStuns   = 16
	BashTimer   = 45

	MobBlock          = 35
	MobBlockPerLevel  = 15
	MobFollow         = 40
	MobFollowPerLevel = 2
	MobTakeChance     = 20 // Percent

	StrCarryMod     = 10 // Per Point
	BaseCarryWeight = 40
	StatDamageMod   = .01 // Per Point

	ReduceSickCon     = 1
	SickConBonus      = 2
	ConBonusHealthDiv = 5
	ConHealRegenMod   = .10
	ConMonkArmor      = 2 // 2 Armor Extra Per Con
	ConFallDamageMod  = 1
	ConArmorMod       = .01

	HitPerDex        = 1
	MissPerDex       = 1
	DexFallDamageMod = 1

	FallDamage = .20

	IntResistMagicBase     = 10
	IntResistMagicPerPoint = 1
	IntManaPool            = 2  // Number of points of mana to add
	IntManaPoolDiv         = 5  // Number to divide by
	IntSpellEffectDuration = 30 // Seconds to add
	IntBroad               = 1  // Number of broadcasts per int point
	IntEvalDivInt          = 3  //Divide int by this number to get eval
	BaseEvals              = 1
	BaseBroads             = 5
	FizzleSave             = 50

	PieRegenMod = .4 // Regen Mana per tick
	PieHealMod  = .7 // Per point

	ArmorReduction       = .007
	ArmorReductionPoints = 10

	MobArmorReduction = .5

	ExperienceReduction = map[int]float64{
		1: .9,
		2: .7,
		3: .6,
		4: .5,
		5: .45,
	}
)

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

var LethalDamage = []int{ // Lethals are 1000000 chance rolls.
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
	1: {125 * 4, 600 * 4, 1200 * 4, 2400 * 4},
	2: {250 * 4, 1000 * 4, 2000 * 4, 4000 * 4},
	3: {500 * 4, 2000 * 4, 4000 * 4, 8000 * 4},
	4: {750 * 4, 3000 * 4, 6000 * 4, 12000 * 4},
	5: {1000 * 4, 4000 * 4, 8000 * 4, 16000 * 4},
	6: {1250 * 4, 5000 * 4, 10000 * 4, 20000 * 4},
	7: {1500 * 4, 6000 * 4, 12000 * 4, 24000 * 4},
	8: {1875 * 4, 7500 * 4, 15000 * 4, 30000 * 4},
	9: {3000 * 4, 12000 * 4, 24000 * 4, 48000 * 4},
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
