package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(tell{},
           "Usage:  send character_name \n \n Send a telepathic message to another player",
           permissions.Player,
           "TELL")
}

type tell cmd

func (tell) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon")
	return


}
