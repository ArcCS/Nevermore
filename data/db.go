//Entry point for the database,  create the connection for the database here

package data

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"strings"
)

var (
	NEOURI    = fmt.Sprintf("bolt://%s:7687", config.Server.NEOAddress)
	DRIVER    neo4j.Driver
	err       error
	PGCONNSTR = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Server.PGAddress,
		config.Server.PGPort,
		config.Server.PGUname,
		config.Server.PGPword,
		config.Server.PGUname)
)

func init() {
	DRIVER, err = neo4j.NewDriver(NEOURI, neo4j.BasicAuth(config.Server.NEOUname, config.Server.NEOPword, ""))
	if err != nil {
		panic(err)
	}
}

// Player, Mob, Object, Room, Quest, ItemInventory
func nextId(dataType string) int {
	session := DRIVER.NewSession(neo4j.SessionConfig{})
	defer func() {
		_ = session.Close()
	}()
	nextId, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			fmt.Sprintf("MATCH (r:%[1]s) RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType),
			nil,
		)
		if err != nil {
			return nil, err
		}
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		log.Println(err)
		return 0
	}

	return int(nextId.(int64)) + 1
}

// Player, Mob, Object, Room, Quest, ItemInventory
func nextLinkId(dataType string) int {
	session := DRIVER.NewSession(neo4j.SessionConfig{})
	defer func() {
		_ = session.Close()
	}()
	nextId, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			fmt.Sprintf("MATCH ()-[r:%[1]s]->() RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType),
			nil,
		)
		if err != nil {
			return nil, err
		}
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		log.Println(err)
		return 0
	}

	return int(nextId.(int64)) + 1
}

func execWrite(query string, params map[string]interface{}) (neo4j.ResultSummary, error) {
	session := DRIVER.NewSession(neo4j.SessionConfig{})
	defer func() {
		_ = session.Close()
	}()
	if results, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(query, params)
		if err != nil {
			return nil, err
		}
		summary, err := result.Consume()
		if err != nil {
			return nil, err
		}
		return summary, nil
	}); err != nil {
		return nil, err
	} else {
		return results.(neo4j.ResultSummary), nil
	}
}

func execRead(query string, params map[string]interface{}) ([]*neo4j.Record, error) {
	session := DRIVER.NewSession(neo4j.SessionConfig{})
	defer func() {
		_ = session.Close()
	}()
	if results, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(query, params)
		if err != nil {
			return nil, err
		}
		summary, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return summary, nil
	}); err != nil {
		return nil, err
	} else {
		return results.([]*neo4j.Record), nil
	}
}
