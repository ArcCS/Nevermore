package data

import (
	"database/sql"
	"log"
	"time"
)

type CombatMetric struct {
	Action       string
	ActionType   int // damage 0, heal 1
	Mode         int // melee 0, ranged 1, spell 2, item 3
	TotalDamage  int
	Resisted     int
	FinalDamage  int
	AttackerType int // player 0, mob 1, npc 2
	AttackerId   int
	AttackerTier int
	VictimType   int // player 0, mob 1, npc 2
	VictimId     int
	CombatTime   time.Time
}

type ChatLog struct {
	ChatType int // say 0 osay 1 sent 2 act 3	gmsay 4, ptell 5
	FromId   int
	ToId     int
	Message  string
	ChatTime time.Time
}

type ItemTotals struct {
	ItemId     int
	TotalSold  int
	TotalValue int
	LastSold   time.Time
}

type ItemSales struct {
	ItemId     int
	TimeSold   time.Time
	SellerId   int
	SellerTier int
	SellValue  int
}

func FlushCombatMetrics() bool {
	if len(CombatMetricsCapture) > 0 {
		db, err := sql.Open("postgres", PGCONNSTR)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return false
		}

		// Prepare the insert statement
		stmt, err := tx.Prepare("INSERT INTO combat_metrics (action, type, mode, total_damage, resisted, final_damage, attacker_type, attacker_id, attacker_tier, victim_type, victim_id, time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)")
		if err != nil {
			log.Println(err)
			return false
		}
		defer stmt.Close()

		// Insert multiple rows using the prepared statement
		for _, combatMetric := range CombatMetricsCapture {
			_, err = stmt.Exec(combatMetric.Action,
				combatMetric.ActionType,
				combatMetric.Mode,
				combatMetric.TotalDamage,
				combatMetric.Resisted,
				combatMetric.FinalDamage,
				combatMetric.AttackerType,
				combatMetric.AttackerId,
				combatMetric.AttackerTier,
				combatMetric.VictimType,
				combatMetric.VictimId,
				combatMetric.CombatTime)
			if err != nil {
				log.Println(err)
				return false
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return false
		}

		ClearCombatMetrics()
		return true
	}
	return false
}

func FlushChatLogs() bool {
	if len(ChatLogsCapture) > 0 {
		db, err := sql.Open("postgres", PGCONNSTR)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return false
		}

		// Prepare the insert statement
		stmt, err := tx.Prepare("INSERT INTO chat (type, from_id, to_id, contents, time) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			log.Println(err)
			return false
		}
		defer stmt.Close()

		// Insert multiple rows using the prepared statement
		log.Println("Running Chatlog Flush")
		for _, chatlog := range ChatLogsCapture {
			_, err = stmt.Exec(chatlog.ChatType, chatlog.FromId, chatlog.ToId, chatlog.Message, chatlog.ChatTime)
			if err != nil {
				log.Println(err)
				return false
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return false
		}

		ClearChatLogs()
		return true
	}
	return false
}

func FlushItemSales() bool {
	if len(ItemSalesCapture) > 0 {
		db, err := sql.Open("postgres", PGCONNSTR)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return false
		}

		// Prepare the insert statement
		stmt, err := tx.Prepare("INSERT INTO item_sales (item_id, time, seller_id, seller_tier, sell_value) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			log.Println(err)
			return false
		}
		defer stmt.Close()

		// Insert multiple rows using the prepared statement
		for _, item := range ItemSalesCapture {
			_, err = stmt.Exec(item.ItemId, item.TimeSold, item.SellerId, item.SellerTier, item.SellValue)
			if err != nil {
				log.Println(err)
				return false
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return false
		}

		ClearItemSales()
		return true
	}
	return false
}

func FlushItemTotals() bool {
	if len(ItemTotalsCapture) > 0 {
		db, err := sql.Open("postgres", PGCONNSTR)
		if err != nil {
			log.Println(err)
		}
		defer db.Close()

		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return false
		}

		// Prepare the insert statement
		stmt, err := tx.Prepare("INSERT INTO item_economy (item_id, total_sold, total_value, last_sold) VALUES ($1, $2, $3, $4) ON CONFLICT (item_id) DO UPDATE SET (total_sold, total_value, last_sold) VALUES total_sold + EXCLUDED.total_sold, total_value + EXCLUDED.total_value, EXCLUDED.last_sold")
		if err != nil {
			log.Println(err)
			return false
		}
		defer stmt.Close()

		// Insert multiple rows using the prepared statement
		for _, item := range ItemTotalsCapture {
			_, err = stmt.Exec(item.ItemId, item.TotalSold, item.TotalValue, item.LastSold)
			if err != nil {
				log.Println(err)
				return false
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return false
		}

		ClearItemTotals()
		return true
	}
	return false
}
