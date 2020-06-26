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
		"attrmoves: a.attrmoves, " +
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

/*
	Object{
		Name:        charData["name"].(string),
		Description: charData["description"].(string),
		Placement:   3,
	},f
	Equipment{},
	ItemInventory{},


 */
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
		"a.attrmoves = 0," +
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
			"a.attrmoves = 0," +
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
