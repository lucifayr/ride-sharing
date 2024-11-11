package main

import (
	"log"

	"ride_sharing_api"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/database/migrations"
	"ride_sharing_api/app/rest"
	"ride_sharing_api/app/simulator"
	sqlc "ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"

	_ "github.com/mattn/go-sqlite3"
)

var S simulator.Simulator

func main() {
	dbFile := simulator.S.DbName()
	err := utils.CreateDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	db, err := initDb(dbFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", dbFile, err)
	}

	handler := rest.NewRESTApi(sqlc.New(db))
	err = simulator.S.HttpListenAndServe(handler, simulator.S.GetEnvRequired(common.ENV_HOST_ADDR))
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Closing HTTP server.")
	}
}

func initDb(dbFile string) (simulator.DB, error) {
	db, err := simulator.S.SqlOpen("sqlite3", "file:"+dbFile)
	if err != nil {
		return nil, err
	}

	m := migrations.FromEmbedFs(embeddings.DbMigrations, "db/migrations")
	m.Up(db)

	return db, nil
}
