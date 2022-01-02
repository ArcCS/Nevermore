// Neo4j account wrapper

package data

import (
	"github.com/ArcCS/Nevermore/config"
	"log"
	"strings"
)

// NewAcct Create a new account from a map of account values
func NewAcct(acctData map[string]interface{}) bool {
	results, err := execWrite("CREATE (a:account) SET "+
		"a.account_id = $acctId, "+
		"a.name = $acctName, "+
		"a.permissions = $permissions,  "+
		"a.password = $acctPass, "+
		"a.active = true",
		map[string]interface{}{
			"acctId":      nextId("account"),
			"acctName":    acctData["name"],
			"acctPass":    acctData["password"],
			"permissions": config.Server.PermissionDefault,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if results.Counters().NodesCreated() > 0 {
		return true
	}
	return false
}

// LoadAcct Retrieve account info based on an acctName
func LoadAcct(acctName string) (map[string]interface{}, bool) {
	results, err := execRead("MATCH (a:account) WHERE toLower(a.name)=toLower($acctName) RETURN "+
		"{account_id: a.account_id, name: a.name, permissions: a.permissions, password: a.password}",
		map[string]interface{}{
			"acctName": acctName,
		},
	)
	if err != nil {
		log.Println(err)
		return nil, true
	}
	if len(results) > 0 {
		return results[0].Values[0].(map[string]interface{}), false
	} else {
		return nil, true
	}
}

// AccountExists Checks whether an account exists based on name
func AccountExists(acctName string) bool {
	results, err := execRead("MATCH (a:account) WHERE toLower(a.name)=toLower($acctName) RETURN a",
		map[string]interface{}{
			"acctName": acctName,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(results) < 1 {
		return false
	} else {
		return true
	}
}

// ListChars Retrieve account characters based on account name
func ListChars(acctName string) []string {
	results, err := execRead("MATCH (a:account)-[o:owns]->(c:character) WHERE toLower(a.name)=toLower($acctName) and "+
		"c.class<>100 and c.active=true RETURN {name: c.name}",
		map[string]interface{}{
			"acctName": acctName,
		},
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(results) < 1 {
		return nil
	} else {
		searchList := make([]string, 0)
		for _, row := range results {
			datum := row.Values[0].(map[string]interface{})["name"].(string)
			if strings.TrimSpace(datum) != "" {
				searchList = append(searchList, datum)
			}
		}
		return searchList
	}
}

// ListPowerChar Retrieve account power characters based on account name
func ListPowerChar(acctName string) (string, bool) {
	results, err := execRead("MATCH (a:account)-[o:owns]->(d:character) WHERE toLower(a.name)=toLower($acctName) and d.class=100 RETURN "+
		"{char: d.name}",
		map[string]interface{}{
			"acctName": acctName,
		},
	)
	if err != nil {
		log.Println(err)
		return "", true
	}
	if len(results) < 1 {
		return "", true
	} else {
		datum := results[0].Values[0].(map[string]interface{})
		return datum["char"].(string), false
	}
}

// Deactivate account based on account name
func Deactivate(acctName string) bool {
	results, err := execWrite("MATCH (a:account) WHERE a.name=$acctName SET a.active=false",
		map[string]interface{}{
			"acctName": acctName,
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

// UpdatePassword Update account
func UpdatePassword(acctName string, acctPass string) bool {
	results, err := execWrite("MATCH (a:account) WHERE toLower(a.name)=toLower($acctName) SET "+
		"a.password = $acctPass ",
		map[string]interface{}{
			"acctName": acctName,
			"acctPass": acctPass,
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

func TogglePermission(acctName string, permission uint32) bool {
	results, err := execWrite("MATCH (a:account) WHERE toLower(a.name)=toLower($acctName) SET "+
		"a.permissions = apoc.bitwise.op(a.permissions,'^',$permission)",
		map[string]interface{}{
			"acctName":   acctName,
			"permission": permission,
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

// DeleteAcct Delete account based on account name
func DeleteAcct(acctName string) bool {
	results, err := execWrite("MATCH (a:account) WHERE a.name=$acctName DELETE a",
		map[string]interface{}{
			"acctName": acctName,
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

func SearchAccountName(searchStr string, skip int) []interface{} {
	data, err := execRead("MATCH (a:account) WHERE toLower(a.name) CONTAINS toLower($search) ORDER BY a.name LIMIT 15 SKIP $skip ",
		map[string]interface{}{
			"search": searchStr,
			"skip":   skip,
		},
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	searchList := make([]interface{}, len(data))
	for _, row := range data {
		searchList = append(searchList, row.Values[0].(map[string]interface{}))
	}
	return searchList
}
