package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
)

func init() {
	addHandler(moveto{},
		"Usage:  moveto character_name room_id \n \n Move a character elsewhere",
		permissions.Builder,
		"moveto", "tpto")
}

type moveto cmd

func (moveto) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Move who where?")
		return
	}
	whoStr := s.words[0]
	who := objects.ActiveCharacters.Find(whoStr)
	if who != nil {
		roomId, _ := strconv.Atoi(s.words[1])
		if to, ok := objects.Rooms[roomId]; ok {
			if !utils.IntIn(to.RoomId, s.cLocks) {
				s.AddCharLock(to.RoomId)
				s.ok = false
				return
			}else if !utils.IntIn(who.ParentId, s.cLocks) {
				s.AddCharLock(who.ParentId)
				s.ok = false
				return
			} else {
				log.Println("Trying to teleport...")
				objects.Rooms[who.ParentId].Chars.Remove(who)
				to.Chars.Add(who)
				who.ParentId = to.RoomId
				go Script(who, "LOOK")
				s.msg.Actor.SendInfo("You teleported " + who.Name + " to " + to.Name + "(" + strconv.Itoa(to.RoomId) + ")")
				s.ok = true
				return
			}
		}else{
			s.msg.Actor.SendBad("Send who where?")
			return
		}
	} else {
		s.msg.Actor.SendBad("Send who where?")
		return
	}
}
