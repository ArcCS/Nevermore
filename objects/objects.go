package objects

// Everyone's an object,  put all the basics here and then everyone gets to ride this train

import (
	"encoding/json"
	"fmt"
	"github.com/ArcCS/Nevermore/prompt"
	"strings"
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
	Placement   int
	Commands    map[string]prompt.MenuItem
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
	o.Commands[strings.ToUpper(cmdItem)] = prompt.MenuItem{
		Command: strings.ToUpper(cmdCmd),
	}
}

func (o *Object) RemoveCommand(cmdItem string) {
	delete(o.Commands, strings.ToUpper(cmdItem))
}

func (o *Object) String() string {
	return fmt.Sprintf("%p %[1]T", o)
}

func (o *Object) ChangePlacement(place int) bool {
	if place < 5 && place > 0 {
		o.Placement = place
		return true
	}
	return false
}

func (o *Object) SerializeCommands() string {
	cmdList := make(map[string]string, 0)

	if len(o.Commands) == 0 {
		return "[]"
	}

	for key, val := range o.Commands {
		cmdList[key] = val.Command
	}

	data, err := json.Marshal(cmdList)
	if err != nil {
		return "[]"
	} else {
		return string(data)
	}
}

func DeserializeCommands(jsonVals string) map[string]prompt.MenuItem {
	commandList := make(map[string]prompt.MenuItem)
	obj := make(map[string]string, 0)
	err := json.Unmarshal([]byte(jsonVals), &obj)
	if err != nil {
		//log.Println("Error deserializing Command list" + err.Error())
		return commandList
	}
	for key, cmdString := range obj {
		commandList[key] = prompt.MenuItem{Command: cmdString}
	}
	return commandList
}
