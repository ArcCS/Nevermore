package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(addOption{}, "",
		permissions.Builder,
		"add_option")
}

type addOption cmd

func (addOption) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Add options to a menu
}
