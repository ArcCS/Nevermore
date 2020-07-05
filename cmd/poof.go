// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import "github.com/ArcCS/Nevermore/permissions"

// Syntax: $POOF
func init() {
	addHandler(poof{},
	"",
	permissions.Anyone,
	"$POOF")
}

type poof cmd

func (poof) process(s *state) {

	name := s.actor.Name
	if s.actor.Flags["invisible"] == false && s.actor.Flags["hidden"] == false {
		s.msg.Observers.SendGood("", name, " appears in a puff of smoke.")
	}

	s.scriptActor("LOOK")
}
