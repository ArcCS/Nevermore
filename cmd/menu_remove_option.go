package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(remove_option{},
	"",
	permissions.Builder,
	"remove_option")
}

type remove_option cmd

func (remove_option) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: List the items in this story inventory
}
