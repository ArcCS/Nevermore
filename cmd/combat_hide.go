package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
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
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.SetTimer("global", config.CombatCooldown)

	// base chance is 15% to hide
	curChance := config.HideChance

	if s.actor.Class == 2 || s.actor.Class == 3 {
		curChance += 30
	}

	curChance += s.actor.Dex.Current * config.HideChancePerPoint

	if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
		s.msg.Actor.SendGood("You slip into the shadows.")
		s.actor.Flags["hidden"] = true
		s.actor.ApplyHook("act", "hide", -1, "10", -1,
			func(){
				s.actor.Flags["hidden"] = false
				s.actor.Write([]byte(text.Info + "You step out of the shadows." + text.Reset + "\n"))
				s.actor.RemoveHook("act", "hide")
				return
			},
			func(){
				s.actor.Flags["hidden"] = false
			},
		)
		s.actor.ApplyHook("say", "hide", -1, "10", -1,
			func(){
				s.actor.Flags["hidden"] = false
				s.actor.Write([]byte(text.Info + "You step out of the shadows." + text.Reset + "\n"))
				s.actor.RemoveHook("say", "hide")
				return
			},
			func(){
				s.actor.Flags["hidden"] = false
				return
			},
		)
		s.actor.ApplyHook("combat", "hide", -1, "10", -1,
			func(){
				s.actor.Flags["hidden"] = false
				s.actor.Write([]byte(text.Info + "You step out of the shadows." + text.Reset + "\n"))
				s.actor.RemoveHook("combat", "hide")
				return
			},
			func(){
				s.actor.Flags["hidden"] = false
				return
			},
			)
		s.actor.ApplyHook("move", "hide", -1, "10", -1,
			func(){
				s.actor.Flags["hidden"] = false
				s.actor.RemoveHook("move", "hide")
				return
			},
			func(){
				s.actor.Flags["hidden"] = false
				return
			},
		)
		s.actor.ApplyHook("gridmove", "hide", -1, "10", -1,
			func(){
				// base chance is 15% to hide
				curChance := config.HideChance

				if s.actor.Class == 2 || s.actor.Class == 3 {
					curChance += 30
				}

				curChance += s.actor.Dex.Current * config.HideChancePerPoint
				if utils.Roll(100, 1, 0) >= curChance {
					s.actor.Flags["hidden"] = false
					s.actor.Write([]byte(text.Bad + "You stumble out of the shadows while changing your position." + text.Reset + "\n"))
					s.actor.RemoveHook("gridmove", "hide")
					return
				}else {
					s.actor.Write([]byte(text.Good + "You stay in the shadows while moving." + text.Reset + "\n"))
					return
				}
			},
			func(){
				s.actor.Flags["hidden"] = false
				return
			},
		)

		s.ok = true
		return
	}else{
		s.msg.Actor.SendBad("Try as you might you fail to find a place to hide.")
	}

	s.ok = true
}
