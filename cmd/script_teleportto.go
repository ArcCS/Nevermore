package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

func init() {
	addHandler(scriptTeleportTo{},
		"",
		permissions.Anyone,
		"$TELEPORTTO")
}

type scriptTeleportTo cmd

func (scriptTeleportTo) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Not enough parameters for teleport_to to execute.")
		return
	}

	var err error
	var ok bool
	var newRoomId int
	var newRoom *objects.Room

	if newRoomId, err = strconv.Atoi(s.words[0]); err != nil {
		s.msg.Actor.SendInfo("Room parameter couldn't be resolved.")
		return
	}

	if newRoom, ok = objects.Rooms[newRoomId]; !ok {
		s.msg.Actor.SendInfo("Room ID not valid.")
		return
	}

	if !utils.IntIn(newRoom.RoomId, s.rLocks) {
		s.AddLocks(newRoom.RoomId)
		s.ok = false
		return
	}

	if !utils.IntIn(newRoom.RoomId, s.rLocks) {
		s.AddLocks(newRoom.RoomId)
		s.ok = false
		return
	}

	s.msg.Actor.SendInfo(utils.Title(strings.ToLower(strings.Join(s.words[1:], " "))))
	s.where.Chars.Remove(s.actor)
	newRoom.Chars.Add(s.actor)
	s.actor.ParentId = newRoom.RoomId
	s.scriptActor("LOOK")
	s.ok = true
	return

}
