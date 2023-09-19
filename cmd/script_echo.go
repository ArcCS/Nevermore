package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"strings"
)

// Syntax: $POOF
func init() {
	addHandler(echo{},
		"",
		permissions.Anyone,
		"$ECHO")
}

type echo cmd

func (echo) process(s *state) {
	s.msg.Actor.SendInfo(strings.Join(s.words, " "))
}
