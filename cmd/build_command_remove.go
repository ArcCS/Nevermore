package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(removeCommand{},
		"Usage: remove_command room|mob|item (name) command_name \n  Deletes a command from the list of commands available \n",
		permissions.Builder,
		"remove_command")
}

type removeCommand cmd

func (removeCommand) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Not enough arguments to process the command.")
		return
	}

	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "room":
		if _, ok := s.where.Commands[s.words[1]]; ok {
			s.where.RemoveCommand(s.words[1])
			s.msg.Actor.SendGood("Script removed from room")
			s.where.Save()
			return
		} else {
			s.msg.Actor.SendBad("The command wasn't found in the rooms commands.")
			return
		}
	// Handle Items
	case "item":
		if len(s.words) < 3 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}
		itemName := s.words[1]
		item := s.actor.Inventory.Search(itemName, 1)

		if item != nil {
			if _, ok := item.Commands[s.words[2]]; ok {
				item.RemoveCommand(s.words[2])
				s.msg.Actor.SendGood("Script removed from item")
				item.Save()
				return
			} else {
				s.msg.Actor.SendBad("The command wasn't found in the items command list.")
			}
		} else {
			s.msg.Actor.SendBad("Item not found.")
			return
		}

	// Handle Mobs
	case "mob":
		if len(s.words) < 3 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}

		mobName := s.words[1]
		mob := s.where.Mobs.Search(mobName, 1, s.actor)

		if mob != nil {
			if _, ok := mob.Commands[s.words[2]]; ok {
				mob.RemoveCommand(s.words[2])
				s.msg.Actor.SendGood("Script removed from mob")
				mob.Save()
				return
			} else {
				s.msg.Actor.SendBad("The command wasn't found on this mob.")
			}
		} else {
			s.msg.Actor.SendBad("Mob not found.")
			return
		}

	default:
		s.msg.Actor.SendBad("Not an object that can be edited.")
	}

	s.ok = true
	return
}
