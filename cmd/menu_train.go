package cmd

import "github.com/ArcCS/Nevermore/permissions"

func init() {
	addHandler(train{}, "", permissions.Player,  "$TRAIN")
}

type train cmd

func (train) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
