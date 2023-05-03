package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(ic{},
		"Usage:  ic # \n \n Switch back to in-character, use OOC to swap out.",
		permissions.Player,
		"ic")
}

type ic cmd

func (ic) process(s *state) {
	if !s.actor.Flags["ooc"] {
		s.msg.Actor.SendBad("You are already in-character.")
		return
	}
	if to, ok := objects.Rooms[s.actor.OOCSwap]; ok {
		if !utils.IntIn(to.RoomId, s.rLocks) {
			s.AddLocks(to.RoomId)
			s.ok = false
			return
		} else {
			s.actor.OOCSwap = 0
			s.actor.Flags["ooc"] = false
			s.actor.Flags["afk"] = false
			s.where.Chars.Remove(s.actor)
			to.Chars.Add(s.actor)
			s.actor.ParentId = to.RoomId
			s.scriptActor("LOOK")
			s.ok = true
			return
		}
	}

}
