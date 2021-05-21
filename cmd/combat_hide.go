package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
)

func init() {
	addHandler(hide{},
		"Usage:  hide <item> # \n\n Hide in the shadows, or attempt to hide an item",
		permissions.Player,
		"hide")
}

type hide cmd

func (hide) process(s *state) {
	if s.actor.Flags["hidden"] {
		s.msg.Actor.SendGood("You're already hidden")
	}

	s.actor.SetTimer("combat", config.CombatCooldown)

	// base chance is 15% to hide
	curChance := config.HideChance

	if s.actor.Class == 2 || s.actor.Class == 3 {
		curChance += 30
	}

	curChance += s.actor.Dex.Current * config.HideChancePerPoint

	if curChance >= 100 {
		s.msg.Actor.SendGood("You slip into the shadows.")
		s.actor.Flags["hidden"] = true
		//TODO:  Add a hook for combat and movement to fall out of the shadows
		s.ok = true
		return
	}

	if utils.Roll(100, 1, 0) <= curChance {
		s.msg.Actor.SendGood("You slip into the shadows.")
		s.actor.Flags["hidden"] = true
		//TODO:  Add a hook for combat and movement to fall out of the shadows
	}else{
		s.msg.Actor.SendBad("You attempt to hide in the shadows but can't find a place for yourself.")
	}
	s.ok = true
}
