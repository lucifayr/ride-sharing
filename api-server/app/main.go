package main

import (
	"database/sql"
	"log"

	embeddings "ride_sharing_api"
	utils "ride_sharing_api/app/database"
	migrations "ride_sharing_api/app/database/migrations"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFile := "./database.db"
	err := utils.CreateDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	_, err = initDb(dbFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", dbFile, err)
	}
}

func initDb(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+dbFile)
	if err != nil {
		return nil, err
	}

	m := migrations.FromEmbedFs(embeddings.DbMigrations, "db/migrations")
	m.Up(db)

	return db, nil
}
