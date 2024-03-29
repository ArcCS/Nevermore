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

const (
	None          Permissions = 0
	Anyone        Permissions = 1 << iota
	Player                    //2
	Builder                   //4
	Dungeonmaster             //8
	Gamemaster                //16
	God
	NPC
	Fighter
	Mage
	Thief
	Paladin
	Cleric
	Ranger
	Barbarian
	Bard
	Monk
)
