package config

import (
	"time"

	"github.com/ArcCS/Nevermore/utils"
)

// CombatModifiers are damage multipliers for special attack results.
var CombatModifiers = map[string]float64{
	// Attack Modifiers
	"critical": 5,
	"double":   2,
	// A barbarian's crushing blow on a regular attack; it also stuns.
	"crushing": 5,

	// Bash (Thunk scales with tier, see ThunkMultiplier)
	"thwomp": 2,
	"thump":  1.5,

	// Sneaky Types (backstab scales with stealth, see BackstabMultiplier)
	"snipe": 4,
	// A critical backstab adds this to the backstab multiplier.
	"backstab_critical": 3,
}

// MultiAttackMultipliers is the damage multiplier for each landed hit of a
// fighter or ranger multi attack, keyed by weapon skill level.  Multipliers
// are applied to landed hits in order, so the first hit that connects always
// takes the 1.0 slot.  Levels not present get a single swing.
var MultiAttackMultipliers = map[int][]float64{
	5:  {1, 0.2},
	6:  {1, 0.4},
	7:  {1, 0.4, 0.1},
	8:  {1, 0.5, 0.2},
	9:  {1, 0.5, 0.3, 0.2},
	10: {1, 0.5, 0.3, 0.2, 0.2},
}

// MultiAttackMissPenaltyFor returns the extra miss chance, in percentage
// points, applied to every swing after the first in a multi attack.  It starts
// at MultiAttackMissPenalty and drops by MultiAttackMissReduction for each
// weapon skill level above MultiAttackMissReductionLevel.
func MultiAttackMissPenaltyFor(skillLevel int) int {
	penalty := MultiAttackMissPenalty - MultiAttackMissReduction*(skillLevel-MultiAttackMissReductionLevel)
	if penalty > MultiAttackMissPenalty {
		penalty = MultiAttackMissPenalty
	}
	if penalty < 0 {
		penalty = 0
	}
	return penalty
}

// FollowUpMissFloorFor is the lowest miss chance a follow-up swing can roll
// at a weapon skill level. The follow-up penalty is added before the clamp,
// so spare dex works it off down to this floor.
func FollowUpMissFloorFor(skillLevel int) int {
	if skillLevel >= FollowUpMissFloorMasterLevel {
		return FollowUpMissFloorMaster
	}
	return FollowUpMissFloor
}

// MajorAbilityTier gates the tier 10 class abilities (berserk, touch of
// death, the long leap, the meditate save). A const so other config
// constants can build on it.
const MajorAbilityTier = 10

var (
	// Follow-up swings add the penalty before the 5-95 clamp, like backstab,
	// then floor at FollowUpMissFloor (FollowUpMissFloorMaster from Master).
	// (2026-10-14: was 25 added after the clamp, so dex past the floor did
	// nothing for follow-ups.)
	MultiAttackMissPenalty        = 35 // Percentage points added to follow-up swing miss chance
	MultiAttackMissReduction      = 5  // Points removed from the penalty per skill level past the threshold
	MultiAttackMissReductionLevel = 7  // Skill level at which the penalty starts falling off: 35 through Expert, then 30/25/20 at Specialist/Master/Grandmaster
	FollowUpMissFloor             = 20
	FollowUpMissFloorMaster       = 15
	FollowUpMissFloorMasterLevel  = 9
	BaselineStatValue             = 10
	ProximityChance               = 80
	ProximityStep                 = 10

	CombatCooldown  = 8
	UnequipCooldown = 2
	// StanceCooldown is the wait between stance toggles (flurry, reckless).
	// Stances ignore every other timer, so this only stops toggle spam.
	StanceCooldown = 1

	// ResumeGraceSeconds is how long after resuming a session a death still
	// counts as a lag death, so a fight that ran while the player was
	// link-dead cannot kill them for real before they have seen the screen.
	ResumeGraceSeconds = 3

	RoomClearTimer            = 3  // Seconds
	// PermMobRegenPctPerMinute is how much of its max HP a permanent mob
	// heals per minute its room sits empty.
	PermMobRegenPctPerMinute = 10
	RoomEffectInvocation      = 18 // Seconds
	RoomDefaultEncounterSpeed = 10 // Seconds
	RoomMaxJigger             = 4
	// Crowding: once a room holds more hostile mobs than characters, double
	// spawns stop and the spawn chance falls linearly, reaching half at
	// chars + RoomEncCrowdHalfAt mobs and staying there.
	RoomEncCrowdHalfAt = 4

	BaseDevicePiety = 8.0

	IntMajorPenalty = 9
	PieMajorPenalty = 5

	MobAugmentPerCharacter = 3

	FreeDeathTier = 4

	SpecialAbilityTier = 7
	MinorAbilityTier   = 5
	SealJusticeTier    = 10 // Courage and faith open at MinorAbilityTier, justice here

	MobVital    = 3
	MobCritical = 4
	MobDouble   = 10
	// A mob critical multiplies damage by MobCriticalMult less
	// MobCriticalPerDex per point of the target's dex: 2.7x at 20 dex.
	MobCriticalMult   = 3.0
	MobCriticalPerDex = 0.015
	MobFollowVital    = 40

	BindCost   = 75000
	RenameCost = 150000

	MissPerLevel = 8 // This is a percentage
	SearchPerInt = 3 // This is a percentage

	SurgeExtraDamage     = .15
	SurgeDamageBonus     = .20 // Percentage added when using surge
	InertialDamageIgnore = .15 // Percentage a mob's inertial barrier ignores; players use InertialBarrierIgnore (2026-10-13: was .20)
	ReflectDamagePerInt  = .02 // Percentage of damage reflected per int point
	ReflectDamageFromMob = .15 // Percentage of damage reflected from mob
	LethalSkillFloor     = .25 // Minimum fraction of a mob's exp granted as weapon skill on a lethal or touch of death kill

	DodgeDamagePerDex     = .01
	FullDodgeChancePerDex = 1.0

	PeekCD                      = 8
	StealCD                     = 8
	HideChance                  = 20
	SneakChance                 = 20
	SneakBonus                  = 10
	StealChance                 = 20
	StealChancePerSkillLevel    = 4
	BackstabMissPenalty         = 30  // points of miss added to the regular weapon roll (before the clamp) at stealth level 0
	BackstabDamageBase          = 2.5 // backstab damage multiplier at stealth level 0
	BackstabDamageSkillModifier = .2  // added to the multiplier per stealth level: 3.9x at specialist, 4.5x at grandmaster
	BackstabGrandmasterBonus    = 1.0 // added on top at grandmaster stealth: 5.5x
	BackstabPenaltyPerStealth   = 2   // points of that penalty removed per stealth level, 10 left at grandmaster
	SnipeChance                 = 15
	HideChancePerPoint          = 3
	SneakChancePerPoint         = 1
	SneakChancePerTier          = 1
	StealChancePerPoint         = 1
	SnipeChancePerPoint         = 1
	SnipeChancePerLevel         = 5
	SnipeFumbleChance           = 20
	MobStealRevengeVitalChance  = 15
	MobBSRevengeVitalChance     = 25
	MobTurnRevengeVitalChance   = 25
	MobTodRevengeVitalChance    = 25
	VitalStrikeScale            = 1.75 // every punishment vital: follow, failed backstab/steal/turn/touch
	QuickdrawCooldown           = 30
	TrackCooldown               = 16
	TrackChance                 = 20
	TrackChancePerLevel         = 5
	TrackChancePerPoint         = 1

	TurnBase            = 40
	TurnScaleDown       = 10
	TurnPiePerChance    = 2
	TurnDivinityPerRank = 1
	TurnMinChance       = 5 // percent; any turn can land
	TurnDivinityExpMult = 0.5
	DisintegrateShare   = 10 // percent of the turn success chance that disintegrates instead
	DisintegrateMin     = 3  // percent of all turns that disintegrate at least, even when a turn would fail
	TurnTimer           = 60
	SlamTimer           = 30
	ShieldSlamStuns     = 16 // seconds

	// Paladin seals persist until dropped, so the effect gets a duration it
	// will never reach.
	SealDuration = 10 * 365 * 24 * 60 * 60 // Seconds
	SealTimer    = 10
	// Seal of courage armor is base + per level, raised by one percent per
	// point of piety (see SealCourageArmor).
	SealCourageArmorBase     = 50
	SealCourageArmorPerLevel = 2
	// Offensive threat a paladin under the seal of courage generates is raised
	// by this percent, so the shield bearer holds mobs off the party.
	SealCourageThreatPercent = 10
	// RescueWindow is how long attacks aimed at a rescued ally are drawn onto
	// the paladin; RescueTimer is the cooldown between rescues.
	RescueWindow = 8  // Seconds
	RescueTimer  = 30 // Seconds
	// A rescue lifts the rescuer to the mob's top threat plus this percent of
	// it (at least 1), so the mob stays on them past a little extra damage
	// from the rest of the party. (2026-10-13: was top + 1.)
	RescueThreatMarginPercent = 10

	ScalePerPiety  = 1
	DurationPerCon = 10

	ParryStuns  = 16
	CircleStuns = 6
	CircleTimer = 16
	BashStuns   = 16

	// Threat bumps, as a percent of the mob's max stamina, added on top of any
	// damage the ability dealt. Abilities pick a size rather than owning a
	// number. Large: the real taunts (circle, hamstring, shield slam) and the
	// cap on a failed turn, which hands over the mob's current stamina. Small:
	// lighter taunts that peel a mob without out-holding the real tanks (bash,
	// the monk's sweep) and a failed backstab. (2026-10-13: was one constant
	// per ability; the monk's taunt moved from feint back to sweep.)
	SmallThreatPercent = 25
	LargeThreatPercent = 50
	// Heal threat is the amount healed split evenly across the mobs
	// attacking the recipient (plus the healer's own target), but no single
	// mob takes more than HealThreatCapPercent of the heal.
	HealThreatCapPercent = 50

	// MobThreatSwitchChance is the percent chance per mob tick (about 8s)
	// that a mob turns on whoever tops its threat table when that is not its
	// current target. Taunts set the target directly and skip this roll.
	MobThreatSwitchChance = 10

	MobBlock          = 35
	MobBlockPerLevel  = 15
	MobFollow         = 40
	MobFollowPerLevel = 2
	// A freshly spawned mob can't follow anyone out of the room for this long.
	MobFollowSpawnLockout = 4 * time.Second
	MobTakeChance         = 20 // Percent

	StrCarryMod     = 10 // Per Point
	BaseCarryWeight = 40
	StatDamageMod   = .01 // Per Point

	ReduceSickCon     = 1
	SickConBonus      = 2
	ConBonusHealthDiv = 5
	ConHealRegenMod   = .10
	ConFallDamageMod  = 1
	ConArmorMod       = .005

	HitPerDex        = 1
	MissPerDex       = 1
	BlessHitBonus    = 5 // percentage points off the miss chance while blessed
	DexFallDamageMod = 1

	FallDamage = .20

	// A resist-fire/air/earth/water spell takes this share off magic damage of
	// its element, for players and mobs alike. (2026-10-13: was 25%.)
	ElementResistPercent = 20

	// Every point of int above IntResistMagicBase blocks IntResistMagicPerPoint
	// magic damage, times max(1, spell tier/IntResistMagicTierDiv); see
	// IntMagicResist. (2026-10-14: was a flat int - 1 for any int above 10.)
	IntResistMagicBase     = 10
	IntResistMagicPerPoint = 1
	IntResistMagicTierDiv  = 2
	ManaPerStatPerTier     = .3 // Mana per point of the pool stat (int, or pie for divine casters and monks), per tier
	IntSpellEffectDuration = 30 // Seconds to add
	IntBroad               = 1  // Number of broadcasts per int point
	IntEvalDivInt          = 3  //Divide int by this number to get eval
	BaseEvals              = 1
	BaseBroads             = 5
	FizzleSave             = 25 // chance to fizzle per int below 9

	PieHealMod           = .7 // Per point
	MinorPieHealMod      = .9 // Per point above MinorHealPieFloor, vigor/mend only
	MinorHealPieFloor    = 10 // Piety up to this adds nothing to vigor/mend
	MinorHealDivinityMod = .3 // Vigor/mend get 30% of the divinity bonus of detraumatize/renewal
	PaladinDivinityMod   = .5 // Paladins get this share of their divinity bonus without the seal of faith

	// Seal of justice damage bonus, in percent: (tier - TierStart)/TierDiv +
	// pie/PieDiv + weaponLevel/WeaponDiv. About 10% at tier 10 and 25.5% at 20,
	// so a justice paladin leads a fighter less while fighters lack multi
	// attack. (2026-10-15: was 5 + tier/2 + pie/4 + weaponLevel/2, about 15.5%
	// at tier 10.)
	SealJusticeTierStart   = 7
	SealJusticeTierDiv     = 1.0
	SealJusticePieDiv      = 3.0
	SealJusticeWeaponDiv   = 2.0
	MinorHealTierDiv       = 4 // Vigor/mend gain +1 base per this many tiers
	CasterMinorHealTierDiv = 3 // MinorHealTierDiv for classes in MinorHealCasterClasses
	// Spellcasting classes; monks, fighters, barbarians and thieves keep MinorHealTierDiv.
	MinorHealCasterClasses = []int{MAGE, CLERIC, PALADIN, BARD, RANGER}
	MinorHealDieSides      = 6 // Vigor/mend roll 1d(sides) + MinorHealDieBonus
	MinorHealDieBonus      = 7
	DetraumatizeBase       = 20
	RenewalBase            = 40
	MajorHealTierDiv       = 3 // Detraumatize gains +1 base per this many tiers
	MajorHealBaseCut       = 2 // Subtracted from detraumatize base to offset the tier bonus
	// Renewal has its own shape: 35 + tier/2 + pie*RenewalPieMod +
	// 2d(5 + tier/3) before divinity, for 25 mana. It stays about 22-24%
	// less mana-efficient than detraumatize for healers with 16+ pie.
	RenewalBaseCut    = 5
	RenewalTierDiv    = 2
	RenewalPieMod     = .8
	RenewalDice       = 2
	RenewalDieBase    = 5
	RenewalDieTierDiv = 3

	ArmorReduction       = .007
	ArmorReductionPoints = 10
	// Player damage taken = constant / (constant + armor). 1000 matches the
	// old linear formula through tier 14 at optimized gear.
	ArmorReductionConstant = 1000

	MobArmorReduction = .5

	ExperienceReduction = map[int]float64{
		1: .9,
		2: .7,
		3: .6,
		4: .5,
		5: .45,
	}
	AttractionChance = 80
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

// CalcMana is the mana pool at a tier. Monks carry chi instead, which scales
// on piety rather than intelligence.
func CalcMana(tier int, intel int, pie int, class int) int {
	if class >= 99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Mana) + int(float64(tier)*float64(ManaPoolStat(intel, pie, class))*ManaPerStatPerTier)
}

// ManaPoolStat is the stat a class's mana pool scales from: piety for the
// divine casters and monks, intelligence for everyone else.
func ManaPoolStat(intel int, pie int, class int) int {
	if class == CLERIC || class == PALADIN || class == MONK {
		return pie
	}
	return intel
}

func CalcHaste(tier int) int {
	switch {
	case tier >= 15:
		return 4
	case tier >= 10:
		return 3
	default:
		return 2
	}
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
	return RollCriticalScaled(skill, 1)
}

// BackstabCriticalChanceMultiplier scales the weapon skill's critical chance
// for a backstab.
const BackstabCriticalChanceMultiplier = 2

// RollCriticalScaled rolls a critical with the table chance multiplied by
// multiplier.
func RollCriticalScaled(skill int, multiplier int) bool {
	if skill > 0 {
		dRoll := utils.Roll(1000, 1, 0)
		if dRoll <= CriticalDamage[skill]*multiplier {
			return true
		}
	}
	return false
}

// WeaponShatterChance is the percent chance that a critical strike, crushing
// blow or Thunk shatters the weapon that landed it, destroying it outright.
const WeaponShatterChance = 5

func RollShatter() bool {
	return utils.Roll(100, 1, 0) <= WeaponShatterChance
}

// Repair pricing. A repair costs RepairCostFraction of the item's value scaled
// by the share of its uses that are gone. Every WeaponRepairsBeforeOverhaul
// repairs a weapon (ArmorRepairsBeforeOverhaul for armor) needs an overhaul instead:
// OverhaulCostMultiplier times the normal price, resetting the count, with an
// OverhaulFailChance percent chance the item is destroyed (and no gold is taken). quest_loot items and rings
// never need an overhaul but pay QuestLootRepairMultiplier times the price.
const (
	RepairCostFraction          = .3
	WeaponRepairsBeforeOverhaul = 10
	ArmorRepairsBeforeOverhaul  = 5
	OverhaulCostMultiplier      = 15
	OverhaulFailChance          = 25
	QuestLootRepairMultiplier   = 1.5
)

// Meld pricing. A meld charges for the larger of the two items' template uses
// (uses melded on earlier are never paid for twice), each at t² + MeldPriceBase, where t is the
// lowest tier any class can cast the item's spell at; spell-less items (ammo)
// fall back to a third of the template value per use. Prices sit near the
// 40th percentile of the old value-based price (people meld their cheapest
// items). Each full MeldPenaltyStep uses of the finished item adds
// MeldPenaltyPerStep to the price once it is past one step: 101–199 uses
// is +50%, 200 is +100%, 300 is +150%.
const (
	MeldPriceBase      = 3
	MeldPenaltyStep    = 100
	MeldPenaltyPerStep = .5
)

// CrushingChance is a barbarian's chance, out of 1000, that the first hit of
// a regular attack round with a blunt or two-handed weapon is a crushing
// blow, by weapon skill level: twice the critical chance.
var CrushingChance = []int{
	0,
	2,
	4,
	6,
	8,
	10,
	12,
	14,
	16,
	20,
	24,
}

// CrushingStuns is how long (seconds) a crushing blow stuns the mob.
const CrushingStuns = 8

func RollCrushing(skill int) bool {
	if skill > 0 {
		if utils.Roll(1000, 1, 0) <= CrushingChance[skill] {
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

// BashChances holds cumulative thresholds out of 1,000,000 for each bash
// result by skill level: Thunk, Thwomp, Thump. Thunk carries the chance the
// old Crushing result had before crushing moved to regular attacks.
var BashChances = map[int][]int{
	0: {0, 0, 0},
	1: {1900, 4300, 9100},
	2: {3000, 7000, 15000},
	3: {6000, 14000, 30000},
	4: {9000, 21000, 45000},
	5: {12000, 28000, 60000},
	6: {15000, 35000, 75000},
	7: {18000, 42000, 90000},
	8: {22500, 52500, 112500},
	9: {36000, 84000, 180000},
}

// RollBash rolls for a special bash result. damModifier multiplies the bash
// damage and stunModifier multiplies BashStuns; both are 1 on a plain bash.
// tier sets the size of a Thunk.
func RollBash(skill int, tier int) (damModifier float64, stunModifier int, output string) {
	damModifier = 1
	stunModifier = 1
	row, ok := BashChances[skill]
	if !ok {
		return
	}
	bashRoll := utils.Roll(1000000, 1, 0)
	if bashRoll <= row[0] { // Thunk
		damModifier = ThunkMultiplier(tier)
		output = "Thunk!!"
	} else if bashRoll <= row[1] { // Thwomp
		damModifier = CombatModifiers["thwomp"]
		output = "Thwomp!!"
	} else if bashRoll <= row[2] { // Thump
		damModifier = CombatModifiers["thump"]
		stunModifier = 2
		output = "Thump!!"
	}
	return
}

func RollLethal(skill int) bool {
	return RollLethalScaled(skill, 1)
}

// RollLethalScaled rolls a lethal blow with the table chance multiplied by
// multiplier (see ExecuteMultiplier for the fighter's execute scaling).
func RollLethalScaled(skill int, multiplier float64) bool {
	if skill > 0 {
		dRoll := utils.Roll(1000000, 1, 0)
		if float64(dRoll) <= float64(LethalDamage[skill])*multiplier {
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

var XP_Modifiers = map[string]float64{
	"poisons":       .03,
	"fast_moving":   .06,
	"block_exit":    .06,
	"follows":       .06,
	"no_stun":       .03,
	"diseases":      .03,
	"spits_acid":    .03,
	"blinds":        .06,
	"no_steal":      .03,
	"ranged_attack": .03,
	"invisible":     .03,
	"steals":        .06,
}

// ThreatPercent is percent of a mob's max stamina, for the threat bumps above.
func ThreatPercent(maxStam int, percent int) int {
	return maxStam * percent / 100
}

// RescueThreat is the threat a rescue lifts the rescuer to on a mob whose top
// threat is top: top plus RescueThreatMarginPercent of it, at least top + 1.
func RescueThreat(top int) int {
	return top + max(1, ThreatPercent(top, RescueThreatMarginPercent))
}

// BackstabMultiplier is the backstab damage multiplier at a stealth level.
func BackstabMultiplier(stealthLevel int) float64 {
	mult := BackstabDamageBase + float64(stealthLevel)*BackstabDamageSkillModifier
	if stealthLevel >= 10 {
		mult += BackstabGrandmasterBonus
	}
	return mult
}

// BackstabMissPenaltyFor is the extra miss a backstab carries over a regular
// swing with the same weapon: BackstabMissPenalty at stealth level 0, falling
// by BackstabPenaltyPerStealth per level. It is added before the 5-95 clamp,
// so dex past a swing's cap still offsets it.
func BackstabMissPenaltyFor(stealthLevel int) int {
	penalty := BackstabMissPenalty - stealthLevel*BackstabPenaltyPerStealth
	if penalty < 0 {
		penalty = 0
	}
	return penalty
}

// TodMissPenaltyFor is the extra miss a touch of death carries over a regular
// swing: TodMissPenalty less tier/TodMissTierDiv, less pie/TodHitPieDiv. Like the
// backstab penalty it is added before the 5-95 clamp, so dex and piety past a
// swing's cap still offset it. It can go negative to soak other penalties.
func TodMissPenaltyFor(tier int, pie int) int {
	return TodMissPenalty - tier/TodMissTierDiv - pie/TodHitPieDiv
}

// FailedTurnThreat is the threat a failed turn hands the caster: the mob's
// current stamina, capped at LargeThreatPercent of its max.
func FailedTurnThreat(current int, max int) int {
	cap := ThreatPercent(max, LargeThreatPercent)
	if current > cap {
		return cap
	}
	return current
}

// MobStunBase is how long (seconds) a mob's stun spell holds a player, and
// MobStunMin is the floor once level difference is subtracted.
const (
	MobStunBase = 20
	MobStunMin  = 1
)

// MobStunDuration is how long a mob of mobLevel stuns a player of playerTier.
// A lower-level mob loses one second per level it is below the player, down
// to MobStunMin; a mob at or above the player's level stuns for the full duration.
func MobStunDuration(mobLevel int, playerTier int) int {
	duration := MobStunBase
	if mobLevel < playerTier {
		duration -= playerTier - mobLevel
	}
	if duration < MobStunMin {
		duration = MobStunMin
	}
	return duration
}

// A player's stun spell holds a mob for StunSpellBase seconds plus one second
// per StunSpellIntPerSecond points of int, so 15s at 12 int. It lands
// StunSpellChance percent of the time against a mob of equal level and int,
// plus StunSpellChancePerLevel per level the caster is above the mob and one
// point per StunSpellIntPerChance int the caster has over the mob.
const (
	StunSpellBase           = 12
	StunSpellIntPerSecond   = 4
	StunSpellChance         = 70
	StunSpellChancePerLevel = 10
	StunSpellIntPerChance   = 2
)

// Combust's burn DoT ticks BurnTicks times, BurnInterval seconds apart, each tick
// rolling 2d(caster tier/2 + 1)+2 scaled by the caster's spell damage bonus.
const (
	BurnTicks    = 3
	BurnInterval = 8
)

// StunSpellDuration is how long a player's stun spell holds a mob.
func StunSpellDuration(intel int) int {
	return StunSpellBase + intel/StunSpellIntPerSecond
}

// StunSpellSuccessChance is the percent chance a player's stun spell lands on a mob.
// Results outside 0-100 are left as-is; the roll treats them as never/always.
func StunSpellSuccessChance(tier int, mobLevel int, intel int, mobInt int) int {
	return StunSpellChance + (tier-mobLevel)*StunSpellChancePerLevel + (intel-mobInt)/StunSpellIntPerChance
}

// TurnSuccessChance is the percent chance a turn lands on an undead mob,
// never below TurnMinChance. Results over 100 are left as-is; the roll treats
// them as always.
func TurnSuccessChance(tier int, mobLevel int, pie int, divinityRank int) int {
	return max(TurnMinChance, TurnBase+(tier-mobLevel)*TurnScaleDown+pie/TurnPiePerChance+divinityRank*TurnDivinityPerRank)
}

// TurnDisintegratePerMille is the chance, out of 1000, that a turn with the
// given success chance destroys the undead outright: DisintegrateShare of the
// (clamped) success chance, never below DisintegrateMin percent. Rolled on the
// same d1000 as the success chance, at the bottom of its band.
func TurnDisintegratePerMille(successChance int) int {
	return max(DisintegrateMin*10, min(100, max(0, successChance))*10*DisintegrateShare/100)
}

// TurnDivinityExp is the divinity exp a turn earns for the damage it did:
// the weapon skill formula (share of the mob's max stamina times its exp)
// at TurnDivinityExpMult in place of a class advancement rate.
func TurnDivinityExp(damage int, mobMaxStam int, mobExp int) int {
	if mobMaxStam <= 0 {
		return 0
	}
	return int(float64(damage) / float64(mobMaxStam) * float64(mobExp) * TurnDivinityExpMult)
}

// Fighting a mob above your level gets harder the further up it is, but every
// character engaged with the mob (the attacker, their party in the room, and
// anyone else on its threat table) buys a level of headroom. Both penalties
// ramp one step per level past their threshold so there is no cliff.
const (
	// MobOutlevelDamagePerLevel is the extra damage a mob's regular attack
	// does per level it sits more than engaged+1 levels above its target.
	MobOutlevelDamagePerLevel = 0.05
)

// OutlevelDamageMult multiplies a mob's regular attack against a target lvlDiff
// levels below it: solo it starts at +5% three levels up, a duo at four.
func OutlevelDamageMult(lvlDiff int, engaged int) float64 {
	over := lvlDiff - engaged - 1
	if over <= 0 {
		return 1
	}
	return 1 + float64(over)*MobOutlevelDamagePerLevel
}

// OutlevelMissPenalty is the miss chance added to a player's attack on a mob
// lvlDiff levels above them: solo it starts at MissPerLevel two levels up, a
// duo at three.
func OutlevelMissPenalty(lvlDiff int, engaged int) int {
	over := lvlDiff - engaged
	if over <= 0 {
		return 0
	}
	return over * MissPerLevel
}

// SealCourageArmor is the armor the seal of courage grants a paladin of the
// given tier and piety, rounded to the nearest point.
func SealCourageArmor(tier int, pie int) int {
	return ((SealCourageArmorBase+SealCourageArmorPerLevel*tier)*(100+pie) + 50) / 100
}

// Inertial barrier, cast by a player on themselves: the share of incoming
// damage ignored is InertialBase percent, plus a point for every
// InertialTierDiv tiers above InertialTierStart, plus a point for every
// InertialIntDiv int. Fractions count; a tier 20 mage with 25 int sits around
// 15% and the theoretical top (tier 25, int 45) is about 19%. Tuned so a robed
// mage with protection and the barrier stays slightly behind a thief in light
// armor. (2026-10-13: was 15 + (tier-9)/4 + int/5, which put the mage 4-6
// points ahead.)
const (
	InertialBase      = 10.0
	InertialTierStart = 9
	InertialTierDiv   = 5.0
	InertialIntDiv    = 8.0
)

// InertialBarrierIgnore is the fraction of damage a player's inertial barrier
// ignores at the given tier and int.
func InertialBarrierIgnore(tier int, intel int) float64 {
	tierPart := 0.0
	if tier > InertialTierStart {
		tierPart = float64(tier-InertialTierStart) / InertialTierDiv
	}
	return (InertialBase + tierPart + float64(intel)/InertialIntDiv) / 100
}

// CombustMagnitude is the magnitude code of combust, which deals damage
// between a tier 4 and a tier 5 spell and counts as tier 4.
const CombustMagnitude = 45

// SpellTier is the tier of a damage spell from its magnitude code, or 0 for
// anything that is not a tiered spell (breath, damage over time, scripts).
func SpellTier(magnitude int) int {
	if magnitude == CombustMagnitude {
		return 4
	}
	if magnitude >= 1 && magnitude <= 7 {
		return magnitude
	}
	return 0
}

// IntMagicResist is the flat magic damage a character's int blocks from a
// hit of the given spell tier: (int - base) * perPoint, times
// max(1, tier/IntResistMagicTierDiv), rounded down. Int at or below the base
// blocks nothing, and tier 0 hits take the plain amount.
func IntMagicResist(intel int, tier int) int {
	if intel <= IntResistMagicBase {
		return 0
	}
	scale := max(IntResistMagicTierDiv, tier)
	return (intel - IntResistMagicBase) * IntResistMagicPerPoint * scale / IntResistMagicTierDiv
}
