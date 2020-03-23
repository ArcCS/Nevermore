package cmd

import (
	"strconv"

	"github.com/ArcCS/Nevermore/stats"
)

// Syntax: WHO
func init() {
	addHandler(who{}, "WHO")
	addHelp("Usage:  who \n \n Display other currently logged in characters.", 0, "who")
}

type who cmd

func (who) process(s *state) {
	var players []string
	if s.actor.Class == 100{
		players = stats.ActiveCharacters.GMList()
	}else{
		players = stats.ActiveCharacters.List()
	}

	if len(players)-1 == 0 {
		s.msg.Actor.SendInfo("You are all alone in this world.")
		return
	}

	for _, player := range players {
		s.msg.Actor.Send(player)
	}

	var (
		plural = len(players) > 1
		start  = "There is currently "
		end    = "."
	)

	if plural {
		start = "There are currently "
		end = "s."
	}

	s.msg.Actor.Send("")
	s.msg.Actor.Send(start, strconv.Itoa(len(players)-1), " other player", end)

	s.ok = true
}
