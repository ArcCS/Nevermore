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
	if s.actor.Flags["ooc"] {
		s.msg.Actor.SendBad("You are already OOC.")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	if to, ok := objects.Rooms[config.OocRoom]; ok {
		if !utils.IntIn(to.RoomId, s.rLocks) {
			s.AddLocks(to.RoomId)
			s.ok = false
			return
		} else {
			for _, mob := range s.where.Mobs.Contents {
				if mob.CheckThreatTable(s.actor.Name) {
					s.msg.Actor.SendBad("You can't do that while in combat!")
					return
				}
			}

			s.actor.OOCSwap = s.actor.ParentId
			s.actor.Flags["ooc"] = true
			s.where.Chars.Remove(s.actor)
			if s.actor.Flags["invisible"] == false && s.actor.Flags["hidden"] == false {
				s.msg.Observers.SendGood("", s.actor.Name, " vanishes in a puff of smoke.")
			}
			to.Chars.Add(s.actor)
			s.actor.ParentId = to.RoomId
			s.scriptActor("LOOK")
			s.ok = true
			return
		}
	}

}
