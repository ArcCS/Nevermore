package cmd

import "github.com/ArcCS/Nevermore/permissions"

// Syntax: SNEEZE
func init() {
	addHandler(suicide{},
           "Usage:  suicide \n \n Permanently kills your character and removes them from the world.",
           permissions.Player,
           "SUICIDE")
}

type suicide cmd

func (suicide) process(s *state) {


	// Notify actor
	s.msg.Actor.SendGood("Oh jeeze; are you sure you want to do that??????")
	return
	// Notify observers in same location
	who := s.actor.Name
	s.msg.Observers.SendInfo(who, " falls to the ground dead and vanishes complete.")


	s.ok = true
}
