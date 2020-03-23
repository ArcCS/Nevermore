package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(teleport{},
		 "teleport",
	)
	addHelp("Usage:  teleport id # \n \n teleports you to the specified room id", 50, "teleport")
}

type teleport cmd

func (teleport) process(s *state) {
	if s.actor.Class < 50 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Teleport where?")
		return
	}

	if s.words[0] == "" {
		s.msg.Actor.SendBad("Teleport where?")
		return
	}
	roomId, _ := strconv.Atoi(s.words[0])
	if to, ok := objects.Rooms[int64(roomId)]; ok {
		if !utils.IntIn(int(to.RoomId), s.cLocks){
			s.AddCharLock(int(to.RoomId))
			return
		}else{
			s.where.Chars.Remove(s.actor)
			to.Chars.Add(s.actor)
			s.actor.ParentId = to.RoomId

			s.scriptActor("LOOK")
			s.ok=true
			return
		}
	}

}