package cmd

func init() {
	addHandler(use{}, "USE")
	addHelp("Usage:  use item # \n\n Use an item", 0, "use")
}

//TODO: Map out the use of items and the effect they map to under spells

type use cmd

func (use) process(s *state) {
	if s.actor.Class == 50 {
		s.msg.Actor.SendInfo("As a builder you can't use these commands.")
		return
	}
	s.msg.Actor.SendGood("Use what?")
	s.ok = true
}
