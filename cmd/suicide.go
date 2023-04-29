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
	addHandler(suicideConfirm{},
		"",
		permissions.Player,
		"$CONFIRM_SUICIDE")
}

type suicide cmd

func (suicide) process(s *state) {
	s.msg.Actor.SendGood("Oh jeeze; are you sure you want to do this?  This action cannot be undone, and your character cannot be restored. \n In order to complete this action you must type \"DELETE" + s.actor.Name + "\" (without quotes, and no space between DELETE and your character name)")
	s.actor.AddCommands("DELETE"+s.actor.Name, "$confirm_suicide")
	s.ok = true
}

type suicideConfirm cmd

func (suicideConfirm) process(s *state) {
	s.msg.Observers.SendInfo(s.actor.Name, " falls to the ground dead and vanishes completely.")
	s.msg.Actor.SendGood("As the life drains from you, the world fades and goes dark")
	s.scriptActor("quit")
	go data.DeleteChar(s.actor.Name)
	s.ok = true
}
