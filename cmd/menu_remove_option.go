package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(removeOption{},
		"",
		permissions.Builder,
		"remove_option")
}

type removeOption cmd

func (removeOption) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")

}
