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
	breathes: m.breathes,
	placement:m.placement,
	commands: m.commands,
	drops: collect({chance: d.chance, item_id: i.item_id}),
	flags:{
	fast_moving: m.fast_moving,
	guard_treasure: m.guard_treasure,
	take_treasure: m.take_treasure,
	steals: m.steals,
	block_exit: m.block_exit,
	follows: m.follows,
	no_steal: m.no_steal,
	detect_invisible: m.detect_invisible,
	no_stun: m.no_stun,
	diseases: m.diseases,
	poisons: m.poisons,
	spits_acid: m.spits_acid,
	ranged_attack: m.ranged_attack,
	flees: m.flees,
	blinds: m.blinds,
	undead: m.undead,
	day_only: m.day_only,
	night_only: m.night_only,
	hide_encounter: m.hide_encounter, 
	invisible:m.invisible, 
	permanent:m.permanent,
    immobile:m.immobile,
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
	breathes: m.breathes,
	commands: m.commands,
	placement:m.placement,
	drops: collect({chance: d.chance, item_id: i.item_id}),
	flags:{
	day_only: m.day_only,
	night_only: m.night_only,
	fast_moving: m.fast_moving,
	guard_treasure: m.guard_treasure,
	take_treasure: m.take_treasure,
	steals: m.steals,
	block_exit: m.block_exit,
	follows: m.follows,
	no_steal: m.no_steal,
	detect_invisible: m.detect_invisible,
	no_stun: m.no_stun,
	diseases: m.diseases,
	poisons: m.poisons,
	spits_acid: m.spits_acid,
	ranged_attack: m.ranged_attack,
	flees: m.flees,
	blinds: m.blinds,
	undead: m.undead,
	hide_encounter: m.hide_encounter, 
	invisible:m.invisible, 
	permanent:m.permanent,
    immobile:m.immobile,
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
		m.commands = '[]',
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
		m.undead= 0, 
		m.invisible= 0, 
		m.permanent= 0,
		m.breathes = "",
		m.fast_moving = 0,
		m.guard_treasure = 0,
		m.take_treasure = 0,
		m.steals = 0,
		m.block_exit = 0,
		m.follows = 0,
		m.no_steal = 0,
		m.detect_invisible = 0,
		m.no_stun = 0,
		m.diseases = 0,
		m.poisons = 0,
		m.spits_acid = 0,
		m.ranged_attack = 0,
		m.flees = 0,
		m.day_only = 0,
		m.night_only =0,
		m.blinds = 0,
		m.placement = 5,
		m.immobile = 0,
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
		return mob_id, false
	} else {
		return -1, true
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
		m.undead=$undead,
		m.hide_encounter=$hide_encounter, 
		m.invisible=$invisible, 
		m.permanent=$permanent,
		m.breathes=$breathes,
		m.commands=$commands,
		m.fast_moving=$fast_moving,
		m.guard_treasure=$guard_treasure,
		m.take_treasure=$take_treasure,
		m.steals=$steals,
		m.block_exit=$block_exit,
		m.follows=$follows,
		m.no_steal=$no_steal,
		m.detect_invisible=$detect_invisible,
		m.no_stun=$no_stun,
		m.diseases=$diseases,
		m.poisons=$poisons,
		m.spits_acid=$spits_acid,
		m.ranged_attack=$ranged_attack,
		m.flees=$flees,
		m.day_only=$day_only,
		m.night_only=$night_only,
		m.blinds=$blinds,
		m.placement=$placement,
		m.immobile=$immobile,
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
			"hpcur":               mobData["hpmax"],
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
			"breathes":            mobData["breathes"],
			"fast_moving":         mobData["fast_moving"],
			"guard_treasure":      mobData["guard_treasure"],
			"take_treasure":       mobData["take_treasure"],
			"steals":              mobData["steals"],
			"block_exit":          mobData["block_exit"],
			"follows":             mobData["follows"],
			"no_steal":            mobData["no_steal"],
			"detect_invisible":    mobData["detect_invisible"],
			"no_stun":             mobData["no_stun"],
			"diseases":            mobData["diseases"],
			"poisons":             mobData["poisons"],
			"spits_acid":          mobData["spits_acid"],
			"ranged_attack":       mobData["ranged_attack"],
			"flees":               mobData["flees"],
			"blinds":              mobData["blinds"],
			"night_only":          mobData["night_only"],
			"day_only":            mobData["day_only"],
			"undead":              mobData["undead"],
			"placement":           mobData["placement"],
			"immobile":            mobData["immobile"],
			"commands":            mobData["commands"],
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

func CopyMob(mobId int) (int, bool) {
	newMobId := nextId("mob")
	results, err := execWrite("MATCH (m:mob{mob_id:$mobId}) CALL apoc.refactor.cloneNodes([m]) YIELD output SET output.mob_id=$newId RETURN output.item_id",
		map[string]interface{}{
			"mobId": mobId,
			"newId": newMobId,
		},
	)
	if err != nil {
		log.Println(err)
		return 0, false
	}
	if results.Counters().ContainsUpdates() {
		return newMobId, true
	} else {
		return 0, false
	}
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
	results, err := execWrite("MATCH (r:room)-[s:spawns]->(m:mob) WHERE r.room_id=$room_id AND m.mob_id=$mob_id DELETE s",
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

func SearchMobRange(loId int, hiId int, skip int) []interface{} {
	results, err := execRead("MATCH (m:mob) WHERE m.mob_id >= $loid AND m.mob_id <= $hiid RETURN {name: m.name, mob_id: m.mob_id, level: m.level} ORDER BY m.mob_id SKIP $skip LIMIT $limit",
		map[string]interface{}{
			"loid":  loId,
			"hiid":  hiId,
			"skip":  skip,
			"limit": config.Server.SearchResults,
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
