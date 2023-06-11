package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
	"strings"
)

func init() {
	addHandler(give{},
		"Usage:  give [person] itemName # \n \n Give the specific person an item. \n\n or \n\n give [person] $[value] to give someone gold",
		permissions.Player,
		"GIVE")
}

type give cmd

func (give) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Give who what???")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	whoStr := s.words[1]

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		if len(s.words) < 3 {
			s.msg.Actor.SendInfo("Who do you want to give it to?")
			return
		} else {
			whoStr = s.words[2]
		}
	} else {
		whoStr = s.words[1]
	}

	var who *objects.Character
	who = s.where.Chars.Search(whoStr, s.actor)
	if who == nil {
		s.msg.Actor.SendInfo("Give who what???")
		return
	}

	if s.actor.Placement != who.Placement {
		s.msg.Actor.SendBad("They are too far away from you to give them anything.")
		return
	}

	s.participant = who

	// We're going to process a money transaction.
	if strings.HasPrefix(targetStr, "$") {
		if amount64, err := strconv.ParseInt(strings.Trim(targetStr, "$"), 10, 64); err == nil {
			amount := int(amount64)
			if s.actor.Gold.CanSubtract(amount) {
				s.actor.RunHook("act")
				s.actor.Gold.SubIfCan(amount)
				who.Gold.Add(amount)
				s.msg.Actor.SendGood("You give ", targetStr, " to ", who.Name, ".")
				if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					s.msg.Participant.SendGood(s.actor.Name + " gives you " + targetStr)
					s.msg.Observers.SendInfo("You see ", s.actor.Name, " give ", who.Name, " some gold.")
				}
				return
			} else {
				s.msg.Actor.SendInfo("You don't have that much gold.")
				return
			}
		}
	}

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to give away?")
		return
	}

	if target.Flags["permament"] {
		s.msg.Actor.SendBad("You cannot get rid of this item.. it is bound to you.")
		return
	}

	if (who.GetCurrentWeight() + target.GetWeight()) <= who.MaxWeight() {
		s.actor.RunHook("act")
		s.actor.Inventory.Remove(target)
		who.Inventory.Add(target)
	} else {
		s.msg.Actor.SendInfo("They can't carry anymore.")
		return
	}

	s.msg.Actor.SendGood("You give ", target.Name, " to ", who.Name, ".")
	if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		s.msg.Participant.SendGood(s.actor.Name + " gives you " + target.Name)
		s.msg.Observers.SendInfo("You see ", s.actor.Name, " give ", target.Name, " to ", who.Name, ".")
	}

	s.ok = true
}
