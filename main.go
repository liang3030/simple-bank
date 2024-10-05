package main

import (
	"database/sql"
	"log"

	"github.com/liang3030/simple-bank/api"
	db "github.com/liang3030/simple-bank/db/sqlc"
	"github.com/liang3030/simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// connnect to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	// create a new store
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
