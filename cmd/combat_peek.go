package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(peek{},
		"Usage:  peek \n \n Display the current items in your inventory.",
		permissions.Thief,
		"peek")
}

type peek cmd

func (peek) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Peek whose inventory?")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
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

	//TODO: Peek players inventory if PvP flag is set

	//TODO: There should be a chance to fail
	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {
		s.actor.SetTimer("peek", config.PeekCD)
		inv := whatMob.Inventory.List()
		s.msg.Actor.SendInfo("In their inventory:")
		if len(inv) == 0 {
			s.msg.Actor.Send("  No items")
		} else {
			s.msg.Actor.Send("  ", strings.Join(whatMob.Inventory.List(), ", "))
		}
	} else {
		s.msg.Actor.SendBad("Peek whose inventory?")
	}
	s.ok = true
	return

}
