package cmd

// Syntax: SNEEZE
func init() {
	addHandler(suicide{}, "SUICIDE")
	addHelp("Usage:  suicide \n \n Permanently kills your character and removes them from the world.", 0, "suicide")
}

type suicide cmd

func (suicide) process(s *state) {


	// Notify actor
	s.msg.Actor.SendGood("Oh jeeze; are you sure you want to do that??????")
	return
	// Notify observers in same location
	who := s.actor.Name
	s.msg.Observer.SendInfo(who, " falls to the ground dead and vanishes complete.")


	s.ok = true
}
