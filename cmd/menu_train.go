package cmd

func init() {
	addHandler(train{}, "$TRAIN")
}

type train cmd

func (train) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
