package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"log"
	"strings"
)

// Syntax: SAY <message> | " <message>
func init() {
	addHandler(pose{},
		"Usage:  pose \n \n Place your character into a passive RP pose!  (Skip the 'IS', auto appended)",
		permissions.Player,
		"POSE")
}

type pose cmd

func (pose) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("How do you want to pose?")
		return
	}

	msg := strings.Join(s.input, " ")

	/*s.actor.ApplyHook("act", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("act", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
		},
	)*/
	s.actor.ApplyHook("say", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("say", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("ooc", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("ooc", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("hide", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("hide", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("combat", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("combat", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("move", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("move", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("gridmove", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("gridmove", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)
	s.actor.ApplyHook("attacked", "pose", -1, "10", -1,
		func() {
			s.actor.Pose = ""
			if _, err := s.actor.Write([]byte(text.Info + "You stop posing." + text.Reset + "\n")); err != nil {
				log.Println("Error writing to player: ", err)
			}
			s.actor.RemoveHook("attacked", "pose")
			return
		},
		func() {
			s.actor.Pose = ""
			return
		},
	)

	s.actor.Pose = msg
	s.msg.Actor.SendGood("You pose: \"", msg, "\"")
	s.msg.Observers.SendGood("You see ", s.actor.Name, " pose: \"", msg, "\"")

	s.ok = true
	return
}
