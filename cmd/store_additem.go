package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(additem{},
	"",
	permissions.Player,
	"additem")
}

type additem cmd

func (additem) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Setup Pawn Shop and let folks barter
}
