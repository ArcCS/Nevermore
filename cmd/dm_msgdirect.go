package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"strings"
)

func init() {
	addHandler(msg{},
		"Usage:  msg character_name \n \n Plant an RP message on a singular character",
		permissions.Dungeonmaster|permissions.Gamemaster,
		"msg")
}

type msg cmd

func (msg) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendBad("Msg what to who?")
		return
	}
	whoStr := s.words[0]
	message := strings.Join(s.input[1:], " ")
	who := objects.ActiveCharacters.Find(whoStr)
	if who != nil {
		objects.ActiveCharacters.Lock()
		who.Write([]byte(text.White + message + "\"" + text.Reset + "\n"))
		objects.ActiveCharacters.Unlock()
		s.msg.Actor.SendGood("GM messaged:, \"" + message + "\", to " + who.Name)
	} else {
		s.msg.Actor.SendBad("Send what to who?")
	}

	s.ok = true
	return
}
