package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(remove_command{}, "Add's a command to a mob.",
	permissions.Builder,
	"remove_command")
}

type remove_command cmd

func (remove_command) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	// TODO: Add Menu/Command to a Mob
}
