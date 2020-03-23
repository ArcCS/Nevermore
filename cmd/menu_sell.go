package cmd

func init() {
	addHandler(sell{}, "$SELL")
}

type sell cmd

func (sell) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
