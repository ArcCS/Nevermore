package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(list_command{},
	"Usage: list_command room|mob|item (name)\n  Lists all of the commands for this object \n example:  list_command mob dragon",
	permissions.Builder,
	"list_commands", "list_command", "lc")
}

type list_command cmd

func (list_command) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Not enough arguments to process the command.")
		return
	}

	//log.Println("Trying to edit: " + strings.ToLower(s.words[0]))
	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "room":
		if len(s.where.Commands) > 0 {
			s.msg.Actor.SendInfo("Commands attached to this room \n ========== \n")
			for key, val := range s.where.Commands {
				s.msg.Actor.SendInfo(key, " | ", val.Command)
			}
			return
		}else{
			s.msg.Actor.SendInfo("There are no commands attached to this room")
		}
	// Handle Items
	case "item":
		if len(s.words) < 2 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}

		itemName := s.words[1]
		item := s.actor.Inventory.Search(itemName, 1)

		if item != nil {
			if len(item.Commands) > 0 {
				s.msg.Actor.SendInfo("Commands attached to this item \n ========== \n")
				for key, val := range item.Commands {
					s.msg.Actor.SendInfo(key, " | ", val.Command)
				}
				return
			}else{
				s.msg.Actor.SendInfo("There are no commands attached to this item")
			}
		} else {
			s.msg.Actor.SendBad("Item not found.")
			return
		}

	// Handle Mobs
	case "mob":
		if len(s.words) < 2 {
			s.msg.Actor.SendBad("There aren't enough arguments to complete this action.")
			return
		}

		mobName := s.input[2]
		mob := s.where.Mobs.Search(mobName, 1, s.actor)

		if mob != nil {
			if len(mob.Commands) > 0 {
				s.msg.Actor.SendInfo("Commands attached to this mob \n ========== \n")
				for key, val := range mob.Commands {
					s.msg.Actor.SendInfo(key, " | ", val.Command)
				}
				return
			}else{
				s.msg.Actor.SendInfo("There are no commands attached to this mob")
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
