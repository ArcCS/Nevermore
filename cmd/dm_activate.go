package cmd

func init() {
	addHandler(activate{}, "activate")
	addHelp("Usage:  activate (room|exit) (name|id) \n \n Activate a room or exit so that it can be seen in the world. ", 60, "activate")
}

type activate cmd

func (activate) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Activate what?")
		return
	}


	s.ok = true
	return
}