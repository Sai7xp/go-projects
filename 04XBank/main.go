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
	"github.com/sai7xp/xbank/utils"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:sumanth123@localhost:5432/xbank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// ) // Configurations moved to app.env file and loaded via Viper

func main() {
	fmt.Println("XBank Main File")
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}
}
