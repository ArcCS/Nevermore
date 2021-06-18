package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"math/rand"
	"strconv"
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

	if !utils.IntIn(newRoom.RoomId, s.cLocks) {
		s.AddCharLock(newRoom.RoomId)
		s.ok = false
		return
	}

	if len(s.words) > 0 {
		target := s.where.Chars.Search(s.input[0], s.actor)
		if target != nil {
			if target.Resist {
				// For every 5 points of int over the target there's an extra 10% chance to teleport
				diff := ((s.actor.Int.Current - target.Int.Current) / 5) * 10
				chance := 30 + diff
				if utils.Roll(100, 1, 0) > chance {
					s.msg.Actor.SendBad("You failed to teleport " + target.Name)
					s.msg.Participant.SendBad(s.actor.Name + " failed to teleport you")
					s.ok = true
					return
				}
			}
			s.msg.Actor.SendGood("You teleport " + target.Name)
			s.msg.Participant.SendBad(s.actor.Name + " teleported you")
			s.where.Chars.Remove(target)
			newRoom.Chars.Add(target)
			target.ParentId = newRoom.RoomId
			Script(target, "LOOK")
			s.ok = true
			return
		}else{
			targetNum := 1
			if len(s.words) > 1 {
				// Try to snag a number off the list
				if val, err := strconv.Atoi(s.words[1]); err == nil {
					targetNum = val
				}
			}
			targetMob := s.where.Mobs.Search(s.input[0], targetNum, s.actor)
			if targetMob != nil {

				diff := (s.actor.Tier - targetMob.Level) * 5
				chance := 10 + diff
				if utils.Roll(100, 1, 0) > chance {
					s.msg.Actor.SendBad("You failed to teleport " + targetMob.Name)
					s.ok = true
					return
				}
				s.msg.Actor.SendGood("You teleport " + targetMob.Name)
				s.where.Mobs.Remove(targetMob)
				for _, char := range s.where.Chars.Contents {
					if char.Victim == target {
						char.Victim = nil
					}
				}
				if len(newRoom.Chars.Contents) > 0 {
					newRoom.Mobs.Add(targetMob, true)
					newRoom.MessageAll(targetMob.Name + " arrives in a puff of smoke.")
					targetMob.StartTicking()
				}
				s.ok = true
				return
			}
			s.msg.Actor.SendBad("Could not target that for teleportation.")
			s.ok = true
			return
		}
	}

	s.where.Chars.Remove(s.actor)
	newRoom.Chars.Add(s.actor)
	s.actor.ParentId = newRoom.RoomId
	s.scriptActor("LOOK")
	s.ok = true
	return

}
