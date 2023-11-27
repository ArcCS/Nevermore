package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(mobInvokeAbility{},
		"",
		permissions.None,
		"$MMOVE")
}

type mobInvokeAbility cmd

func (mobInvokeAbility) process(s *state) {
	return
}
