package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(hide{},
		"Usage:  hide <item> # \n\n Hide in the shadows, or attempt to hide an item",
		permissions.Player,
		"hide")
}

type hide cmd

func (hide) process(s *state) {
	//TODO Finish hide
	s.msg.Actor.SendInfo("Where ya gonna hide???")
	s.ok = true
}
