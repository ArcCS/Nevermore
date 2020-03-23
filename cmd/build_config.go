package cmd

func init() {
	addHandler(config_var{}, "config")
	addHelp("Usage:  config list  \n Show all configurable variables \n" +
		                      "config var_name value \n Change the configuration value" ,60, "config")
}

type config_var cmd

func (config_var) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 55 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Add what?")
		return
	}

	// TODO Various configuration options

	s.ok = true
	return
}
