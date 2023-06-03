package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(add_command{},
		"Usage: add_command room|mob|item (name) command_name command_string \n  Adds a command to the list of commands available \n"+
			"example:  add_command mob dragon talk $TELEPORTTO room_id "+
			"$TELEPORTTO is from the list of script commands"+
			"The word talk is what will be processed as a command when the user, and the room_id is where the teleport will send them. "+
			"Use of the room option, will add the command to the current room",
		permissions.Builder,
		"add_command")
}

type add_command cmd

func (add_command) process(s *state) {
	if len(s.words) < 3 {
		s.msg.Actor.SendBad("Not enough arguments to process the command.")
		return
	}

	//log.Println("Trying to edit: " + strings.ToLower(s.words[0]))
	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "room":
		if _, ok := ScriptList[s.words[2]]; ok {
			s.where.AddCommands(s.words[1], strings.Join(s.words[2:], " "))
			s.msg.Actor.SendGood("Script set on room")
			s.where.Save()
			return
		} else {
			s.msg.Actor.SendBad("The inputted script was not recognized: " + s.words[1])
			return
		}
	// Handle Items
	case "item":
		if len(s.words) < 4 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}
		itemName := s.words[1]
		item := s.actor.Inventory.Search(itemName, 1)

		if item != nil {
			if _, ok := ScriptList[s.words[3]]; ok {
				item.AddCommands(s.words[2], strings.Join(s.words[3:], " "))
				s.msg.Actor.SendGood("Script set on item")
				item.Save()
				return
			} else {
				s.msg.Actor.SendBad("The inputted script was not recognized.")
			}
		} else {
			s.msg.Actor.SendBad("Item not found.")
			return
		}

	// Handle Mobs
	case "mob":
		if len(s.words) < 4 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}

		mobName := s.input[2]
		mob := s.where.Mobs.Search(mobName, 1, s.actor)

		if mob != nil {
			if _, ok := ScriptList[s.words[3]]; ok {
				mob.AddCommands(s.words[2], strings.Join(s.words[3:], " "))
				s.msg.Actor.SendGood("Script set on mob")
				mob.Save()
				return
			} else {
				s.msg.Actor.SendBad("The inputted script was not recognized.")
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
