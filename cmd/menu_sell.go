package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(sell{}, "", permissions.Player, "$SELL")
}

type sell cmd

func (sell) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")

}
