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
	21: 115,
	22: 120,
	23: 125,
	24: 130,
	25: 135,
	26: 140,
}

func CanWield(tier int, class int, max int) bool {
	if class == 0 {
		tier += 1
	}
	if max < MaxWeaponDamage[tier] {
		return true
	}
	return false
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
	var weaponLevel = CalculateLevel(exp, WeaponExpLevels)
	if weaponLevel == 10 {
		if class == 0 {
			return WeaponTitles[10]
		} else {
			return WeaponTitles[9]
		}
	} else {
		return WeaponTitles[weaponLevel]
	}
}

func AffinityExpTitle(exp int) string {
	return AffinityTitles[CalculateLevel(exp, WeaponExpLevels)]
}

func DivinityExpTitle(exp int) string {
	return DivinityTitles[CalculateLevel(exp, WeaponExpLevels)]
}

func StealthExpTitle(exp int) string {
	return StealthTitles[CalculateLevel(exp, StealthExpLevels)]
}

func StealthLevel(exp int) int {
	return CalculateLevel(exp, StealthExpLevels)
}

func StealthExpNext(exp int) int {
	var currentLevel = CalculateLevel(exp, StealthExpLevels)
	if currentLevel == 10 {
		return 0
	} else {
		return StealthExpLevels[currentLevel+1]
	}
}

func WeaponLevel(exp int, class int) int {
	var currentLevel = CalculateLevel(exp, WeaponExpLevels)
	if currentLevel == 10 {
		if class == 0 || class == 4 || class == 5 || class == 6 {
			return 10
		} else {
			return 9
		}
	} else {
		return currentLevel
	}
}

func WeaponExpNext(exp int, class int) int {
	var currentLevel = CalculateLevel(exp, WeaponExpLevels)
	if currentLevel >= 9 {
		if currentLevel == 9 && (class == 0 || class == 4 || class == 5 || class == 6) {
			return WeaponExpLevels[10]
		} else {
			return 0
		}
	} else {
		return WeaponExpLevels[currentLevel+1]
	}
}

func WeaponMissChance(exp int) int {
	var currentLevel = CalculateLevel(exp, WeaponExpLevels)
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
