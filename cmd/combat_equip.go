package cmd

func init() {
	addHandler(equip{}, "equip")
	addHelp("Usage:  equip item # \n\n Try to equip an item from your inventory", 0, "equip")
}

type equip cmd

func (equip) process(s *state) {
	if s.actor.Class == 50 {
		s.msg.Actor.SendInfo("As a builder you can't use these commands.")
		return
	}
	s.msg.Actor.SendInfo("Mighty fine air you want to put on")
	s.ok = true
}
