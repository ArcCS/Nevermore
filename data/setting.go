//Neo4j setting wrappers

package data

import (
	"fmt"
	"log"
)

// Retrieve character information
func LoadSetting(settingName string) (string, bool) {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap := conn.QueryNeoAll("MATCH (s:setting) WHERE toLower(s.name)=toLower({settingName}) RETURN {"+
		"value: s.value}",
		map[string]interface{}{
			"settingName": settingName,
		},
	)
	if len(data) < 1 {
		log.Println(rtrap)
		return "", true
	} else {
		dataR := data[0][0].(map[string]interface{})
		return dataR["value"].(string), false
	}
}

// Update setting
func UpdateSetting(settingName string, setting string) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo(
		"MATCH (a:setting) WHERE a.name={name} SET "+
			"a.value = {value}",
		map[string]interface{}{
			"name":  settingName,
			"value": setting,
		},
	)

	fmt.Println(rtrap)
	numResult, _ := result.RowsAffected()
	if numResult > 0 {
		return false
	} else {
		return true
	}
}
