package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
	"strings"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(peek{},
		"Usage:  peek \n \n Display the current items in your inventory.",
		permissions.Thief,
		"peek")
}

type peek cmd

func (peek) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Peek whose inventory?")
		return
	}

	if s.actor.CheckFlag("blind") {
		s.msg.Actor.SendBad("You can't see anything!")
		return
	}

	if s.actor.Stam.Current <= 0 {
		s.msg.Actor.SendBad("You are far too tired to do that.")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("peek")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {
		curChance := config.StealChance + (s.actor.Dex.Current * config.StealChancePerPoint) + (config.StealthLevel(s.actor.Skills[11].Value) * config.StealChancePerSkillLevel)
		lvlDiff := float64(whatMob.Level - s.actor.Tier)
		if lvlDiff > 2 {
			lvlDiff = (lvlDiff - 2) * 0.2
			curChance -= int(float64(curChance) * lvlDiff)
		}

		curChance = curChance + 10

		//s.msg.Actor.SendInfo("Steal chance = " + strconv.Itoa(curChance))

		if curChance > 95 {
			curChance = 95
		}

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		//log.Println(s.actor.Name+"Peek Chance Roll: ", curChance)
		if utils.Roll(100, 1, 0) > curChance {
			s.msg.Actor.SendBad("You fail to peek into their inventory.")
			s.msg.Observers.SendInfo(s.actor.Name + " tries to peek into " + whatMob.Name + "'s inventory.")
			s.actor.SetTimer("peek", config.PeekCD*2)
			s.ok = true
			return
		}

		s.actor.SetTimer("peek", config.PeekCD)
		inv := whatMob.Inventory.List()
		s.msg.Actor.SendInfo("In their inventory:")
		if len(inv) == 0 {
			s.msg.Actor.Send("  No items")
			return
		} else {
			s.msg.Actor.Send("  ", strings.Join(whatMob.Inventory.List(), ", "))
			return
		}
	}

	// Try searching through players if gm
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		var whatChar *objects.Character
		whatChar = s.where.Chars.Search(name, s.actor)
		if whatChar != nil {
			inv := whatChar.Inventory.List()
			s.msg.Actor.SendInfo("In their inventory:")
			if len(inv) == 0 {
				s.msg.Actor.Send("  No items")
				return
			} else {
				s.msg.Actor.Send("  ", strings.Join(whatChar.Inventory.List(), ", "))
				return
			}
		} else {
			s.msg.Actor.SendBad("Peek whose inventory?")
			return
		}
	}

}
