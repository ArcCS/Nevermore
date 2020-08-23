package spells

import "strconv"

type mobnpcspells interface {
	ChangePlacement(place int) bool
	ApplyEffect(effectname string, length string, interval string,  effect func(), effectoff func())
	RemoveEffect(effectname string)
	ReceiveDamage(damage int) (int, int)
	ReceiveVitalDamage(damage int) int
	Heal(damage int) (int, int)
	HealVital(damage int)
	HealStam(damage int)
	RestoreMana(damage int)
}

var Effects = map[string]func(caller mobnpcspells, target mobnpcspells, modifier int) string {
	"heal-stam": healstam,
	"heal-vit": healvit,
	"heal": heal,
	"heal-all": healall,
	"fire-damage": firedamage,
	"earth-damage": earthdamage,
	"air-damage": airdamage,
	"water-damage": waterdamage,
	"light": light,
	"curepoison": curepoison,
	"bless": bless,
	"protection": protection,
	"invisibility": invisibility,
	"detect-invisible": detectinvisible,
	"teleport": teleport,
	"stun": stun,
	"enchant": enchant,
	"recall": recall,
	"summon": summon,
	"wizard-walk": wizardwalk,
	"levitate": levitate,
	"resist-fire": resistfire,
	"resist-magic": resistmagic,
	"remove-curse": removecurse,
	"resist-air": resistair,
	"resist-water": resistwater,
	"resist-earth": resistearth,
	"clairvoyance": clairvoyance,
	"remove-disease": removedisease,
	"cure-blindness": cureblindness,
	"polymorph": polymorph,
	"attraction": attraction,
	"inertial-barrier": inertialbarrier,
	"surge": surge,
	"resist-poison": resistpoison,
	"resilient-aura": resilientaura,
	"resist-disease": resistdisease,
	"disrupt-magic": disruptmagic,
	"reflection": reflection,
	"dodge": dodge,
	"resist-acid": resistacid,
	"embolden": embolden,
}

func healstam(caller mobnpcspells, target mobnpcspells, modifier int) string{
	/*vigor devices seemed to be very low 4-10?
	mend-wounds devices were around 5-15.
	detraum devices seemed around 38-50.
	renew devices seemed around 60-80 (less sure on this)
	casting vigor with a l14 int 25 priest was around 45-49, whilst casting vigor with a l14 int 30 wiz was around 20ish
	casting detraum with l14 priest was  80-100ish
	casting renewal with l14 priest was approx 150
	casting detraum with bard was around 45?
	casting renewal with bard was approx 60-80? */
	target.HealStam(10)
	return "You restored 10 stamina"
}

func healvit(caller mobnpcspells, target mobnpcspells, modifier int) string{
	target.HealVital(10)
	return "You restored 10 vitality"
}

func heal(caller mobnpcspells, target mobnpcspells, modifier int) string{
	if modifier == 1 {
		target.Heal(50)
		return "you healed for 50 damage"
	}else{
		target.Heal(100)
		return "You healed for 100 damage"
	}

}

func healall(caller mobnpcspells, target mobnpcspells, modifier int) string{
	target.Heal(50000)
	return "You healed all of their damage"
}

func firedamage(caller mobnpcspells, target mobnpcspells, modifier int) string{
	damage := 0
	if modifier == 1 {
		damage = 10
	}else if modifier == 2 {
		damage = 30
	}else if modifier == 3 {
		damage = 60
	}else if modifier == 4 {
		damage = 120
	}else if modifier == 5 {
		damage = 250
	}else if modifier == 6 {
		damage = 400
	}else if modifier == 7 {
		damage = 600
	}
	//todo: get magical strength to do a x2 or x4 for spells or other modifiers as necessary
	target.ReceiveDamage(damage)
	return "Your spell struck for " + strconv.Itoa(damage)
}

func earthdamage(caller mobnpcspells, target mobnpcspells, modifier int) string{ 	damage := 0
	if modifier == 1 {
		damage = 10
	}else if modifier == 2 {
		damage = 30
	}else if modifier == 3 {
		damage = 60
	}else if modifier == 4 {
		damage = 120
	}else if modifier == 5 {
		damage = 250
	}else if modifier == 6 {
		damage = 400
	}else if modifier == 7 {
		damage = 600
	}
	//todo: get magical strength to do a x2 or x4 for spells or other modifiers as necessary
	target.ReceiveDamage(damage)
	return "Your spell struck for " + strconv.Itoa(damage) }

func airdamage(caller mobnpcspells, target mobnpcspells, modifier int) string{ 	damage := 0
	if modifier == 1 {
		damage = 10
	}else if modifier == 2 {
		damage = 30
	}else if modifier == 3 {
		damage = 60
	}else if modifier == 4 {
		damage = 120
	}else if modifier == 5 {
		damage = 250
	}else if modifier == 6 {
		damage = 400
	}else if modifier == 7 {
		damage = 600
	}
	//todo: get magical strength to do a x2 or x4 for spells or other modifiers as necessary
	target.ReceiveDamage(damage)
	return "Your spell struck for " + strconv.Itoa(damage) }

func waterdamage(caller mobnpcspells, target mobnpcspells, modifier int) string{ 	damage := 0
	if modifier == 1 {
		damage = 10
	}else if modifier == 2 {
		damage = 30
	}else if modifier == 3 {
		damage = 60
	}else if modifier == 4 {
		damage = 120
	}else if modifier == 5 {
		damage = 250
	}else if modifier == 6 {
		damage = 400
	}else if modifier == 7 {
		damage = 600
	}
	//todo: get magical strength to do a x2 or x4 for spells or other modifiers as necessary
	target.ReceiveDamage(damage)
	return "Your spell struck for " + strconv.Itoa(damage) }

func light(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func curepoison(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func bless(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func protection(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func invisibility(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func detectinvisible(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func teleport(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func stun(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func enchant(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func recall(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func summon(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func wizardwalk(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func levitate(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistfire(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistmagic(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func removecurse(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistair(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistwater(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistearth(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func clairvoyance(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func removedisease(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func cureblindness(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func polymorph(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func attraction(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func inertialbarrier(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func surge(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistpoison(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resilientaura(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistdisease(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func disruptmagic(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func reflection(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func dodge(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func resistacid(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }

func embolden(caller mobnpcspells, target mobnpcspells, modifier int) string{ return "" }
