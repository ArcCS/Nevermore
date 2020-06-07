// Neo4j account wrapper

package data

import (
	"fmt"
	"log"
)

// New account
func NewAcct(acctData map[string]interface{}) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, _ := conn.ExecNeo("CREATE (a:account) SET " +
		"a.account_id = {acctId}, " +
		"a.name = {acctName}, " +
		"a.permissions = 3,  " +
		"a.password = {acctPass}, " +
		"a.active = true",
		map[string]interface{}{
			"acctId": nextId("account"),
			"acctName": acctData["name"],
			"acctPass": acctData["password"],
		},
	)
	numResult, _ := result.RowsAffected()
	if numResult < 0 {
		return true
	}else {
		return false
	}
}

// Retrieve account info
func LoadAcct(acctName string) (map[string]interface{}, bool) {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (a:account) WHERE toLower(a.name)=toLower({acctName}) RETURN " +
		"{account_id: a.account_id, name: a.name, permissions: a.permissions, password: a.password}",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	if len(data) < 1 {
		return nil, true
	}else {
		return data[0][0].(map[string]interface{}), false
	}
}

// Retrieve whether the account exists\
func AccountExists(acctName string) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (a:account) WHERE toLower(a.name)=toLower({acctName}) RETURN a",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	if len(data) < 1 {
		return false
	}else {
		return true
	}
}

// Retrieve account characters
func ListChars(acctName string) ([]string, bool){
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (a:account)-[o:owns]->(c:character) WHERE toLower(a.name)=toLower({acctName}) and " +
		"c.class<>100 and c.active=true RETURN c.name",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	if len(data) < 1 {
		return nil, true
	}else {
		charList := make([]string, len(data))
		for _, k := range data[0] {
			charList = append(charList, k.(string))
		}
		return charList, false
	}
}

// Retrieve account power characters
func ListPowerChar(acctName string) (string, bool) {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (a:account)-[o:owns]->(d:character) WHERE toLower(a.name)=toLower({acctName}) and d.class=100 RETURN " +
		"{char: d.name}",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	if len(data) < 1 {
		return "", true
	}else {
		datum := data[0][0].(map[string]interface{})
		return datum["char"].(string), false
	}
}


// Deactivate account
func Deactivate(acctName string) bool {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll("MATCH (a:account) WHERE a.name={acctName} SET a.active=false",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	if len(data) < 1 {
		return false
	}else {
		return true
	}
}

// Update account
func UpdatePassword(acctName string, acctPass string) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo("MATCH (a:account) WHERE toLower(a.name)=toLower({acctName}) SET " +
		"a.password = {acctPass} ",
		map[string]interface {}{
			"acctName": acctName,
			"acctPass": acctPass,
		},
	)
	fmt.Println(rtrap)
	numResult, _ := result.RowsAffected()
	if numResult < 0 {
		return true
	}else {
		return false
	}
}

func TogglePermission(acctName string, permission uint32) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, rtrap := conn.ExecNeo("MATCH (a:account) WHERE toLower(a.name)=toLower({acctName}) SET " +
		"a.permissions = apoc.bitwise.op(a.permissions,'^',{permission})",
		map[string]interface {}{
			"acctName": acctName,
			"permission": permission,
		},
	)
	log.Println(rtrap)
	numResult, _ := result.RowsAffected()
	if numResult < 0 {
		return true
	}else {
		return false
	}
}

// Delete account
func DeleteAcct(acctName string) bool {
	conn, _ := getConn()
	defer conn.Close()
	result, _ := conn.ExecNeo("MATCH (a:account) WHERE a.name={acctName} DELETE a",
		map[string]interface {}{
			"acctName": acctName,
		},
	)
	numResult, _ := result.RowsAffected()
	if numResult < 0 {
		return true
	}else {
		return false
	}
}

func SearchAccountName(searchStr string, skip int64) []interface{} {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, rtrap:= conn.QueryNeoAll("MATCH (a:account) WHERE toLower(a.name) CONTAINS toLower({search}) ORDER BY a.name LIMIT 15 SKIP {skip} ",
		map[string]interface {}{
			"search": searchStr,
			"skip": skip,
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
