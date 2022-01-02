package cmd

import "github.com/ArcCS/Nevermore/permissions"

// Syntax: WHO
func init() {
	addHandler(save{},
		"Usage:  Commit your current character state to the db.",
		permissions.Player,
		"SAVE")
}

type save cmd

func (save) process(s *state) {
	s.msg.Actor.SendGood("Saving....")
	s.actor.Save()
	s.msg.Actor.SendGood("Saved.")
}
