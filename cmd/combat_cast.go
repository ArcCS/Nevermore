package cmd

func init() {
	addHandler(cast{}, "cast")
	addHelp("Usage:  cast spell_name target # \n\n Attempts to cast a known spell from your spellbook", 0, "cast")
}

type cast cmd

func (cast) process(s *state) {
	if s.actor.Class == 50 {
		s.msg.Actor.SendInfo("As a builder you can't use these commands.")
		return
	}
	s.msg.Actor.SendInfo("You focus really hard but...  couldn't cast anything...")
	s.ok = true
}
