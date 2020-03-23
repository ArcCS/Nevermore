package cmd

func init() {
	addHandler(shutdown{}, "shutdown")
	addHelp("Usage:  shutdown \n \n Safely shutdown the game", 60, "shutdown")
}

type shutdown cmd

func (shutdown) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("WIP.. remind me to do this.. ")
		return
	}


	s.ok = true
	return
}