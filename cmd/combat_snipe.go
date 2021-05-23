package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"github.com/ArcCS/Nevermore/utils"
	"math"
	"strconv"
)

func init() {
	addHandler(snipe{},
		"Usage:  snipe target # \n\n Snipe the target, can only be done while hidden",
		permissions.Ranger,
		"snipe")
}

type snipe cmd

func (snipe) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Snipe what exactly?")
		return
	}

	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You must be at least tier 10 to use this skill.")
		return
	}

	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You must be hidden to snipe.")
	}

	ready, msg := s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.RunHook("combat")

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(name, nameNum, true)
	if whatMob != nil {

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Main.ItemType != 4 {
			s.msg.Actor.SendBad("You can only snipe with a ranged weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement == whatMob.Placement {
			s.msg.Actor.SendBad("You are too close to snipe them.")
			return
		}

		_, ok := whatMob.ThreatTable[s.actor.Name]
		if ok {
			s.msg.Actor.SendBad("You have already engaged ", whatMob.Name, " in combat!")
			return
		}

		s.actor.RunHook("combat")

		curChance := config.SnipeChance

		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			curChance = 100
		}

		curChance += s.actor.Dex.Current * config.SnipeChancePerPoint
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor.Name)
		if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
			actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(config.CombatModifiers["snipe"]))))
			s.msg.Actor.SendInfo("You sniped the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
			s.msg.Observers.SendInfo(s.actor.Name + " bashes " + whatMob.Name)
			if whatMob.Stam.Current <= 0 {
				s.msg.Actor.SendInfo("You killed " + whatMob.Name + text.Reset)
				s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name + text.Reset)
				//TODO Calculate experience
				stringExp := strconv.Itoa(whatMob.Experience)
				for k := range whatMob.ThreatTable {
					s.where.Chars.Search(k, true).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
					s.where.Chars.Search(k, true).Experience.Add(whatMob.Experience)
				}
				s.msg.Observers.SendInfo(whatMob.Name + " dies.")
				s.msg.Actor.SendInfo(whatMob.DropInventory())
				objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
				whatMob = nil
			}
			s.actor.SetTimer("combat", config.CombatCooldown)
			return
		}else{
			s.msg.Actor.SendBad("You failed to snipe ", whatMob.Name , ", you fumble your weapon!")
			s.msg.Observer.SendInfo(s.actor.Name + " fails to snipe and fumbles their weapon. ")
			s.actor.SetTimer("global", 25)
			_, what := s.actor.Equipment.Unequip(s.actor.Equipment.Main.Name)
			if what != nil {
				s.actor.Inventory.Lock()
				s.actor.Inventory.Add(what)
				s.actor.Inventory.Unlock()
				s.ok = true
				return
			}
		}
	}

	s.msg.Actor.SendInfo("Bash what?")
	s.ok = true
}
