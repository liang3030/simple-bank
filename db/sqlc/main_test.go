package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/liang3030/simple-bank/util"
	_ "github.com/lib/pq" // Make sure the postgres enginee driver is imported
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
