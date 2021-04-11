//Entry point for the database,  create the connection for the database here

package data

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"strings"
)

var (
	URI = fmt.Sprintf("bolt://%s:%s@%s:7687", config.Server.DBUname, config.Server.DBPword, config.Server.DBAddress)
)

func getConn() (bolt.Conn, error) {
	driver := bolt.NewDriver()
	return driver.OpenNeo(URI)
}

// Player, Mob, Object, Room, Quest, ItemInventory
func nextId(dataType string) int {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll(fmt.Sprintf("MATCH (r:%[1]s) RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType), nil)
	return int(data[0][0].(int64)) + 1
}

// Player, Mob, Object, Room, Quest, ItemInventory
func nextLinkId(dataType string) int {
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll(fmt.Sprintf("MATCH ()-[r:%[1]s]->() RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType), nil)
	return int(data[0][0].(int64)) + 1
}
