package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
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

	s.msg.Actor.SendInfo("In your inventory, ", strconv.Itoa(len(inv)) , " items weighing approximately ", strconv.Itoa(s.actor.Inventory.TotalWeight) ,"lbs. (Equipment: " + strconv.Itoa(s.actor.Equipment.Weight) + "lbs)" )

	if len(inv) == 0 {
		s.msg.Actor.Send("  No items")
	}else {
		s.msg.Actor.Send("  ", strings.Join(s.actor.Inventory.List(), ", "))
	}
	s.ok = true
	return
}
