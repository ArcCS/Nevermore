package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(buy{}, "", permissions.Player, "$BUY")
}

type buy cmd

func (buy) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")

}
