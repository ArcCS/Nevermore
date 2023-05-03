package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(encounter{},
		"Usage:  encounters off/on /n Suspend or restart rooms encounters",
		permissions.Dungeonmaster,
		"encounter")
}

type encounter cmd

func (encounter) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("on or off??")
		return
	}

	if !utils.StringIn(s.words[0], []string{"ON", "OFF"}) {
		s.msg.Actor.SendInfo("on or off!??!")
		return
	}

	if s.words[0] == "ON" {
		s.where.Flags["encounters_on"] = true
		s.msg.Actor.SendGood("Encounters restarted in this room.")
		return
	}

	if s.words[0] == "OFF" {
		s.where.Flags["encounters_on"] = false
		s.msg.Actor.SendGood("Encounters suspended in this room.")
		s.ok = true

		return
	}
	return
}
