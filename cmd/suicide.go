package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(suicide{},
		"Usage:  suicide \n \n Permanently kills your character and removes them from the world.",
		permissions.Player,
		"SUICIDE")
}

type suicide cmd

func (suicide) process(s *state) {
	s.msg.Actor.SendGood("Oh jeeze; are you sure you want to do this?  This action cannot be undone, and your character cannot be restored. \n In order to complete this action you must type \"DELETE" + s.actor.Name + "\" (without quotes)")
	s.actor.AddCommands("DELETE" + s.actor.Name, "$confirm_suicide")
	s.ok = true
}

func init() {
	addHandler(suicide_confirm{},
		"",
		permissions.Player,
		"$CONFIRM_SUICIDE")
}

type suicide_confirm cmd

func (suicide_confirm) process(s *state) {
	s.msg.Observers.SendInfo(s.actor.Name, " falls to the ground dead and vanishes complete.")
	s.msg.Actor.SendGood("As the life drains from you, the world fades and goes dark")
	data.DeleteChar(s.actor.Name)
	s.scriptActor("quit")
	s.ok = true
}