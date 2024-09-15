package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Make sure the postgres enginee driver is imported
)

const (
	dbDriver = "postgres"
	// TODO: add environment variable for test database
	dbSource = "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
