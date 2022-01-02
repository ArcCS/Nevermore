package data

import (
	"log"
)

func LoadSetting(settingName string) (string, bool) {
	results, err := execRead("MATCH (s:setting) WHERE toLower(s.name)=toLower($settingName) RETURN {"+
		"value: s.value}",
		map[string]interface{}{
			"settingName": settingName,
		},
	)
	if err != nil {
		log.Println(err)
		return "", false
	}
	return results[0].Values[0].(map[string]interface{})["value"].(string), false
}

func UpdateSetting(settingName string, setting string) bool {
	results, err := execWrite(
		"MATCH (a:setting) WHERE a.name=$name SET "+
			"a.value = $value",
		map[string]interface{}{
			"name":  settingName,
			"value": setting,
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
