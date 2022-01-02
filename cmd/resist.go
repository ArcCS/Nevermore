package cmd

import "github.com/ArcCS/Nevermore/permissions"

// Syntax: WHO
func init() {
	addHandler(resist{},
		"Usage:  resist,  toggles whether you are resisting teleport and summons.",
		permissions.Player,
		"resist")
}

type resist cmd

func (resist) process(s *state) {
	s.actor.Resist = !s.actor.Resist
	if s.actor.Resist {
		s.msg.Actor.SendInfo("You are now resisting teleport and summon attempts.")
	}else{
		s.msg.Actor.SendInfo("You are no longer resisting teleport or summon attempts.")
	}
}
