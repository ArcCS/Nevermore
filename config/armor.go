package config

// Max Damage
var MaxArmor = map[int]map[int]int{
	// Level int ref is cap
	4:  {5: 50, 20: 40, 21: 15, 19: 15, 25: 30, 22: 15, 26: 15, 23: 20, 24: 4},
	9:  {5: 70, 20: 50, 21: 25, 19: 25, 25: 40, 22: 25, 26: 25, 23: 40, 24: 6},
	14: {5: 90, 20: 60, 21: 35, 19: 35, 25: 50, 22: 35, 26: 35, 23: 60, 24: 8},
	19: {5: 110, 20: 70, 21: 45, 19: 45, 25: 60, 22: 45, 26: 45, 23: 80, 24: 10},
	25: {5: 130, 20: 80, 21: 55, 19: 55, 25: 70, 22: 55, 26: 55, 23: 100, 24: 12},
}

func CheckArmor(aType, tier, val int) bool {
	switch {
	case tier <= 4:
		if val <= MaxArmor[4][aType] {
			return true
		}
	case tier <= 9:
		if val <= MaxArmor[9][aType] {
			return true
		}
	case tier <= 14:
		if val <= MaxArmor[14][aType] {
			return true
		}
	case tier <= 19:
		if val <= MaxArmor[19][aType] {
			return true
		}
	case tier <= 25:
		if val <= MaxArmor[25][aType] {
			return true
		}
	}
	return false
}
