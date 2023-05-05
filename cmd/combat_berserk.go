package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(berserk{},
		"Usage:  berserk \n\n Begin an uncontrollable rage with enhanced strength",
		permissions.Barbarian,
		"berserk", "rage", "berz")
}

type berserk cmd

func (berserk) process(s *state) {
	// Check some timers
	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	if s.actor.Tier < config.SpecialAbilityTier {
		s.msg.Actor.SendBad("You must be at least tier " + strconv.Itoa(config.SpecialAbilityTier) + " to use this skill.")
		return
	}

	berz, ok := s.actor.Flags["berserk"]
	if ok {
		if berz {
			s.msg.Actor.SendBad("You're already in the grips of the red rage!")
			return
		}
	}
	ready, msg := s.actor.TimerReady("combat_berserk")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.RunHook("combat")
	objects.Effects["berserk"](s.actor, s.actor, 0)
	s.msg.Observers.SendInfo(s.actor.Name + " goes berserk!")
	s.actor.SetTimer("combat_berserk", 60*10)

	s.ok = true
}
