package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(bonus{},
           "Usage:  bonus player/all ## Bonus a player or use bonus all to bonus the entire room",
           permissions.Dungeonmaster,
           "bonus")
}

type bonus cmd

func (bonus) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Bonus who with what?")
		return
	}

	s.msg.Observer.SendInfo("WIP")
	s.ok = true
	return
}
