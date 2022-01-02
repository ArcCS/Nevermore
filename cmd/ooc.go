package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(ooc{},
		"Usage:  ooc # \n \n Switch to OOC lounge, use IC to swap out.",
		permissions.Player,
		"ooc")
}

type ooc cmd

func (ooc) process(s *state) {
	if s.actor.Flags["ooc"]{
		s.msg.Actor.SendBad("You are already OOC.")
		return
	}
	if to, ok := objects.Rooms[config.OocRoom]; ok {
		if !utils.IntIn(to.RoomId, s.cLocks) {
			s.AddCharLock(to.RoomId)
			return
		} else {
			s.actor.OOCSwap = s.actor.ParentId
			s.actor.Flags["ooc"] = true
			s.where.Chars.Remove(s.actor)
			to.Chars.Add(s.actor)
			s.actor.ParentId = to.RoomId
			s.scriptActor("LOOK")
			s.ok = true
			return
		}
	}

}
