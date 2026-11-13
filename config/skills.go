package config

// Combat values

var MaxWeaponDamage = map[int]int{
	1:  15,
	2:  20,
	3:  25,
	4:  30,
	5:  35,
	6:  40,
	7:  45,
	8:  50,
	9:  55,
	10: 60,
	11: 65,
	12: 70,
	13: 75,
	14: 80,
	15: 85,
	16: 90,
	17: 95,
	18: 100,
	19: 105,
	20: 110,
	21: 120,
	22: 130,
	23: 140,
	24: 150,
	25: 160,
	26: 170,
}

// MaxTierFighterBonus is the flat damage a fighter at LevelCap adds to every
// weapon roll, standing in for reading one tier ahead once caps no longer apply.
const MaxTierFighterBonus = 5

// IgnoresGearCaps reports whether a tier is past the weapon and armor caps.
func IgnoresGearCaps(tier int) bool {
	return tier >= LevelCap
}

// CanWieldTwoHanded reports whether a class may wield two-handed weapons.
func CanWieldTwoHanded(class int) bool {
	return class == FIGHTER || class == BARBARIAN || class == PALADIN
}

// TwoHandedCapPercent is the share of the tier's weapon damage cap that a
// two-handed weapon adds to it; TwoHandedCapStep is the multiple the bonus is
// rounded to, and also the least it can be.
const TwoHandedCapPercent = 15
const TwoHandedCapStep = 5

// TwoHandedCapBonus is how much higher the damage cap sits for a two-handed
// weapon: TwoHandedCapPercent of the cap, rounded to the nearest
// TwoHandedCapStep (halves up), never less than one step.
func TwoHandedCapBonus(cap int) int {
	steps := (cap*TwoHandedCapPercent + TwoHandedCapStep*50) / (TwoHandedCapStep * 100)
	if steps < 1 {
		steps = 1
	}
	return steps * TwoHandedCapStep
}

// WeaponDamageCap is the max roll a class at a tier must stay under to wield a
// weapon.  Fighters read one tier ahead; two-handed weapons get TwoHandedCapBonus.
func WeaponDamageCap(tier int, class int, twoHanded bool) int {
	if class == FIGHTER {
		tier += 1
	}
	cap := MaxWeaponDamage[tier]
	if twoHanded {
		cap += TwoHandedCapBonus(cap)
	}
	return cap
}

func CanWield(tier int, class int, max int, twoHanded bool) bool {
	return max < WeaponDamageCap(tier, class, twoHanded)
}

func CalculateLevel(exp int, expTable map[int]int) int {
	switch {
	case exp >= expTable[1] && exp < expTable[2]:
		return 1
	case exp >= expTable[2] && exp < expTable[3]:
		return 2
	case exp >= expTable[3] && exp < expTable[4]:
		return 3
	case exp >= expTable[4] && exp < expTable[5]:
		return 4
	case exp >= expTable[5] && exp < expTable[6]:
		return 5
	case exp >= expTable[6] && exp < expTable[7]:
		return 6
	case exp >= expTable[7] && exp < expTable[8]:
		return 7
	case exp >= expTable[8] && exp < expTable[9]:
		return 8
	case exp >= expTable[9] && exp < expTable[10]:
		return 9
	case exp >= expTable[10]:
		return 10
	default:
		return 0
	}
}

// SkillExpLevels is the exp needed for each level of every skill: weapons,
// elemental affinity, divinity and stealth.
var SkillExpLevels = map[int]int{
	0:  0,
	1:  3000,
	2:  20000,
	3:  100000,
	4:  500000,
	5:  1000000,
	6:  1600000,
	7:  2400000,
	8:  6000000,
	9:  12000000,
	10: 36000000,
}

var WeaponTitles = []string{
	"Unskilled",
	"Basic",
	"Skilled",
	"Experienced",
	"Refined",
	"Ace",
	"Adept",
	"Expert",
	"Specialist",
	"Master",
	"Grandmaster",
}

var AffinityTitles = []string{
	"Unattuned",
	"Neophyte",
	"Novice",
	"Channeler",
	"Artisan",
	"Specialist",
	"Attuned",
	"Elementalist",
	"Savant",
	"Virtuoso",
	"Ascended",
}

var DivinityTitles = []string{
	"Agnostic",
	"Novice",
	"Graceful",
	"Blessed",
	"Radiant",
	"Sanctified",
	"Sacred",
	"Exhalted",
	"Supernal",
	"Angelic",
	"Transcendent",
}

var StealthTitles = []string{
	"Street Urchin",
	"Footpad",
	"Cutpurse",
	"Burguler",
	"Prowler",
	"Infiltrator",
	"Elusive",
	"Shadow Dancer",
	"Phantom Blade",
	"Assassin",
	"Master Assassin",
}

var HealingSkill = map[int]int{
	0:  0,
	1:  20,
	2:  40,
	3:  60,
	4:  80,
	5:  100,
	6:  115, // tapers from rank 6 so top healers don't run away
	7:  130,
	8:  140,
	9:  150,
	10: 160,
}

// SpellTierDamagePercent is the damage bonus, in percent per caster tier,
// every player's damage spells get; a mage's affinity bonus adds to it.
const SpellTierDamagePercent = 1

// SpellIntDamagePercent is the damage bonus, in percent per point of int above
// BaselineStatValue, every player's damage spells get.
const SpellIntDamagePercent = 1

// SpellDamageBonus is the percent a player's damage spell is raised by: the
// int and tier bonuses, plus the affinity bonus for mages. They add together
// rather than compounding.
func SpellDamageBonus(tier int, intel int, class int, affinityExp int, slot int) int {
	bonus := tier * SpellTierDamagePercent
	if intel > BaselineStatValue {
		bonus += (intel - BaselineStatValue) * SpellIntDamagePercent
	}
	if class == MAGE {
		bonus += SpellDmgSkill[WeaponLevel(affinityExp, class, slot)]
	}
	return bonus
}

var SpellDmgSkill = map[int]int{
	0:  0,
	1:  5,
	2:  10,
	3:  15,
	4:  20,
	5:  30,
	6:  40,
	7:  50,
	8:  60,
	9:  70,
	10: 80,
}

func WeaponExpTitle(exp int, class int, slot int) string {
	var weaponLevel = CalculateLevel(exp, SkillExpLevels)
	if weaponLevel == 10 {
		if ReachesGrandmaster(class, slot) {
			return WeaponTitles[10]
		} else {
			return WeaponTitles[9]
		}
	} else {
		return WeaponTitles[weaponLevel]
	}
}

func AffinityExpTitle(exp int) string {
	return AffinityTitles[CalculateLevel(exp, SkillExpLevels)]
}

func DivinityExpTitle(exp int) string {
	return DivinityTitles[CalculateLevel(exp, SkillExpLevels)]
}

func StealthExpTitle(exp int) string {
	return StealthTitles[CalculateLevel(exp, SkillExpLevels)]
}

func StealthLevel(exp int) int {
	return CalculateLevel(exp, SkillExpLevels)
}

func StealthExpNext(exp int) int {
	var currentLevel = CalculateLevel(exp, SkillExpLevels)
	if currentLevel == 10 {
		return 0
	} else {
		return SkillExpLevels[currentLevel+1]
	}
}

// MissileSkill is the skill slot for missile weapons (item type 4).
const MissileSkill = 4

// DivinitySkill is the skill slot for clerical healing skill.
const DivinitySkill = 10

// HandSkill is the skill slot for unarmed hand-to-hand fighting.
const HandSkill = 5

// RangerMeleeAdvancement is the ranger's weapon exp multiplier for anything
// but a missile weapon; their class value applies to bows alone.
const RangerMeleeAdvancement = .7

// WeaponAdvancementFor is the weapon exp multiplier a class earns in a skill
// slot. Rangers advance at full pace only with missile weapons.
func WeaponAdvancementFor(class int, slot int) float64 {
	if class == RANGER && slot != MissileSkill {
		return RangerMeleeAdvancement
	}
	return Classes[AvailableClasses[class]].WeaponAdvancement
}

// FirstElementSkill and LastElementSkill bound the elemental affinity skill
// slots (fire, air, earth, water).
const FirstElementSkill, LastElementSkill = 6, 9

// ReachesGrandmaster reports whether a class can reach skill level 10 in a
// slot. Only clerics grandmaster divinity and only mages an elemental
// affinity; with weapons it is fighters in anything, monks with hand-to-hand,
// and rangers with missile weapons.
func ReachesGrandmaster(class int, slot int) bool {
	if slot == DivinitySkill {
		return class == CLERIC
	}
	if slot >= FirstElementSkill && slot <= LastElementSkill {
		return class == MAGE
	}
	switch class {
	case FIGHTER:
		return true
	case MONK:
		return slot == HandSkill
	case RANGER:
		return slot == MissileSkill
	}
	return false
}

func WeaponLevel(exp int, class int, slot int) int {
	var currentLevel = CalculateLevel(exp, SkillExpLevels)
	if currentLevel == 10 {
		if ReachesGrandmaster(class, slot) {
			return 10
		} else {
			return 9
		}
	} else {
		return currentLevel
	}
}

func WeaponExpNext(exp int, class int, slot int) int {
	var currentLevel = CalculateLevel(exp, SkillExpLevels)
	if currentLevel >= 9 {
		if currentLevel == 9 && ReachesGrandmaster(class, slot) {
			return SkillExpLevels[10]
		} else {
			return 0
		}
	} else {
		return SkillExpLevels[currentLevel+1]
	}
}

func WeaponMissChance(exp int) int {
	var currentLevel = CalculateLevel(exp, SkillExpLevels)
	switch {
	case currentLevel == 0:
		return 30
	case currentLevel == 1:
		return 28
	case currentLevel == 2:
		return 26
	case currentLevel == 3:
		return 24
	case currentLevel == 4:
		return 22
	case currentLevel == 5:
		return 20
	case currentLevel == 6:
		return 15
	case currentLevel == 7:
		return 10
	case currentLevel == 8:
		return 5
	case currentLevel >= 9:
		return 0
	default:
		return 50
	}
}
