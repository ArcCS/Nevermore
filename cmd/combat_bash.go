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
	addHandler(bash{},
		"Usage:  bash target # \n\n Bash the target",
		permissions.Barbarian,
		"bash")
}

type bash cmd

func (bash) process(s *state) {
	if len(s.input) < 1 {
		s.msg.Actor.SendBad("Bash what exactly?")
		return
	}
	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You must be at least tier 5 to use this skill.")
		return
	}
	// Check some timers
	ready, msg := s.actor.TimerReady("combat_bash")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
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
	whatMob = s.where.Mobs.Search(name, nameNum, s.actor)
	if whatMob != nil {

		// Shortcut a missing weapon:
		if s.actor.Equipment.Main == nil {
			s.msg.Actor.SendBad("You have no weapon to attack with.")
			return
		}

		// Shortcut weapon not being blunt
		if s.actor.Equipment.Main.ItemType != 2 {
			s.msg.Actor.SendBad("You can only bash with a blunt weapon.")
			return
		}

		// Shortcut target not being in the right location, check if it's a missile weapon, or that they are placed right.
		if s.actor.Placement != whatMob.Placement {
			s.msg.Actor.SendBad("You are too far away to bash them.")
			return
		}

		//TODO: Parry/Miss/Resist being bashed?

		// Check the rolls in reverse order from hardest to lowest for bash rolls.
		damageModifier := 1
		stunModifier := 1
		if utils.Roll(config.ThunkRoll, 1, 0) == 1 { // Thunk
			damageModifier = config.CombatModifiers["thunk"]
			s.msg.Actor.SendGood("Thunk!!")
		} else if utils.Roll(config.CrushingRoll, 1, 0) == 1 { // Crushing
			damageModifier = config.CombatModifiers["crushing"]
			s.msg.Actor.SendGood("Craaackk!!")
		} else if utils.Roll(config.ThwompRoll, 1, 0) == 1 { // Thwomp
			damageModifier = config.CombatModifiers["thwomp"]
			s.msg.Actor.SendGood("Thwomp!!")
		} else if utils.Roll(config.ThumpRoll, 1, 0) == 1 { // Thump
			stunModifier = 2
			s.msg.Actor.SendGood("Thump!!")
		}
		whatMob.MobStunned = config.BashStuns * stunModifier
		actualDamage, _ := whatMob.ReceiveDamage(int(math.Ceil(float64(s.actor.InflictDamage()) * float64(damageModifier))))
		whatMob.AddThreatDamage(whatMob.Stam.Max/10, s.actor.Name)
		s.msg.Actor.SendInfo("You bashed the " + whatMob.Name + " for " + strconv.Itoa(actualDamage) + " damage!" + text.Reset)
		s.msg.Observers.SendInfo(s.actor.Name + " bashes " + whatMob.Name)
		if whatMob.Stam.Current <= 0 {
			s.msg.Actor.SendInfo("You killed " + whatMob.Name + text.Reset)
			s.msg.Observers.SendInfo(s.actor.Name + " killed " + whatMob.Name + text.Reset)
			//TODO Calculate experience
			stringExp := strconv.Itoa(whatMob.Experience)
			for k := range whatMob.ThreatTable {
				s.where.Chars.Search(k, s.actor).Write([]byte(text.Cyan + "You earn " + stringExp + " exp for the defeat of the " + whatMob.Name + "\n" + text.Reset))
				s.where.Chars.Search(k, s.actor).Experience.Add(whatMob.Experience)
			}
			s.msg.Observers.SendInfo(whatMob.Name + " dies.")
			s.msg.Actor.SendInfo(whatMob.DropInventory())
			objects.Rooms[whatMob.ParentId].Mobs.Remove(whatMob)
			whatMob = nil
		}
		s.actor.SetTimer("combat_bash", config.BashTimer)
		s.actor.SetTimer("combat", config.CombatCooldown)
		return
	}

	s.msg.Actor.SendInfo("Bash what?")
	s.ok = true
}
