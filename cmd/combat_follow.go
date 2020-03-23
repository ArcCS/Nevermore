package cmd

func init() {
	addHandler(follow{}, "follow")
	addHelp("Usage:  follow player # \n\n Become a part of another players party", 0, "follow")
}

type follow cmd

func (follow) process(s *state) {
	if s.actor.Class == 50 {
		s.msg.Actor.SendInfo("As a builder you can't use these commands.")
		return
	}
	s.msg.Actor.SendInfo("Who ya followin'??")
	s.ok = true
}

