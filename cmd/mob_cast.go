package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobCast{},
		"",
		permissions.Anyone,
		"$MCAST")
}

type mobCast cmd

func (mobCast) process(s *state) {
	return
}
