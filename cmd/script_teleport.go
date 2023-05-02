package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"math/rand"
	"strings"
	"time"
)

func init() {
	addHandler(scriptTeleport{},
		"",
		permissions.Anyone,
		"$TELEPORT")
}

type scriptTeleport cmd

func (scriptTeleport) process(s *state) {

	rand.Seed(time.Now().Unix())
	newRoom := objects.Rooms[objects.TeleportTable[rand.Intn(len(objects.TeleportTable))]]

	if !utils.IntIn(newRoom.RoomId, s.rLocks) {
		s.AddLocks(newRoom.RoomId)
		s.ok = false
		return
	}

	if len(s.words) != 0 {
		s.msg.Actor.Send(strings.Join(s.input[1:], " "))
	} else {
		s.msg.Actor.Send("You were teleported!!")
	}
	s.where.Chars.Remove(s.actor)
	newRoom.Chars.Add(s.actor)
	s.actor.ParentId = newRoom.RoomId
	s.msg.Observers.Send(s.actor.Name + " was teleported away!")
	s.scriptActor("LOOK")
	s.ok = true
	return

}
