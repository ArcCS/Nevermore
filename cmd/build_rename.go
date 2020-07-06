package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(rename{},"Usage:  rename name new_name \n \n Use this command to change the instance item name. (Does not overwrite DB) \n",
		permissions.Builder,
		"rename")
}

type rename cmd

func (rename) process(s *state) {
	// Check arguments
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Rename what to what?")
		return
	}

	// Toggle Flags
	itemName := s.input[0]
	item := s.actor.Inventory.Search(itemName, 1)

	if item != nil {
		item.Name = strings.Join(s.input[1:], "")
		s.msg.Actor.SendGood("Item name changed to " + item.Name)
	}else{
		s.msg.Actor.SendBad("Not an object that can be edited.")
	}

	s.ok = true
	return
}
