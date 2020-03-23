package cmd

/*

Eat, drink, laugh, nod, shrug, bow, sneeze, cough, dance, wave, shake head, whistle,  thumbs up/down, cross fingers/arms, wink, blink, frown, smile, taunt, cook, lean, snap, clap/applaud, love, angry, confused
Poke. Slap. Kick. Burp.
Jump, glare
Hiccup
Stare
Tickle
Hug
Kiss

func init() {
	addHandler(act{}, "act", "emote")
	addHelp("Usage:  act performs for all to see \n \n Perform actions.", 0, "act", "emote")
}

type act cmd

func (act) process(s *state) {

	// Did they send an action?
	if len(s.words) == 0 {
		s.msg.Actor.SendBad("... what were you trying to do???")
		return
	}

	action := strings.Join(s.input, " ")
	who := s.actor.Name

	s.msg.Actor.SendInfo("You ", action)
	s.msg.Observer.SendInfo(who, " ", action)

	s.ok = true
}

 */
