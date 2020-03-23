package cmd

func init() {
	addHandler(create{}, "create", "new",)
	addHelp("Usage:  create (room|mob|object) name description \n \n Create a brand new object with a name and description. \n Note:  Use the modify command to add modify traits of the object.", 50, "create", "new")
}

type create cmd

func (create) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 50 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Delete what?")
		return
	}


	s.ok = true
	return
}