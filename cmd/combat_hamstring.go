package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(hamstring{},
		"Usage:  hamstring target # \n\n Try to hamstring a mob to generate a large amount of threat",
		permissions.Fighter,
		"hamstring", "ham")
}

type hamstring cmd

func (hamstring) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Hamstring what exactly?")
		return
	}

	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("hamstring")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {
		s.actor.Victim = whatMob
		s.actor.RunHook("combat")

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == (*objects.Item)(nil) {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Equipment.Main.ItemType != 4 && (s.actor.Placement != whatMob.Placement) {
			s.msg.Actor.SendBad("You are too far away to hamstring them.")
			return
		}

		// Check for a miss
		if utils.Roll(100, 1, 0) <= DetermineMissChance(s, whatMob.Level-s.actor.Tier) {
			s.msg.Actor.SendBad("You missed!!")
			s.msg.Observers.SendBad(s.actor.Name + " fails to hamstring " + whatMob.Name)
			whatMob.AddThreatDamage(1, s.actor)
			whatMob.CurrentTarget = s.actor.Name
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		}

		whatMob.AddThreatDamage(whatMob.Stam.Max, s.actor)
		whatMob.CurrentTarget = s.actor.Name
		s.actor.SetTimer("hamstring", config.HamTimer)
		s.msg.Actor.SendInfo("You hamstring " + whatMob.Name)
		s.msg.Observers.SendInfo(s.actor.Name + " hamstrings " + whatMob.Name)
		return

	}

	s.msg.Actor.SendInfo("Hamstring what?")
	s.ok = true
}
