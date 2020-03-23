package cmd

import (
	"github.com/ArcCS/Nevermore/stats"
	"strings"
)

func init() {
	addHandler(appeal{}, "APPEAL")
	addHelp("Usage:  appeal HELP ME OH GODS OR CREATORS\n \n Appeal a message to higher powers.  Note: Append OOC to clarify it's a non RP issue", 0, "appeal")
}

type appeal cmd

func (appeal) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("What did you want to say?")
		return
	}

	stats.ActiveCharacters.MessageGM("### " + s.actor.Name + " appeals: " + strings.Join(s.input, " "))

	s.ok = true
	return
}
