package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strings"
)

func init() {
	addHandler(sneak{},
		"Usage:  sneak (ExitName) # \n\n Attempt to sneak into another area",
		permissions.Thief|permissions.Ranger|permissions.Monk,
		"sneak")
}

type sneak cmd

func (sneak) process(s *state) {
	var exitName string
	from := s.where
	// Does this place even have exits?
	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You must be hiding to sneak.")
		return
	}

	if len(from.Exits) == 0 {
		s.msg.Actor.SendInfo("You can't see anywhere to sneak from here.")
		return
	}

	// Decide what exit we are going to
	if utils.StringIn(s.cmd, directionals) {
		exitName = directionIndex[s.cmd]
	} else {
		if len(s.words) > 0 {
			// Join the strings together for exits with spaces
			exitName = strings.Join(s.words, " ")
		} else {
			s.msg.Actor.SendBad("Sneak where?")
		}
	}

	// Test for partial exit names
	exitTxt := strings.ToLower(exitName)
	if !utils.StringIn(strings.ToUpper(exitTxt), directionals) {
		for txtE := range from.Exits {
			if strings.Contains(txtE, exitTxt) {
				exitTxt = txtE
			}
		}
	}
	if toE, ok := from.Exits[exitTxt]; ok {
		s.actor.RunHook("sneak")
		s.actor.SetTimer("global", config.CombatCooldown)

		// base chance is 15% to hide
		curChance := config.SneakChance

		if s.actor.Class == 2 || s.actor.Class == 3 {
			curChance += 30
		}
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}
		curChance += s.actor.Dex.Current * config.SneakChancePerPoint

		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
			// Check that the room ID exists
			if to, ok := objects.Rooms[toE.ToId]; ok {
				// Apply a lock
				if !utils.IntIn(toE.ToId, s.cLocks) {
					s.AddCharLock(toE.ToId)
					return
				} else {
					if !toE.Flags["placement_dependent"] {
						if !objects.Rooms[toE.ToId].Crowded() {
							from.Chars.Remove(s.actor)
							to.Chars.Add(s.actor)
							s.actor.Placement = 3
							s.actor.ParentId = toE.ToId
							s.scriptActor("LOOK")
							s.ok = true
							return
						} else {
							s.msg.Actor.SendInfo("That area is crowded.")
							s.ok = true
							return
						}
					}
				}
			}
		}else{
			s.msg.Actor.SendBad("You stumble out of the shadows while leaving..")
			s.actor.Flags["hidden"] = false
			s.scriptActor("GO", exitName)
		}
	} else {
		s.msg.Actor.SendInfo("You can't go that direction.")
		s.ok = true
		return
	}

}
