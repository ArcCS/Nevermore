// Neo4j character wrapper

package data

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"log"
	"strings"
)

// Retrieve character information
func LoadChar(charName string) (map[string]interface{}, bool) {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (a:character) WHERE toLower(a.name)=toLower({charName}) RETURN {" +
		"gender: a.gender, " +
		"character_id: a.character_id, " +
		"name: a.name, " +
		"class: a.class, " +
		"race: a.race, " +
		"passages: a.passages, " +
		"bonuspoints: a.bonuspoints, " +
		"title: a.title, " +
		"tier: a.tier, " +
		"strcur: a.strcur, " +
		"concur: a.concur, " +
		"dexcur: a.dexcur, " +
		"piecur: a.piecur, " +
		"intcur: a.intcur, " +
		"strmod: a.strmod, " +
		"conmod: a.conmod, " +
		"dexmod: a.dexmod, " +
		"piemod: a.piemod, " +
		"intmod: a.intmod, " +
		"manacur: a.manacur, " +
		"vitcur: a.vitcur, " +
		"stamcur: a.stamcur, " +
		"manamax: a.manamax, " +
		"vitmax: a.vitmax, " +
		"stammax: a.stammax, " +
		"description: a.description, " +
		"parentid: a.parentid, " +
		"birthday: a.birthday, " +
		"played: a.played, " +
		"broadcasts: a.broadcasts, " +
		"evals: a.evals, " +
		"gold: a.gold, " +
		"bankgold: a.bankgold, " +
		"experience: a.experience, " +
		"sharpexp: a.sharpexp, " +
		"thrustexp: a.thrustexp, " +
		"bluntexp: a.bluntexp, " +
		"poleexp: a.poleexp, " +
		"missileexp: a.missileexp, " +
		"spells: a.spells, " +
		"equipment: a.equipment, " +
		"inventory: a.inventory, " +
		"flags:{invisible: a.invisible, darkvision: a.darkvision, hidden: a.hidden}}",

		map[string]interface {}{
			"charName": charName,
		},
	)
	if len(data) < 1 {
		log.Println(rtrap)
		return nil, true
	}else {
		return data[0][0].(map[string]interface{}), false
	}
}

// New character
func CreateChar(charData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"CREATE (a:character) SET " +
		"a.character_id = {characterId}, " +
		"a.gender = {gender}, " +
		"a.name = {name}, " +
		"a.class = {class}, " +
		"a.race = {race}, " +
		"a.active = true, " +
		"a.passages = 0," +
		"a.bonuspoints = 0," +
		"a.title = '', " +
		"a.tier = 1, " +
		"a.strcur = {strcur}, " +
		"a.concur = {concur}, " +
		"a.dexcur = {dexcur}, " +
		"a.piecur = {piecur}, " +
		"a.intcur = {intcur}, " +
		"a.strmod = 0, " +
		"a.conmod = 0, " +
		"a.dexmod = 0, " +
		"a.piemod = 0, " +
		"a.intmod = 0, " +
		"a.manacur = {curr_mana}, " +
		"a.vitcur = {curr_vit}, " +
		"a.stamcur = {curr_stam}, " +
		"a.manamax = {curr_mana}, " +
		"a.vitmax = {curr_vit}, " +
		"a.stammax = {curr_stam}, " +
		"a.description = '', " +
		"a.parentid = {start_room}," +
		"a.birthday = 0," +
		"a.played = 0," +
		"a.broadcasts = 5," +
		"a.evals = 5," +
		"a.gold = 0," +
		"a.bankgold = 0," +
		"a.sharpexp = 0," +
		"a.bluntexp = 0," +
		"a.poleexp = 0," +
		"a.thrustexp = 0," +
		"a.missileexp = 0," +
		"a.spells = ''," +
		"a.equipment = '[]'," +
		"a.inventory = '[]'," +
		"a.experience = 0, " +
		"a.invisible = 0, " +
		"a.darkvision = 0, " +
		"a.hidden = 0 ",
		map[string]interface {}{
			"characterId": nextId("character"),
			"gender":      	charData["gender"],
			"name":        	strings.Title(charData["name"].(string)),
			"class":       	charData["class"],
			"race":        	charData["race"],
			"strcur":       charData["str"],
			"concur":       charData["con"],
			"dexcur":       charData["dex"],
			"piecur":       charData["pie"],
			"intcur":       charData["intel"],
			"curr_mana":	config.CalcMana(1, charData["con"].(int), charData["class"].(int)),
			"curr_vit":		config.CalcHealth(1, charData["intel"].(int), charData["class"].(int)),
			"curr_stam":	config.CalcStamina(1, charData["con"].(int), charData["class"].(int)),
			"start_room": 	config.StartingRoom,
		},
	)

	owner, otrap := conn.ExecNeo(
		"MATCH (a:account), (c:character) WHERE " +
			"a.name = {aname} AND c.name = {cname} " +
			"CREATE (a)-[o:owns]->(c) RETURN o",
		map[string]interface {}{
			"aname":        charData["account"],
			"cname":       	charData["name"],
		},
	)

	fmt.Println(rtrap)
	fmt.Println(otrap)
	ownResult, _ := owner.RowsAffected()
	numResult, _ := result.RowsAffected()
	if numResult > 0 && ownResult > 0 {
		return false
	}else {
		return true
	}
}

// Update character
func SaveChar(charData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"MATCH (a:character) WHERE a.character_id={characterid} SET " +
			"a.name = {name}, " +
			"a.passages = 0," +
			"a.bonuspoints = 0," +
			"a.title = {title}, " +
			"a.tier = {tier},  " +
			"a.strcur = {strcur}, " +
			"a.concur = {concur}, " +
			"a.dexcur = {dexcur}, " +
			"a.piecur = {piecur}, " +
			"a.intcur = {intcur}, " +
			"a.manacur = {curr_mana}, " +
			"a.vitcur = {curr_vit}, " +
			"a.stamcur = {curr_stam}, " +
			"a.manamax = {max_mana}, " +
			"a.vitmax = {max_vit}, " +
			"a.stammax = {max_stam}, " +
			"a.description = {description}, " +
			"a.parentid = {parent_id}," +
			"a.played = {played}," +
			"a.broadcasts = {broadcasts}," +
			"a.evals = {evals}," +
			"a.gold = {gold}," +
			"a.bankgold = {bankgold}," +
			"a.sharpexp = {sharpexp}," +
			"a.bluntexp = {bluntexp}," +
			"a.poleexp = {poleexp}," +
			"a.thrustexp = {thrustexp}," +
			"a.missileexp = {missileexp}," +
			"a.spells = {spells}," +
			"a.equipment = {equipment}," +
			"a.inventory = {inventory}," +
			"a.experience = {experience}",
		map[string]interface {}{
			"characterid": charData["character_id"],
			"name": charData["name"],
			"title": charData["title"],
			"tier": charData["tier"],
			"experience":	charData["experience"],
			"spells":		charData["spells"],
			"thrustexp":     charData["thrustexp"],
			"bluntexp":     charData["bluntexp"],
			"missileexp":     charData["missileexp"],
			"poleexp":     charData["poleexp"],
			"sharpexp":     charData["sharpexp"],
			"bankgold":     charData["bankgold"],
			"gold":			charData["gold"],
			"evals":		charData["evals"],
			"broadcasts":	charData["broadcasts"],
			"played":      	charData["played"],
			"description":  charData["description"],
			"parent_id":    charData["parent_id"],
			"strcur":       charData["str"],
			"concur":       charData["con"],
			"dexcur":       charData["dex"],
			"piecur":       charData["pie"],
			"intcur":       charData["intel"],
			"curr_mana":	charData["manacur"],
			"curr_vit":		charData["vitcurr"],
			"curr_stam":	charData["stamcurr"],
			"max_mana":	charData["manamax"],
			"max_vit":		charData["vitmax"],
			"max_stam":	charData["stammax"],
			"equipment": charData["equipment"],
			"inventory": charData["inventory"],
		},
	)


	fmt.Println(rtrap)
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

func CharacterExists(charName string) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (c:character) WHERE toLower(c.name)=toLower({charName}) RETURN c",
		map[string]interface {}{
			"charName": charName,
		},
	)
	if len(data) < 1 {
		return false
	}else {
		return true
	}
}

// Delete character
func DeleteChar(charName string) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, _ := conn.ExecNeo("MATCH (a:character) WHERE a.name={charName} DELETE a",
		map[string]interface {}{
			"charName": charName,
		},
	)
	numResult, _ := result.RowsAffected()
	if numResult < 0 {
		return true
	}else {
		return false
	}
}

func SearchCharName(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (c:character) WHERE toLower(c.name) CONTAINS toLower({search}) RETURN c ORDER BY c.name SKIP {skip} LIMIT {limit}",
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

func SearchCharDesc(searchStr string, skip int) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (c:character) WHERE toLower(c.description) CONTAINS toLower({search}) RETURN c ORDER BY c.name SKIP {skip} LIMIT {limit}",
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

/*
// Get an inventory_id
func NextInventoryId(characterName string) int {
conn, _ := getConn()
defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (c:character)-[h:has]->(item)"+
	"WHERE c.name={characterName}" +
"WITH COLLECT(h.has_id) as has_ids" +
"WITH MAX(has_ids)+1 AS new_id" +
"RETURN new_id",
		map[string]interface {}{
			"characterName":  characterName,
		},
	)
return data[0][0].(int)
}

// Create Inventory Item
func CreateInventory(inventoryData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	hasInventory, rtrap := conn.ExecNeo(
		"MATCH (c:character), (i:item) WHERE " +
			"c.character_id = {characterName} AND i.item_id = {itemId} " +
			`CREATE (c)-[h:has]->(i) SET
				h.has_id={hasId},
				h.name={itemName},
				h.uses={uses},
				h.magic={magic},
				h.spell={spell},
				h.armor={armor},
				h.equipped={equipped}`,
		map[string]interface {}{
			"characterName":        inventoryData["characterName"],
			"itemId":       inventoryData["itemId"],
			"hasId":		NextInventoryId(inventoryData["characterName"].(string)),
			"itemName":		inventoryData["itemName"],
			"uses": 		inventoryData["uses"],
			"magic":		inventoryData["magic"],
			"spell":		inventoryData["spell"],
			"armor":		inventoryData["armor"],
			"equipped":		inventoryData["equipped"],
		},
	)
	if rtrap != nil{
		log.Println(rtrap)
	}

	numResult, _ := hasInventory.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

// Delete inventory
func DeleteInventoryItem(characterName string, hasId int) bool {
	conn, _ := getConn()
	defer conn.Close()
	toExit, rtrap := conn.ExecNeo(
		"MATCH (c:character)-[h:has_id]->(item) WHERE " +
			"c.name = {characterName} AND h.has_id = {hasId} " +
			`DELETE (c)-[s:spawns]->(m) SET 
	s.chance={chance}`,
		map[string]interface {}{
			"hasId":        hasId,
			"characterName":       characterName,
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

func ClearAllInventory(characterName string){

}

// Create Inventory Item
func UpdateInventory(inventoryData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	hasInventory, rtrap := conn.ExecNeo(
		"MATCH (c:character)-[h:has]->(i:item) WHERE " +
			"c.character_id = {characterName} AND i.item_id = {itemId} and h.has_id =  {hasId} " +
			`CREATE (c)-[h:has]->(i) SET
				h.has_id={hasId},
				h.name={itemName},
				h.uses={uses},
				h.magic={magic},
				h.spell={spell},
				h.armor={armor},
				h.equipped={equipped}`,
		map[string]interface {}{
			"characterName":        inventoryData["characterName"],
			"itemId":       inventoryData["itemId"],
			"hasId":		inventoryData["hasId"],
			"itemName":		inventoryData["itemName"],
			"uses": 		inventoryData["uses"],
			"magic":		inventoryData["magic"],
			"spell":		inventoryData["spell"],
			"armor":		inventoryData["armor"],
			"equipped":		inventoryData["equipped"],
		},
	)
	if rtrap != nil{
		log.Println(rtrap)
	}

	numResult, _ := hasInventory.RowsAffected()
	if numResult > 0 {
		return false
	}else {
		return true
	}
}

 */