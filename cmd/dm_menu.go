package cmd

func init() {
	addHandler(menu{}, "menu")
	addHelp("Usage:  menu \n \n WIP:  Create/Modify a menu prompt ", 60, "menu")
}

type menu cmd

func (menu) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Menu's.. mumble.. mumble..")
		return
	}


	s.ok = true
	return
}