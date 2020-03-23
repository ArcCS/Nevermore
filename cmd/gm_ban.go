package cmd

func init() {
	addHandler(ban{}, "ban")
	addHelp("Usage:  ban (character|ip) \n \n Deactivate a character or perma ban an IP", 60, "ban")
}

type ban cmd

func (ban) process(s *state) {
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