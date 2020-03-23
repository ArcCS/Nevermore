package cmd

func init() {
	addHandler(about{}, "about")
	addHelp("Display version and codebase information.", 0, "about")
}

type about cmd

func (about) process(s *state) {

	// Echo some stuff
	s.msg.Actor.SendInfo("We're running Nevermore for Aalynor's Nexus. " +
		"©2019 \n" +
		"Some components used from WolfMUD (https://www.wolfmud.org/)")

	s.ok = true
}
