package cmd

// Syntax: WHO
func init() {
	addHandler(who{}, "SAVE")
	addHelp("Usage:  Commit your current character state to the db.", 0, "save")
}

type save cmd

func (save) process(s *state) {
	s.msg.Actor.SendGood("Saving....")
	s.actor.Save()
	s.msg.Actor.SendGood("Saved.")
}
