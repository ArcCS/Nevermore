package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(questmode{},
		"Usage:  questmode off/on /n Change the realms to quest mode or turn it off.",
		permissions.Dungeonmaster,
		"questmode")
}

type questmode cmd

func (questmode) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("on or off??")
		return
	}

	if !utils.StringIn(s.words[0], []string{"ON", "OFF"}) {
		s.msg.Actor.SendInfo("on or off!??!")
		return
	}

	if s.words[0] == "ON" {
		config.QuestMode = true
		objects.ActiveCharacters.MessageAll("A barely perceptible mist fills the air and the realms grow disturbingly quiet. \n "+
			"(Quest Mode Has been activated: Death Loss capped at 10%, max loss to current level base/Tiered Exp gains are lifted)", config.BroadcastChannel)
		s.msg.Actor.SendGood("Realm Wide Quest mode has been activated")
		return
	}

	if s.words[0] == "OFF" {
		config.QuestMode = false
		objects.ActiveCharacters.MessageAll("The mist dissipates and the realms return to normal. (Quest Mode has been deactivated)", config.BroadcastChannel)
		s.msg.Actor.SendGood("Realm Wide Quest mode has been deactivated")
		s.ok = true

		return
	}
	return
}
