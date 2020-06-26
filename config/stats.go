package config

// Str Mods
var StrCarryMod = 10 // Per Point
var BaseCarryWeight = 40
var StrDamageMod = .01 // Per Point

// Con Mods
var ConArmorMod = .01
var ConBonusHealth = 1
var ConBonusHealthDiv = 5

// Dex Mods
var DexDodgeMod = .0025  //Chance to dodge
var DexGlobalMod = .05 // Seconds to subtract from global ticker

// Int Mods
var IntOffensiveMod = .01
var IntManaPool = 2 // Number of points of mana to add
var IntManaPoolDiv = 5  // Number to divide by
var IntSpellEffectDuration = 30 // Seconds to add
var IntBroadDaily = 1
var IntEvalDaily = 1
var IntEvalDailyDiv = 3
var BaseEvals = 3

// Piety Mods
var PieRegenMod = .01 // Regen Mana per tick
var PieHealMod = .01 // Per point

func MaxWeight(str int) int {
	return BaseCarryWeight + (str*StrCarryMod)
}

func CalcHealth(tier int, con int, class int) int {
	if class>=99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Health) + (tier * ((con/ConBonusHealthDiv)*ConBonusHealth))
}

func CalcStamina(tier int, con int, class int) int {
	if class>=99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Stamina) + (tier * ((con/ConBonusHealthDiv)*ConBonusHealth))
}

func CalcMana(tier int, intel int, class int) int {
	if class>=99 {
		return 800
	}
	return (tier * Classes[AvailableClasses[class]].Mana) + (tier * ((intel/IntManaPoolDiv)*IntManaPool))
}