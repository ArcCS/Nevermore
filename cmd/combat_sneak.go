package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(sneak{},
           "Usage:  sneak (ExitName) # \n\n Attempt to sneak into another area",
           permissions.Thief & permissions.Ranger & permissions.Monk,
           "sneak")
}

type sneak cmd

func (sneak) process(s *state) {
	//TODO Finish sneak
	s.msg.Actor.SendInfo("Where ya sneakin'??")
	s.ok = true
}

