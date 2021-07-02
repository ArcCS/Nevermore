package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
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

	if newRoomId, err = strconv.Atoi(s.words[1]); err != nil {
		s.msg.Actor.SendInfo("Room parameter couldn't be resolved.")
		return
	}

	if newRoom, ok = objects.Rooms[newRoomId]; !ok {
		s.msg.Actor.SendInfo("Room ID not valid.")
		return
	}

	if !utils.IntIn(newRoom.RoomId, s.cLocks) {
		s.AddCharLock(newRoom.RoomId)
		s.ok = false
		return
	}

	target := stats.ActiveCharacters.Find(s.words[0])
	if target != nil {
		if !utils.IntIn(newRoom.RoomId, s.cLocks) {
			s.AddCharLock(newRoom.RoomId)
			s.ok = false
			return
		}
		if s.actor != target {
			s.participant = target
		}
		if target.Resist && target != s.actor {
			// For every 5 points of int over the target there's an extra 10% chance to teleport
			diff := ((s.actor.GetStat("int") - target.GetStat("int")) / 5) * 10
			chance := 30 + diff
			if utils.Roll(100, 1, 0) > chance {
				s.msg.Actor.SendBad("You failed to magically transport " + target.Name)
				s.msg.Participant.SendBad(s.actor.Name + " failed to magically transport you")
				s.ok = true
				return
			}
		}
		s.msg.Actor.SendGood("You magically transported " + target.Name)
		s.msg.Participant.SendBad(s.actor.Name + " magically transported you.")
		s.where.Chars.Remove(target)
		newRoom.Chars.Add(target)
		target.ParentId = newRoom.RoomId
		go Script(target, "LOOK")
		s.ok = true
		return
	}else{
		s.msg.Actor.SendBad("Could not find dthat person to cast the spell on.")
	}

}
