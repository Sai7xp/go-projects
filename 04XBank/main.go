/*
* Created on 06 May 2024
* @author Sai Sumanth
 */
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sai7xp/xbank/api"
	db "github.com/sai7xp/xbank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:sumanth123@localhost:5432/xbank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	fmt.Println("XBank Main File")
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}
}
