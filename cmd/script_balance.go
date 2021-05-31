package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(balance{},
	"",
	permissions.Player,
	"$BALANCE")
}

type balance cmd

func (balance) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Setup Pawn Shop and let folks barter
}
