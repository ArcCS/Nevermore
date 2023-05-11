package config

// Allow these to be modified in the future if these are too aggressive
var IntMajorPenalty = 5
var IntMinorPenalty = 9
var StrMajorPenalty = 5
var StrMinorPenalty = 9
var ConMajorPenalty = 5
var ConMinorPenalty = 9
var DexMajorPenalty = 5
var DexMinorPenalty = 9
var PieMajorPenalty = 5
var PieMinorPenalty = 9

func IntUsePenalty(intel int) bool {
	//failChance := 0
	if intel <= MinorAbilityTier {
		return true
	}
	return false
}
