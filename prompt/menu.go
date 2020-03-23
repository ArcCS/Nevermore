package prompt

// Structure to contain dynamic menus
type Menu struct {
	Text string
	Options map[string]MenuItem
}

// Structure to contain dynamic menus
type MenuItem struct {
	Command string
}

func (m *Menu) DisplayMenu() {
	return
}