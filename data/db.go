//Entry point for the database,  create the connection for the database here

package data

import (
	"database/sql"
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"strings"
	"time"
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
	ChatLogsCapture      []ChatLog
	ItemSalesCapture     []ItemSales
	ItemTotalsCapture    map[int]ItemTotals
	CombatMetricsCapture []CombatMetric
)

func init() {
	DRIVER, err = neo4j.NewDriver(NEOURI, neo4j.BasicAuth(config.Server.NEOUname, config.Server.NEOPword, ""))
	if err != nil {
		panic(err)
	}

	ChatLogsCapture = make([]ChatLog, 0)
	ItemSalesCapture = make([]ItemSales, 0)
	ItemTotalsCapture = make(map[int]ItemTotals)
	CombatMetricsCapture = make([]CombatMetric, 0)

}

func StoreChatLog(chatType int, fromId int, toId int, message string) {
	ChatLogsCapture = append(ChatLogsCapture, ChatLog{ChatType: chatType, FromId: fromId, ToId: toId, Message: message, ChatTime: time.Now()})
}

func ClearChatLogs() {
	ChatLogsCapture = make([]ChatLog, 0)
}

func StoreItemSale(ItemId int, SellerId int, SellerTier int, SellValue int) {
	ItemSalesCapture = append(ItemSalesCapture, ItemSales{ItemId: ItemId, TimeSold: time.Now(), SellerId: SellerId, SellerTier: SellerTier, SellValue: SellValue})
}

func ClearItemSales() {
	ItemSalesCapture = make([]ItemSales, 0)
}

func StoreCombatMetric(Action string, ActionType int, Mode int, TotalDamage int, Resisted int, FinalDamage int, AttackerType int, AttackerId int, AttackerTier int, VictimType int, VictimId int) {
	CombatMetricsCapture = append(CombatMetricsCapture, CombatMetric{Action: Action, ActionType: ActionType, Mode: Mode, TotalDamage: TotalDamage, Resisted: Resisted, FinalDamage: FinalDamage, AttackerType: AttackerType, AttackerId: AttackerId, AttackerTier: AttackerTier, VictimType: VictimType, VictimId: VictimId, CombatTime: time.Now()})
}

func ClearCombatMetrics() {
	CombatMetricsCapture = make([]CombatMetric, 0)
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

//*** POST GRES FUNCTIONS ***//

func pgExecRead(selectStmt string) (*sql.Rows, error) {
	db, err := sql.Open("postgres", PGCONNSTR)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	// Select rows from the table
	rows, err := db.Query(selectStmt)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return rows, err
}

func pgExec(execStmt string, params ...interface{}) (success bool, err error) {
	db, err := sql.Open("postgres", PGCONNSTR)
	success = true
	if err != nil {
		success = false
		log.Println(err)
	}
	defer db.Close()
	/*
		insertStmt := "INSERT INTO mytable (col1, col2) VALUES ($1, $2)"
		updateStmt := "UPDATE mytable SET col2 = $1 WHERE col1 = $2"
		deleteStmt := "DELETE FROM mytable WHERE col1 = $1"
	*/
	_, err = db.Exec(execStmt, params)
	if err != nil {
		success = false
		log.Println(err)
	}
	return success, err
}
