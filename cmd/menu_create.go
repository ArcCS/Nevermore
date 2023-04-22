package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(create_menu{}, "",
		permissions.Builder,
		"create_menu")
}

type create_menu cmd

func (create_menu) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")

}
