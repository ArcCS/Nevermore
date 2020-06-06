package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(activate{},
           "Usage:  activate (room|exit) (name|id) \n \n Activate a room or exit so that it can be seen in the world. ",
           permissions.Dungeonmaster,
           "activate")
}

type activate cmd

func (activate) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Activate what?")
		return
	}


	s.ok = true
	return
}