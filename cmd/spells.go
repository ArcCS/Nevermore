package cmd

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

// Syntax: spells
func init() {
	addHandler(spells{}, "SPELLS", "SPELL", "SPL")
	addHelp("Usage:  spells \n \n List the spells currently bound to your character, and your spellbook", 0, "spells")
}

type spells cmd

func (spells) process(s *state) {
spell_template := `Your spellbook contains the following spell incantations:
----------------------------------------------------------------------
{{.Spells}}

You sense the following enchantments bound to your lifeforce:
----------------------------------------------------------------------
{{.SpellEffects}}
`
	spell_effects := ""
	for k, _ := range s.actor.Effects {
		spell_effects += k + " "
	}

	data := struct {
		Spells string
		SpellEffects string
	}{
		strings.Join(s.actor.Spells, ", "),
		spell_effects,
	}
	tmpl, _ := template.New("stat_info").Parse(spell_template)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	}else{
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
