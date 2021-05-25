package spells

import (
	"fmt"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
)

var (
	teleportTable = []int{117}
	defaultDuration = 600
)

var CharEffects = map[string]func(target *objects.Character, modifiers map[string]interface{}) string{
	"heal-stam":        healstam,
	"heal-vit":         healvit,
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
	"detect-invisible": detectinvisible,
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

var MobEffects = map[string]func(target *objects.Mob, modifiers map[string]interface{}) string{
	"heal-stam":        mobheal,
	"heal-vit":         mobheal,
	"heal":             mobheal,
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
	"detect-invisible": mobdetectinvisible,
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

func MobCast(caller *objects.Mob, target interface{}, spell string, modifiers map[string]interface{}) {
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
	case *objects.Character:
		CharEffects[spell](target.(*objects.Character), modifiers)
	case *objects.Mob:
		MobEffects[spell](target.(*objects.Mob), modifiers)
	default:
		fmt.Printf("Strange behavior attempting to player cast a spell on %T!\n", v)
	}
}

func PlayerCast(caller *objects.Character, target interface{}, spell string, modifiers map[string]interface{}) string {
	// Pass some of the player data to the spell
	modifiers["name"] = caller.Name
	modifiers["tier"] = caller.Tier
	modifiers["int"] = caller.Int.Current
	modifiers["str"] = caller.Str.Current
	modifiers["dex"] = caller.Dex.Current
	modifiers["pie"] = caller.Pie.Current
	modifiers["con"] = caller.Con.Current
	modifiers["multiplier"] = 1

	switch v := target.(type) {
	case *objects.Character:
		return CharEffects[spell](target.(*objects.Character), modifiers)
	case *objects.Mob:
		modifiers["multiplier"] = caller.GetSpellMultiplier()
		return MobEffects[spell](target.(*objects.Mob), modifiers)
	default:
		fmt.Printf("Strange behavior attempting to player cast a spell on %T!\n", v)
		return "The spell fizzles."
	}
}

func healstam(target *objects.Character, modifiers map[string]interface{}) string {
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

func healvit(target *objects.Character, modifiers map[string]interface{}) string {
	damage := 10
	target.HealVital(damage)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(damage)  +" vitality."
}

func heal(target *objects.Character, modifiers map[string]interface{}) string {
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

func healall(target *objects.Character, modifiers map[string]interface{}) string {
	stamDam, vitDam := target.Heal(2000)
	target.Write([]byte(text.Info + "You now have " + strconv.Itoa(target.Stam.Current) + " stamina and " + strconv.Itoa(target.Vit.Current) + " vitality." + text.Reset + "\n"))
	return "You heal " + target.Name + " for " + strconv.Itoa(stamDam) + " stamina and " + strconv.Itoa(vitDam) + " vitality."

}

func firedamage(target *objects.Character, modifiers map[string]interface{}) string {
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

func earthdamage(target *objects.Character, modifiers map[string]interface{}) string {
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

func airdamage(target *objects.Character, modifiers map[string]interface{}) string {
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

func waterdamage(target *objects.Character, modifiers map[string]interface{}) string {
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

func light(target *objects.Character, modifiers map[string]interface{}) string {

	return "true"
}

func curepoison(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func bless(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func protection(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func invisibility(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func detectinvisible(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func teleport(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func stun(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func enchant(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func recall(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func summon(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func wizardwalk(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func levitate(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistfire(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistmagic(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func removecurse(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistair(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistwater(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistearth(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func clairvoyance(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func removedisease(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func cureblindness(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func polymorph(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func attraction(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func inertialbarrier(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func surge(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistpoison(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resilientaura(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistdisease(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func disruptmagic(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func reflection(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func dodge(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func resistacid(target *objects.Character, modifiers map[string]interface{}) string { return "" }

func embolden(target *objects.Character, modifiers map[string]interface{}) string {
	target.RemoveEffect("fear")
	return "true"
}

func mobheal(target *objects.Mob, modifiers map[string]interface{}) string {
	damage := 0
	if modifiers["magnitude"].(int) == 1 {
		damage = 50
	} else {
		damage = 100
	}
	target.Heal(damage)
	return "You heal " + target.Name + " for " + strconv.Itoa(damage) + " health"
}

func mobhealall(target *objects.Mob, modifiers map[string]interface{}) string {
	target.Heal(2000)
	return "You heal " + target.Name + " for 2000 health"
}

func mobfiredamage(target *objects.Mob, modifiers map[string]interface{}) string {
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
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " fire damage."
}

func mobearthdamage(target *objects.Mob, modifiers map[string]interface{}) string {
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
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " earth damage."
}

func mobairdamage(target *objects.Mob, modifiers map[string]interface{}) string {
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
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " air damage."
}

func mobwaterdamage(target *objects.Mob, modifiers map[string]interface{}) string {
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
	return "Your spell struck "+ target.Name +" for " + strconv.Itoa(damage) + " water damage."
}

func moblight(target *objects.Mob, modifiers map[string]interface{}) string {
	return ""
}

func mobcurepoison(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobbless(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobprotection(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobinvisibility(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobdetectinvisible(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobteleport(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobstun(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobenchant(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobrecall(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobsummon(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobwizardwalk(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func moblevitate(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistfire(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistmagic(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobremovecurse(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistair(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistwater(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistearth(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobclairvoyance(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobremovedisease(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobcureblindness(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobpolymorph(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobattraction(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobinertialbarrier(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobsurge(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistpoison(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresilientaura(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistdisease(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobdisruptmagic(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobreflection(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobdodge(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobresistacid(target *objects.Mob, modifiers map[string]interface{}) string { return "" }

func mobembolden(target *objects.Mob, modifiers map[string]interface{}) string {
	target.RemoveEffect("fear")
	return "true"
}
