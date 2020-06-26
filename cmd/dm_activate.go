package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(activate{},
           "Usage:  activate (id) \n \n Activate a room so that it can be seen in the world. ",
           permissions.Dungeonmaster,
           "activate")
}

type activate cmd

func (activate) process(s *state) {
	if len(s.words) == 0 {
		s.where.Flags["active"] = true
		s.where.Save()
		s.msg.Actor.SendGood("Current room activated")
	}else {
		objectRef, _ := strconv.Atoi(s.input[1])
		room, rErr := objects.Rooms[objectRef]
		if rErr {
			room.Flags["active"] = true
			room.Save()
			s.msg.Actor.SendGood("Room activated")
		} else {
			s.msg.Actor.SendBad("Couldn't find room.")
		}
	}

	s.ok = true
	return
}