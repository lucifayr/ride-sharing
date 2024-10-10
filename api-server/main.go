package main

import (
	"database/sql"
	"embed"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	migrations "ride_sharing_api/app/database"
)

//go:embed db/migrations/*
var dbMigrations embed.FS

func main() {
	dbFile := "./database.db"
	err := createDbFileIfNotExists(dbFile)
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

	m, err := migrations.FromEmbedFs(dbMigrations, "db/migrations")
	m.Up(db)

	return db, nil
}

func createDbFileIfNotExists(path string) error {
	_, err := os.Stat(path)

	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Create(path)
	f.Close()

	return err
}
