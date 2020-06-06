package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(menu{},
           "Usage:  menu \n \n WIP:  Create/Modify a menu prompt ",
           permissions.Dungeonmaster,
           "menu")
}

type menu cmd

func (menu) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Menu's.. mumble.. mumble..")
		return
	}


	s.ok = true
	return
}