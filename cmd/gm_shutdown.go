package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(shutdown{},
           "Usage:  shutdown \n \n Safely shutdown the game",
           permissions.Gamemaster,
           "shutdown")
}

type shutdown cmd

func (shutdown) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("WIP.. remind me to do this.. ")
		return
	}


	s.ok = true
	return
}