package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"strconv"
)

func init() {
	addHandler(show{},
		"Usage:  show itemName [person] # \n \n Show your item off to someone else.",
		permissions.Player,
		"SHOW")
}

type show cmd

func (show) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Show who what?")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	whoStr := s.words[1]

	if val, err := strconv.Atoi(s.words[1]); err == nil {
		targetNum = val
		if len(s.words) < 3 {
			s.msg.Actor.SendInfo("Who do you want to show it to?")
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

	target := s.actor.Inventory.Search(targetStr, targetNum)

	if target == nil {
		s.msg.Actor.SendInfo("What're you trying to show off?")
		return
	}

	s.msg.Actor.SendGood("You show ", target.Name, " to ", who.Name, ".")
	s.msg.Participant.SendGood(s.actor.Name + " shows you " + target.Name)
	s.msg.Participant.SendInfo(target.Look())
	if !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		s.msg.Observers.SendInfo("You see ", s.actor.Name, " show ", target.Name, " to ", who.Name, ".")
	}

	s.ok = true
}
