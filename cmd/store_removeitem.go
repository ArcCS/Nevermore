package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(removeitem{},
	"",
	permissions.Player,
	"removeitem")
}

type removeitem cmd

func (removeitem) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Setup Pawn Shop and let folks barter
}
