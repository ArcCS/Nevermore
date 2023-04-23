package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(about{},
		"Display version and codebase information.",
		permissions.Player,
		"about")
}

type about cmd

func (about) process(s *state) {

	s.msg.Actor.SendInfo("We're running Nevermore for Aalynor's Nexus. Copyright 2020-2023 \n" +
		"World is a fan restoration of out of publication original Aalynor's Nexus, 1996-2013 \n" +
		"Some components used from WolfMUD (https://www.wolfmud.org/)")

	s.ok = true
}
