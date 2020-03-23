package cmd

func init() {
	addHandler(bonus{}, "bonus")
	addHelp("Usage:  bonus player/all ## Bonus a player or use bonus all to bonus the entire room", 60, "bonus")
}

type bonus cmd

func (bonus) process(s *state) {
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Bonus who with what?")
		return
	}

	s.msg.Observer.SendInfo("WIP")
	s.ok = true
	return
}
