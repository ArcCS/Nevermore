package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"strconv"
)

func init() {
	addHandler(perms{}, "perms")
	addHelp("Usage:  perms account_name ## \n \n Apply a specific privilege to an account", 60, "perms")
}

type perms cmd

func (perms) process(s *state) {
	// Handle Permissions
	if s.actor.Class < 60 {
		s.msg.Actor.SendInfo("Unknown command, type HELP to get a list of commands")
		return
	}
	if len(s.words) < 2 {
		s.msg.Actor.SendInfo("Change who to what?")
		return
	}

	// Update the DB
	acctLevel, _ := strconv.Atoi(s.words[1])
	data.ChangeAcctType(s.words[0], acctLevel)
	s.msg.Actor.SendInfo("Sent account modification")

	s.ok = true
	return
}