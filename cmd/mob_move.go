package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobMove{},
		"Usage:  balance \n \n Displays the ",
		permissions.Player,
		"$MMOVE")
}

type mobMove cmd

func (mobMove) process(s *state) {
	return
}
