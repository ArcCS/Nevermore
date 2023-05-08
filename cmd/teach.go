package cmd

/*
import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strings"
)

func init() {
	addHandler(teach{},
		"Usage:  teach character spell_name # \n \n Teach a spell to someone else.",
		permissions.Player,
		"TEACH")
}

type teach cmd

func (teach) process(s *state) {

	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("You need to specify the person and the spell you want to teach.")
		return
	}

	name := s.words[0]
	spell := strings.ToLower(s.words[1])

	// Try searching inventory where we are
	whatChar := s.where.Chars.Search(name, s.actor)

	// Was item to read found?
	if whatChar == nil {
		s.msg.Actor.SendBad("We can't find that player.")
		return
	}
	s.participant = whatChar
	if utils.StringIn(spell, []string{"light", "curepoison", "hurt", "burn", "blister", "rumble"}) {
		s.msg.Actor.SendGood("You teach ", spell, " to  "+whatChar.Name)
		whatChar.Spells = append(s.actor.Spells, spell)
		s.msg.Participant.SendGood("You learn ", spell, " from "+s.actor.Name)
		return
	} else if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		if _, ok := objects.Spells[spell]; ok {
			s.msg.Actor.SendGood("You teach ", spell, " to  "+whatChar.Name)
			whatChar.Spells = append(s.actor.Spells, spell)
			s.msg.Participant.SendGood("You learn ", spell, " from "+s.actor.Name)
			return
		} else {
			s.msg.Actor.SendBad("That's not a known spell.")
		}
	} else {
		s.msg.Actor.SendBad("That's not a spell that you can teach.")
	}

	s.ok = true
}

*/
