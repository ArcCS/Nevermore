package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"github.com/jinzhu/copier"
	"math/rand"
	"strings"
)

func init() {
	addHandler(scriptTeleport{},
		"",
		permissions.Anyone,
		"$TELEPORT")
}

type scriptTeleport cmd

func (scriptTeleport) process(s *state) {

	var modTeleportTable []int
	copier.Copy(&modTeleportTable, &objects.TeleportTable)
	if utils.IntIn(s.actor.ParentId, modTeleportTable) {
		modTeleportTable = utils.RemoveInt(modTeleportTable, s.actor.ParentId)
	}
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
	s.msg.Observers[newRoom.RoomId].SendInfo(s.actor.Name, " arrives in a puff of smoke.")
	s.scriptActor("LOOK")
	s.ok = true
	return

}
