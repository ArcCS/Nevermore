package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(build_copy{}, "Usage: copy (room|mob|item) (SubjectID) \n \n Use this to copy an existing item in the database \n",
		permissions.Builder,
		"copy", "duplicate")
}

type build_copy cmd

func (build_copy) process(s *state) {
	// TODO: Finish this up to copy something from the database to make a new db entry
	return
	s.ok = true
	return
}
