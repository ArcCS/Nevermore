package cmd

func init() {
	addHandler(whisper{}, "WHISPER")
	addHelp("Usage:  whisper \n \n Get the specified item.", 0, "whisper")
}

type whisper cmd

func (whisper) process(s *state) {
		s.msg.Actor.SendInfo("WIP, coming soon.")

}
