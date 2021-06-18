package objects

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
)

var (
	TeleportTable = []int{117}
	defaultDuration = 600
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
	/*vigor devices seemed to be very low 4-10?
	mend-wounds devices were around 5-15.
	detraum devices seemed around 38-50.
	renew devices seemed around 60-80 (less sure on this)
	casting vigor with a l14 int 25 priest was around 45-49, whilst casting vigor with a l14 int 30 wiz was around 20ish
	casting detraum with l14 priest was  80-100ish
	casting renewal with l14 priest was approx 150
	casting detraum with bard was around 45?
	casting renewal with bard was approx 60-80? */
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

func stun(caller interface{}, target interface{}, magnitude int) string { return "" }

func enchant(caller interface{}, target interface{}, magnitude int) string { return "" }

func recall(caller interface{}, target interface{}, magnitude int) string { return "" }

func summon(caller interface{}, target interface{}, magnitude int) string { return "" }

func wizardwalk(caller interface{}, target interface{}, magnitude int) string { return "" }

func levitate(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistfire(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistmagic(caller interface{}, target interface{}, magnitude int) string { return "" }

func removecurse(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistair(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistwater(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistearth(caller interface{}, target interface{}, magnitude int) string { return "" }

func clairvoyance(caller interface{}, target interface{}, magnitude int) string { return "" }

func removedisease(caller interface{}, target interface{}, magnitude int) string { return "" }

func cureblindness(caller interface{}, target interface{}, magnitude int) string { return "" }

func polymorph(caller interface{}, target interface{}, magnitude int) string { return "" }

func attraction(caller interface{}, target interface{}, magnitude int) string { return "" }

func inertialbarrier(caller interface{}, target interface{}, magnitude int) string { return "" }

func surge(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistpoison(caller interface{}, target interface{}, magnitude int) string { return "" }

func resilientaura(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistdisease(caller interface{}, target interface{}, magnitude int) string { return "" }

func disruptmagic(caller interface{}, target interface{}, magnitude int) string { return "" }

func reflection(caller interface{}, target interface{}, magnitude int) string { return "" }

func dodge(caller interface{}, target interface{}, magnitude int) string { return "" }

func resistacid(caller interface{}, target interface{}, magnitude int) string { return "" }

func embolden(caller interface{}, target interface{}, magnitude int) string { return ""
}