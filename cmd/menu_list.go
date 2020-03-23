package cmd

func init() {
	addHandler(list{}, "$LIST")
}

type list cmd

func (list) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
