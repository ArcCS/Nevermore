package cmd

import (
	"bytes"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"log"
	"strings"
	"text/template"
)

// Syntax: spells
func init() {
	addHandler(spell{},
		"Usage:  spells \n \n List the spells currently bound to your character, and your spellbook",
		permissions.Player,
		"SPELLS", "SPELL", "SP")
}

type spell cmd

func (spell) process(s *state) {
	spellTemplate := `Your spellbook contains the following spell incantations:
----------------------------------------------------------------------
{{.Spells}}

You sense the following enchantments bound to your lifeforce:
----------------------------------------------------------------------
{{.SpellEffects}}
`

	var spellEffects []string
	for k := range s.actor.Effects {
		if _, ok := objects.Spells[k]; ok {
			spellEffects = append(spellEffects, k)
		}
	}

	data := struct {
		Spells       string
		SpellEffects string
	}{
		strings.Join(s.actor.Spells, ", "),
		strings.Join(spellEffects, ", "),
	}
	tmpl, _ := template.New("stat_info").Parse(spellTemplate)
	var output bytes.Buffer
	err := tmpl.Execute(&output, data)
	if err != nil {
		log.Println(err)
	} else {
		s.msg.Actor.SendGood(output.String())
	}

	s.ok = true
}
