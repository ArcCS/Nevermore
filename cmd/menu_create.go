package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(createMenu{}, "",
		permissions.Builder,
		"create_menu")
}

type createMenu cmd

func (createMenu) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")

}
