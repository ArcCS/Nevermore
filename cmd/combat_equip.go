package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(equip{},
           "Usage:  equip item # \n\n Try to equip an item from your inventory",
           permissions.Player,
           "equip")
}

type equip cmd

func (equip) process(s *state) {
	s.msg.Actor.SendInfo("Mighty fine air you want to put on")
	s.ok = true
}
