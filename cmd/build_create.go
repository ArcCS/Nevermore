package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(create{},
	"Usage:  create (room|mob|object) name description \n \n Create a brand new object with a name and description. \n Note:  Use the modify command to add modify traits of the object.",
	permissions.Builder,
	"create", "new")
}

type create cmd

func (create) process(s *state) {
	if len(s.words) == 0 {
		s.msg.Actor.SendInfo("Delete what?")
		return
	}


	s.ok = true
	return
}