package cmd

import (
	"strings"

	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(act{},
		"Usage:  act performs for all to see \n \n Perform actions.",
		permissions.Player,
		"act", "blink", "blush", "bow", "burp", "cackle", "cheer", "chuckle", "clap", "confused", "cough", "crossarms", "crossfingers", "cry", "dance", "emote", "flex", "flinch", "frown", "gasp", "giggle", "grin", "groan", "hiccup", "hug", "jump", "kneel", "laugh", "nod", "ponder", "salute", "shake", "shiver", "shrug", "sigh", "sneeze", "snap", "smile", "smirk", "spit", "stare", "stretch", "tap", "thumbsdown", "thumbsup", "wave", "whistle", "wink", "yawn")
}

type act cmd

func (act) process(s *state) {
	cmdStr := strings.ToLower(s.cmd)
	action := ""
	s.actor.RunHook("act")
	if cmdStr == "act" || cmdStr == "emote" {
		// Did they send an action?
		if len(s.words) == 0 {
			s.msg.Actor.SendBad("... what were you trying to do???")
			return
		}
		action = strings.Join(s.input, " ")
	} else if cmdStr == "blink" {
		action = "blinks slowly."
	} else if cmdStr == "burp" {
		action = "burps."
	} else if cmdStr == "cackle" {
		action = "cackles maniacally!"
	} else if cmdStr == "clap" {
		action = "claps " + config.TextPosPronoun[s.actor.Gender] + " hands."
	} else if cmdStr == "confused" {
		action = "appears confused."
	} else if cmdStr == "cough" {
		action = "coughs."
	} else if cmdStr == "crossarms" {
		action = "crosses " + config.TextPosPronoun[s.actor.Gender] + " arms."
	} else if cmdStr == "crossfingers" {
		action = "crosses " + config.TextPosPronoun[s.actor.Gender] + " fingers."
	} else if cmdStr == "dance" {
		action = "does a little jig."
	} else if cmdStr == "frown" {
		action = "frowns."
	} else if cmdStr == "grin" {
		action = "grins widely."
	} else if cmdStr == "hiccup" {
		action = "hiccups!"
	} else if cmdStr == "jump" {
		action = "jumps up and down."
	} else if cmdStr == "laugh" {
		action = "laughs."
	} else if cmdStr == "nod" {
		action = "nods."
	} else if cmdStr == "shake" {
		action = "shakes " + config.TextPosPronoun[s.actor.Gender] + " head."
	} else if cmdStr == "shrug" {
		action = "shrugs " + config.TextPosPronoun[s.actor.Gender] + " shoulders."
	} else if cmdStr == "smile" {
		action = "smiles broadly."
	} else if cmdStr == "smirk" {
		action = "smirks."
	} else if cmdStr == "snap" {
		action = "snaps " + config.TextPosPronoun[s.actor.Gender] + " fingers."
	} else if cmdStr == "sneeze" {
		action = "sneezes! ACHOO!"
	} else if cmdStr == "spit" {
		action = "spits on the ground."
	} else if cmdStr == "tap" {
		action = "taps" + config.TextPosPronoun[s.actor.Gender] + " foot."
	} else if cmdStr == "thumbsdown" {
		action = "gives a thumbs down!"
	} else if cmdStr == "thumbsup" {
		action = "gives a thumbs up!"
	} else if cmdStr == "whistle" {
		action = "whistles a tune."
	}
	s.msg.Actor.SendInfo(s.actor.Name, " ", action)
	s.msg.Observers.SendInfo(s.actor.Name, " ", action)

	s.ok = true
}
