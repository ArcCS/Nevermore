package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(deactivate{},
		"Usage:  deactivate (id) \n \n Activate a room so that it can be seen in the world. ",
		permissions.Dungeonmaster,
		"deactivate")
}

type deactivate cmd

func (deactivate) process(s *state) {
	if len(s.words) == 0 {
		s.where.Flags["active"] = false
		s.where.Save()
		s.msg.Actor.SendGood("Current room deactivated")
	} else {
		objectRef, _ := strconv.Atoi(s.input[1])
		room, rErr := objects.Rooms[objectRef]
		if rErr {
			room.Flags["active"] = false
			room.Save()
			s.msg.Actor.SendGood("Room deactivated")
		} else {
			s.msg.Actor.SendBad("Couldn't find room.")
		}
	}

	s.ok = true
	return
}
