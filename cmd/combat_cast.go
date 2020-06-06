package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(cast{},
           "Usage:  cast spell_name target # \n\n Attempts to cast a known spell from your spellbook",
           permissions.Player,
           "cast")
}

type cast cmd

func (cast) process(s *state) {
	s.msg.Actor.SendInfo("You focus really hard but...  couldn't cast anything...")
	s.ok = true
}
