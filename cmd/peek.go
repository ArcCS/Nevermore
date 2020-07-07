package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(peek{},
           "Usage:  inventory \n \n Display the current items in your inventory.",
           permissions.Thief,
           "peek")
}

type peek cmd

func (peek) process(s *state) {
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

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, false)
	if whatMob != nil {
		inv := whatMob.Inventory.List()
		s.msg.Actor.SendInfo("In their inventory:")
		if len(inv) == 0 {
			s.msg.Actor.Send("  No items")
		}else {
			s.msg.Actor.Send("  ", strings.Join(s.actor.Inventory.List(), ", "))
		}
	}
	return
}
