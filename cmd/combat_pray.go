package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(pray{},
		"Usage:  pray \n\n Focus your attention on your faith and be overwhelmed with piousness.",
		permissions.Cleric|permissions.Paladin,
		"pray")
}

type pray cmd

func (pray) process(s *state) {
	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	haste, ok := s.actor.Flags["pray"]
	if ok {
		if haste {
			s.msg.Actor.SendBad("You've recently prayed.'")
			return
		}
	}
	ready, msg := s.actor.TimerReady("combat_pray")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	objects.CharEffects["pray"](s.actor, map[string]interface{}{})
	s.msg.Observers.SendInfo(s.actor.Name + " prays.")
	s.actor.SetTimer("combat_pray", 60*10)

	s.ok = true
}

