// Neo4j wrapper for all of the room data

package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadRooms() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (r:room) OPTIONAL MATCH (r)-[e:exit]->(d:room) OPTIONAL MATCH (r)-[s:spawns]->(m:mob) RETURN " +
		`{room_id: r.room_id, creator: r.creator, name: r.name, description: r.description, encounter_rate: r.encounter_rate, 
	encounters: collect({chance: s.chance, mob_id: m.mob_id}),
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
	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		roomList = append(roomList, datum)
	}
	return roomList
}

func LoadRoom(roomId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (r:room {room_id: {roomId}}) OPTIONAL MATCH (r)-[e:exit]->(d:room) OPTIONAL MATCH (r)-[s:spawns]->(m:mob) RETURN " +
		`{room_id: r.room_id, creator: r.creator, name: r.name, description: r.description, encounter_rate: r.encounter_rate, 
	encounters: collect({chance: s.change, mob_id: m.mob_id}), exits: collect({direction:e.name, description: e.description, placement: e.placement, key_id: e.key_id, dest: d.room_id, 
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
	map[string]interface {}{
		"roomId": roomId,
	})
	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	return data[0][0].(map[string]interface{})
}

func LoadExit(exitName string, roomId int, toId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (r:room)-[e:exit]->(d:room) WHERE r.room_id={fromId} AND e.name={exitname} AND d.room_id={toId} RETURN " +
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
	map[string]interface {}{
		"exitname": exitName,
		"fromId":  roomId,
		"toId":	 toId,
	})
	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	return data[0][0].(map[string]interface{})
}

// Create Room
func CreateRoom(roomName string, creator string) (int, bool) {
	conn, _ := getConn()
	defer conn.Close()
	room_id := nextId("room")
	result, rtrap := conn.ExecNeo(
		"CREATE (r:room) SET " +
			"r.room_id = {roomId}, " +
			"r.name = {name}, " +
			"r.description = 'This is a nice room you made here... needs a description though.', " +
			"r.encounter_rate = 0," +
			"r.creator = {creator}, " +
			"r.repair = 0, " +
			"r.mana_drain = 0, " +
			"r.no_summon = 0, " +
			"r.heal_fast = 0, " +
			"r.no_teleport = 0, " +
			"r.lo_level = 0, " +
			"r.no_scry = 0, " +
			"r.shielded = 0, " +
			"r.dark_always = 0, " +
			"r.light_always = 0, " +
			"r.natural_light = 1, " +
			"r.indoors = 0, " +
			"r.fire = 0, " +
			"r.encounters_on = 0, " +
			"r.no_word_of_recall = 0, " +
			"r.water = 0, " +
			"r.no_magic = 0, " +
			"r.urban = 0, " +
			"r.underground = 0, " +
			"r.hilevel = 0, " +
			"r.earth = 0, " +
			"r.active = 0, " +
			"r.train = 0, " +
			"r.wind = 0",
		map[string]interface {}{
			"roomId": room_id,
			"name":   roomName,
			"creator": creator,
		},
	)

	if rtrap != nil{
		log.Println(rtrap)
	}
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return room_id, false
	}else {
		return -1, true
	}
}

// Create Exits
func CreateExit(exitData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, rtrap := conn.ExecNeo(
		"MATCH (r:room), (r2:room) WHERE " +
			"r.room_id = {baseRoom} AND r2.room_id = {toRoom} " +
			`CREATE (r)-[e:exit]->(r2) SET 
	e.name={exitName}, 
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
		map[string]interface {}{
			"exitName":        exitData["name"],
			"baseRoom":       exitData["fromId"],
			"toRoom":		exitData["toId"],
		},
	)
	if rtrap != nil{
		log.Println(rtrap)
	}

	numResult, _ := toExit.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

// Does this exit exist?
func ExitExists(exitName string, roomId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (r:room)-(e:exit)->() WHERE r.room_id={room_id} AND e.name={name} RETURN e",
		map[string]interface {}{
			"room_id": roomId,
			"name": exitName,
		},
	)
	if len(data) <= 1 {
		return false
	}else {
		return true
	}
}

// Delete Room
func DeleteRoom(roomId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _ := conn.ExecNeo("MATCH ()-[e:exit]->(r:room)-[e2:exit]->() WHERE r.room_id={room_id} DELETE r, e, e2",
		map[string]interface {}{
			"room_id": roomId,
		},
	)

	numResult, _ := data.RowsAffected()
	if numResult < 1 {
		return false
	}else {
		return true
	}
}

// Update Room
func UpdateRoom(roomData map[string]interface{})  bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"MATCH (r:room) WHERE r.room_id={room_id} SET " +
			"r.name = {name}, " +
			"r.description = {description}, " +
			"r.repair = {repair}, " +
			"r.mana_drain = {mana_drain}, " +
			"r.no_summon = {no_summon}, " +
			"r.encounter_rate = {encounter_rate}," +
			"r.heal_fast = {heal_fast}, " +
			"r.no_teleport = {no_teleport}, " +
			"r.lo_level = {lo_level}, " +
			"r.no_scry = {no_scry}, " +
			"r.shielded = {shielded}, " +
			"r.dark_always = {dark_always}, " +
			"r.light_always = {light_always}, " +
			"r.natural_light = {natural_light}, " +
			"r.indoors = {indoors}, " +
			"r.fire = {fire}, " +
			"r.encounters_on = {encounters_on}, " +
			"r.no_word_of_recall = {no_word_of_recall}, " +
			"r.water = {water}, " +
			"r.no_magic = {no_magic}, " +
			"r.urban = {urban}, " +
			"r.underground = {underground}, " +
			"r.hilevel = {hilevel}, " +
			"r.earth = {earth}, " +
			"r.active = {active}, " +
			"r.train = {train}," +
			"r.wind = {wind}",
		map[string]interface {}{
			"room_id": 		  roomData["room_id"],
			"name":      	  roomData["name"],
			"description":    roomData["description"],
			"encounter_rate": roomData["encounter_rate"],
			"repair": roomData["repair"],
			"mana_drain": roomData["mana_drain"],
			"no_summon": roomData["no_summon"],
			"heal_fast": roomData["heal_fast"],
			"no_teleport": roomData["no_teleport"],
			"lo_level": roomData["lo_level"],
			"no_scry": roomData["no_scry"],
			"shielded": roomData["shielded"],
			"dark_always": roomData["dark_always"],
			"light_always": roomData["light_always"],
			"natural_light": roomData["natural_light"],
			"indoors": roomData["indoors"],
			"fire": roomData["fire"],
			"encounters_on": roomData["encounters_on"],
			"no_word_of_recall": roomData["no_word_of_recall"],
			"water": roomData["water"],
			"no_magic": roomData["no_magic"],
			"urban": roomData["urban"],
			"underground": roomData["underground"],
			"hilevel": roomData["hilevel"],
			"earth": roomData["earth"],
			"wind": roomData["wind"],
			"active": roomData["active"],
			"train": roomData["train"],
		},
	)

	if rtrap != nil{
		log.Println(rtrap)
	}
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

// Rename Exit
func RenameExit(exitName string, oldName string,  baseRoom int, toRoom int) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, etrap := conn.ExecNeo(
		"MATCH (r:room)-[e:exit]->(r2:room) WHERE " +
			"r.room_id={baseRoom} AND r.room_id={toRoom} AND e.name={oldexit} SET " +
			`e.name={exitname}`,
		map[string]interface {}{
			"exitname":      exitName,
			"oldexit":       oldName,
			"baseRoom":		baseRoom,
			"toRoom":		toRoom,
		},
	)
	if etrap != nil{
		log.Println(etrap)
	}
	numResult, _ := toExit.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

// Update Exit
func UpdateExit(exitData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, etrap := conn.ExecNeo(
		"MATCH (r:room)-[e:exit]->(r2:room) WHERE " +
			"r.room_id={baseRoom} AND r.room_id={toRoom} SET " +
			`e.name={exitname}, 
	e.placement={placement}, 
	e.description={description},
	e.key_id={key_id}, 
	e.closeable={closeable},
	e.closed={closed},
	e.autoclose={autoclose},
	e.lockable={lockable},
	e.unpickable={unpickable},
	e.locked={locked},
	e.hidden={hidden},
	e.invisible={invisible},
	e.levitate={levitate},
	e.day_only={day_only},
	e.night_only={night_only},
	e.placement_dependent={placement_dependent}`,
		map[string]interface {}{
			"exitname":        exitData["exitname"],
			"baseRoom":       exitData["fromId"],
			"toRoom":	exitData["toId"],
			"description": exitData["description"],
			"placement": exitData["placement"],
			"key_id": exitData["key_id"],
			"closeable": exitData["closeable"],
			"closed": exitData["closed"],
			"autoclose": exitData["autoclose"],
			"lockable": exitData["lockable"],
			"unpickable": exitData["unpickable"],
			"locked": exitData["locked"],
			"hidden": exitData["hidden"],
			"invisible": exitData["invisible"],
			"levitate": exitData["levitate"],
			"day_only": exitData["day_only"],
			"night_only": exitData["night_only"],
			"placement_dependent": exitData["placement_dependent"],
		},
	)
	if etrap != nil{
		log.Println(etrap)
	}
	numResult, _ := toExit.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}


// Delete Exit
func DeleteExit(exitName string, roomId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, rtrap := conn.ExecNeo("MATCH (r:room)-[e:exit]->() WHERE r.room_id={room_id} AND e.name={name} DELETE e",
		map[string]interface {}{
			"room_id": roomId,
			"name": exitName,
		},
	)
	if rtrap != nil{
		log.Println(rtrap)
	}
	numResult, _ := data.RowsAffected()
	if numResult < 1 {
		return false
	}else {
		return true
	}
}

func SearchRoomName(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (r:room) WHERE toLower(r.name) CONTAINS toLower({search}) RETURN {room_id: r.room_id, creator: r.creator, name: r.name} ORDER BY r.name SKIP {skip} LIMIT {limit}",
		map[string]interface {}{
			"search": searchStr,
			"skip": skip,
			"limit": config.Server.SearchResults,
		},
	)

	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		roomList = append(roomList, datum)
	}
	return roomList
}

func SearchRoomDesc(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (r:room) WHERE toLower(r.description) CONTAINS toLower({search}) RETURN {room_id: r.room_id, creator: r.creator, name: r.name} ORDER BY r.name SKIP {skip} LIMIT {limit}",
		map[string]interface {}{
			"search": searchStr,
			"skip": skip,
			"limit": config.Server.SearchResults,
		},
	)

	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	roomList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		roomList = append(roomList, datum)
	}
	return roomList
}