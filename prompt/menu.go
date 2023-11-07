package prompt

type Menu struct {
	Text    string
	Options map[string]MenuItem
}

type MenuItem struct {
	Command string
}

func (m *Menu) DisplayMenu() {
	return
}
