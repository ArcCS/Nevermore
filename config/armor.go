package config

// Armor Values
var ArmorReduction = .01
var ArmorReductionPoints = 10

var MobArmorReduction = .03
var MobArmorReductionPoints = 10

// Max Damage
var MaxArmor = map[int]map[string]int{
	// Level int ref is cap
	4:  {"max": 20, "chest": 50, "legs": 40, "wrist": 15, "boot": 15, "head": 30, "back": 15, "hand": 15, "shield": 20, "shield_block": 5},
	9:  {"max": 30, "chest": 70, "legs": 50, "wrist": 25, "boot": 25, "head": 40, "back": 25, "hand": 25, "shield": 40, "shield_block": 10},
	14: {"max": 40, "chest": 90, "legs": 60, "wrist": 35, "boot": 35, "head": 50, "back": 35, "hand": 35, "shield": 60, "shield_block": 15},
	19: {"max": 50, "chest": 110, "legs": 80, "wrist": 45, "boot": 45, "head": 60, "back": 45, "hand": 45, "shield": 80, "shield_block": 20},
	25: {"max": 60, "chest": 130, "legs": 80, "wrist": 55, "boot": 55, "head": 70, "back": 55, "hand": 55, "shield": 100, "shield_block": 25},
}

func ReturnReduction(totalArmor int) float64 {
	return float64(totalArmor/ArmorReductionPoints) * ArmorReduction
}

func CheckArmor(aType string, tier int, val int) bool {
	for lev, vals := range MaxArmor {
		if tier <= lev {
			if val <= vals[aType] {
				return true
			}
		}
	}
	return false
}
