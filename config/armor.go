package config

// Max Damage
var MaxArmor = map[int]map[int]int{
	// Level int ref is cap
	4:  {9: 50, 20: 40, 21: 15, 19: 15, 25: 30, 22: 15, 26: 15, 23: 20},
	9:  {9: 70, 20: 50, 21: 25, 19: 25, 25: 40, 22: 25, 26: 25, 23: 40},
	14: {9: 90, 20: 60, 21: 35, 19: 35, 25: 50, 22: 35, 26: 35, 23: 60},
	19: {9: 110, 20: 70, 21: 45, 19: 45, 25: 60, 22: 45, 26: 45, 23: 80},
	25: {9: 130, 20: 80, 21: 55, 19: 55, 25: 70, 22: 55, 26: 55, 23: 100},
}

var MaxTotals = map[int]map[string]int{
	// Level int ref is cap
	4:  {"max": 20, "shield_block": 5},
	9:  {"max": 30, "shield_block": 10},
	14: {"max": 40, "shield_block": 15},
	19: {"max": 50, "shield_block": 20},
	25: {"max": 60, "shield_block": 25},
}

func ReturnReduction(totalArmor int) float64 {
	return float64(totalArmor/ArmorReductionPoints) * ArmorReduction
}

func CheckArmor(aType int, tier int, val int) bool {
	for lev, vals := range MaxArmor {
		if tier <= lev {
			if val <= vals[aType] {
				return true
			}
		}
	}
	return false
}

func CheckMaxArmor(maxType string, tier int, newTotal int) bool {
	for lev, vals := range MaxTotals {
		if tier <= lev {
			if newTotal <= vals[maxType] {
				return true
			}
		}
	}
	return false
}
