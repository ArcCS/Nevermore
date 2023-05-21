package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

func init() {
	addHandler(actwith{},
		"Usage:  'action-with' place holder, these actions require someone else to ",
		permissions.Player,
		"taunt", "bug", "bow", "hug", "angry", "glare", "stare", "tickle", "poke", "slap", "kick", "wave", "wink")
}

var actWithMap = map[string]string{
	"taunt":  "taunts",
	"bug":    "lightly harasses",
	"bow":    "bows before",
	"hug":    "hugs",
	"angry":  "appears furious at",
	"glare":  "glare at",
	"stare":  "stares at",
	"tickle": "tickles",
	"poke":   "pokes",
	"slap":   "slaps",
	"kick":   "kicks",
	"wave":   "waves at",
	"wink":   "winks at",
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
	var ok bool
	s.actor.RunHook("act")
	if action, ok = actWithMap[cmdStr]; !ok {
		s.msg.Actor.SendBad("Action not available")
		s.ok = true
		return
	}

	s.msg.Actor.SendInfo("You " + action + " " + whoWith)
	if s.participant != nil {
		s.msg.Participant.SendInfo(s.actor.Name + " " + action + " you.")
	}
	s.msg.Observers.SendInfo(s.actor.Name + " " + action + " " + whoWith)

	s.ok = true
}
