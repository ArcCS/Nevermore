// Neo4j character wrapper

package data

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strings"
	"time"
)

// LoadChar Retrieve character information
func LoadChar(charName string) (map[string]interface{}, bool) {
	results, err := execRead("MATCH (a:character) WHERE toLower(a.name)=toLower($charName) RETURN {"+
		"gender: a.gender, "+
		"character_id: a.character_id, "+
		"name: a.name, "+
		"class: a.class, "+
		"race: a.race, "+
		"passages: a.passages, "+
		"bonuspoints: a.bonuspoints, "+
		"title: a.title, "+
		"tier: a.tier, "+
		"strcur: a.strcur, "+
		"concur: a.concur, "+
		"dexcur: a.dexcur, "+
		"piecur: a.piecur, "+
		"intcur: a.intcur, "+
		"strmod: a.strmod, "+
		"conmod: a.conmod, "+
		"dexmod: a.dexmod, "+
		"piemod: a.piemod, "+
		"intmod: a.intmod, "+
		"manacur: a.manacur, "+
		"vitcur: a.vitcur, "+
		"stamcur: a.stamcur, "+
		"manamax: a.manamax, "+
		"vitmax: a.vitmax, "+
		"stammax: a.stammax, "+
		"description: a.description, "+
		"parentid: a.parentid, "+
		"birthday: a.birthday, "+
		"birthdate: a.birthdate, "+
		"birthmonth: a.birthmonth, "+
		"played: a.played, "+
		"broadcasts: a.broadcasts, "+
		"evals: a.evals, "+
		"gold: a.gold, "+
		"bankgold: a.bankgold, "+
		"experience: a.experience, "+
		"sharpexp: a.sharpexp, "+
		"thrustexp: a.thrustexp, "+
		"bluntexp: a.bluntexp, "+
		"poleexp: a.poleexp, "+
		"missileexp: a.missileexp, "+
		"handexp: a.handexp, "+
		"fireexp: a.fireexp, "+
		"airexp: a.airexp, "+
		"earthexp: a.earthexp, "+
		"waterexp: a.waterexp, "+
		"divinity: a.divinity, "+
		"stealthexp: a.stealthexp, "+
		"spells: a.spells, "+
		"equipment: a.equipment, "+
		"inventory: a.inventory, "+
		"lastrefresh: a.lastrefresh, "+
		"effects: a.effects, "+
		"timers: a.timers, "+
		"enchants: a.enchants, "+
		"heals: a.heals, "+
		"restores: a.restores, "+
		"rerolls: a.rerolls, "+
		"oocswap: a.oocswap, "+
		"flags:{invisible: a.invisible, darkvision: a.darkvision, hidden: a.hidden, ooc: a.ooc}}",

		map[string]interface{}{
			"charName": charName,
		},
	)
	if err != nil {
		log.Println(err)
		return nil, true
	}
	if len(results) < 1 {
		return nil, false
	} else {
		return results[0].Values[0].(map[string]interface{}), false
	}
}

// CreateChar New character
func CreateChar(charData map[string]interface{}) bool {
	results, err := execWrite(
		"CREATE (a:character) SET "+
			"a.character_id = $characterId, "+
			"a.gender = $gender, "+
			"a.name = $name, "+
			"a.class = $class, "+
			"a.race = $race, "+
			"a.active = true, "+
			"a.passages = 0,"+
			"a.bonuspoints = 0,"+
			"a.title = '', "+
			"a.tier = 1, "+
			"a.strcur = $strcur, "+
			"a.concur = $concur, "+
			"a.dexcur = $dexcur, "+
			"a.piecur = $piecur, "+
			"a.intcur = $intcur, "+
			"a.strmod = 0, "+
			"a.conmod = 0, "+
			"a.dexmod = 0, "+
			"a.piemod = 0, "+
			"a.intmod = 0, "+
			"a.manacur = $curr_mana, "+
			"a.vitcur = $curr_vit, "+
			"a.stamcur = $curr_stam, "+
			"a.manamax = $curr_mana, "+
			"a.vitmax = $curr_vit, "+
			"a.stammax = $curr_stam, "+
			"a.description = '', "+
			"a.parentid = $start_room,"+
			"a.birthday = $birth_day,"+
			"a.birthdate = $birth_date,"+
			"a.birthmonth = $birth_month,"+
			"a.played = 0,"+
			"a.broadcasts = 5,"+
			"a.evals = 5,"+
			"a.gold = 0,"+
			"a.bankgold = 0,"+
			"a.sharpexp = 0,"+
			"a.bluntexp = 0,"+
			"a.poleexp = 0,"+
			"a.thrustexp = 0,"+
			"a.missileexp = 0,"+
			"a.handexp = 0,"+
			"a.fireexp = 0,"+
			"a.airexp = 0,"+
			"a.earthexp = 0,"+
			"a.waterexp = 0,"+
			"a.divinity = 0,"+
			"a.stealthexp = 0,"+
			"a.spells = '',"+
			"a.equipment = '[]',"+
			"a.inventory = '[]',"+
			"a.experience = 0, "+
			"a.invisible = 0, "+
			"a.lastrefresh = $lastrefresh, "+
			"a.darkvision = $darkvision, "+
			"a.effects = '[]', "+
			"a.timers = '[]', "+
			"a.enchants = 0, "+
			"a.heals = 0, "+
			"a.restores = 0, "+
			"a.rerolls = 0, "+
			"a.ooc = 0, "+
			"a.oocswap = 0, "+
			"a.hidden = 0 ",
		map[string]interface{}{
			"characterId": nextId("character"),
			"gender":      charData["gender"],
			"name":        utils.Title(charData["name"].(string)),
			"class":       charData["class"],
			"race":        charData["race"],
			"strcur":      charData["str"],
			"concur":      charData["con"],
			"dexcur":      charData["dex"],
			"piecur":      charData["pie"],
			"intcur":      charData["intel"],
			"birth_day":   charData["birthday"],
			"birth_date":  charData["birthdate"],
			"birth_month": charData["birthmonth"],
			"lastrefresh": time.Now().String(),
			"darkvision":  utils.Btoi(charData["darkvision"].(bool)),
			"curr_mana":   config.CalcMana(1, charData["con"].(int), charData["class"].(int)),
			"curr_vit":    config.CalcHealth(1, charData["intel"].(int), charData["class"].(int)),
			"curr_stam":   config.CalcStamina(1, charData["con"].(int), charData["class"].(int)),
			"start_room":  config.StartingRoom,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	owner, oerr := execWrite(
		"MATCH (a:account), (c:character) WHERE "+
			"a.name = $aname AND c.name = $cname "+
			"CREATE (a)-[o:owns]->(c) RETURN o",
		map[string]interface{}{
			"aname": charData["account"],
			"cname": utils.Title(charData["name"].(string)),
		},
	)
	if oerr != nil {
		log.Println(oerr)
		return false
	}
	if results.Counters().NodesCreated() > 0 && owner.Counters().RelationshipsCreated() > 0 {
		return true
	} else {
		return false
	}
}

// SaveChar Update character information from a map
func SaveChar(charData map[string]interface{}) bool {
	results, err := execWrite(
		"MATCH (a:character) WHERE a.character_id=$characterid SET "+
			"a.name = $name, "+
			"a.passages = 0,"+
			"a.bonuspoints = $bonuspoints,"+
			"a.title = $title, "+
			"a.class = $class, "+
			"a.tier = $tier,  "+
			"a.strcur = $strcur, "+
			"a.concur = $concur, "+
			"a.dexcur = $dexcur, "+
			"a.piecur = $piecur, "+
			"a.intcur = $intcur, "+
			"a.manacur = $curr_mana, "+
			"a.vitcur = $curr_vit, "+
			"a.stamcur = $curr_stam, "+
			"a.manamax = $max_mana, "+
			"a.vitmax = $max_vit, "+
			"a.stammax = $max_stam, "+
			"a.description = $description, "+
			"a.parentid = $parent_id,"+
			"a.played = $played,"+
			"a.broadcasts = $broadcasts,"+
			"a.evals = $evals,"+
			"a.gold = $gold,"+
			"a.bankgold = $bankgold,"+
			"a.sharpexp = $sharpexp,"+
			"a.bluntexp = $bluntexp,"+
			"a.poleexp = $poleexp,"+
			"a.thrustexp = $thrustexp,"+
			"a.missileexp = $missileexp,"+
			"a.handexp = $handexp, "+
			"a.fireexp = $fireexp, "+
			"a.airexp = $airexp,"+
			"a.earthexp = $earthexp,"+
			"a.waterexp = $waterexp,"+
			"a.divinity = $divinity,"+
			"a.stealthexp = $stealthexp,"+
			"a.spells = $spells,"+
			"a.equipment = $equipment,"+
			"a.inventory = $inventory,"+
			"a.effects = $effects, "+
			"a.timers = $timers, "+
			"a.lastrefresh = $lastrefresh, "+
			"a.oocswap = $oocswap, "+
			"a.enchants = $enchants, "+
			"a.heals = $heals, "+
			"a.restores = $restores, "+
			"a.rerolls = $rerolls, "+
			"a.ooc = $ooc, "+
			"a.experience = $experience",
		map[string]interface{}{
			"characterid": charData["character_id"],
			"name":        charData["name"],
			"title":       charData["title"],
			"class":       charData["class"],
			"tier":        charData["tier"],
			"experience":  charData["experience"],
			"spells":      charData["spells"],
			"thrustexp":   charData["thrustexp"],
			"bluntexp":    charData["bluntexp"],
			"missileexp":  charData["missileexp"],
			"poleexp":     charData["poleexp"],
			"sharpexp":    charData["sharpexp"],
			"handexp":     charData["handexp"],
			"fireexp":     charData["fireexp"],
			"airexp":      charData["airexp"],
			"earthexp":    charData["earthexp"],
			"waterexp":    charData["waterexp"],
			"divinity":    charData["divinity"],
			"stealthexp":  charData["stealthexp"],
			"bankgold":    charData["bankgold"],
			"gold":        charData["gold"],
			"evals":       charData["evals"],
			"broadcasts":  charData["broadcasts"],
			"played":      charData["played"],
			"description": charData["description"],
			"parent_id":   charData["parent_id"],
			"strcur":      charData["str"],
			"concur":      charData["con"],
			"dexcur":      charData["dex"],
			"piecur":      charData["pie"],
			"intcur":      charData["intel"],
			"curr_mana":   charData["manacur"],
			"curr_vit":    charData["vitcurr"],
			"curr_stam":   charData["stamcurr"],
			"max_mana":    charData["manamax"],
			"max_vit":     charData["vitmax"],
			"max_stam":    charData["stammax"],
			"equipment":   charData["equipment"],
			"inventory":   charData["inventory"],
			"effects":     charData["effects"],
			"lastrefresh": charData["lastrefresh"],
			"timers":      charData["timers"],
			"oocswap":     charData["oocswap"],
			"ooc":         charData["ooc"],
			"enchants":    charData["enchants"],
			"heals":       charData["heals"],
			"restores":    charData["restores"],
			"rerolls":     charData["rerolls"],
			"bonuspoints": charData["bonuspoints"],
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().ContainsUpdates() {
		return false
	} else {
		return true
	}
}

// SaveCharField Update character
func SaveCharField(charName string, field string, fieldVal interface{}) bool {
	results, err := execWrite(
		"MATCH (a:character) WHERE toLower(a.name)=toLower($charName) SET "+
			fmt.Sprintf("a.%[1]s", strings.ToLower(field))+" = $field_val",
		map[string]interface{}{
			"charName":  charName,
			"field_val": fieldVal,
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

func CharacterExists(charName string) bool {
	results, err := execRead("MATCH (c:character) WHERE toLower(c.name)=toLower($charName) RETURN c",
		map[string]interface{}{
			"charName": charName,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(results) > 0 {
		return true
	} else {
		return false
	}
}

// DeleteChar Delete character
func DeleteChar(charName string) bool {
	results, err := execWrite("MATCH ()-[o:owns]->(a:character) WHERE a.name=$charName DELETE o,a",
		map[string]interface{}{
			"charName": charName,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().NodesDeleted() > 0 {
		return true
	} else {
		return false
	}
}

// PuppetChar Remove current ownership and set ownership to the PuppetMaster account
func PuppetChar(charName string) bool {
	results, err := execWrite("MATCH ()-[o:owns]->(c:character) WHERE toLower(c.name)=toLower($charName) DELETE o",
		map[string]interface{}{
			"charName": charName,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().RelationshipsDeleted() > 0 {
		puppetRes, err := execWrite("MATCH (a:account), (c:character) WHERE toLower(c.name)=toLower($charName) AND a.name='PuppetMaster' CREATE (a)-[o:owns]->(c) RETURN o",
			map[string]interface{}{
				"charName": charName,
			},
		)
		if err != nil {
			log.Println(err)
			return false
		}
		if puppetRes.Counters().RelationshipsCreated() > 0 {
			return true
		}
		return true
	} else {
		return false
	}
}

func SearchCharName(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (c:character) WHERE toLower(c.name) CONTAINS toLower($search) RETURN c ORDER BY c.name SKIP $skip LIMIT $limit",
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

func SearchCharDesc(searchStr string, skip int) []interface{} {
	results, err := execRead("MATCH (c:character) WHERE toLower(c.description) CONTAINS toLower($search) RETURN c ORDER BY c.name SKIP $skip LIMIT $limit",
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
