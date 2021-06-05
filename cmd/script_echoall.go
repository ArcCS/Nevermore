// Copyright 2017 Andrew 'Diddymus' Rolfe. All rights reserved.
//
// Use of this source code is governed by the license in the LICENSE file
// included with the source code.

package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: $POOF
func init() {
	addHandler(echoall{},
		"",
		permissions.Anyone,
		"$ECHOALL")
}

type echoall cmd

func (echoall) process(s *state) {
		s.msg.Actor.SendInfo(strings.Join(s.words, " "))
		s.msg.Observers.SendInfo(strings.Join(s.words, " "))
}
