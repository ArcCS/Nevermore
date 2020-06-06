package permissions

type Permissions uint32

func (f Permissions) HasFlags(flags ...Permissions) bool {
	groupCheck := true
	for _, flag := range flags {
		if !groupCheck {
			return groupCheck
		}
		if f&flag == 0 {
			groupCheck = false
		}
	}
	return groupCheck
}

func (f Permissions) HasAnyFlags(flags ...Permissions) bool {
	groupCheck := false
	for _, flag := range flags {
		if groupCheck {
			return true
		}
		if f&flag != 0 {
			groupCheck = true
		}
	}
	return groupCheck
}

func (f Permissions) HasFlag(flag Permissions) bool { return f&flag != 0 }
func (f *Permissions) AddFlag(flag Permissions)     { *f |= flag }
func (f *Permissions) ClearFlag(flag Permissions)   { *f &= ^flag }
func (f *Permissions) ToggleFlag(flag Permissions)  { *f ^= flag }

// List of constants to use for bit wise operations
const (
	Anyone Permissions = 1 << iota
	Player
	Builder
	Dungeonmaster
	Gamemaster
	God
	NPC
	Fighter
	Mage
	Thief
	Paladin
	Cleric
	Ranger
	Berserker
	Bard
	Monk
)