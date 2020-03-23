package objects
// Everyone's an object,  put all the basics here and then everyone gets to ride this train

import (
	"fmt"
)

/*
Object is a base level instantiation of any world item.
// It includes a name, a description, and a handful of properties
// to determine how it functions in the world.
The rest of the the world constructs from this
*/
type Object struct {
	Name string
	Description string

	Placement int64
	// Commands and Script
	// Todo: Commands are a list of bound menu items that dispatch hangs on to
	// EG. "touch orb"  "damage 50 $player, teleport player 540"
}

var ObjectCount chan uint64

func init() {
	ObjectCount = make(chan uint64)
}


func (o *Object) String() string {
	return fmt.Sprintf("%p %[1]T", o)
}

func (o *Object) Free() {

}

func (o *Object) ChangePlacement(place int64) bool{
	if place < 5 && place > 0 {
		o.Placement = place
		return true
	}
	return false
}
