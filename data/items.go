package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
)

func LoadItems() []interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (i:item) RETURN "+
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
	flags: {always_crit: i.always_crit, permanent:i.permanent,
	magic:i.magic,
	no_take: i.no_take,
	light: i.light,
	weightless_chest: i.weightless_chest}
	}`, nil)
	if rtrap != nil {
		log.Println(rtrap)
		return nil
	}
	itemList := make([]interface{}, len(data))
	for _, row := range data {
		datum := row[0].(map[string]interface{})
		itemList = append(itemList, datum)
	}
	return itemList
}

func LoadItem(itemId int) map[string]interface{} {
	// Return all of the rooms to be pushed into the room stack
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (i:item) WHERE i.item_id={itemId} RETURN "+
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
	flags: {always_crit: i.always_crit,permanent:i.permanent,
	magic:i.magic,
	no_take: i.no_take,
	light: i.light,
	weightless_chest: i.weightless_chest}
	}`,
		map[string]interface{}{
			"itemId": itemId,
		})
	if rtrap != nil {
		log.Println(rtrap)
		return nil
	}
	return data[0][0].(map[string]interface{})
}

// Create Room
func CreateItem(itemData map[string]interface{}) (int, bool) {
	conn, _ := getConn()
	defer conn.Close()
	item_id := nextId("item")
	result, rtrap := conn.ExecNeo(
		"CREATE (i:item) SET "+
			`i.creator = {creator},
		i.item_id = {item_id},
		i.ndice = 1,
		i.weight = 1,
		i.description = "Your new shiny item!",
		i.weapon_speed = 0,
		i.type = {type},
		i.pdice = 1,
		i.armor = 0,
		i.max_uses = 1,
		i.name = {name},
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

	if rtrap != nil {
		log.Println(rtrap)
	}
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return item_id, false
	} else {
		return -1, true
	}
}

// Update Room
func UpdateItem(itemData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"MATCH (i:item) WHERE i.item_id={item_id} SET "+
			`i.ndice = {ndice},
		i.weight = {weight},
		i.description = {description},
		i.weapon_speed = {weapon_speed},
		i.type = {type},
		i.pdice = {pdice},
		i.armor = {armor},
		i.max_uses = {max_uses},
		i.name = {name},
		i.sdice = {sdice},
		i.value = {value},
		i.spell = {spell},
		i.always_crit = {always_crit},
		i.permanent = {permanent},
		i.no_take = {no_take},
		i.light = {light},
		i.weightless_chest = {weightless_chest},
		i.magic = {magic}`,
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
		},
	)

	if rtrap != nil {
		log.Println(rtrap)
	}
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return false
	} else {
		return true
	}
}

// Delete Item
func DeleteItem(roomId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _ := conn.ExecNeo("MATCH ()-[e:exit]->(r:room)-[e2:exit]->() WHERE r.room_id={room_id} DELETE r, e, e2",
		map[string]interface{}{
			"room_id": roomId,
		},
	)

	numResult, _ := data.RowsAffected()
	if numResult < 1 {
		return false
	} else {
		return true
	}
}

// Create Drop
func CreateDrop(dropData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, rtrap := conn.ExecNeo(
		"MATCH (m:mob), (i:item) WHERE "+
			"i.item_id = {itemId} AND m.mob_id = {mobId} "+
			`CREATE (m)-[d:drops]->(i) SET 
	s.chance={chance}`,
		map[string]interface{}{
			"mobId":  dropData["mobId"],
			"itemId": dropData["itemId"],
			"chance": dropData["chance"],
		},
	)
	if rtrap != nil {
		log.Println(rtrap)
	}

	numResult, _ := toExit.RowsAffected()
	if numResult > 0 {
		return false
	} else {
		return true
	}
}

// Update Drop
func UpdateDrop(mobData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, etrap := conn.ExecNeo(
		"MATCH (m:mob)-[d:drops]->(i:item) WHERE "+
			"m.mob_id={mob_id} AND i.item_id={item_id} SET "+
			`s.chance={chance}`,
		map[string]interface{}{
			"item_id": mobData["item_id"],
			"mob_id":  mobData["mob_id"],
			"chance":  mobData["chance"],
		},
	)
	if etrap != nil {
		log.Println(etrap)
	}
	numResult, _ := toExit.RowsAffected()
	if numResult > 0 {
		return false
	} else {
		return true
	}
}

// Delete Drop
func DeleteDrop(mobId int, itemId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, rtrap := conn.ExecNeo("MATCH (m:mob)-[d:drops]->(i:item) WHERE m.mob_id={mob_id} AND i.item_id={item_id} DELETE d",
		map[string]interface{}{
			"item_id": itemId,
			"mob_id":  mobId,
		},
	)
	if rtrap != nil {
		log.Println(rtrap)
	}
	numResult, _ := data.RowsAffected()
	if numResult < 1 {
		return false
	} else {
		return true
	}
}

func SearchItemName(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (o:item) WHERE toLower(o.name) CONTAINS toLower({search}) RETURN {name:o.name, type:o.type, item_id: o.item_id} ORDER BY o.name  SKIP {skip} LIMIT {limit}",
		map[string]interface{}{
			"search": searchStr,
			"skip":   skip,
			"limit":  config.Server.SearchResults,
		},
	)

	if rtrap != nil {
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

func SearchItemDesc(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (o:item) WHERE toLower(o.description) CONTAINS toLower({search}) RETURN {name:o.name, type:o.type, item_id: o.item_id} ORDER BY o.name  SKIP {skip} LIMIT {limit}",
		map[string]interface{}{
			"search": searchStr,
			"skip":   skip,
			"limit":  config.Server.SearchResults,
		},
	)

	if rtrap != nil {
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
