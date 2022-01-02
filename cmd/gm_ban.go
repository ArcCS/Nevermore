package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(ban{},
		"Usage:  ban (character|ip) \n \n Deactivate a character or perma ban an IP",
		permissions.Gamemaster,
		"ban")
}

type ban cmd

func (ban) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Activate what?")
		return
	}

	s.ok = true
	return
}
