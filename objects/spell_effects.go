package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"math/rand"
	"strconv"
)

var (
	TeleportTable = []int{117, 2567,
		2538, // Rymek Sewers
		2530, //  Slimes
		2761, //  Snake Pit
		2740, //  Sink Hole
		2508, //  Gars
		2698, //  Beavers
		1963, //  Plains
		369,  //  Untended Crop Field
		2483, //  Wild Garden
		2913, //  Nexus Sewers
		2130, //  Trolls
		2253, //  NorthEast Eldane Marsh
		2699, //  Crabs

		//Non Hostile Teleport Rooms:
		2696, //  Twisted Maze in Mandrake Farm
		848,  //  SouthEast side of Rymek Island
		1750, //  Odd Bush in Woods East of Rymek
		1831, //  Woods between Rymek & Nexus
		1864, //  Nexus Town Square
		1858, //  Nexus Park middle
		2104, //  Game Trail East of nexus
		2182, //  Eldane Forest Oak Grove
		2216, //  Clearing by IceWine river
		2375, //  Beach NorthWest of Nexus
		513,  //  BlackWoods boarder
		2015, //  NorthWest Crystal Mountains Path
		1991, //  Along Tilnar's Veign (right before Drow)
		2868, //  UnderCity Grotto
	}
	RecallRoom = "77"
)

var Effects = map[string]func(caller interface{}, target interface{}, magnitude int) string{
	"poison":           poison,
	"disease":          disease,
	"blind":            blind,
	"berserk":          berserk,
	"haste":            haste,
	"pray":             pray,
	"heal-stam":        healstam,
	"heal-vit":         healvit,
	"restore":          restore,
	"heal":             heal,
	"heal-all":         healall,
	"fire-damage":      firedamage,
	"earth-damage":     earthdamage,
	"air-damage":       airdamage,
	"water-damage":     waterdamage,
	"light":            light,
	"curepoison":       curepoison,
	"bless":            bless,
	"protection":       protection,
	"invisibility":     invisibility,
	"detect-invisible": detectInvisible,
	"teleport":         teleport,
	"stun":             stun,
	"recall":           recall,
	"summon":           summon,
	"wizard-walk":      wizardwalk,
	"levitate":         levitate,
	"resist-fire":      resistfire,
	"resist-magic":     resistmagic,
	//"remove-curse":     removecurse,
	"resist-air":       resistair,
	"resist-water":     resistwater,
	"resist-earth":     resistearth,
	"remove-disease":   removedisease,
	"remove-blindness": cureblindness,
	//"polymorph":        polymorph,
	"attraction":       attraction,
	"inertial-barrier": inertialbarrier,
	"surge":            surge,
	"resist-poison":    resistpoison,
	"resilient-aura":   resilientaura,
	"resist-disease":   resistdisease,
	"disrupt-magic":    disruptmagic,
	"reflection":       reflection,
	"dodge":            dodge,
	"resist-acid":      resistacid,
	//"embolden":         embolden,
}

func Cast(caller interface{}, target interface{}, spell string, magnitude int) string {
	return Effects[spell](caller, target, magnitude)
}

func berserk(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("berserk", "60", 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("berserk", "berserk", text.Red+"The red rage grips you!!!\n")
				target.SetModifier("str", 5)
				target.SetModifier("base_damage", target.GetStat("str")*config.CombatModifiers["berserk"])
			},
			func() {
				target.FlagOffAndMsg("berserk", "berserk", text.Cyan+"The tension releases and your rage fades...\n")
				target.SetModifier("base_damage", -target.GetStat("str")*config.CombatModifiers["berserk"])
				target.SetModifier("str", -5)
			})
	case *Mob:
		return ""
	}
	return ""
}

func blind(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("blind", "300", 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("blind", "blind", text.Red+"You've been blinded!!!!\n")
			},
			func() {
				target.FlagOffAndMsg("blind", "blind", text.Cyan+"Your vision returns!\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func poison(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		if !target.CheckFlag("resist-poison") {
			target.FlagOn("poisoned", "mob_poisoned")
			target.ApplyEffect("poison", strconv.Itoa(magnitude*5), 8, magnitude, // magnitude maps to level of mob
				func(triggers int) {
					damage := magnitude
					switch {
					case triggers <= 1:
						return
					case triggers <= 3:
						damage *= 2
					default:
						damage *= 2
					}
					// Reduce Damage by Con
					damage -= (target.GetStat("con") / config.SickConBonus) * config.ReduceSickCon
					if damage >= 1 {
						target.ReceiveDamageNoArmor(damage)
						if _, err := target.Write([]byte(text.Red + "The poison courses through your veins for " + strconv.Itoa(damage) + " damage!\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
					} else {
						if _, err := target.Write([]byte(text.Cyan + "You persevere through the poison and it has no effect!\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
					}
				},
				func() {
					target.FlagOff("poisoned", "mob_poisoned")
					if _, err := target.Write([]byte(text.Cyan + "The effects of the poison subside...\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				})
		} else {
			if _, err := target.Write([]byte(text.Cyan + "The poison has no effect on you!\n")); err != nil {
				log.Println("Error writing to player:", err)
			}
		}
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func disease(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		if !target.CheckFlag("resist-disease") {
			target.ApplyEffect("disease", strconv.Itoa(magnitude*4), 8, magnitude,
				func(triggers int) {
					damage := magnitude
					switch {
					case triggers <= 1:
						return
					case triggers <= 3:
						damage *= 2
					default:
						damage *= 2
					}
					// Reduce Damage by Con
					damage -= (target.GetStat("con") / config.SickConBonus) * config.ReduceSickCon
					if damage >= 1 {
						target.ReceiveDamageNoArmor(damage)
						target.FlagOn("disease", "mob_disease")
						if _, err := target.Write([]byte(text.Red + "The disease progresses, racking your body for " + strconv.Itoa(damage) + " damage!\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
					} else {
						if _, err := target.Write([]byte(text.Cyan + "You persevere through the disease and it has no effect!\n")); err != nil {
							log.Println("Error writing to player:", err)
						}
					}
				},
				func() {
					target.FlagOff("disease", "mob_disease")
					if _, err := target.Write([]byte(text.Cyan + "The disease subsides...\n")); err != nil {
						log.Println("Error writing to player:", err)
					}
				})
		} else {
			if _, err := target.Write([]byte(text.Cyan + "You resist the disease!\n")); err != nil {
				log.Println("Error writing to player:", err)
			}
		}
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func haste(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("haste", "60", 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("haste", "haste", text.Info+"Your muscles tighten and your reflexes hasten!!!\n")
				target.SetModifier("dex", 5)
			},
			func() {
				target.FlagOffAndMsg("haste", "haste", text.Cyan+"Your reflexes return to normal.\n")
				target.SetModifier("dex", -5)
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func pray(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("pray", "300", 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("pray", "pray", text.Red+"Your faith fills your being.\n")
				target.SetModifier("pie", 5)
			},
			func() {
				target.FlagOffAndMsg("pray", "pray", text.Cyan+"Your piousness returns to normal.\n")
				target.SetModifier("pie", -5)
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func healstam(caller interface{}, target interface{}, magnitude int) string {
	switch caller := caller.(type) {
	case *Character:
		damage := 0
		if caller.CheckFlag("casting") {
			divinityLevel := config.HealingSkill[config.WeaponLevel(caller.Skills[10].Value, caller.Class)]
			damage = int((float64(caller.Pie.Current) * config.PieHealMod) + float64(utils.Roll(10, 1, 0))*(1+float64(divinityLevel)*.01))
			damage = caller.CalcHealPenalty(damage)
		} else {
			damage = int((config.BaseDevicePiety * config.PieHealMod) + float64(utils.Roll(10, 1, 0)))
		}
		switch target := target.(type) {
		case *Character:
			healAmount := target.HealStam(damage)
			mode := 3
			if caller.CheckFlag("casting") {
				caller.AdvanceDivinity(healAmount*2, caller.Class)
				mode = 2
			}
			data.StoreCombatMetric("vigor", 1, mode, damage, damage-healAmount, healAmount, 0, caller.CharId, caller.Tier, 0, target.CharId)
			if utils.IntIn(caller.Class, []int{5, 6, 7}) {
				/*
					for _, mob := range Rooms[target.ParentId].Mobs.Contents {
						if mob.Flags["hostile"] {
							mob.AddThreatDamage(healAmount/10, caller)
						}
					}*/
				if target.Victim != nil {
					switch victim := target.Victim.(type) {
					case *Mob:
						victim.AddThreatDamage(healAmount/3, caller)
					}
				}
			}

			return text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"
		case *Mob:
			target.HealStam(damage)
			return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " stamina"
		}
		return ""
	}
	return ""
}

func healvit(caller interface{}, target interface{}, magnitude int) string {
	switch caller := caller.(type) {
	case *Character:
		damage := 0
		if caller.CheckFlag("casting") {
			divinityLevel := config.HealingSkill[config.WeaponLevel(caller.Skills[10].Value, caller.Class)]
			damage = int((float64(caller.Pie.Current) * config.PieHealMod) + float64(utils.Roll(10, 1, 0))*(1+float64(divinityLevel)*.01))
			damage = caller.CalcHealPenalty(damage)
		} else {
			damage = int((config.BaseDevicePiety * config.PieHealMod) + float64(utils.Roll(10, 1, 0)))
		}
		switch target := target.(type) {
		case *Character:
			healAmount := target.HealVital(damage)
			mode := 3
			if caller.CheckFlag("casting") {
				caller.AdvanceDivinity(healAmount*2, caller.Class)
				mode = 2
			}
			data.StoreCombatMetric("mend", 1, mode, damage, damage-healAmount, healAmount, 0, caller.CharId, caller.Tier, 0, target.CharId)
			if utils.IntIn(caller.Class, []int{5, 6, 7}) {
				/*
					for _, mob := range Rooms[target.ParentId].Mobs.Contents {
						if mob.Flags["hostile"] {
							mob.AddThreatDamage(healAmount/10, caller)
						}
					}*/
				if target.Victim != nil {
					switch victim := target.Victim.(type) {
					case *Mob:
						victim.AddThreatDamage(healAmount/3, caller)
					}
				}
			}
			return text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"
		case *Mob:
			target.HealVital(damage)
			return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " vitality."
		}
		return ""
	}
	return ""
}

func heal(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	action := ""
	if magnitude == 1 {
		damage = 20
		action = "detraumatize"
	} else {
		damage = 40
		action = "renewal"
	}
	switch caller := caller.(type) {
	case *Character:
		if caller.CheckFlag("casting") {

			divinityLevel := config.HealingSkill[config.WeaponLevel(caller.Skills[10].Value, caller.Class)]
			damage = int((float64(damage) + (float64(caller.Pie.Current) * config.PieHealMod) + float64(utils.Roll(10, 1, 0))) * (1 + float64(divinityLevel)*.01))
			damage = caller.CalcHealPenalty(damage)
		} else {
			damage += int((config.BaseDevicePiety * config.PieHealMod) + float64(utils.Roll(10, 1, 0)))
		}
		switch target := target.(type) {
		case *Character:
			stam, vit := target.Heal(damage)
			mode := 3
			if caller.CheckFlag("casting") {
				caller.AdvanceDivinity((stam*2)+(vit*2), caller.Class)
				mode = 2
			}
			data.StoreCombatMetric(action, 1, mode, damage, damage-(stam+vit), stam+vit, 0, caller.CharId, caller.Tier, 0, target.CharId)
			if utils.IntIn(caller.Class, []int{5, 6, 7}) {
				/*
					for _, mob := range Rooms[target.ParentId].Mobs.Contents {
						if mob.Flags["hostile"] {
							mob.AddThreatDamage(healAmount/10, caller)
						}
					}*/
				switch victim := target.Victim.(type) {
				case *Mob:
					victim.AddThreatDamage(stam+vit/3, caller)
				}
			}
			return text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"
		case *Mob:
			stamDam, vitDam := target.Heal(damage)
			return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
		}
		return ""
	}
	return ""
}

func restore(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.Mana.Current = target.Mana.Max
		return text.Info + "You now have " + strconv.Itoa(target.Mana.Current) + " mana" + text.Reset + "\n"
	case *Mob:
		target.Mana.Current = target.Mana.Max
		return "You cast a restore on " + target.Name + " and replenish their mana stores."
	}
	return ""

}

func healall(caller interface{}, target interface{}, magnitude int) string {
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			stam, vit := target.Heal(2000)
			mode := 3
			if caller.CheckFlag("casting") {
				caller.AdvanceDivinity((stam*2)+(vit*2), caller.Class)
				mode = 2
			}
			data.StoreCombatMetric("heal", 1, mode, stam+vit, 0, stam+vit, 0, caller.CharId, caller.Tier, 0, target.CharId)
			return text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"
		case *Mob:
			stamDam, vitDam := target.Heal(2000)
			return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
		}
		return ""
	}
	return ""
}

var magicSkillMap = map[string]int{
	"fire":  6,
	"air":   7,
	"earth": 8,
	"water": 9,
}

func spellDamage(caller interface{}, target interface{}, magnitude int, magicType string) string {
	var name string
	var intel int
	var level int
	var id int
	var callerType int
	var spellType = 2
	actualDamage := 0
	damage := 0
	switch caller := caller.(type) {
	case *Character:
		name = caller.Name
		intel = caller.Int.Current
		level = caller.Tier
		id = caller.CharId
		callerType = 0
		if !caller.CheckFlag("casting") {
			spellType = 3
		}
		actualDamage = elementalDamage(magnitude, intel)
		damage = int(float64(actualDamage) + float64(actualDamage)*float64(caller.Int.Current)*config.StatDamageMod)
		if caller.Class == 4 {
			affinityLevel := config.SpellDmgSkill[config.WeaponLevel(caller.Skills[magicSkillMap[magicType]].Value, caller.Class)]
			damage = int(float64(damage) * (1 + float64(affinityLevel)*.01))
		}
	case *Mob:
		name = caller.Name
		level = caller.Level
		id = caller.MobId
		callerType = 1
		intel = caller.Int.Current
		actualDamage = elementalDamage(magnitude, intel)
		damage = actualDamage
	}
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam, resisted := target.ReceiveMagicDamage(damage, magicType)
		data.StoreCombatMetric(magicType+"spell_"+strconv.Itoa(magnitude), 0, spellType, actualDamage+resisted, resisted, actualDamage, callerType, id, level, 0, target.CharId)
		returnString := text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality. You resisted " + strconv.Itoa(resisted) + "damage." + text.Reset
		// Reflect
		switch caller := caller.(type) {
		case *Character:
			if target.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * (float64(target.GetStat("int")) * config.ReflectDamagePerInt))
				stamDamage, vitDamage, resisted := caller.ReceiveMagicDamage(reflectDamage, magicType)
				data.StoreCombatMetric("spell_player_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 0, target.CharId, target.Tier, 0, caller.CharId)
				returnString += "\n" + text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + caller.Name + "!\n" + text.Reset
				if _, err := caller.Write([]byte(text.Red + target.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back to you!\n" + text.Reset)); err != nil {
					log.Println("Error writing to player:", err)
				}
				caller.DeathCheck(" was slain by reflection!")
			}
			return returnString
		case *Mob:
			if _, err := target.Write([]byte(text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality. You resisted " + strconv.Itoa(resisted) + "damage." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player:", err)
			}

			if target.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * (float64(target.GetStat("int")) * config.ReflectDamagePerInt))
				stamDamage, _, resisted := caller.ReceiveMagicDamage(reflectDamage, magicType)
				data.StoreCombatMetric("spell_player_reflect", 0, 0, stamDamage+resisted, resisted, stamDamage, 0, target.CharId, target.Tier, 1, caller.MobId)
				if _, err := target.Write([]byte(text.Cyan + "You reflect " + strconv.Itoa(reflectDamage) + " damage back to " + caller.Name + "!\n" + text.Reset)); err != nil {
					log.Println("Error writing to player:", err)
				}
				caller.DeathCheck(target)
			}
			return ""
		}

	case *Mob:
		damage, _, resisted := target.ReceiveMagicDamage(damage, magicType)
		data.StoreCombatMetric(magicType+"spell_"+strconv.Itoa(magnitude), 0, spellType, actualDamage+resisted, resisted, actualDamage, callerType, id, level, 0, target.MobId)
		switch caller := caller.(type) {
		case *Character:
			target.AddThreatDamage(damage, caller)
			caller.AdvanceElementalExp(int(float64(damage)/float64(target.Stam.Max)*float64(target.Experience)), magicType, caller.Class)
		}
		returnString := "Your spell struck " + target.Name + " for " + strconv.Itoa(damage) + " " + magicType + " damage. They resisted " + strconv.Itoa(resisted) + "."
		// Reflect
		switch caller := caller.(type) {
		case *Character:
			if target.CheckFlag("reflection") {
				reflectDamage := int(float64(actualDamage) * config.ReflectDamageFromMob)
				stamDamage, vitDamage, resisted := caller.ReceiveMagicDamage(reflectDamage, magicType)
				data.StoreCombatMetric("spell_mob_reflect", 0, 0, stamDamage+vitDamage+resisted, resisted, stamDamage+vitDamage, 1, target.MobId, target.Level, 0, caller.CharId)
				returnString += "\n" + text.Red + target.Name + " reflects " + strconv.Itoa(reflectDamage) + " damage back to you!"
				caller.DeathCheck(" was slain by reflection!")
			}
		case *Mob:
			log.Println("mob on mob violence not implemented yet")
		}

		return returnString
	}
	return ""
}

func firedamage(caller interface{}, target interface{}, magnitude int) string {
	return spellDamage(caller, target, magnitude, "fire")
}

func earthdamage(caller interface{}, target interface{}, magnitude int) string {
	return spellDamage(caller, target, magnitude, "earth")
}

func airdamage(caller interface{}, target interface{}, magnitude int) string {
	return spellDamage(caller, target, magnitude, "air")
}

func waterdamage(caller interface{}, target interface{}, magnitude int) string {
	return spellDamage(caller, target, magnitude, "water")
}

func elementalDamage(magnitude int, intel int) (damage int) {
	power := 0
	if magnitude == 1 {
		power = utils.Roll(3, 2, 0)
		damage = 7 + power
	} else if magnitude == 2 {
		power = utils.Roll(3, 4, 0)
		damage = 21 + power
	} else if magnitude == 3 {
		power = utils.Roll(3, 7, 0)
		damage = 42 + power
	} else if magnitude == 4 {
		power = utils.Roll(4, 10, 0)
		damage = 84 + power
	} else if magnitude == 5 {
		power = utils.Roll(5, 16, 0)
		damage = 175 + power
	} else if magnitude == 6 {
		power = utils.Roll(5, 18, 0)
		damage = 280 + power
	} else if magnitude == 7 {
		power = utils.Roll(6, 35, 0)
		damage = 350 + power
	}
	return damage
}

func light(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration = 300 + config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("light", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("light", "light_spell", text.Info+"A small orb of light flits next to you.\n")
			},
			func() {
				target.FlagOffAndMsg("light", "light_spell", text.Cyan+"The orb of light fades away\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""

}

func curepoison(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.RemoveEffect("poison")
		if _, err := target.Write([]byte(text.Bad + "Your fever subsides." + text.Reset + "\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
		return ""
	case *Mob:
		target.RemoveEffect("poison")
	}
	return ""
}

func bless(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration = 300 + config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("bless", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("bless", "bless_spell", text.Info+"The devotion to Gods fills your soul.\n")
			},
			func() {
				target.FlagOffAndMsg("bless", "bless_spell", text.Cyan+"The blessing fades from you.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func protection(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration = 300 + config.IntSpellEffectDuration*caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("protection", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("protection", "protection_spell", text.Info+"Your aura flows from you, protecting you. \n")
				target.SetModifier("armor", 25)
			},
			func() {
				target.FlagOffAndMsg("protection", "protection_spell", text.Cyan+"Your aura returns to normal.\n")
				target.SetModifier("armor", -25)
			})
		return ""
	case *Mob:
		target.ApplyEffect("protection", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("protection", "protection_spell")
				target.Armor += 25
			},
			func() {
				target.FlagOff("protection", "protection_spell")
				target.Armor -= 25
			})
		return ""
	}
	return ""
}

func invisibility(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("invisibility", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("invisibility", "invisibility_spell", text.Info+"Light flows around you. \n")
			},
			func() {
				target.FlagOffAndMsg("invisibility", "invisibility_spell", text.Cyan+"The cloak falls and you become visible.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func detectInvisible(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("detect-invisible", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("detect-invisible", "detectinvisible_spell", text.Info+"Your senses are magnified, detecting the unseen.\n")
			},
			func() {
				target.FlagOffAndMsg("detect-invisible", "detectinvisible_spell", text.Cyan+"Your invisibility detection fades away.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func teleport(caller interface{}, target interface{}, magnitude int) string {
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			if !Rooms[caller.ParentId].Flags["no_teleport"] {
				if caller == target {
					return "$CRIPT $TELEPORT"

				} else {
					return "$CRIPT $TELEPORT " + target.Name
				}
			} else {
				return "Oppressive magical energies prevent you from teleporting."
			}
		case *Mob:
			return "$CRIPT $TELEPORT " + target.Name + " " + strconv.Itoa(Rooms[caller.ParentId].Mobs.GetNumber(target))
		}

	case *Mob:
		switch target := target.(type) {
		case *Character:
			return "$CRIPT $TELEPORT " + target.Name
		case *Mob:
			// TODO: Should mobs really bother teleporting other mobs?
			return ""
		}
	}
	return ""
}

func stun(caller interface{}, target interface{}, magnitude int) string {
	duration := 15
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			return "No PVP yet"
			/*
				diff := ((caller.GetStat("int") - target.GetStat("int")) / 5) * 10
				chance := 30 + diff
				if utils.Roll(100, 1, 0) > chance {
				}
			*/
		case *Mob:
			diff := (caller.Tier - target.Level) * 5
			chance := 10 + diff
			if utils.Roll(100, 1, 0) > chance {
				return "You failed to stun " + target.Name
			} else {
				target.Stun(duration)
				return "You stunned " + target.Name
			}
		}

	case *Mob:
		switch target := target.(type) {
		case *Character:
			diff := (caller.Level - target.Tier) * 5
			chance := 10 + diff
			if utils.Roll(100, 1, 0) > chance {
				if _, err := target.Write([]byte(text.Info + caller.Name + " failed to stun you." + text.Reset + "\n")); err != nil {
					log.Println("Error writing to player:", err)
				}
			} else {
				if _, err := target.Write([]byte(text.Bad + caller.Name + " stunned you." + text.Reset + "\n")); err != nil {
					log.Println("Error writing to player:", err)
				}
				target.SetTimer("global", 20)
			}
		case *Mob:
			// Mobs stun mobs?
			return ""
		}
	}
	return ""
}

func recall(caller interface{}, target interface{}, magnitude int) string {
	switch caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			return "$CRIPT $TELEPORTTO " + target.Name + " " + RecallRoom
		case *Mob:
			return "Cannot be cast on a mob."
		}

	case *Mob:
		return ""
	}
	return ""
}

func summon(caller interface{}, target interface{}, magnitude int) string {
	if caller == target {
		return "You cannot cast summon on yourself."
	}
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			return "$CRIPT $TELEPORTTO " + target.Name + " " + strconv.Itoa(caller.ParentId)
		case *Mob:
			return "Cannot be cast on a mob."
		}

	case *Mob:
		return ""
	}
	return ""
}

func wizardwalk(caller interface{}, target interface{}, magnitude int) string {
	if caller == target {
		return "Why would you walk to yourself?"
	}
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			return "$CRIPT $TELEPORTTO " + caller.Name + " " + strconv.Itoa(target.ParentId)
		case *Mob:
			return "Cannot be cast on a mob."
		}
	case *Mob:
		return ""
	}
	return ""
}

func levitate(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("levitate", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("levitate", "levitate_spell", text.Info+"You lift off of your feet. \n")
			},
			func() {
				target.FlagOffAndMsg("levitate", "levitate_spell", text.Cyan+"Your feet touch the ground as the spell fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistfire(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-fire", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-fire", "resistfire_spell", text.Info+"Magical shielding springs up around you protecting you from fire. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-fire", "resistfire_spell", text.Cyan+"The magical cloak protecting you from fire fades. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("resist-fire", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("resist-fire", "resistfire_spell")
			},
			func() {
				target.FlagOff("resist-fire", "resistfire_spell")
			})
	}
	return ""
}

func resistmagic(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-magic", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-magic", "resistmagic_spell", text.Info+"Magical shielding springs up around you protecting you from magic. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-magic", "resistmagic_spell", text.Cyan+"The magical cloak protecting you from magic fades. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("resist-magic", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("resist-magic", "resistmagic_spell")
			},
			func() {
				target.FlagOff("resist-magic", "resistmagic_spell")
			})
		return ""
	}
	return ""
}

func resistair(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-air", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-air", "resistair_spell", text.Info+"Magical shielding springs up around you protecting you from air. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-air", "resistair_spell", text.Cyan+"The magical cloak protecting you from air fades. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("resist-air", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("resist-air", "resistair_spell")
			},
			func() {
				target.FlagOff("resist-air", "resistair_spell")
			})
		return ""
	}
	return ""
}

func resistwater(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-water", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-water", "resistwater_spell", text.Info+"Magical shielding springs up around you protecting you from water. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-water", "resistwater_spell", text.Cyan+"The magical cloak protecting you from water fades. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("resist-water", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("resist-water", "resistwater_spell")
			},
			func() {
				target.FlagOff("resist-water", "resistwater_spell")
			})
		return ""
	}
	return ""
}

func resistearth(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-earth", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-earth", "resistearth_spell", text.Info+"Magical shielding springs up around you protecting you from earth. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-earth", "resistearth_spell", text.Cyan+"The magical cloak protecting you from earth fades. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("resist-water", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("resist-earth", "resistearth_spell")
			},
			func() {
				target.FlagOff("resist-earth", "resistearth_spell")
			})
		return ""
	}
	return ""
}

func removedisease(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.RemoveEffect("disease")
		if _, err := target.Write([]byte(text.Bad + "The affliction is purged." + text.Reset + "\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
		return ""
	case *Mob:
		target.RemoveEffect("disease")
	}
	return ""

}

func cureblindness(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.RemoveEffect("blind")
		if _, err := target.Write([]byte(text.Bad + "Your vision returns." + text.Reset + "\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
		return ""
	case *Mob:
		target.RemoveEffect("blind")
	}
	return ""

}

func inertialbarrier(caller interface{}, target interface{}, magnitude int) string {
	if caller != target {
		return "You can only cast this spell on yourself."
	}
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}

	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("inertial-barrier", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("inertial-barrier", "inertialbarrier_spell", text.Info+"A dampening barrier forms around you.\n")
			},
			func() {
				target.FlagOffAndMsg("inertial-barrier", "inertialbarrier_spell", text.Cyan+"The dampening barrier falls away. \n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("inertial-barrier", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("inertial-barrier", "inertialbarrier_spell")
			},
			func() {
				target.FlagOff("inertial-barrier", "inertialbarrier_spell")
			})
		return ""
	}
	return ""
}

func surge(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("surge", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("surge", "surge_spell", text.Info+"You feel the power flow into you.\n")
			},
			func() {
				target.FlagOffAndMsg("surge", "surge_spell", text.Cyan+"The surge of power fades from you.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistpoison(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-poison", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-poison", "resistpoison_spell", text.Info+"Your blood thickens, protecting you from poison. \n")
			},
			func() {
				target.FlagOffAndMsg("resist-poison", "resistpoison_spell", text.Cyan+"Your blood returns to normal, your magical protection from poison fading. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resilientaura(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resilient-aura", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resilient-aura", "resilientaura_spell", text.Info+"A magical shield forms around your gear protecting it from damage.\n")
			},
			func() {
				target.FlagOffAndMsg("resilient-aura", "resilientaura_spell", text.Cyan+"The magical shield around your equipment fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistdisease(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-disease", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-disease", "resistdisease_spell", text.Info+"Your blood heats, protecting you from disease.\n")
			},
			func() {
				target.FlagOffAndMsg("resist-disease", "resistdisease_spell", text.Cyan+"Your magical fever fades, removing your resistance to disease.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func reflection(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	case *Mob:
		duration += 300
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("reflection", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("reflection", "reflect_spell", text.Info+"A mirrored shell forms around you and fades from view.\n")
			},
			func() {
				target.FlagOffAndMsg("reflection", "reflect_spell", text.Cyan+"The mirrored shell shatters, and falls away.\n")
			})
		return ""
	case *Mob:
		target.ApplyEffect("reflection", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOn("reflection", "reflect_spell")
			},
			func() {
				target.FlagOff("reflection", "reflect_spell")
			})
		return ""
	}
	return ""
}

func dodge(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("dodge", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("dodge", "dodge_spell", text.Info+"Your reflexes quicken.\n")
			},
			func() {
				target.FlagOffAndMsg("dodge", "dodge_spell", text.Cyan+"Your magically quickened reflexes return to normal.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistacid(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration * caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resist-acid", strconv.Itoa(duration), 0, 0,
			func(triggers int) {
				target.FlagOnAndMsg("resist-acid", "resistacid_spell", text.Info+"A thick mucous coats your skin protecting you from acid damage.\n")
			},
			func() {
				target.FlagOffAndMsg("resist-acid", "resistacid_spell", text.Cyan+"The mucous falls away.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func disruptmagic(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		if target.CheckFlag("resist-magic") {
			// 50:50 chance to resist the disrupt spell
			if utils.Roll(100, 1, 0) > 50 {
				if _, err := target.Write([]byte(text.Info + "You resist the disruption to your magic.\n")); err != nil {
					log.Println("Error writing to player:", err)
				}
				return ""
			}
		}
		var spellEffects []string
		for k := range target.Effects {
			if utils.StringIn(k, SupportSpells) {
				spellEffects = append(spellEffects, k)
			}
		}
		if len(spellEffects) > 0 {
			chosenSpell := spellEffects[rand.Intn(len(spellEffects))]
			target.RemoveEffect(chosenSpell)
			if _, err := target.Write([]byte(text.Bad + "The disruptive magic removes " + chosenSpell + " from you.\n")); err != nil {
				log.Println("Error writing to player:", err)
			}
			return ""
		}
		if _, err := target.Write([]byte(text.Info + "The disruptive magic has no effect on you.\n")); err != nil {
			log.Println("Error writing to player:", err)
		}
	case *Mob:
		if target.CheckFlag("resist-magic") {
			// 50:50 chance to resist the disrupt spell
			if utils.Roll(100, 1, 0) > 50 {
				return text.Bad + target.Name + " resisted the disruption from your spell.\n"
			}
		}
		var spellEffects []string
		for k := range target.Effects {
			if utils.StringIn(k, SupportSpells) {
				spellEffects = append(spellEffects, k)
			}
		}
		if len(spellEffects) > 0 {
			chosenSpell := spellEffects[rand.Intn(len(spellEffects))]
			target.RemoveEffect(chosenSpell)
			switch caller := caller.(type) {
			case *Character:
				if _, err := caller.Write([]byte(text.Bad + "Your disruptive magic removes " + chosenSpell + " from " + target.Name + " .\n")); err != nil {
					log.Println("Error writing to player:", err)
				}
			}
			return ""
		}
		switch caller := caller.(type) {
		case *Character:
			if _, err := caller.Write([]byte(text.Bad + "Your disruptive magic has no effect on " + target.Name + " .\n")); err != nil {
				log.Println("Error writing to player:", err)
			}
		}
	}

	return ""

}

func attraction(caller interface{}, target interface{}, magnitude int) string {
	switch caller := caller.(type) {
	case *Character:
		go Script(caller, "$ATTRACT")
		return text.Cyan + "Light coalesces into a vaguely sprite shape and darts around the area creating as much commotion as possible, then fades away.\n"
	}
	return ""
}

/*
func embolden(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		if target.HasEffect("fear") {
			target.RemoveEffect("fear")
			target.Write([]byte(text.Bad + "Your fears subside, and your resolve itensifies." + text.Reset + "\n"))
			return ""
		}
		target.Write([]byte(text.Bad + "You are unaffected by the embolden spell." + text.Reset + "\n"))
	case *Mob:
		target.RemoveEffect("fear")
	}
	return ""
}

func polymorph(caller interface{}, target interface{}, magnitude int) string { return "" }

func removecurse(caller interface{}, target interface{}, magnitude int) string {
	return ""
}
*/
