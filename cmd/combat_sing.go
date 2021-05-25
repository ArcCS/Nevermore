package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
)

func init() {
	addHandler(sing{},
		"Usage:  sing song_name # \n\n Use your instrument to sing a song!",
		permissions.Bard,
		"sing")
}

type sing cmd

func (sing) process(s *state) {
	//TODO Finish Sing Command
	return
}
