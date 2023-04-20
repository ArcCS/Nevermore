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

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
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
		// Check that the room ID exists
		if to, ok := objects.Rooms[toE.ToId]; ok {
			// Apply a lock
			if !utils.IntIn(toE.ToId, s.cLocks) {
				s.AddCharLock(toE.ToId)
				s.ok = false
				return
			} else {
				s.actor.RunHook("sneak")
				// Reactivate this if the thieves are getting too spicy
				//s.actor.SetTimer("global", config.CombatCooldown)

				// base chance is 15% to hide
				curChance := config.SneakChance

				if s.actor.Class == 2 || s.actor.Class == 3 {
					curChance += config.SneakBonus
				}
				if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					curChance = 100
				}
				// Calculate bonus based on dex
				curChance += s.actor.Dex.Current * config.SneakChancePerPoint
				// Calculate bonus based on tier
				curChance += s.actor.Tier * config.SneakChancePerTier

				if utils.Roll(100, 1, 0) > curChance {
					s.msg.Actor.SendBad("You stumble out of the shadows while leaving..")
					s.actor.Flags["hidden"] = false
					s.ok = true
					s.scriptActor("GO", exitName)
					return
				}

				if !objects.Rooms[toE.ToId].Flags["active"] {
					s.msg.Actor.SendBad("Go where?")
					return
				}

				if toE.Flags["invisible"] && !s.actor.CheckFlag("detect_invisible") {
					s.msg.Actor.SendBad("Go where?")
					return
				}

				if toE.Flags["placement_dependent"] && s.actor.Placement != toE.Placement {
					s.msg.Actor.SendBad("You must be next to the exit to use it.")
					return
				}

				if toE.Flags["closed"] {
					s.msg.Actor.SendBad("The way is closed.")
					return
				}

				if toE.Flags["day_only"] && !objects.DayTime {
					s.msg.Actor.SendBad("You can only go there at night.")
					return
				}

				if toE.Flags["night_only"] && objects.DayTime {
					s.msg.Actor.SendBad("You can only go there during the day.")
					return
				}

				if s.actor.Equipment.Weight > s.actor.MaxWeight() {
					s.msg.Actor.SendBad("You are carrying too much to move.")
					return
				}

				if toE.Flags["levitate"] && !s.actor.CheckFlag("levitate") {
					s.msg.Actor.Send("You fall while trying to go that way!  You take 20 points of damage!")
					s.actor.ReceiveDamage(20)
					return
				}

				if objects.Rooms[toE.ToId].Crowded() {
					s.msg.Actor.SendInfo("That area is crowded.")
					s.ok = true
					return
				}

				from.Chars.Remove(s.actor)
				to.Chars.Add(s.actor)
				s.actor.Placement = 3
				s.actor.ParentId = toE.ToId
				s.scriptActor("LOOK")
				s.ok = true
				return
			}
		} else {
			s.msg.Actor.SendInfo("You can't go that direction.")
			s.ok = true
			return
		}
	} else {
		s.msg.Actor.SendInfo("You can't go that direction.")
		s.ok = true
		return
	}

}
