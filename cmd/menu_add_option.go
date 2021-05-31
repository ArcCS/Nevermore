package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(add_option{}, "",
	permissions.Builder,
	"add_option")
}

type add_option cmd

func (add_option) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Add options to a menu
}
