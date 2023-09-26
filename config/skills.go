package config

// Combat values

// Max Damage
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
	21: 115,
	22: 120,
	23: 125,
	24: 130,
	25: 135,
	26: 140,
}

// Quick Function to check if character can wield
func CanWield(tier int, class int, max int) bool {
	if class == 0 {
		tier += 1
	}
	if max < MaxWeaponDamage[tier] {
		return true
	}
	return false
}

// Exp to level weapon classes
var WeaponExpLevels = map[int]int{
	0:  0,
	1:  3000,
	2:  30000,
	3:  300000,
	4:  750000,
	5:  1500000,
	6:  2250000,
	7:  3000000,
	8:  7500000,
	9:  15000000,
	10: 45000000,
}

var SkillAdvancement = map[int]float32{
	0: 1,
	1: .7,
	2: .7,
	3: .7,
	6: .7,
	8: .5,
	7: .5,
	5: .4,
	4: .4,
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
	"Moonshadow",
}

var StealthExpLevels = map[int]int{
	0:  0,
	1:  3000,
	2:  30000,
	3:  300000,
	4:  750000,
	5:  1500000,
	6:  2250000,
	7:  3000000,
	8:  7500000,
	9:  15000000,
	10: 45000000,
}

var HealingSkill = map[int]int{
	0:  0,
	1:  20,
	2:  40,
	3:  60,
	4:  80,
	5:  100,
	6:  120,
	7:  140,
	8:  160,
	9:  180,
	10: 200,
}

var SpellDmgSkill = map[int]int{
	0:  0,
	1:  5,
	2:  10,
	3:  15,
	4:  25,
	5:  35,
	6:  45,
	7:  70,
	8:  85,
	9:  100,
	10: 120,
}

func WeaponExpTitle(exp int, class int) string {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return WeaponTitles[0]
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return WeaponTitles[1]
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return WeaponTitles[2]
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return WeaponTitles[3]
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return WeaponTitles[4]
	case exp >= WeaponExpLevels[5] && exp < WeaponExpLevels[6]:
		return WeaponTitles[5]
	case exp >= WeaponExpLevels[6] && exp < WeaponExpLevels[7]:
		return WeaponTitles[6]
	case exp >= WeaponExpLevels[7] && exp < WeaponExpLevels[8]:
		return WeaponTitles[7]
	case exp >= WeaponExpLevels[8] && exp < WeaponExpLevels[9]:
		return WeaponTitles[8]
	case exp >= WeaponExpLevels[9] && exp < WeaponExpLevels[10]:
		return WeaponTitles[9]
	case exp >= WeaponExpLevels[10]:
		if class == 0 {
			return WeaponTitles[10]
		} else {
			return WeaponTitles[9]
		}
	default:
		return WeaponTitles[0]
	}
}

func AffinityExpTitle(exp int) string {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return AffinityTitles[0]
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return AffinityTitles[1]
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return AffinityTitles[2]
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return AffinityTitles[3]
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return AffinityTitles[4]
	case exp >= WeaponExpLevels[5] && exp < WeaponExpLevels[6]:
		return AffinityTitles[5]
	case exp >= WeaponExpLevels[6] && exp < WeaponExpLevels[7]:
		return AffinityTitles[6]
	case exp >= WeaponExpLevels[7] && exp < WeaponExpLevels[8]:
		return AffinityTitles[7]
	case exp >= WeaponExpLevels[8] && exp < WeaponExpLevels[9]:
		return AffinityTitles[8]
	case exp >= WeaponExpLevels[9] && exp < WeaponExpLevels[10]:
		return AffinityTitles[9]
	case exp >= WeaponExpLevels[10]:
		return AffinityTitles[10]
	default:
		return AffinityTitles[0]
	}
}

func DivinityExpTitle(exp int) string {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return DivinityTitles[0]
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return DivinityTitles[1]
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return DivinityTitles[2]
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return DivinityTitles[3]
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return DivinityTitles[4]
	case exp >= WeaponExpLevels[5] && exp < WeaponExpLevels[6]:
		return DivinityTitles[5]
	case exp >= WeaponExpLevels[6] && exp < WeaponExpLevels[7]:
		return DivinityTitles[6]
	case exp >= WeaponExpLevels[7] && exp < WeaponExpLevels[8]:
		return DivinityTitles[7]
	case exp >= WeaponExpLevels[8] && exp < WeaponExpLevels[9]:
		return DivinityTitles[8]
	case exp >= WeaponExpLevels[9] && exp < WeaponExpLevels[10]:
		return DivinityTitles[9]
	case exp >= WeaponExpLevels[10]:
		return DivinityTitles[10]
	default:
		return DivinityTitles[0]
	}
}

func StealthExpTitle(exp int) string {
	switch {
	case exp >= StealthExpLevels[0] && exp < StealthExpLevels[1]:
		return StealthTitles[0]
	case exp >= StealthExpLevels[1] && exp < StealthExpLevels[2]:
		return StealthTitles[1]
	case exp >= StealthExpLevels[2] && exp < StealthExpLevels[3]:
		return StealthTitles[2]
	case exp >= StealthExpLevels[3] && exp < StealthExpLevels[4]:
		return StealthTitles[3]
	case exp >= StealthExpLevels[4] && exp < StealthExpLevels[5]:
		return StealthTitles[4]
	case exp >= StealthExpLevels[5] && exp < StealthExpLevels[6]:
		return StealthTitles[5]
	case exp >= StealthExpLevels[6] && exp < StealthExpLevels[7]:
		return StealthTitles[6]
	case exp >= StealthExpLevels[7] && exp < StealthExpLevels[8]:
		return StealthTitles[7]
	case exp >= StealthExpLevels[8] && exp < StealthExpLevels[9]:
		return StealthTitles[8]
	case exp >= StealthExpLevels[9] && exp < StealthExpLevels[10]:
		return StealthTitles[9]
	case exp >= StealthExpLevels[10]:
		return StealthTitles[10]
	default:
		return StealthTitles[0]
	}
}

func StealthLevel(exp int) int {
	switch {
	case exp >= StealthExpLevels[0] && exp < StealthExpLevels[1]:
		return 0
	case exp >= StealthExpLevels[1] && exp < StealthExpLevels[2]:
		return 1
	case exp >= StealthExpLevels[2] && exp < StealthExpLevels[3]:
		return 2
	case exp >= StealthExpLevels[3] && exp < StealthExpLevels[4]:
		return 3
	case exp >= StealthExpLevels[4] && exp < StealthExpLevels[5]:
		return 4
	case exp >= StealthExpLevels[5] && exp < StealthExpLevels[6]:
		return 5
	case exp >= StealthExpLevels[6] && exp < StealthExpLevels[7]:
		return 6
	case exp >= StealthExpLevels[7] && exp < StealthExpLevels[8]:
		return 7
	case exp >= StealthExpLevels[8] && exp < StealthExpLevels[9]:
		return 8
	case exp >= StealthExpLevels[9] && exp < StealthExpLevels[10]:
		return 9
	case exp >= StealthExpLevels[10]:
		return 10
	default:
		return 0
	}
}

func StealthExpNext(exp int) int {
	switch {
	case exp >= StealthExpLevels[0] && exp < StealthExpLevels[1]:
		return StealthExpLevels[1]
	case exp >= StealthExpLevels[1] && exp < StealthExpLevels[2]:
		return StealthExpLevels[2]
	case exp >= StealthExpLevels[2] && exp < StealthExpLevels[3]:
		return StealthExpLevels[3]
	case exp >= StealthExpLevels[3] && exp < StealthExpLevels[4]:
		return StealthExpLevels[4]
	case exp >= StealthExpLevels[4] && exp < StealthExpLevels[5]:
		return StealthExpLevels[5]
	case exp >= StealthExpLevels[5] && exp < StealthExpLevels[6]:
		return StealthExpLevels[6]
	case exp >= StealthExpLevels[6] && exp < StealthExpLevels[7]:
		return StealthExpLevels[7]
	case exp >= StealthExpLevels[7] && exp < StealthExpLevels[8]:
		return StealthExpLevels[8]
	case exp >= StealthExpLevels[8] && exp < StealthExpLevels[9]:
		return StealthExpLevels[9]
	case exp >= StealthExpLevels[9] && exp < StealthExpLevels[10]:
		return StealthExpLevels[10]
	case exp >= StealthExpLevels[10]:
		return 0
	default:
		return StealthExpLevels[1]
	}
}

func WeaponLevel(exp int, class int) int {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return 0
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return 1
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return 2
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return 3
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return 4
	case exp >= WeaponExpLevels[5] && exp < WeaponExpLevels[6]:
		return 5
	case exp >= WeaponExpLevels[6] && exp < WeaponExpLevels[7]:
		return 6
	case exp >= WeaponExpLevels[7] && exp < WeaponExpLevels[8]:
		return 7
	case exp >= WeaponExpLevels[8] && exp < WeaponExpLevels[9]:
		return 8
	case exp >= WeaponExpLevels[9] && exp < WeaponExpLevels[10]:
		return 9
	case exp >= WeaponExpLevels[10]:
		if class == 0 || class == 4 || class == 5 || class == 6 {
			return 10
		} else {
			return 9
		}
	default:
		return 0
	}
}

func WeaponExpNext(exp int, class int) int {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return WeaponExpLevels[1]
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return WeaponExpLevels[2]
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return WeaponExpLevels[3]
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return WeaponExpLevels[4]
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return WeaponExpLevels[5]
	case exp >= WeaponExpLevels[5] && exp < WeaponExpLevels[6]:
		return WeaponExpLevels[6]
	case exp >= WeaponExpLevels[6] && exp < WeaponExpLevels[7]:
		return WeaponExpLevels[7]
	case exp >= WeaponExpLevels[7] && exp < WeaponExpLevels[8]:
		return WeaponExpLevels[8]
	case exp >= WeaponExpLevels[8] && exp < WeaponExpLevels[9]:
		return WeaponExpLevels[9]
	case exp >= WeaponExpLevels[9] && exp < WeaponExpLevels[10]:
		if class == 0 || class == 4 || class == 5 || class == 6 {
			return WeaponExpLevels[10]
		} else {
			return 0
		}
	case exp >= WeaponExpLevels[10]:
		return 0
	default:
		return WeaponExpLevels[1]
	}
}

func WeaponMissChance(exp int, class int) int {
	switch {
	case exp >= WeaponExpLevels[0] && exp < WeaponExpLevels[1]:
		return 30
	case exp >= WeaponExpLevels[1] && exp < WeaponExpLevels[2]:
		return 25
	case exp >= WeaponExpLevels[2] && exp < WeaponExpLevels[3]:
		return 20
	case exp >= WeaponExpLevels[3] && exp < WeaponExpLevels[4]:
		return 15
	case exp >= WeaponExpLevels[4] && exp < WeaponExpLevels[5]:
		return 10
	case exp >= WeaponExpLevels[5]:
		return 0
	default:
		return 50
	}
}
