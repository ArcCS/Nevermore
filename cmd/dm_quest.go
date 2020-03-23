package cmd

func init() {
	addHandler(quest{}, "quest")
	addHelp("Usage:  quest \n \n WIP:  Create/Modify a quest ", 60, "quest")
}

type quest cmd

func (quest) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Quest's mumble mumble...")
		return
	}


	s.ok = true
	return
}