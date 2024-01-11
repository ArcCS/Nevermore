package cmd

import (
	"regexp"
	"strings"

	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(act{},
		"Usage:  act performs for all to see \n \n Perform actions.",
		permissions.Player,
		"act", "blink", "blush", "bow", "burp", "cackle", "cheer", "chuckle", "clap", "confused", "cough", "crossarms", "crossfingers", "cry", "dance", "emote", "flex", "flinch", "frown", "gasp", "giggle", "grin", "groan", "hiccup", "jump", "kneel", "laugh", "me", "nod", "ponder", "salute", "shake", "shiver", "shrug", "sigh", "sneeze", "snap", "smile", "smirk", "snicker", "spit", "stare", "stretch", "tap", "thumbsdown", "thumbsup", "wave", "whistle", "wink", "yawn")
}

var actDict = map[string]string{
	"act":          "acts out a scene.",
	"blink":        "blinks slowly.",
	"blush":        "blushes.",
	"bow":          "bows respectfully.",
	"burp":         "burps.",
	"cackle":       "cackles with insane glee!",
	"cheer":        "cheers!",
	"chuckle":      "chuckles politely.",
	"clap":         "claps enthusiastically.",
	"confused":     "looks very confused.",
	"cough":        "coughs.",
	"crossarms":    "crosses their arms.",
	"crossfingers": "crosses their fingers.",
	"cry":          "cries.",
	"dance":        "dances around.",
	"emote":        "acts out a scene.",
	"flex":         "flexes their muscles.",
	"flinch":       "flinches.",
	"frown":        "frowns.",
	"gasp":         "gasps.",
	"giggle":       "giggles.",
	"grin":         "grins.",
	"groan":        "groans.",
	"hiccup":       "hiccups.",
	"jump":         "jumps up and down.",
	"kneel":        "kneels down.",
	"laugh":        "laughs.",
	"nod":          "nods.",
	"ponder":       "ponders the situation.",
	"salute":       "salutes.",
	"shake":        "shakes their head.",
	"shiver":       "shivers.",
	"shrug":        "shrugs.",
	"sigh":         "sighs.",
	"sneeze":       "sneezes, ACHOOO!",
	"snap":         "snaps their fingers.",
	"smile":        "smiles.",
	"smirk":        "smirks.",
	"snicker":      "snickers.",
	"spit":         "spits.",
	"stare":        "stares off into space.",
	"stretch":      "stretches.",
	"tap":          "taps their foot impatiently.",
	"thumbsdown":   "gives a thumbs down.",
	"thumbsup":     "gives a thumbs up.",
	"wave":         "waves.",
	"whistle":      "whistles.",
	"wink":         "winks.",
	"yawn":         "yawns.",
}

type act cmd

func (act) process(s *state) {
	cmdStr := strings.ToLower(s.cmd)
	action := ""
	var ok bool
	s.actor.RunHook("act")
	if cmdStr == "act" || cmdStr == "emote" || cmdStr == "me" {
		// Did they send an action?
		if len(s.words) == 0 {
			s.msg.Actor.SendBad("... what were you trying to do???")
			return
		}
		action = strings.Join(s.input, " ")
		match, _ := regexp.MatchString("([?.,\"'()!;:])", action[len(action)-1:])
		if !match {
			action = action + "."
		}
	} else {
		if action, ok = actDict[cmdStr]; !ok {
			s.msg.Actor.SendBad("Action not available")
			s.ok = true
			return
		}
	}
	data.StoreChatLog(3, s.actor.CharId, 0, action)
	s.msg.Actor.SendInfo(s.actor.Name, " ", action)
	s.msg.Observers.SendInfo(s.actor.Name, " ", action)

	s.ok = true
}
