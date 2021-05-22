package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(whisper{},
		"Usage:  whisper \n \n Get the specified item.",
		permissions.Player,
		"WHISPER")
}

type whisper cmd

func (whisper) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon.")
	//TODO Need to do the whispering
}
