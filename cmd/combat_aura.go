package cmd

import (
	"strconv"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(aura{},
		"Usage:  seal [courage|faith|justice|off] \n\n Invoke a holy seal upon yourself; only one seal can be held at a time.\n"+
			" Typing seal on its own tells you which seal you currently bear.\n\n"+
			" courage: you stand firm behind your shield, turning aside blows meant for you and your allies and drawing your foes' wrath.\n"+
			" faith:   your devotion is answered in full; your healing grows stronger and your body and spirit mend more swiftly.\n"+
			" justice: your weapon is guided by righteous fury, striking truer and harder against the wicked.",
		permissions.Paladin,
		"seal")
}

type aura cmd

// sealDef is what a seal demands before it can be invoked.
type sealDef struct {
	tier   int
	shield bool
}

var seals = map[string]sealDef{
	"courage": {tier: config.MinorAbilityTier, shield: true},
	"faith":   {tier: config.MinorAbilityTier},
	"justice": {tier: config.SealJusticeTier},
}

// activeSeal returns the name of the seal the character currently holds, or "".
func activeSeal(c *objects.Character) string {
	for name := range seals {
		if c.CheckFlag("seal-" + name) {
			return name
		}
	}
	return ""
}

func (aura) process(s *state) {
	if !s.requireTier(config.MinorAbilityTier) {
		return
	}
	current := activeSeal(s.actor)
	if len(s.input) < 1 {
		if current == "" {
			s.msg.Actor.SendInfo("You bear no seal.")
		} else {
			s.msg.Actor.SendInfo("You bear the seal of " + current + ".")
		}
		return
	}

	which := s.input[0]

	if which == "off" || which == "none" {
		if current == "" {
			s.msg.Actor.SendBad("You have no seal to release.")
			return
		}
		s.actor.RemoveEffect("seal-" + current)
		s.msg.Observers.SendInfo(s.actor.Name + " releases " + config.PosPronoun(s.actor.Gender) + " seal.")
		s.ok = true
		return
	}

	def, ok := seals[which]
	if !ok {
		s.msg.Actor.SendBad("I don't know that seal.")
		return
	}
	if s.actor.Tier < def.tier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(def.tier) + " to invoke the seal of " + which + ".")
		return
	}
	if current == which {
		s.msg.Actor.SendBad("You already bear the seal of " + which + ".")
		return
	}
	if def.shield && !s.actor.HasShield() {
		s.msg.Actor.SendBad("The seal of " + which + " requires a shield in your offhand.")
		return
	}

	ready, msg := s.actor.TimerReady("combat_seal")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	if current != "" {
		s.actor.RemoveEffect("seal-" + current)
	}
	objects.Effects["seal-"+which](s.actor, s.actor, 0)
	s.msg.Observers.SendInfo(s.actor.Name + " invokes the seal of " + which + ".")
	s.actor.SetTimer("combat_seal", config.SealTimer)
	s.actor.SetTimer("combat", 3)
	s.ok = true
}
