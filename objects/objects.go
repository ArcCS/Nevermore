package objects

// Everyone's an object,  put all the basics here and then everyone gets to ride this train

import (
	"fmt"
	"github.com/ArcCS/Nevermore/prompt"
)

/*
Object is a base level instantiation of any world item.
It includes a name, a description, and a handful of properties
to determine how it functions in the world.
The rest of the the world constructs from this
*/
type Object struct {
	Name        string
	Description string
	Placement int
	Commands map[string]prompt.MenuItem
}

var ObjectCount chan uint

func init() {
	ObjectCount = make(chan uint)
}

func (o *Object) EmptyCommands() {
	for k := range o.Commands {
		delete(o.Commands, k)
	}
}

func (o *Object) AddCommands(cmdItem string, cmdCmd string) {
	o.Commands[cmdItem] = prompt.MenuItem{
		Command: cmdCmd,
	}
}

func (o *Object) String() string {
	return fmt.Sprintf("%p %[1]T", o)
}

func (o *Object) Free() {

}

func (o *Object) ChangePlacement(place int) bool {
	if place < 5 && place > 0 {
		o.Placement = place
		return true
	}
	return false
}
