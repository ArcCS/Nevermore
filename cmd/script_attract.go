// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import "github.com/ArcCS/Nevermore/permissions"

// Syntax: $ATTRACT
func init() {
	addHandler(attract{},
		"",
		permissions.Anyone,
		"$ATTRACT")
}

type attract cmd

func (attract) process(s *state) {

	s.where.Encounter()

}
