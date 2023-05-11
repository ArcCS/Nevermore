package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(inventory{},
		"Usage:  inventory \n \n Display the current items in your inventory.",
		permissions.Player,
		"inv", "inventory")
}

type inventory cmd

func (inventory) process(s *state) {

	// Try and find out if we are carrying anything
	inv := s.actor.Inventory.List()

	s.msg.Actor.SendInfo("You are carrying " + strconv.Itoa(s.actor.GetCurrentWeight()) + ", (Inventory: " + strconv.Itoa(len(inv)) + " items at " + strconv.Itoa(s.actor.Inventory.GetTotalWeight()) + "lbs; Equipment: " + strconv.Itoa(s.actor.Equipment.GetWeight()) + "lbs)")

	if len(inv) == 0 {
		s.msg.Actor.Send("  No items")
	} else {
		s.msg.Actor.Send("  ", s.actor.Inventory.ReducedList())
	}
	s.ok = true
	return
}
