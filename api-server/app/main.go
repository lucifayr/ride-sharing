package main

import (
	"log"

	"ride_sharing_api"
	"ride_sharing_api/app/database"
	"ride_sharing_api/app/database/migrations"
	"ride_sharing_api/app/simulator"
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

	_, err = initDb(dbFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", dbFile, err)
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
