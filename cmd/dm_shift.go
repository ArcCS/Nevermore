package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
)

func init() {
	addHandler(shift{},
		"Usage:  shft room_id \n \n Move all characters in the room to a new room",
		permissions.Dungeonmaster,
		"shft")
}

type shift cmd

func (shift) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Shift everyone where?")
		return
	}

	var to *objects.Room
	var ok bool
	var roomId int
	var err error
	if roomId, err = strconv.Atoi(s.words[0]); err == nil {
		if to, ok = objects.Rooms[roomId]; ok {
			if !utils.IntIn(to.RoomId, s.rLocks) {
				s.AddLocks(to.RoomId)
				s.ok = false
				return
			} else if !utils.IntIn(s.actor.ParentId, s.rLocks) {
				s.AddLocks(s.actor.ParentId)
				s.ok = false
				return
			}
		} else {
			s.msg.Actor.SendBad("Send who where?")
			return
		}

		charList := s.where.Chars.ListAll()

		for _, char := range charList {
			log.Println("Shifting: ", char.Name)
			s.where.Chars.Remove(char)
			to.Chars.Add(char)
			char.ParentId = to.RoomId
			go Script(char, "LOOK")
			s.msg.Actor.SendInfo("You teleported " + char.Name + " to " + to.Name + "(" + strconv.Itoa(to.RoomId) + ")")
			s.ok = true
		}
	}
}
