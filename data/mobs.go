package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadMobs() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (m:mob) OPTIONAL MATCH (m)-[d:drops]->(i:item) RETURN " +
	`{mob_id:m.mob_id, 
	name:m.name, 
	description:m.description, 
	experience:m.experience, 
	level:m.level, 
	gold:m.gold, 
	constitution:m.constitution, 
	strength:m.strength, 
	intelligence:m.intelligence, 
	dexterity:m.dexterity, 
	piety:m.piety, 
	mpmax:m.mpmax, 
	mpcur:m.mpcur, 
	hpcur:m.hpcur, 
	hpmax:m.hpmax, 
	sdice:m.sdice, 
	ndice:m.ndice, 
	pdice:m.pdice, 
	spells:m.spells, 
	casting_probability:m.casting_probability, 
	armor:m.armor, 
	numwander:m.numwander, 
	wimpyvalue:m.wimpyvalue, 
	air_resistance:m.air_resistance, 
	fire_resistance:m.fire_resistance, 
	earth_resistance:m.earth_resistance, 
	water_resistance:m.water_resistance, 
	drops: collect({chance: d.chance, item_id: i.item_id}),
	flags:{
	hide_encounter: m.hide_encounter, 
	invisible:m.invisible, 
	permanent:m.permanent }}`, nil)
	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	mobList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		mobList = append(mobList, datum)
	}
	return mobList
}

func LoadMob(mobId int64) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (m:mob {mob_id: {mobId}) OPTIONAL MATCH (m)-[d:drops]->(i:item) RETURN " +
		`{mob_id:m.mob_id, 
	name:m.name, 
	description:m.description, 
	experience:m.experience, 
	level:m.level, 
	gold:m.gold, 
	constitution:m.constitution, 
	strength:m.strength, 
	intelligence:m.intelligence, 
	dexterity:m.dexterity, 
	piety:m.piety, 
	mpmax:m.mpmax, 
	mpcur:m.mpcur, 
	hpcur:m.hpcur, 
	hpmax:m.hpmax, 
	sdice:m.sdice, 
	ndice:m.ndice, 
	pdice:m.pdice, 
	spells:m.spells, 
	casting_probability:m.casting_probability, 
	armor:m.armor, 
	numwander:m.numwander, 
	wimpyvalue:m.wimpyvalue, 
	air_resistance:m.air_resistance, 
	fire_resistance:m.fire_resistance, 
	earth_resistance:m.earth_resistance, 
	water_resistance:m.water_resistance, 
	drops: collect({chance: d.chance, item_id: i.item_id}),
	flags:{
	hide_encounter: m.hide_encounter, 
	invisible:m.invisible, 
	permanent:m.permanent }}`,
		map[string]interface {}{
			"mobId": mobId,
		})
	if rtrap != nil{
		log.Println(rtrap)
		return nil
	}
	return data[0][0].(map[string]interface{})
}

// Create Room
func CreateMob(mobName string, creator string) (int64, bool) {
	conn, _ := getConn()
	defer conn.Close()
	mob_id := nextId("mob")
	result, rtrap := conn.ExecNeo(
		"CREATE (m:mob) SET " +
		`mob_id = {mobId}, 
		name = {name}, 
		description=  "A shiny new mob!", 
		experience=  0, 
		level= 1, 
		gold= 0, 
		constitution= 0, 
		strength= 0, 
		intelligence= 0, 
		dexterity= 0, 
		piety= 0, 
		mpmax= 0, 
		mpcur= 0, 
		hpcur= 0, 
		hpmax= 0, 
		sdice= 0, 
		ndice= 0, 
		pdice= 0, 
		spells= "", 
		casting_probability= 0, 
		armor= 0, 
		numwander= 0, 
		wimpyvalue= 0, 
		air_resistance= 0, 
		fire_resistance= 0, 
		earth_resistance= 0, 
		water_resistance= 0, 
		hide_encounter=  0, 
		invisible= 0, 
		permanent= 0`,
		map[string]interface {}{
			"mobId": mob_id,
			"name":   mobName,
			"creator": creator,
		},
	)

	if rtrap != nil{
		log.Println(rtrap)
	}
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return mob_id, false
	}else {
		return -1, true
	}
}

// Update Room
func UpdateMob(mobData map[string]interface{})  bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"MATCH (m:mob) WHERE m.mob_id={mob_id} SET " +
		`m.name={name},
		m.description={description},
		m.experience={experience}, 
		m.level={level}, 
		m.gold={gold}, 
		m.constitution={constitution}, 
		m.strength={strength}, 
		m.intelligence={intelligence}, 
		m.dexterity={dexterity}, 
		m.piety={piety}, 
		m.mpmax={mpmax}, 
		m.mpcur={mpcur}, 
		m.hpcur={hpcur}, 
		m.hpmax={hpmax}, 
		m.sdice={sdice}, 
		m.ndice={ndice}, 
		m.pdice={pdice}, 
		m.spells={spells}, 
		m.casting_probability={casting_probability}, 
		m.armor={armor}, 
		m.numwander={numwander}, 
		m.wimpyvalue={wimpyvalue}, 
		m.air_resistance={air_resistance}, 
		m.fire_resistance={fire_resistance}, 
		m.earth_resistance={earth_resistance}, 
		m.water_resistance={water_resistance}, 
		m.hide_encounter={hide_encounter}, 
		m.invisible={invisible}, 
		m.permanent={permanent}`,
		map[string]interface {}{
			"mob_id": mobData["mob_id"],
			"name": mobData["name"],
			"description": mobData["description"],
			"experience": mobData["experience"],
			"level": mobData["level"],
			"gold": mobData["gold"],
			"constitution": mobData["constitution"],
			"strength": mobData["strength"],
			"intelligence": mobData["intelligence"],
			"dexterity": mobData["dexterity"],
			"piety": mobData["piety"],
			"mpmax": mobData["mpmax"],
			"mpcur": mobData["mpcur"],
			"hpcur": mobData["hpcur"],
			"hpmax": mobData["hpmax"],
			"sdice": mobData["sdice"],
			"ndice": mobData["ndice"],
			"pdice": mobData["pdice"],
			"spells": mobData["spells"],
			"casting_probability": mobData["casting_probability"],
			"armor": mobData["armor"],
			"numwander": mobData["numwander"],
			"wimpyvalue": mobData["wimpyvalue"],
			"air_resistance": mobData["air_resistance"],
			"fire_resistance": mobData["fire_resistance"],
			"earth_resistance": mobData["earth_resistance"],
			"water_resistance": mobData["water_resistance"],
			"hide_encounter": mobData["hide_encounter"],
			"invisible": mobData["invisible"],
			"permanent": mobData["permanent"],
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


// Create Encounter
func CreateEncounter(encounterData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, rtrap := conn.ExecNeo(
		"MATCH (r:room), (m:mob) WHERE " +
			"r.room_id = {roomId} AND m.mob_id = {mobId} " +
			`CREATE (r)-[s:spawns]->(m) SET 
	s.chance={chance}`,
		map[string]interface {}{
			"mobId":        encounterData["mobId"],
			"roomId":       encounterData["roomId"],
			"chance":		encounterData["chance"],
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

// Does a room already have too many encounters?
func SumEncounters(roomId int64) int64 {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (r:room)-[s:spawns]->() WHERE r.room_id={room_id} RETURN {rate_sum: sum(s.chance)}",
		map[string]interface {}{
			"room_id": roomId,
		},
	)
	datum := data[0][0].(map[string]interface{})
	return datum["rate_sum"].(int64)
}

// Delete Mob
func DeleteMob(mobId int64) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _ := conn.ExecNeo("MATCH ()-[s:spawns]->(m:mob)-[d:drops]->() WHERE m.mob_id={mob_id} DELETE s, m, d",
		map[string]interface {}{
			"room_id": mobId,
		},
	)

	numResult, _ := data.RowsAffected()
	if numResult < 1 {
		return false
	}else {
		return true
	}
}

// Update Encounter
func UpdateEncounter(mobData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, etrap := conn.ExecNeo(
		"MATCH (r:room)-[s:spawns]->(m:mob) WHERE " +
			"r.room_id={roomId} AND m.mob_id={mobId} SET " +
			"s.chance={chance}",
		map[string]interface {}{
			"roomId":  mobData["roomId"],
			"mobId":   mobData["mobId"],
			"chance":	mobData["chance"],

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

// Delete encounter
func DeleteEncounter(mobId int64, roomId int64) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, rtrap := conn.ExecNeo("MATCH (r:room)-[s:spawns]->(m:mob) WHERE r.room_id={room_id} AND s.mob_id={mob_id} DELETE s",
		map[string]interface {}{
			"room_id": roomId,
			"mob_id": mobId,
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

// Search mobs
func SearchMobName(searchStr string, skip int64) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (m:mob) WHERE toLower(m.name) CONTAINS toLower({search}) RETURN {name: m.name, mob_id: m.mob_id, level: m.level} ORDER BY m.name SKIP {skip} LIMIT {limit}",
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
	searchList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		searchList = append(searchList, datum)
	}
	return searchList
}

func SearchMobDesc(searchStr string, skip int64) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (m:mob) WHERE toLower(m.description) CONTAINS toLower({search}) RETURN {name: m.name, mob_id: m.mob_id, level: m.level} ORDER BY m.name SKIP {skip} LIMIT {limit}",
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
	searchList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		searchList = append(searchList, datum)
	}
	return searchList
}