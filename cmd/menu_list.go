package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(list{}, "", permissions.Player, "$LIST")
}

type list cmd

func (list) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: List the items in this story inventory
}
