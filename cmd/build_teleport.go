package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(teleport{},
	"Usage:  teleport id # \n \n teleports you to the specified room id",
	permissions.Builder,
	"teleport")
}

type teleport cmd

func (teleport) process(s *state) {
	if len(s.words) < 1 {
		s.msg.Actor.SendBad("Teleport where?")
		return
	}

	if s.words[0] == "" {
		s.msg.Actor.SendBad("Teleport where?")
		return
	}
	roomId, _ := strconv.Atoi(s.words[0])
	if to, ok := objects.Rooms[roomId]; ok {
		if !utils.IntIn(to.RoomId, s.cLocks){
			s.AddCharLock(to.RoomId)
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