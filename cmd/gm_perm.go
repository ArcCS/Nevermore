package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/stats"
	"github.com/ArcCS/Nevermore/utils"
	"strings"
)

func init() {
	addHandler(togglePerm{},
		"Usage:  toggle_perm account_name dungeonmaster \n \n Apply a specific privilege to an account",
		permissions.Gamemaster,
		"toggleperm")
}

type togglePerm cmd

func (togglePerm) process(s *state) {
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Change who to what?")
		return
	}

	// Update the DB
	acctLevel := strings.ToLower(s.words[1])
	if acctLevel == "gamemaster" || acctLevel == "gm" {
		if utils.StringIn(s.words[0], stats.ActiveCharacters.List()) {
			character := stats.ActiveCharacters.Find(s.words[0])
			character.Permission.ToggleFlag(permissions.Gamemaster)
		}
		data.TogglePermission(s.words[0], uint32(permissions.Gamemaster))
	} else if acctLevel == "dungeonmaster" || acctLevel == "dm" {
		if utils.StringIn(s.words[0], stats.ActiveCharacters.List()) {
			character := stats.ActiveCharacters.Find(s.words[0])
			character.Permission.ToggleFlag(permissions.Dungeonmaster)
		}
		data.TogglePermission(s.words[0], uint32(permissions.Dungeonmaster))
	} else if acctLevel == "builder" || acctLevel == "build" {
		if utils.StringIn(s.words[0], stats.ActiveCharacters.List()) {
			character := stats.ActiveCharacters.Find(s.words[0])
			character.Permission.ToggleFlag(permissions.Builder)
		}
		data.TogglePermission(s.words[0], uint32(permissions.Builder))
	} else {
		s.msg.Actor.SendInfo("Appropriate permission toggle not found: gamemaster, gm, dungeonmaster, dm, builder, build")
	}

	s.ok = true
	return
}
