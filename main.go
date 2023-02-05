package main

import (
	"database/sql"
	"log"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/routes"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("app", ".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB : ", err)
	}

	store := db.NewStore(conn)
	server, err := api.SetupServer(config, store)
	if err != nil {
		log.Fatal("cannot create server : ", err)
	}
	routes.ApplyAllPublicRoutes(server)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server : ", err)
	}

}
