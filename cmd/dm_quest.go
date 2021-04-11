package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(quest{},
		"Usage:  quest \n \n WIP:  Create/Modify a quest ",
		permissions.Dungeonmaster,
		"quest")
}

type quest cmd

func (quest) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Quest's mumble mumble...")
		return
	}

	s.ok = true
	return
}
