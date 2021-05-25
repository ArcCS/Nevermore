package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
)

func init() {
	addHandler(meditate{},
		"Usage:  meditate \n\n Enter a meditative trance to recover your health and chi",
		permissions.Monk,
		"meditate")
}

type meditate cmd

func (meditate) process(s *state) {
	//TODO: Finish Meditate Command
	// Check some timers
	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	berz, ok := s.actor.Flags["berserk"]
	if ok {
		if berz {
			s.msg.Actor.SendBad("You're already in the grips of the red rage!")
			return
		}
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

	s.actor.ApplyEffect("berserk", "60", "0",
		func() {
			s.actor.ToggleFlagAndMsg("meditate", "meditate", text.Red+"The red rage grips you!!!\n")
			_, ok := s.actor.Modifiers["base_damage"]
			if ok {
				s.actor.Modifiers["base_damage"] += s.actor.Str.Current * config.CombatModifiers["berserk"]
			} else {
				s.actor.Modifiers["base_damage"] = s.actor.Str.Current * config.CombatModifiers["berserk"]
			}
			s.actor.Str.Current += 5
		},
		func() {
			s.actor.ToggleFlagAndMsg("meditate", "meditate", text.Cyan+"The tension releases and your rage fades...\n")
			s.actor.Str.Current -= 5
			s.actor.Modifiers["base_damage"] -= s.actor.Str.Current * config.CombatModifiers["berserk"]
		})
	s.msg.Observers.SendInfo(s.actor.Name + " goes berserk!")
	s.actor.SetTimer("combat_berserk", 60*10)

	s.ok = true
}
