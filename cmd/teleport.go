package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

func init() {
	addHandler(scriptTeleport{},
		"",
		permissions.Anyone,
		"$TELEPORT")
}

type scriptTeleport cmd

func (scriptTeleport) process(s *state) {
	roomId, _ := strconv.Atoi(s.words[0])
	if to, ok := objects.Rooms[roomId]; ok {
		if !utils.IntIn(to.RoomId, s.cLocks) {
			s.AddCharLock(to.RoomId)
			return
		} else {
			s.where.Chars.Remove(s.actor)
			to.Chars.Add(s.actor)
			s.actor.ParentId = to.RoomId
			s.scriptActor("LOOK")

			s.ok = true
			return
		}
	}

}
