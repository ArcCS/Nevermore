package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(withdraw{},
	"",
	permissions.Player,
	"$WITHDRAW")
}

type withdraw cmd

func (withdraw) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Setup Pawn Shop and let folks barter
}
