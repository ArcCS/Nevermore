// Neo4j wrapper for all of the room data

package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadRooms() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	data, err := execRead("MATCH (r:room) OPTIONAL MATCH (r)-[e:exit]->(d:room) OPTIONAL MATCH (r)-[s:spawns]->(m:mob) RETURN "+
		`{room_id: r.room_id, creator: r.creator, name: r.name, description: r.description, encounter_rate: r.encounter_rate, 
	encounters: collect({chance: s.chance, mob_id: m.mob_id}),
	mobs: r.mobs,
	inventory: r.inventory,
	exits: collect({direction:e.name, description: e.description, placement: e.placement, key_id: e.key_id, dest: d.room_id, 
	flags:{closeable: e.closeable,
	closed: e.closed,
	autoclose: e.autoclose,
	lockable: e.lockable,
	unpickable: e.unpickable,
	locked: e.locked,
	hidden: e.hidden,
	invisible: e.invisible,
	levitate: e.levitate,
	day_only: e.day_only,
	night_only: e.night_only,
	placement_dependent: e.placement_dependent}}), flags:{train: r.train, active: r.active, repair: r.repair,
	mana_drain: r.mana_drain,
	no_summon: r.no_summon,
	heal_fast:  r.heal_fast,
	no_teleport: r.no_teleport,
	lo_level: r.lolevel,
	no_scry: r.no_scry,
	shielded: r.shielded,
	dark_always: r.dark_always,
	light_always: r.light_always,
	natural_light: r.natural_light,
	indoors: r.indoors,
	fire: r.fire,
	encounters_on: r.encounters_on,
	no_word_of_recall: r.no_word_of_recall,
	water: r.water,
	no_magic: r.no_magic,
	urban: r.urban,
	underground: r.underground,
	hilevel: r.hilevel,
	earth: r.earth,
	wind: r.wind}}`, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		roomList = append(roomList, row.Values[0].(map[string]interface{}))
	}
	return roomList
}

func LoadRoom(room_id int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	data, err := execRead("MATCH (r:room {room_id: $room_id}) OPTIONAL MATCH (r)-[e:exit]->(d:room) OPTIONAL MATCH (r)-[s:spawns]->(m:mob) RETURN "+
		`{room_id: r.room_id, creator: r.creator, name: r.name, description: r.description, encounter_rate: r.encounter_rate, 
	encounters: collect({chance: s.change, mob_id: m.mob_id}),
	mobs: r.mobs,
	inventory: r.inventory,
	exits: collect({direction:e.name, description: e.description, placement: e.placement, key_id: e.key_id, dest: d.room_id, 
	flags:{closeable: e.closeable,
	closed: e.closed,
	autoclose: e.autoclose,
	lockable: e.lockable,
	unpickable: e.unpickable,
	locked: e.locked,
	hidden: e.hidden,
	invisible: e.invisible,
	levitate: e.levitate,
	day_only: e.day_only,
	night_only: e.night_only,
	placement_dependent: e.placement_dependent}}), flags:{train: r.train, active: r.active, repair: r.repair,
	mana_drain: r.mana_drain,
	no_summon: r.no_summon,
	heal_fast:  r.heal_fast,
	no_teleport: r.no_teleport,
	lo_level: r.lolevel,
	no_scry: r.no_scry,
	shielded: r.shielded,
	dark_always: r.dark_always,
	light_always: r.light_always,
	natural_light: r.natural_light,
	indoors: r.indoors,
	fire: r.fire,
	encounters_on: r.encounters_on,
	no_word_of_recall: r.no_word_of_recall,
	water: r.water,
	no_magic: r.no_magic,
	urban: r.urban,
	underground: r.underground,
	hilevel: r.hilevel,
	earth: r.earth,
	wind: r.wind}}`,
		map[string]interface{}{
			"room_id": room_id,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	return data[0].Values[0].(map[string]interface{})
}

func LoadExit(exitName string, room_id int, toId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	data, err := execRead("MATCH (r:room)-[e:exit]->(d:room) WHERE r.room_id=$fromId AND e.name=$exitname AND d.room_id=$toId RETURN "+
		`{direction:e.name, description: e.description, placement: e.placement, key_id: e.key_id, dest: d.room_id, 
	flags:{closeable: e.closeable,
	closed: e.closed,
	autoclose: e.autoclose,
	lockable: e.lockable,
	unpickable: e.unpickable,
	locked: e.locked,
	hidden: e.hidden,
	invisible: e.invisible,
	levitate: e.levitate,
	day_only: e.day_only,
	night_only: e.night_only,
	placement_dependent: e.placement_dependent}}`,
		map[string]interface{}{
			"exitname": exitName,
			"fromId":   room_id,
			"toId":     toId,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return data[0].Values[0].(map[string]interface{})
}

// CreateRoom will create a new room from a roomname and a creator
func CreateRoom(roomName string, creator string) (int, bool) {
	room_id := nextId("room")
	results, err := execWrite(
		"CREATE (r:room) SET "+
			"r.room_id = $room_id, "+
			"r.name = $name, "+
			"r.description = 'This is a nice room you made here... needs a description though.', "+
			"r.encounter_rate = 0,"+
			"r.creator = $creator, "+
			"r.repair = 0, "+
			"r.mana_drain = 0, "+
			"r.no_summon = 0, "+
			"r.heal_fast = 0, "+
			"r.no_teleport = 0, "+
			"r.lo_level = 0, "+
			"r.no_scry = 0, "+
			"r.shielded = 0, "+
			"r.dark_always = 0, "+
			"r.light_always = 0, "+
			"r.natural_light = 1, "+
			"r.indoors = 0, "+
			"r.fire = 0, "+
			"r.encounters_on = 0, "+
			"r.no_word_of_recall = 0, "+
			"r.water = 0, "+
			"r.no_magic = 0, "+
			"r.urban = 0, "+
			"r.underground = 0, "+
			"r.hilevel = 0, "+
			"r.earth = 0, "+
			"r.active = 0, "+
			"r.train = 0, "+
			"r.mobs = '[]', "+
			"r.inventory = '[]', "+
			"r.wind = 0",
		map[string]interface{}{
			"room_id":  room_id,
			"name":    roomName,
			"creator": creator,
		},
	)
	if err != nil {
		log.Println(err)
	}
	if 	results.Counters().NodesCreated() > 0 {
		return room_id, false
	} else {
		return -1, true
	}
}

// CreateExit Create Exits from a map of exitData
func CreateExit(exitData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room), (r2:room) WHERE "+
			"r.room_id = $baseRoom AND r2.room_id = $toRoom "+
			`CREATE (r)-[e:exit]->(r2) SET 
	e.name=$exitName, 
	e.placement=3, 
	e.key_id=0, 
	e.closeable=0,
	e.closed=0,
	e.autoclose=0,
	e.lockable=0,
	e.unpickable=0,
	e.locked=0,
	e.hidden=0,
	e.invisible=0,
	e.levitate=0,
	e.day_only=0,
	e.night_only=0,
	e.placement_dependent=0`,
		map[string]interface{}{
			"exitName": exitData["name"],
			"baseRoom": exitData["fromId"],
			"toRoom":   exitData["toId"],
		},
	)
	if err != nil {
		log.Println(err)
		return true
	}
	if 	results.Counters().RelationshipsCreated() > 0 {
		return false
	} else {
		return true
	}
}

// ExitExists Does this exit exist?
func ExitExists(exitName string, room_id int) bool {
	data, err := execRead("MATCH (r:room)-[e:exit]->() WHERE r.room_id=$room_id AND e.name=$name RETURN e",
		map[string]interface{}{
			"room_id": room_id,
			"name":    exitName,
		},
	)
	if err != nil {
		log.Println(err)
		return true
	}
	if len(data) <= 1 {
		return false
	} else {
		return true
	}
}

// DeleteRoom Delete Room
func DeleteRoom(room_id int) bool {
	results, err := execWrite("MATCH ()-[e:exit]->(r:room)-[e2:exit]->() WHERE r.room_id=$room_id DELETE r, e, e2",
		map[string]interface{}{
			"room_id": room_id,
		},
	)
	if err != nil {
		log.Println(err)
		return true
	}
	if results.Counters().NodesDeleted() > 0 {
		return true
	} else {
		return false
	}
}

// UpdateRoom Update Room
func UpdateRoom(roomData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room) WHERE r.room_id=$room_id SET "+
			"r.name = $name, "+
			"r.description = $description, "+
			"r.repair = $repair, "+
			"r.mana_drain = $mana_drain, "+
			"r.no_summon = $no_summon, "+
			"r.encounter_rate = $encounter_rate,"+
			"r.heal_fast = $heal_fast, "+
			"r.no_teleport = $no_teleport, "+
			"r.lo_level = $lo_level, "+
			"r.no_scry = $no_scry, "+
			"r.shielded = $shielded, "+
			"r.dark_always = $dark_always, "+
			"r.light_always = $light_always, "+
			"r.natural_light = $natural_light, "+
			"r.indoors = $indoors, "+
			"r.fire = $fire, "+
			"r.encounters_on = $encounters_on, "+
			"r.no_word_of_recall = $no_word_of_recall, "+
			"r.water = $water, "+
			"r.no_magic = $no_magic, "+
			"r.urban = $urban, "+
			"r.underground = $underground, "+
			"r.hilevel = $hilevel, "+
			"r.earth = $earth, "+
			"r.active = $active, "+
			"r.train = $train,"+
			"r.mobs = $mobs, "+
			"r.inventory = $inventory, "+
			"r.wind = $wind",
		map[string]interface{}{
			"room_id":           roomData["room_id"],
			"name":              roomData["name"],
			"description":       roomData["description"],
			"encounter_rate":    roomData["encounter_rate"],
			"repair":            roomData["repair"],
			"mana_drain":        roomData["mana_drain"],
			"no_summon":         roomData["no_summon"],
			"heal_fast":         roomData["heal_fast"],
			"no_teleport":       roomData["no_teleport"],
			"lo_level":          roomData["lo_level"],
			"no_scry":           roomData["no_scry"],
			"shielded":          roomData["shielded"],
			"dark_always":       roomData["dark_always"],
			"light_always":      roomData["light_always"],
			"natural_light":     roomData["natural_light"],
			"indoors":           roomData["indoors"],
			"fire":              roomData["fire"],
			"encounters_on":     roomData["encounters_on"],
			"no_word_of_recall": roomData["no_word_of_recall"],
			"water":             roomData["water"],
			"no_magic":          roomData["no_magic"],
			"urban":             roomData["urban"],
			"underground":       roomData["underground"],
			"hilevel":           roomData["hilevel"],
			"earth":             roomData["earth"],
			"wind":              roomData["wind"],
			"active":            roomData["active"],
			"train":             roomData["train"],
			"mobs":             roomData["mobs"],
			"inventory":       	roomData["inventory"],
		},
	)

	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	}
	return false
}

// RenameExit Renames an exit based on it's previous name, the new name, and the connecting room information
func RenameExit(exitName string, oldName string, baseRoom int, toRoom int) bool {
	results, err := execWrite(
		"MATCH (r:room)-[e:exit]->(r2:room) WHERE "+
			"r.room_id=$baseRoom AND r2.room_id=$toRoom AND e.name=$oldexit SET "+
			`e.name=$exitname`,
		map[string]interface{}{
			"exitname": exitName,
			"oldexit":  oldName,
			"baseRoom": baseRoom,
			"toRoom":   toRoom,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	}
	return false
}

// UpdateExit Update Exit based on
func UpdateExit(exitData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room)-[e:exit]->(r2:room) WHERE "+
			"r.room_id=$baseRoom AND r2.room_id=$toRoom AND e.name=$exitname SET "+
			`e.placement=$placement, 
	e.description=$description,
	e.key_id=$key_id, 
	e.closeable=$closeable,
	e.closed=$closed,
	e.autoclose=$autoclose,
	e.lockable=$lockable,
	e.unpickable=$unpickable,
	e.locked=$locked,
	e.hidden=$hidden,
	e.invisible=$invisible,
	e.levitate=$levitate,
	e.day_only=$day_only,
	e.night_only=$night_only,
	e.placement_dependent=$placement_dependent`,
		map[string]interface{}{
			"exitname":            exitData["exitname"],
			"baseRoom":            exitData["fromId"],
			"toRoom":              exitData["toId"],
			"description":         exitData["description"],
			"placement":           exitData["placement"],
			"key_id":              exitData["key_id"],
			"closeable":           exitData["closeable"],
			"closed":              exitData["closed"],
			"autoclose":           exitData["autoclose"],
			"lockable":            exitData["lockable"],
			"unpickable":          exitData["unpickable"],
			"locked":              exitData["locked"],
			"hidden":              exitData["hidden"],
			"invisible":           exitData["invisible"],
			"levitate":            exitData["levitate"],
			"day_only":            exitData["day_only"],
			"night_only":          exitData["night_only"],
			"placement_dependent": exitData["placement_dependent"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	}
	return false
}

// DeleteExit Delete Exit based on the exit name and the containing room
func DeleteExit(exitName string, room_id int) bool {
	results, err := execWrite("MATCH (r:room)-[e:exit]->() WHERE r.room_id=$room_id AND e.name=$name DELETE e",
		map[string]interface{}{
			"room_id": room_id,
			"name":    exitName,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	}
	return false
}

func SearchRoomName(searchStr string, skip int) []interface{} {
	data, err := execRead("MATCH (r:room) WHERE toLower(r.name) CONTAINS toLower($search) RETURN {room_id: r.room_id, creator: r.creator, name: r.name} ORDER BY r.name SKIP $skip LIMIT $limit",
		map[string]interface{}{
			"search": searchStr,
			"skip":   skip,
			"limit":  config.Server.SearchResults,
		},
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		roomList = append(roomList, row.Values[0].(map[string]interface{}))
	}
	return roomList
}

func SearchRoomDesc(searchStr string, skip int) []interface{} {
	data, err := execRead("MATCH (r:room) WHERE toLower(r.description) CONTAINS toLower($search) RETURN {room_id: r.room_id, creator: r.creator, name: r.name} ORDER BY r.name SKIP $skip LIMIT $limit",
		map[string]interface{}{
			"search": searchStr,
			"skip":   skip,
			"limit":  config.Server.SearchResults,
		},
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		roomList = append(roomList, row.Values[0].(map[string]interface{}))
	}
	return roomList
}

func CreateNarrative(narrData map[string]interface{}) (bool, error) {
	results, err := execWrite(
		"MATCH (r:room) WHERE "+
			"r.room_id = $room_id "+
			`CREATE (r)-[nar:narr]->(n:narrative) SET 
	n.text=$narrText, n.title=$narrTitle`,
		map[string]interface{}{
			"room_id":    narrData["room_id"],
			"narrTitle": narrData["narrTitle"],
			"narrText":  narrData["narrText"],
		},
	)
	if err != nil {
		log.Println(err)
	}
	if 	results.Counters().NodesCreated() > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func UpdateNarrative(narrData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room)-[nar:narr]-> WHERE "+
			"r.room_id = $room_id "+
			`CREATE (r)-[nar:narr]->(n:narrative) SET 
	n.text=$narrText, n.title=$narrTitle`,
		map[string]interface{}{
			"room_id":    narrData["room_id"],
			"narrTitle": narrData["narrTitle"],
			"narrText":  narrData["narrText"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	}
	return false
}
