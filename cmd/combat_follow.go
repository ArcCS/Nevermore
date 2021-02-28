package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(follow{},
           "Usage:  follow player # \n\n Become a part of another players party",
           permissions.Player,
           "follow")
}

type follow cmd

func (follow) process(s *state) {
	//TODO Finish follow and party dynamics
	s.msg.Actor.SendInfo("Who ya followin'??")
	s.ok = true
}

