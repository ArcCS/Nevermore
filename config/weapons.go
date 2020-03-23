package config
// Combat values

// Max Damage
var MaxWeaponDamage = map[int]int{
	1: 5,
	2: 10,
	3: 15,
	4: 20,
	5: 25,
	6: 30,
	7: 35,
	8: 40,
	9: 45,
	10: 50,
	11: 55,
	12: 60,
	13: 65,
	14: 70,
	15: 75,
	16: 80,
	17: 85,
	18: 90,
	19: 95,
	20: 100,
	21: 105,
	22: 110,
	23: 115,
	24: 120,
	25: 125,
	26: 130,
}

// Exp to level weapon classes
var WeaponExpLevels = map[int64]int64{
	2: 1000,
	3: 10000,
	4: 100000,
	5: 250000,
	6: 500000,
	7: 750000,
	8: 1000000,
	9: 2500000,
	10: 5000000,
	11: 15000000,
}

type WeaponClass struct {
	Title string
	DoubleDamage float64
	CriticalStrike float64
	LethalRate float64
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

func WeaponExpTitle(exp int64) string {
	switch {
	case exp < WeaponExpLevels[0] && exp > 0:
		return WeaponTitles[0]
	case exp < WeaponExpLevels[1]:
		return WeaponTitles[1]
	case exp < WeaponExpLevels[2]:
		return WeaponTitles[2]
	case exp < WeaponExpLevels[3]:
		return WeaponTitles[3]
	case exp < WeaponExpLevels[4]:
		return WeaponTitles[4]
	case exp < WeaponExpLevels[5]:
		return WeaponTitles[5]
	case exp < WeaponExpLevels[6]:
		return WeaponTitles[6]
	case exp < WeaponExpLevels[7]:
		return WeaponTitles[7]
	case exp < WeaponExpLevels[8]:
		return WeaponTitles[8]
	case exp < WeaponExpLevels[9]:
		return WeaponTitles[9]
	default:
		return WeaponTitles[0]
	}
}