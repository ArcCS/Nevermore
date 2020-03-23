package cmd

func init() {
	addHandler(equipment{}, "equipment")
	addHelp("Usage:  equipment \n\n Display currently equipped gear", 0, "equipment")
}

type equipment cmd

func (equipment) process(s *state) {
	/*
	{{.Sub_pronoun}} {{.Isare}} wearing {{.Chest}} about {{.Pos_pronoun}} body, some
	  magic Darkplate Leggings (15) on {{.Pos_pronoun}} legs and a Cloak of Fog (15) around
	  {{.Pos_pronoun}} neck.
	{{.Sub_pronoun}} {{.Isare}} holding a {{.Main}} and {{.Offhand}}.
	{{.Sub_pronoun}} {{.Isare}} wearing some {{.Arms}} on {{.Pos_pronoun}} arms.
	{{.Sub_pronoun}} {{.HasHave}} a {{.Finger1}} on {{.Pos_pronoun}} finger.
	{{.Sub_pronoun}} {{.HasHave}} a {{.Finger2}} on {{.Pos_pronoun}} finger.
	{{.Sub_pronoun}} {{.HasHave}} a {{.Finger3}} on {{.Pos_pronoun}} finger.
	{{.Sub_pronoun}} {{.HasHave}} {{.Legs}} on {{.Pos_pronoun}} legs.
	{{.Sub_pronoun}} {{.HasHave}} {{.Feet}} on {{.Pos_pronoun}} feet.
	{{.Sub_pronoun}} {{.Isare}} wearing a {{.Head}}.
	*/
	s.msg.Actor.SendInfo("Yyeeeaaahh look at you wearing.. a whole lotta nothin'")
	s.ok = true
}
