package objects

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/text"
	"math/rand"
	"strconv"
	"time"
)

var (
	teleportTable = []int{117}
	defaultDuration = 600
)

var CharEffects = map[string]func(target *Character, modifiers map[string]interface{}) string{
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
	"cure-blindness":   cureblindness,
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

var MobEffects = map[string]func(target *Mob, modifiers map[string]interface{}) string{
	"heal-stam":        mobheal,
	"heal-vit":         mobheal,
	"heal":             mobheal,
	"restore":          mobrestore,
	"heal-all":         mobhealall,
	"fire-damage":      mobfiredamage,
	"earth-damage":     mobearthdamage,
	"air-damage":       mobairdamage,
	"water-damage":     mobwaterdamage,
	"light":            moblight,
	"curepoison":       mobcurepoison,
	"bless":            mobbless,
	"protection":       mobprotection,
	"invisibility":     mobinvisibility,
	"detect-invisible": mobdetect_invisible,
	"teleport":         mobteleport,
	"stun":             mobstun,
	"enchant":          mobenchant,
	"recall":           mobrecall,
	"summon":           mobsummon,
	"wizard-walk":      mobwizardwalk,
	"levitate":         moblevitate,
	"resist-fire":      mobresistfire,
	"resist-magic":     mobresistmagic,
	"remove-curse":     mobremovecurse,
	"resist-air":       mobresistair,
	"resist-water":     mobresistwater,
	"resist-earth":     mobresistearth,
	"clairvoyance":     mobclairvoyance,
	"remove-disease":   mobremovedisease,
	"cure-blindness":   mobcureblindness,
	"polymorph":        mobpolymorph,
	"attraction":       mobattraction,
	"inertial-barrier": mobinertialbarrier,
	"surge":            mobsurge,
	"resist-poison":    mobresistpoison,
	"resilient-aura":   mobresilientaura,
	"resist-disease":   mobresistdisease,
	"disrupt-magic":    mobdisruptmagic,
	"reflection":       mobreflection,
	"dodge":            mobdodge,
	"resist-acid":      mobresistacid,
	"embolden":         mobembolden,
}

/*
A robust casting system requires multiple entry and pass around points for casting on mobs and players
as they are technically handled differently.  It's easier to redirect and cast, and then as needed
create a spell invocation for both target types
 */

func MobCast(caller *Mob, target interface{}, spell string, modifiers map[string]interface{}) {
	// Pass some of the player data to the spell
	modifiers["name"] = caller.Name
	modifiers["tier"] = caller.Level
	modifiers["int"] = caller.Int.Current
	modifiers["str"] = caller.Str.Current
	modifiers["dex"] = caller.Dex.Current
	modifiers["pie"] = caller.Pie.Current
	modifiers["con"] = caller.Con.Current
	modifiers["multiplier"] = 1

	switch v := target.(type) {
	case *Character:
		CharEffects[spell](target.(*Character), modifiers)
	case *Mob:
		MobEffects[spell](target.(*Mob), modifiers)
	default:
		fmt.Printf("Strange behavior attempting to player cast a spell on %T!\n", v)
	}
}

func PlayerCast(caller *Character, target interface{}, spell string, modifiers map[string]interface{}) string {
	// Pass some of the player data to the spell
	modifiers["name"] = caller.Name
	modifiers["tier"] = caller.Tier
	modifiers["int"] = caller.GetStat("int")
	modifiers["str"] = caller.GetStat("str")
	modifiers["dex"] = caller.GetStat("dex")
	modifiers["pie"] = caller.GetStat("pie")
	modifiers["con"] = caller.GetStat("con")
	modifiers["multiplier"] = 1

	switch v := target.(type) {
	case *Character:
		return CharEffects[spell](target.(*Character), modifiers)
	case *Mob:
		modifiers["multiplier"] = caller.GetSpellMultiplier()
		return MobEffects[spell](target.(*Mob), modifiers)
	default:
		fmt.Printf("Strange behavior attempting to player cast a spell on %T!\n", v)
		return "The spell fizzles."
	}
}


func berserk(target *Character, modifiers map[string]interface{}) string{
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
	return ""
}

func haste(target *Character, modifiers map[string]interface{}) string {
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
}

func pray(target *Character, modifiers map[string]interface{}) string {
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
}

func healstam(target *Character, modifiers map[string]interface{}) string {
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
	target.HealStam(damage)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " stamina"
}

func healvit(target *Character, modifiers map[string]interface{}) string {
	damage := 10
	target.HealVital(damage)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(damage)  +" vitality."
}

func heal(target *Character, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 50
	} else {
		damage = 100
	}
	stamDam, vitDam := target.Heal(damage)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."
}

func restore(target *Character, modifiers map[string]interface{}) string {
	target.Mana.Current = target.Mana.Max
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Mana.Current) + " mana" + text.Reset + "\n"))
	return "You cast a restore on " + target.Name + " and replenish their mana stores."
}

func healall(target *Character, modifiers map[string]interface{}) string {
	stamDam, vitDam := target.Heal(2000)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."

}

func firedamage(target *Character, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		// 7 + (45/int)*(roll 2*3)

		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	stamDam, vitDam := target.ReceiveDamage(damage)
	target.Write([]byte(text.Bad + modifiers["name"].(string) + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
}

func earthdamage(target *Character, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	stamDam, vitDam := target.ReceiveDamage(damage)
	target.Write([]byte(text.Bad + modifiers["name"].(string) + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " earth damage."
}

func airdamage(target *Character, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	stamDam, vitDam := target.ReceiveDamage(damage)
	target.Write([]byte(text.Bad + modifiers["name"].(string) + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " air damage."
}

func waterdamage(target *Character, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	stamDam, vitDam := target.ReceiveDamage(damage)
	target.Write([]byte(text.Bad + modifiers["name"].(string) + "'s spell struck you for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."  + text.Reset + "\n"))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
}

func light(target *Character, modifiers map[string]interface{}) string {
	duration := 300 + config.IntSpellEffectDuration*modifiers["int"].(int)
	target.ApplyEffect("light", strconv.Itoa(duration), "0",
		func() {
			target.ToggleFlagAndMsg("light", "light_spell", text.Info +"A small orb of light flits next to you.\n")
		},
		func() {
			target.ToggleFlagAndMsg("light", "light_spell", text.Cyan+"The orb of light fades away\n")
		})
	return ""
}

func curepoison(target *Character, modifiers map[string]interface{}) string {
	target.RemoveEffect("poison")
	return ""
}

func bless(target *Character, modifiers map[string]interface{}) string {
	duration := 300 + config.IntSpellEffectDuration*modifiers["int"].(int)
	target.ApplyEffect("bless", strconv.Itoa(duration),"0",
		func() {
			target.ToggleFlagAndMsg("bless", "bless_spell", text.Info +"The devotion to Gods fills your soul.\n")
		},
		func() {
			target.ToggleFlagAndMsg("bless", "bless_spell", text.Cyan+"The blessing fades from you.\n")
		})
	return ""
}

func protection(target *Character, modifiers map[string]interface{}) string {
	duration := 300 + config.IntSpellEffectDuration*modifiers["int"].(int)
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
}

func invisibility(target *Character, modifiers map[string]interface{}) string {
	duration := 30 + (config.IntSpellEffectDuration/2)*modifiers["int"].(int)
	target.ApplyEffect("invisibility", strconv.Itoa(duration), "0",
		func() {
			target.ToggleFlagAndMsg("invisibility", "invisibility_spell", text.Info +"Light flows around you. \n")
		},
		func() {
			target.ToggleFlagAndMsg("invisibility", "invisibility_spell", text.Cyan+"The cloak falls and you become visible.\n")
		})
	return ""
}

func detect_invisible(target *Character, modifiers map[string]interface{}) string {
	duration := 30 + (config.IntSpellEffectDuration/2)*modifiers["int"].(int)
	target.ApplyEffect("detectinvisibile", strconv.Itoa(duration), "0",
		func() {
			target.ToggleFlagAndMsg("detectinvisibile", "detectinvisibile_spell", text.Info +"Your senses are magnified, detecting the unseen.\n")
		},
		func() {
			target.ToggleFlagAndMsg("detectinvisibile", "detectinvisibile_spell", text.Cyan+"Your invisibility detection fades away.\n")
		})
	return ""
}

func teleport(target *Character, modifiers map[string]interface{}) string {
	rand.Seed(time.Now().Unix())
	newRoom := teleportTable[rand.Intn(len(teleportTable))]
	return "$CRIPT $TELEPORT " + strconv.Itoa(newRoom)
}

func stun(target *Character, modifiers map[string]interface{}) string { return "" }

func enchant(target *Character, modifiers map[string]interface{}) string { return "" }

func recall(target *Character, modifiers map[string]interface{}) string { return "" }

func summon(target *Character, modifiers map[string]interface{}) string { return "" }

func wizardwalk(target *Character, modifiers map[string]interface{}) string { return "" }

func levitate(target *Character, modifiers map[string]interface{}) string { return "" }

func resistfire(target *Character, modifiers map[string]interface{}) string { return "" }

func resistmagic(target *Character, modifiers map[string]interface{}) string { return "" }

func removecurse(target *Character, modifiers map[string]interface{}) string { return "" }

func resistair(target *Character, modifiers map[string]interface{}) string { return "" }

func resistwater(target *Character, modifiers map[string]interface{}) string { return "" }

func resistearth(target *Character, modifiers map[string]interface{}) string { return "" }

func clairvoyance(target *Character, modifiers map[string]interface{}) string { return "" }

func removedisease(target *Character, modifiers map[string]interface{}) string { return "" }

func cureblindness(target *Character, modifiers map[string]interface{}) string { return "" }

func polymorph(target *Character, modifiers map[string]interface{}) string { return "" }

func attraction(target *Character, modifiers map[string]interface{}) string { return "" }

func inertialbarrier(target *Character, modifiers map[string]interface{}) string { return "" }

func surge(target *Character, modifiers map[string]interface{}) string { return "" }

func resistpoison(target *Character, modifiers map[string]interface{}) string { return "" }

func resilientaura(target *Character, modifiers map[string]interface{}) string { return "" }

func resistdisease(target *Character, modifiers map[string]interface{}) string { return "" }

func disruptmagic(target *Character, modifiers map[string]interface{}) string { return "" }

func reflection(target *Character, modifiers map[string]interface{}) string { return "" }

func dodge(target *Character, modifiers map[string]interface{}) string { return "" }

func resistacid(target *Character, modifiers map[string]interface{}) string { return "" }

func embolden(target *Character, modifiers map[string]interface{}) string {
	target.RemoveEffect("fear")
	return "true"
}

func mobheal(target *Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 50
	} else {
		damage = 100
	}
	target.Heal(damage)
	return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " health"
}

func mobhealall(target *Mob, modifiers map[string]interface{}) string {
	target.Heal(target.Stam.Max)
	return "You heal " + target.Name + " for 2000 health"
}

func mobrestore(target *Mob, modifiers map[string]interface{}) string {
	target.Mana.Current = target.Mana.Max
	return "You cast a restore on " + target.Name + " and replenish their mana stores."
}

func mobfiredamage(target *Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	target.ReceiveDamage(damage*modifiers["multiplier"].(int))
	target.AddThreatDamage(damage*modifiers["multiplier"].(int), modifiers["name"].(string))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
}

func mobearthdamage(target *Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	target.ReceiveDamage(damage*modifiers["multiplier"].(int))
	target.AddThreatDamage(damage*modifiers["multiplier"].(int), modifiers["name"].(string))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " earth damage."
}

func mobairdamage(target *Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	target.ReceiveDamage(damage*modifiers["multiplier"].(int))
	target.AddThreatDamage(damage*modifiers["multiplier"].(int), modifiers["name"].(string))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " air damage."
}

func mobwaterdamage(target *Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 10
	} else if modifiers["magnitude"].(int) == 2 {
		damage = 30
	} else if modifiers["magnitude"].(int) == 3 {
		damage = 60
	} else if modifiers["magnitude"].(int) == 4 {
		damage = 120
	} else if modifiers["magnitude"].(int) == 5 {
		damage = 250
	} else if modifiers["magnitude"].(int) == 6 {
		damage = 400
	} else if modifiers["magnitude"].(int) == 7 {
		damage = 600
	}
	target.ReceiveDamage(damage*modifiers["multiplier"].(int))
	target.AddThreatDamage(damage*modifiers["multiplier"].(int), modifiers["name"].(string))
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " water damage."
}

func moblight(target *Mob, modifiers map[string]interface{}) string {
	return ""
}

func mobcurepoison(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobbless(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobprotection(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobinvisibility(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobdetect_invisible(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobteleport(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobstun(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobenchant(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobrecall(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobsummon(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobwizardwalk(target *Mob, modifiers map[string]interface{}) string { return "" }

func moblevitate(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistfire(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistmagic(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobremovecurse(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistair(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistwater(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistearth(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobclairvoyance(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobremovedisease(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobcureblindness(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobpolymorph(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobattraction(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobinertialbarrier(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobsurge(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistpoison(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresilientaura(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistdisease(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobdisruptmagic(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobreflection(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobdodge(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobresistacid(target *Mob, modifiers map[string]interface{}) string { return "" }

func mobembolden(target *Mob, modifiers map[string]interface{}) string {
	target.RemoveEffect("fear")
	return "true"
}
