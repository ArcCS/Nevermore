package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobAttack{},
		"",
		permissions.Anyone,
		"$MATTACK")
}

type mobAttack cmd

func (mobAttack) process(s *state) {
	return
}
