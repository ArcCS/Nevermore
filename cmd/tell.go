package cmd

func init() {
	addHandler(tell{}, "TELL", "SEND")
	addHelp("Usage:  send character_name \n \n Send a telepathic message to another player", 0, "send", "tell")
}

type tell cmd

func (tell) process(s *state) {
	s.msg.Actor.SendInfo("WIP, coming soon")
	return


}
