package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(haste{},
		"Usage:  haste \n\n Hasten your actions temporarily increasing your dex and your combat actions",
		permissions.Ranger,
		"haste")
}

type haste cmd

func (haste) process(s *state) {
	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	haste, ok := s.actor.Flags["haste"]
	if ok {
		if haste {
			s.msg.Actor.SendBad("You're already moving quickly!")
			return
		}
	}
	ready, msg := s.actor.TimerReady("combat_haste")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	objects.Effects["haste"](s.actor, s.actor, 0)
	s.msg.Observers.SendInfo(s.actor.Name + " begins moving faster!")
	s.actor.SetTimer("combat_haste", 60*10)

	s.ok = true
}
