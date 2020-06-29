package objects

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/utils"
)

type Exit struct {
	Object
	ParentId int
	ToId int
	Flags map[string]bool
	KeyId int

}

func NewExit(room_id int, exitData map[string]interface{}) *Exit {
	placement := 3
	if exitData["placement"] == nil{
		placement = 3
	}
	ok := false
	description := ""
	if description, ok = exitData["description"].(string); !ok {
		description = ""
	}
	if exitData["key_id"] == nil {
		exitData["key_id"] = -1
	}
	newExt := &Exit{
		Object{
			exitData["direction"].(string),
			description,
			placement,
		},
		room_id,
		int(exitData["dest"].(int64)),
		make(map[string]bool),
			int(exitData["key_id"].(int64)),
	}
	for k, v := range exitData["flags"].(map[string]interface{}){
		if v == nil{
			newExt.Flags[k] = false
		}else {
			newExt.Flags[k] = int(v.(int64)) != 0
		}
	}
	return newExt
}

func (e *Exit) Look() string{
	return e.Description
}

func (e *Exit) Close() bool {
	if e.Flags["closeable"] == true {
		e.Flags["closed"] = true
		return true
	}
	return false
}

func (e *Exit) Open() bool {
	if e.Flags["locked"] == false {
		if e.Flags["closeable"] == true {
			e.Flags["closed"] = true
			return true
		}
	}
	return false
}
func (e *Exit) ToggleFlag(flagName string) bool {
	if val, exists := e.Flags[flagName]; exists{
		e.Flags[flagName] = !val
		return true
	}else{
		return false
	}
}

func (e *Exit) Save() {
	exitData := make(map[string]interface{})
	exitData["exitname"] = e.Name
	exitData["description"] = e.Description
	exitData["fromId"] = e.ParentId
	exitData["toId"] = e.ToId
	exitData["placement"] = e.Placement
	exitData["key_id"] = e.KeyId
	exitData["closeable"] = utils.Btoi(e.Flags["closeable"])
	exitData["closed"] = utils.Btoi(e.Flags["closed"])
	exitData["autoclose"] = utils.Btoi(e.Flags["autoclose"])
	exitData["lockable"] = utils.Btoi(e.Flags["lockable"])
	exitData["unpickable"] = utils.Btoi(e.Flags["unpickable"])
	exitData["locked"] = utils.Btoi(e.Flags["locked"])
	exitData["hidden"] = utils.Btoi(e.Flags["hidden"])
	exitData["invisible"] = utils.Btoi(e.Flags["invisible"])
	exitData["levitate"] = utils.Btoi(e.Flags["levitate"])
	exitData["day_only"] = utils.Btoi(e.Flags["day_only"])
	exitData["night_only"] = utils.Btoi(e.Flags["night_only"])
	exitData["placement_dependent"] = utils.Btoi(e.Flags["placement_dependent"])
	data.UpdateExit(exitData)
}