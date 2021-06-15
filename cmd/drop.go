package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

// Syntax: DROP item
func init() {
	addHandler(drop{},
		"Usage:  drop itemName # \n \n Drop the specified item name and number.",
		permissions.Player,
		"drop")
}

type drop cmd

func (drop) process(s *state) {

	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to drop?")
		return
	}

	// We have at least 2 items here so lets move forward with that
	targetStr := s.words[0]
	targetNum := 1

	if len(s.words) < 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			s.msg.Actor.SendInfo("What did you want to drop?")
			return
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to drop?")
		return
	}

	s.actor.RunHook("act")
	where := s.where.Items

	s.actor.Inventory.Lock()
	s.actor.Inventory.Remove(target)
	where.Add(target)
	s.actor.Inventory.Unlock()
	if target.Flags["permanent"] {
		s.where.Save()
	}

	s.msg.Actor.SendGood("You drop ", target.Name, ".")
	if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		s.msg.Observers.SendInfo(s.actor.Name, " drops ", target.Name, ".")
	}

	s.ok = true
}
