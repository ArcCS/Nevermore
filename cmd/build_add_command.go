package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(add_command{}, "Add's a command to a mob.",
	permissions.Builder,
	"$BUY")
}

type add_command cmd

func (add_command) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Add Menu/Command to a Mob
}
