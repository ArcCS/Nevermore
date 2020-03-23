package cmd

import (
	"github.com/ArcCS/Nevermore/stats"
	"strings"
)

func init() {
	addHandler(broadcast{}, "BROADCAST")
	addHelp("Usage:  broadcast I have so many things to talk about! \n \n Broadcast messages to the entire realm at the cost of a broadcast point.", 0, "broadcast")
}

type broadcast cmd

func (broadcast) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	/* TODO: Uncomment this block after Beta
	if s.actor.Tier	< 50 {
		if s.actor.Broadcasts < 1 {
			s.msg.Actor.SendBad("You're out of broadcasts today.")
		}else{
			s.actor.Broadcasts -= 1
			stats.ActiveCharacters.MessageAll("### " + s.actor.Name + ": " + strings.Join(s.input, " "))
		}
	}else {
	 */
		stats.ActiveCharacters.MessageAll("### " + s.actor.Name + ": " + strings.Join(s.input, " "))
	//}

	s.ok = true
	return
}
