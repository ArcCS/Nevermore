package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

var (
	TeleportTable = []int{117}
	defaultDuration = 600
	lvl5RecallRoom = "77"
	RecallRoom = "77"
)

var Effects = map[string]func(caller interface{}, target interface{}, magnitude int) string{
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
	"detect-invisible": detect_invisible,
	"teleport":         teleport,
	"stun":             stun,
	"enchant":          enchant,
	"recall":           recall,
	"summon":           summon,
	"wizard-walk":      wizardwalk,
	"levitate":         levitate,
	"resist-fire":      resistfire,
	"resist-magic":     resistmagic,
	"remove-curse":     removecurse,
	"resist-air":       resistair,
	"resist-water":     resistwater,
	"resist-earth":     resistearth,
	"clairvoyance":     clairvoyance,
	"remove-disease":   removedisease,
	"remove-blindness":   cureblindness,
	"polymorph":        polymorph,
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
	"embolden":         embolden,
}

func Cast(caller interface{}, target interface{}, spell string, magnitude int) string {
		return Effects[spell](caller, target, magnitude)
}


func berserk(caller interface{}, target interface{}, magnitude int) string{
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("berserk", "60", "0",
			func() {
				target.ToggleFlagAndMsg("berserk", "berserk", text.Red+"The red rage grips you!!!\n")
				target.SetModifier("str", 5)
				target.SetModifier("base_damage",  target.GetStat("str") * config.CombatModifiers["berserk"])
			},
			func() {
				target.ToggleFlagAndMsg("berserk", "berserk", text.Cyan+"The tension releases and your rage fades...\n")
				target.SetModifier("base_damage",  -target.GetStat("str") * config.CombatModifiers["berserk"])
				target.SetModifier("str", -5)
			})
	case *Mob:
		return ""
	}
	return ""
}

func haste(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("haste", "60", "0",
			func() {
				target.ToggleFlagAndMsg("haste", "haste", text.Info+"Your muscles tighten and your reflexes hasten!!!\n")
				target.SetModifier("dex", 5)
			},
			func() {
				target.ToggleFlagAndMsg("haste", "haste", text.Cyan+"Your reflexes return to normal.\n")
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
		target.ApplyEffect("pray", "300", "0",
			func() {
				target.ToggleFlagAndMsg("pray", "pray", text.Red+"Your faith fills your being.\n")
				target.SetModifier("pie", 5)
			},
			func() {
				target.ToggleFlagAndMsg("pray", "pray", text.Cyan+"Your piousness returns to normal.\n")
				target.SetModifier("pie", -5)
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func healstam(caller interface{}, target interface{}, magnitude int) string {
	damage := 10
	switch caller := caller.(type) {
	case *Character:
		damage = damage + int(float64(caller.Pie.Current)*config.PieHealMod)
	}

	switch target := target.(type) {
	case *Character:
		target.HealStam(damage)
		target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
		return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " stamina"
	case *Mob:
		target.HealStam(damage)
		return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " stamina"
	}
	return ""
}

func healvit(caller interface{}, target interface{}, magnitude int) string {
	damage := 10
	switch caller := caller.(type) {
	case *Character:
		damage = damage + int(float64(caller.Pie.Current)*config.PieHealMod)
	}

	switch target := target.(type) {
	case *Character:
		target.HealVital(damage)
		target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
		return "You heal " + target.Name + " for " + strconv.Itoa(damage)  +" vitality."
	case *Mob:
		target.HealVital(damage)
		return "You heal " + target.Name + " for " + strconv.Itoa(damage)  +" vitality."
	}
	return ""

}

func heal(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	if magnitude == 1 {
		damage = 50
	} else {
		damage = 100
	}
	switch caller := caller.(type) {
	case *Character:
		damage = damage + int(float64(caller.Pie.Current)*config.PieHealMod)
	}

	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.Heal(damage)
		target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
		return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
	case *Mob:
		stamDam, vitDam := target.Heal(damage)
		return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
	}
	return ""

}

func restore(caller interface{}, target interface{}, magnitude int) string {
	if caller == target {
		return "You can only cast this spell on others."
	}
	/*TODO: Restore this after class props are implemented
	switch caller := caller.(type) {

	case *Character:
		caller.ClassProps["restore"] -= 1
	}
	*/
	switch target := target.(type) {
	case *Character:
		target.Mana.Current = target.Mana.Max
		target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Mana.Current) + " mana" + text.Reset + "\n"))
		return "You cast a restore on " + target.Name + " and replenish their mana stores."
	case *Mob:
		target.Mana.Current = target.Mana.Max
		return "You cast a restore on " + target.Name + " and replenish their mana stores."
	}
	return ""

}

func healall(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.Heal(2000)
		target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
		return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
	case *Mob:
		stamDam, vitDam := target.Heal(2000)
		return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
	}
	return ""
}

func firedamage(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	if magnitude == 1 {
		// 7 + (45/int)*(roll 2*3)
		damage = 10
	} else if magnitude == 2 {
		damage = 30
	} else if magnitude == 3 {
		damage = 60
	} else if magnitude == 4 {
		damage = 120
	} else if magnitude == 5 {
		damage = 250
	} else if magnitude == 6 {
		damage = 400
	} else if magnitude == 7 {
		damage = 600
	}
	name := ""
	switch caller := caller.(type) {
	case *Character:
		name = caller.Name
	case *Mob:
		name = caller.Name
	}
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.ReceiveDamage(damage)
		target.Write([]byte(text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
	case *Mob:
		target.ReceiveDamage(damage)
		switch caller := caller.(type) {
		case *Character:
			target.AddThreatDamage(damage, caller)
		}
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
	}
	return ""
}

func earthdamage(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	if magnitude == 1 {
		// 7 + (45/int)*(roll 2*3)
		damage = 10
	} else if magnitude == 2 {
		damage = 30
	} else if magnitude == 3 {
		damage = 60
	} else if magnitude == 4 {
		damage = 120
	} else if magnitude == 5 {
		damage = 250
	} else if magnitude == 6 {
		damage = 400
	} else if magnitude == 7 {
		damage = 600
	}
	name := ""
	switch caller := caller.(type) {
	case *Character:
		name = caller.Name
	case *Mob:
		name = caller.Name
	}
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.ReceiveDamage(damage)
		target.Write([]byte(text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " earth damage."
	case *Mob:
		target.ReceiveDamage(damage)
		switch caller := caller.(type) {
		case *Character:
			target.AddThreatDamage(damage, caller)
		}
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " earth damage."
	}
	return ""
}

func airdamage(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	if magnitude == 1 {
		// 7 + (45/int)*(roll 2*3)
		damage = 10
	} else if magnitude == 2 {
		damage = 30
	} else if magnitude == 3 {
		damage = 60
	} else if magnitude == 4 {
		damage = 120
	} else if magnitude == 5 {
		damage = 250
	} else if magnitude == 6 {
		damage = 400
	} else if magnitude == 7 {
		damage = 600
	}
	name := ""
	switch caller := caller.(type) {
	case *Character:
		name = caller.Name
	case *Mob:
		name = caller.Name
	}
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.ReceiveDamage(damage)
		target.Write([]byte(text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " air damage."
	case *Mob:
		target.ReceiveDamage(damage)
		switch caller := caller.(type) {
		case *Character:
			target.AddThreatDamage(damage, caller)
		}
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " air damage."
	}
	return ""
}

func waterdamage(caller interface{}, target interface{}, magnitude int) string {
	damage := 0
	if magnitude == 1 {
		// 7 + (45/int)*(roll 2*3)
		damage = 10
	} else if magnitude == 2 {
		damage = 30
	} else if magnitude == 3 {
		damage = 60
	} else if magnitude == 4 {
		damage = 120
	} else if magnitude == 5 {
		damage = 250
	} else if magnitude == 6 {
		damage = 400
	} else if magnitude == 7 {
		damage = 600
	}
	name := ""
	switch caller := caller.(type) {
	case *Character:
		name = caller.Name
	case *Mob:
		name = caller.Name
	}
	switch target := target.(type) {
	case *Character:
		stamDam, vitDam := target.ReceiveDamage(damage)
		target.Write([]byte(text.Bad + name + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " water damage."
	case *Mob:
		target.ReceiveDamage(damage)
		switch caller := caller.(type) {
		case *Character:
			target.AddThreatDamage(damage, caller)
		}
		return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " water damage."
	}
	return ""
}

func light(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration = 300 + config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("light", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("light", "light_spell", text.Info +"A small orb of light flits next to you.\n")
			},
			func() {
				target.ToggleFlagAndMsg("light", "light_spell", text.Cyan+"The orb of light fades away\n")
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
		target.Write([]byte(text.Bad + "Your fever subsides." + text.Reset + "\n"))
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
		target.ApplyEffect("bless", strconv.Itoa(duration),"0",
			func() {
				target.ToggleFlagAndMsg("bless", "bless_spell", text.Info +"The devotion to Gods fills your soul.\n")
			},
			func() {
				target.ToggleFlagAndMsg("bless", "bless_spell", text.Cyan+"The blessing fades from you.\n")
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
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("protection", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("protection", "protection_spell", text.Info +"Your aura flows from you, protecting you. \n")
				target.SetModifier("armor", 25)
			},
			func() {
				target.ToggleFlagAndMsg("protection", "protection_spell", text.Cyan+"Your aura returns to normal.\n")
				target.SetModifier("armor", -25)
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func invisibility(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("invisibility", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("invisibility", "invisibility_spell", text.Info +"Light flows around you. \n")
			},
			func() {
				target.ToggleFlagAndMsg("invisibility", "invisibility_spell", text.Cyan+"The cloak falls and you become visible.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func detect_invisible(caller interface{}, target interface{}, magnitude int) string {
	duration := 300
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("detectinvisibile", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("detectinvisibile", "detectinvisibile_spell", text.Info +"Your senses are magnified, detecting the unseen.\n")
			},
			func() {
				target.ToggleFlagAndMsg("detectinvisibile", "detectinvisibile_spell", text.Cyan+"Your invisibility detection fades away.\n")
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
				if caller == target {
					return "$CRIPT $TELEPORT"

				}else {
					return "$CRIPT $TELEPORT " + target.Name
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			return "No PVP yet"
			diff := ((caller.GetStat("int") - target.GetStat("int")) / 5) * 10
			chance := 30 + diff
			if utils.Roll(100, 1, 0) > chance {
			}
		case *Mob:
			diff := (caller.Tier - target.Level) * 5
			chance := 10 + diff
			if utils.Roll(100, 1, 0) > chance {
				return "You failed to teleport " + target.Name
			}else{
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
				target.Write([]byte(text.Info + caller.Name + " failed to stun you." + text.Reset + "\n"))
			}else{
				target.Write([]byte(text.Bad + caller.Name + " stunned you." + text.Reset + "\n"))
				target.SetTimer("global", 20)
			}
		case *Mob:
			// Mobs stun mobs?  Meh maybe
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("levitate", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("levitate", "levitate_spell", text.Info +"You lift off of your feet. \n")
			},
			func() {
				target.ToggleFlagAndMsg("levitate", "levitate_spell", text.Cyan+"Your feet touch the ground as the spell fades. \n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistfire", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistfire", "resistfire_spell", text.Info +"Magical shielding springs up around you protecting you from fire. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistfire", "resistfire_spell", text.Cyan+"The magical cloak protecting you from fire fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistmagic(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistmagic", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistmagic", "resistmagic_spell", text.Info +"Magical shielding springs up around you protecting you from magic. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistmagic", "resistmagic_spell", text.Cyan+"The magical cloak protecting you from magic fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistair(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistair", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistair", "resistair_spell", text.Info +"Magical shielding springs up around you protecting you from air. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistair", "resistair_spell", text.Cyan+"The magical cloak protecting you from air fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistwater(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistwater", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistwater", "resistwater_spell", text.Info +"Magical shielding springs up around you protecting you from water. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistwater", "resistwater_spell", text.Cyan+"The magical cloak protecting you from water fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func resistearth(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistearth", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistearth", "resistearth_spell", text.Info +"Magical shielding springs up around you protecting you from earth. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistearth", "resistearth_spell", text.Cyan+"The magical cloak protecting you from earth fades. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}


func clairvoyance(caller interface{}, target interface{}, magnitude int) string {
	if caller == target {
		return "You cannot cast clairvoyance on yourself."
	}
	switch caller := caller.(type) {
	case *Character:
		switch target := target.(type) {
		case *Character:
			if target.Resist {
				// For every 5 points of int over the target there's an extra 10% chance to clairvoyance
				diff := ((caller.GetStat("int") - target.GetStat("int")) / 5) * 10
				chance := 30 + diff
				if utils.Roll(100, 1, 0) > chance {
					target.Write([]byte(text.Info + caller.Name + " failed to cast clairvoyance on you. \n" + text.Reset))
					return "You failed to cast clairvoyance on " + target.Name
				}else{
					target.Write([]byte(text.Info + caller.Name + " sees through your eyes. \n" + text.Reset))
					return Rooms[target.ParentId].Look(caller)
				}
			}
		case *Mob:
			return "Cannot be cast on a mob."
		}

	case *Mob:
		return ""
	}
	return ""
}

func removedisease(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.RemoveEffect("disease")
		target.Write([]byte(text.Bad + "The affliction is purged." + text.Reset + "\n"))
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
		target.Write([]byte(text.Bad + "Your vision returns." + text.Reset + "\n"))
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("inertialbarrier", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("inertialbarrier", "inertialbarrier_spell", text.Info +"A dampening barrier forms around you.\n")
			},
			func() {
				target.ToggleFlagAndMsg("inertialbarrier", "inertialbarrier_spell", text.Cyan+"The dampening barrier falls away. \n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func surge(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("surge", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("surge", "surge_spell", text.Info +"You feel the power flow into you.\n")
			},
			func() {
				target.ToggleFlagAndMsg("surge", "surge_spell", text.Cyan+"The surge of power fades from you.\n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistpoison", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistpoison", "resistpoison_spell", text.Info +"Your blood thickens, protecting you from poison. \n")
			},
			func() {
				target.ToggleFlagAndMsg("resistpoison", "resistpoison_spell", text.Cyan+"Your blood returns to normal, your magical protection from poison fading. \n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resilientaura", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resilientaura", "resilientaura_spell", text.Info +"A magical shield forms around your gear protecting it from damage.\n")
			},
			func() {
				target.ToggleFlagAndMsg("resilientaura", "resilientaura_spell", text.Cyan+"The magical shield around your equipment fades. \n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistdisease", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistdisease", "resistdisease_spell", text.Info +"Your blood heats, protecting you from disease.\n")
			},
			func() {
				target.ToggleFlagAndMsg("resistdisease", "resistdisease_spell", text.Cyan+"Your magical fever fades, removing your resistance to disease.\n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("reflect", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("reflect", "reflect_spell", text.Info +"A mirrored shell forms around you and fades from view.\n")
			},
			func() {
				target.ToggleFlagAndMsg("reflect", "reflect_spell", text.Cyan+"The mirrored shell shatters, and falls away.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func dodge(caller interface{}, target interface{}, magnitude int) string {
	duration := 30
	switch caller := caller.(type) {
	case *Character:
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("dodge", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("dodge", "dodge_spell", text.Info +"Your reflexes quicken.\n")
			},
			func() {
				target.ToggleFlagAndMsg("dodge", "dodge_spell", text.Cyan+"Your magically quickened reflexes return to normal.\n")
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
		duration += config.IntSpellEffectDuration*caller.Int.Current
	}
	switch target := target.(type) {
	case *Character:
		target.ApplyEffect("resistacid", strconv.Itoa(duration), "0",
			func() {
				target.ToggleFlagAndMsg("resistacid", "resistacid_spell", text.Info +"A thick mucous coats your skin protecting you from acid damage.\n")
			},
			func() {
				target.ToggleFlagAndMsg("resistacid", "resistacid_spell", text.Cyan+"The mucous falls away.\n")
			})
		return ""
	case *Mob:
		return ""
	}
	return ""
}

func embolden(caller interface{}, target interface{}, magnitude int) string {
	switch target := target.(type) {
	case *Character:
		target.RemoveEffect("fear")
		target.Write([]byte(text.Bad + "Your irrational fear vanishes." + text.Reset + "\n"))
		return ""
	case *Mob:
		target.RemoveEffect("fear")
	}
	return ""
}

func disruptmagic(caller interface{}, target interface{}, magnitude int) string {
	//TODO: make a list of disruptable spells
	//TODO: Remove one of those spells
	return "" }

func polymorph(caller interface{}, target interface{}, magnitude int) string { return "" }

func attraction(caller interface{}, target interface{}, magnitude int) string { return "" }

func removecurse(caller interface{}, target interface{}, magnitude int) string {
	//TODO: Remove Curse?
	return "" }

func enchant(caller interface{}, target interface{}, magnitude int) string { return "" }