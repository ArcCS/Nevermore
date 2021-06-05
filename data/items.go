package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadItems() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	results, err := execRead("MATCH (i:item) RETURN "+
		`{creator:i.creator,
	item_id:i.item_id,
	ndice:i.ndice,
	weight:i.weight,
	description:i.description,
	weapon_speed:i.weapon_speed,
	type:i.type,
	pdice:i.pdice,
	armor:i.armor,
	max_uses:i.max_uses,
	name:i.name,
	sdice:i.sdice,
	value:i.value,
	spell:i.spell,
	commands: i.commands,
	flags: {always_crit: i.always_crit, permanent:i.permanent,
	magic:i.magic,
	no_take: i.no_take,
	light: i.light,
	weightless_chest: i.weightless_chest}
	}`, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	itemList := make([]interface{}, len(results))
	for _, row := range results {
		itemList = append(itemList, row.Values[0].(map[string]interface{}))
	}
	return itemList
}

func LoadItem(itemId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	results, err := execRead("MATCH (i:item) WHERE i.item_id=$itemId RETURN "+
		`{creator:i.creator,
	item_id:i.item_id,
	ndice:i.ndice,
	weight:i.weight,
	description:i.description,
	weapon_speed:i.weapon_speed,
	type:i.type,
	pdice:i.pdice,
	armor:i.armor,
	max_uses:i.max_uses,
	name:i.name,
	sdice:i.sdice,
	value:i.value,
	spell:i.spell,
	commands: i.commands,
	flags: {always_crit: i.always_crit,permanent:i.permanent,
	magic:i.magic,
	no_take: i.no_take,
	light: i.light,
	weightless_chest: i.weightless_chest}
	}`,
		map[string]interface{}{
			"itemId": itemId,
		})
	if err != nil {
		log.Println(err)
		return nil
	}
	return results[0].Values[0].(map[string]interface{})
}

func CreateItem(itemData map[string]interface{}) (int, bool) {
	item_id := nextId("item")
	results, err := execWrite(
		"CREATE (i:item) SET "+
			`i.creator = $creator,
		i.item_id = $item_id,
		i.ndice = 1,
		i.weight = 1,
		i.description = "Your new shiny item!",
		i.weapon_speed = 0,
		i.type = $type,
		i.pdice = 1,
		i.armor = 0,
		i.max_uses = 1,
		i.commands = '[]',
		i.name = $name,
		i.sdice = 1,
		i.spell = "",
		i.value = 1,
		i.always_crit = 0,
		i.permanent = 0,
		i.magic = 0,
		i.no_take = 0,
		i.light = 0,
		i.weightless_chest = 0`,
		map[string]interface{}{
			"item_id": item_id,
			"name":    itemData["name"],
			"creator": itemData["creator"],
			"type":    itemData["type"],
		},
	)
	if err != nil {
		log.Println(err)
		return -1, true
	}
	if results.Counters().ContainsUpdates() {
		return item_id, false
	} else {
		return -1, true
	}
}

func UpdateItem(itemData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (i:item) WHERE i.item_id=$item_id SET "+
			`i.ndice = $ndice,
		i.weight = $weight,
		i.description = $description,
		i.weapon_speed = $weapon_speed,
		i.type = $type,
		i.pdice = $pdice,
		i.armor = $armor,
		i.max_uses = $max_uses,
		i.name = $name,
		i.sdice = $sdice,
		i.value = $value,
		i.spell = $spell,
		i.commands = $commands,
		i.always_crit = $always_crit,
		i.permanent = $permanent,
		i.no_take = $no_take,
		i.light = $light,
		i.weightless_chest = $weightless_chest,
		i.magic = $magic`,
		map[string]interface{}{
			"item_id":          itemData["item_id"],
			"ndice":            itemData["ndice"],
			"weight":           itemData["weight"],
			"description":      itemData["description"],
			"weapon_speed":     itemData["weapon_speed"],
			"type":             itemData["type"],
			"pdice":            itemData["pdice"],
			"armor":            itemData["armor"],
			"max_uses":         itemData["max_uses"],
			"name":             itemData["name"],
			"sdice":            itemData["sdice"],
			"value":            itemData["value"],
			"spell":            itemData["spell"],
			"always_crit":      itemData["always_crit"],
			"permanent":        itemData["permanent"],
			"magic":            itemData["magic"],
			"light":            itemData["light"],
			"no_take":          itemData["no_take"],
			"weightless_chest": itemData["weightless_chest"],
			"commands": 		itemData["commands"],
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

func DeleteItem(roomId int) bool {
	results, err := execWrite("MATCH ()-[e:exit]->(r:room)-[e2:exit]->() WHERE r.room_id=$room_id DELETE r, e, e2",
		map[string]interface{}{
			"room_id": roomId,
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

func CreateDrop(dropData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (m:mob), (i:item) WHERE "+
			"i.item_id = $itemId AND m.mob_id = $mobId "+
			`CREATE (m)-[d:drops]->(i) SET 
	d.chance=$chance`,
		map[string]interface{}{
			"mobId":  dropData["mobId"],
			"itemId": dropData["itemId"],
			"chance": dropData["chance"],
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

func UpdateDrop(mobData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (m:mob)-[d:drops]->(i:item) WHERE "+
			"m.mob_id=$mob_id AND i.item_id=$item_id SET "+
			`d.chance=$chance`,
		map[string]interface{}{
			"item_id": mobData["item_id"],
			"mob_id":  mobData["mob_id"],
			"chance":  mobData["chance"],
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

func DeleteDrop(mobId int, itemId int) bool {
	results, err := execWrite("MATCH (m:mob)-[d:drops]->(i:item) WHERE m.mob_id=$mob_id AND i.item_id=$item_id DELETE d",
		map[string]interface{}{
			"item_id": itemId,
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

func SearchItemName(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (o:item) WHERE toLower(o.name) CONTAINS toLower($search) RETURN {name:o.name, type:o.type, item_id: o.item_id} ORDER BY o.name  SKIP $skip LIMIT $limit",
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

func SearchItemDesc(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (o:item) WHERE toLower(o.description) CONTAINS toLower($search) RETURN {name:o.name, type:o.type, item_id: o.item_id} ORDER BY o.name  SKIP $skip LIMIT $limit",
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

func SearchItemMaxDamage(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (i:item) WHERE i.ndice*i.sdice+i.pdice <= toInteger($search) and i.type <=4 RETURN {name:i.name, type:i.type, item_id: i.item_id, max_damage: i.ndice*i.sdice+i.pdice}  SKIP $skip LIMIT $limit",
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

func SearchItemRange(loId int, hiId int, skip int) []interface{} {
	results, err := execRead("MATCH (i:item) WHERE i.item_id >= $loid AND i.item_id <= $hiid RETURN {name: i.name, item_id: i.item_id} ORDER BY m.name SKIP $skip LIMIT $limit",
		map[string]interface{}{
			"loid": loId,
			"hiId": hiId,
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
