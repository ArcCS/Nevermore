package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadMobs() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	results, err := execRead("MATCH (m:mob) OPTIONAL MATCH (m)-[d:drops]->(i:item) RETURN "+
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
	permanent:m.permanent,
	hostile:m.hostile}}`, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	mobList := make([]interface{}, len(results))
	for _, row := range results {
		mobList = append(mobList, row.Values[0].(map[string]interface{}))
	}
	return mobList
}

func LoadMob(mobId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	results, err := execRead("MATCH (m:mob {mob_id: $mobId}) OPTIONAL MATCH (m)-[d:drops]->(i:item) RETURN "+
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
	permanent:m.permanent,
	hostile:m.hostile }}`,
		map[string]interface{}{
			"mobId": mobId,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	return results[0].Values[0].(map[string]interface{})
}

func CreateMob(mobName string, creator string) (int, bool) {
	mob_id := nextId("mob")
	results, err := execWrite(
		"CREATE (m:mob) SET "+
			`m.mob_id = $mobId, 
		m.name = $name, 
		m.creator = $creator, 
		m.description=  "A shiny new mob!", 
		m.experience=  0, 
		m.level= 1, 
		m.gold= 0, 
		m.constitution= 0, 
		m.strength= 0, 
		m.intelligence= 0, 
		m.dexterity= 0, 
		m.piety= 0, 
		m.mpmax= 0, 
		m.mpcur= 0, 
		m.hpcur= 0, 
		m.hpmax= 0, 
		m.sdice= 1, 
		m.ndice= 1, 
		m.pdice= 0, 
		m.spells= "", 
		m.casting_probability= 0, 
		m.armor= 0, 
		m.numwander= 20, 
		m.wimpyvalue= 0, 
		m.air_resistance= 0, 
		m.fire_resistance= 0, 
		m.earth_resistance= 0, 
		m.water_resistance= 0, 
		m.hide_encounter=  0, 
		m.invisible= 0, 
		m.permanent= 0,
		m.hostile=0`,
		map[string]interface{}{
			"mobId":   mob_id,
			"name":    mobName,
			"creator": creator,
		},
	)
	if err != nil {
		log.Println(err)
		return -1, false
	}
	if results.Counters().ContainsUpdates() {
		return mob_id, true
	} else {
		return -1, false
	}
}

func UpdateMob(mobData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (m:mob) WHERE m.mob_id=$mob_id SET "+
			`m.name=$name,
		m.description=$description,
		m.experience=$experience, 
		m.level=$level, 
		m.gold=$gold, 
		m.constitution=$constitution, 
		m.strength=$strength, 
		m.intelligence=$intelligence, 
		m.dexterity=$dexterity, 
		m.piety=$piety, 
		m.mpmax=$mpmax, 
		m.mpcur=$mpcur, 
		m.hpcur=$hpcur, 
		m.hpmax=$hpmax, 
		m.sdice=$sdice, 
		m.ndice=$ndice, 
		m.pdice=$pdice, 
		m.spells=$spells, 
		m.casting_probability=$casting_probability, 
		m.armor=$armor, 
		m.numwander=$numwander, 
		m.wimpyvalue=$wimpyvalue, 
		m.air_resistance=$air_resistance, 
		m.fire_resistance=$fire_resistance, 
		m.earth_resistance=$earth_resistance, 
		m.water_resistance=$water_resistance, 
		m.hide_encounter=$hide_encounter, 
		m.invisible=$invisible, 
		m.permanent=$permanent,
		m.hostile=$hostile`,
		map[string]interface{}{
			"mob_id":              mobData["mob_id"],
			"name":                mobData["name"],
			"description":         mobData["description"],
			"experience":          mobData["experience"],
			"level":               mobData["level"],
			"gold":                mobData["gold"],
			"constitution":        mobData["constitution"],
			"strength":            mobData["strength"],
			"intelligence":        mobData["intelligence"],
			"dexterity":           mobData["dexterity"],
			"piety":               mobData["piety"],
			"mpmax":               mobData["mpmax"],
			"mpcur":               mobData["mpcur"],
			"hpcur":               mobData["hpcur"],
			"hpmax":               mobData["hpmax"],
			"sdice":               mobData["sdice"],
			"ndice":               mobData["ndice"],
			"pdice":               mobData["pdice"],
			"spells":              mobData["spells"],
			"casting_probability": mobData["casting_probability"],
			"armor":               mobData["armor"],
			"numwander":           mobData["numwander"],
			"wimpyvalue":          mobData["wimpyvalue"],
			"air_resistance":      mobData["air_resistance"],
			"fire_resistance":     mobData["fire_resistance"],
			"earth_resistance":    mobData["earth_resistance"],
			"water_resistance":    mobData["water_resistance"],
			"hide_encounter":      mobData["hide_encounter"],
			"invisible":           mobData["invisible"],
			"permanent":           mobData["permanent"],
			"hostile":             mobData["hostile"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	} else {
		return false
	}
}

func CreateEncounter(encounterData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room), (m:mob) WHERE "+
			"r.room_id = $roomId AND m.mob_id = $mobId "+
			`CREATE (r)-[s:spawns]->(m) SET 
        s.chance=$chance`,
		map[string]interface{}{
			"mobId":  encounterData["mobId"],
			"roomId": encounterData["roomId"],
			"chance": encounterData["chance"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	} else {
		return false
	}
}

func SumEncounters(roomId int) int {
	results, err := execRead("MATCH (r:room)-[s:spawns]->() WHERE r.room_id=$room_id RETURN {rate_sum: sum(s.chance)}",
		map[string]interface{}{
			"room_id": roomId,
		},
	)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int(results[0].Values[0].(map[string]interface{})["rate_sum"].(int64))
}

func DeleteMob(mobId int) bool {
	results, err := execWrite("MATCH ()-[s:spawns]->(m:mob)-[d:drops]->() WHERE m.mob_id=$mob_id DELETE s, m, d",
		map[string]interface{}{
			"room_id": mobId,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	} else {
		return false
	}
}

func UpdateEncounter(mobData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (r:room)-[s:spawns]->(m:mob) WHERE "+
			"r.room_id=$roomId AND m.mob_id=$mobId SET "+
			"s.chance=$chance",
		map[string]interface{}{
			"roomId": mobData["roomId"],
			"mobId":  mobData["mobId"],
			"chance": mobData["chance"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	} else {
		return false
	}
}

func DeleteEncounter(mobId int, roomId int) bool {
	results, err := execWrite("MATCH (r:room)-[s:spawns]->(m:mob) WHERE r.room_id=$room_id AND s.mob_id=$mob_id DELETE s",
		map[string]interface{}{
			"room_id": roomId,
			"mob_id":  mobId,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return true
	} else {
		return false
	}
}

func SearchMobName(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (m:mob) WHERE toLower(m.name) CONTAINS toLower($search) RETURN {name: m.name, mob_id: m.mob_id, level: m.level} ORDER BY m.name SKIP $skip LIMIT $limit",
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
	searchList := make([]interface{}, len(results))
	for _, row := range results {
		searchList = append(searchList, row.Values[0].(map[string]interface{}))
	}
	return searchList
}

func SearchMobDesc(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (m:mob) WHERE toLower(m.description) CONTAINS toLower($search) RETURN {name: m.name, mob_id: m.mob_id, level: m.level} ORDER BY m.name SKIP $skip LIMIT $limit",
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
	searchList := make([]interface{}, len(results))
	for _, row := range results {
		searchList = append(searchList, row.Values[0].(map[string]interface{}))
	}
	return searchList
}
