package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(configVar{},
		"Usage:  config list  \n Show all configurable variables \n" +
			"config var_name value \n Change the configuration value",
		permissions.Builder,
		"config")

}

type configVar cmd

func (configVar) process(s *state) {
	// Handle necessary arguements
	if len(s.words) < 1 {
		s.msg.Actor.SendInfo("Edit which config with what value?")
		return
	}

	// TODO Various configuration options

	s.ok = true
	return
}
