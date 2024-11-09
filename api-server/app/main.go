package main

import (
	"log"
	"net/http"

	"ride_sharing_api"
	"ride_sharing_api/app/database"
	"ride_sharing_api/app/database/migrations"
	"ride_sharing_api/app/rest"
	"ride_sharing_api/app/simulator"
	sqlc "ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"

	_ "github.com/mattn/go-sqlite3"
)

var S simulator.Simulator

func main() {
	dbFile := database.NAME
	err := utils.CreateDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	db, err := initDb(dbFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", dbFile, err)
	}

	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: rest.NewRESTApi(sqlc.New(db)),
	}
	log.Println("Listening on", server.Addr)
	err = server.ListenAndServe()
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
