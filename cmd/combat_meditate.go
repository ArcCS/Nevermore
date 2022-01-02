package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(meditate{},
		"Usage:  meditate \n\n Enter a meditative trance to recover your health and chi",
		permissions.Monk,
		"meditate")
}

type meditate cmd

func (meditate) process(s *state) {
	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	ready, msg := s.actor.TimerReady("combat_meditate")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.Stam.Current = s.actor.Stam.Max
	s.actor.Vit.Current = s.actor.Vit.Max
	s.actor.Mana.Current = s.actor.Mana.Max
	s.msg.Actor.SendGood("You slow your thoughts and enter a brief trance restoring your health and chi.")
	s.msg.Observers.SendInfo(s.actor.Name + " meditates!")
	s.actor.SetTimer("combat_meditate", config.MeditateTime)

	s.ok = true
}
