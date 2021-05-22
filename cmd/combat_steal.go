package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(steal{},
		"Usage:  steal target item \n \n Try to steal an item from a targets inventory",
		permissions.Thief,
		"steal")
}

type steal cmd

func (steal) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Attack what exactly?")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("peek")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	//TODO: Steal from players inventory if PvP flag is set

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, false)
	if whatMob != nil {
		s.actor.RunHook("steal")
		inv := whatMob.Inventory.List()
		s.msg.Actor.SendInfo("In their inventory:")
		if len(inv) == 0 {
			s.msg.Actor.Send("  No items")
		} else {
			s.msg.Actor.Send("  ", strings.Join(whatMob.Inventory.List(), ", "))
		}
	}
	return
}
