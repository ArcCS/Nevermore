package cmd

func init() {
	addHandler(buy{}, "$BUY")
}

type buy cmd

func (buy) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
