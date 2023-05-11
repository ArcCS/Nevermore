package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(actwith{},
		"Usage:  'action-with' place holder, these actions require someone else to ",
		permissions.Player,
		"taunt", "bug", "bow", "love", "angry", "glare", "stare", "tickle", "poke", "slap", "kick", "wave", "wink")
}

type actwith cmd

func (actwith) process(s *state) {
	whoWith := ""
	if len(s.words) > 0 {
		targetPlayer := s.where.Chars.Search(s.words[0], s.actor)
		if targetPlayer != nil {
			s.participant = targetPlayer
			whoWith = targetPlayer.Name
		} else {
			targetMob := s.where.Mobs.Search(s.words[0], 1, s.actor)
			if targetMob != nil {
				whoWith = targetMob.Name
			} else {
				s.msg.Actor.SendBad("Who did you want to do that with?")
				s.ok = true
				return
			}
		}
	} else {
		s.msg.Actor.SendBad("Who did you want to do that with?")
		s.ok = true
		return
	}

	cmdStr := strings.ToLower(s.cmd)
	action := ""
	s.actor.RunHook("act")
	if cmdStr == "taunt" {
		// Did they send an action?
		if len(s.words) == 0 {
			s.msg.Actor.SendBad("... what were you trying to do???")
			return
		}
		action = strings.Join(s.input, " ")
	} else if cmdStr == "bow" {
		action = "bows before"
	} else if cmdStr == "love" {
		action = "loves"
	} else if cmdStr == "angry" {
		action = "appears furious at"
	} else if cmdStr == "glare" {
		action = "glares at"
	} else if cmdStr == "bug" {
		action = "lightly harasses"
	} else if cmdStr == "stare" {
		action = "stares at"
	} else if cmdStr == "tickle" {
		action = "tickles"
	} else if cmdStr == "poke" {
		action = "pokes"
	} else if cmdStr == "slap" {
		action = "slaps"
	} else if cmdStr == "kick" {
		action = "kicks"
	} else if cmdStr == "wave" {
		action = "waves at"
	} else if cmdStr == "wink" {
		action = "winks at"
	}

	s.msg.Actor.SendInfo("You " + action + " " + whoWith)
	if s.participant != nil {
		s.msg.Participant.SendInfo(s.actor.Name + " " + action + " you.")
	}
	s.msg.Observers.SendInfo(s.actor.Name + " " + action + " " + whoWith)

	s.ok = true
}
