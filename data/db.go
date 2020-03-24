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
func nextId(dataType string) int64{
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll(fmt.Sprintf("MATCH (r:%[1]s) RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType), nil)
	return data[0][0].(int64)+1
}

// Player, Mob, Object, Room, Quest, ItemInventory
func nextLinkId(dataType string) int64{
	conn, _ := getConn()
	defer conn.Close()
	data, _, _, _ := conn.QueryNeoAll(fmt.Sprintf("MATCH ()-[r:%[1]s]->() RETURN COALESCE(MAX(r.%[2]s_id), 0)", strings.ToLower(dataType), dataType), nil)
	return data[0][0].(int64)+1
}

//
// Need to increment ID's
// MATCH (r.Room) MAX(r.room_id)+1
// MATCH (r.Player) MAX(r.player_id+1

// CREATE CONSTRAINT ON (personj:Person) ASSERT person.id IS UNIQUE

// Here we prepare a new statement. This gives us the flexibility to
// cancel that statement without any request sent to Neo
/*
func prepareStatement(query string) bolt.Stmt {
	st, err := conn.PrepareNeo(query)
	handleError(err)
	return st
}

// Here we prepare a new pipeline statement for running multiple
// queries concurrently
func preparePipeline(stmts ...string) bolt.PipelineStmt {
	pipeline, err := conn.PreparePipeline(
		stmts...
	)
	handleError(err)
	return pipeline
}

func executePipeline(pipeline bolt.PipelineStmt) {
	pipelineResults, err := pipeline.ExecPipeline(nil, nil, nil, nil, nil, nil)
	handleError(err)

	for _, result := range pipelineResults {
		numResult, _ := result.RowsAffected()
		fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 2 (per each iteration)
	}

	err = pipeline.Close()
	handleError(err)
}

func queryStatement(st bolt.Stmt) bolt.Rows {
	// Even once I get the rows, if I do not consume them and close the
	// rows, Neo will discard and not send the data
	rows, err := st.QueryNeo(nil)
	handleError(err)
	return rows
}

func consumeMetadata(rows bolt.Rows, st bolt.Stmt) {
	// Here we loop through the rows until we get the metadata object
	// back, meaning the row stream has been fully consumed

	var err error
	err = nil

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {
			fmt.Printf("PATH: %#v\n", row[0].(graph.Path)) // Prints all paths
		}
	}
	st.Close()
}

func consumeRows(rows bolt.Rows, st bolt.Stmt) {
	// This interface allows you to consume rows one-by-one, as they
	// come off the bolt stream. This is more efficient especially
	// if you're only looking for a particular row/set of rows, as
	// you don't need to load up the entire dataset into memory
	data, _, err := rows.NextNeo()
	handleError(err)

	// This query only returns 1 row, so once it's done, it will return
	// the metadata associated with the query completion, along with
	// io.EOF as the error
	_, _, err = rows.NextNeo()
	handleError(err)
	fmt.Printf("COLUMNS: %#v\n", rows.Metadata()["fields"].([]interface{})) // COLUMNS: n.foo,n.bar
	fmt.Printf("FIELDS: %d %f\n", data[0].(int64), data[1].(float64))       // FIELDS: 1 2.2

	st.Close()
}

// Executing a statement just returns summary information
func executeStatement(st bolt.Stmt) {
	result, err := st.ExecNeo(map[string]interface{}{"foo": 1, "bar": 2.2})
	handleError(err)
	numResult, err := result.RowsAffected()
	handleError(err)
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	st.Close()
}

// Here we create a simple function that will take care of errors, helping with some code clean up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
*/