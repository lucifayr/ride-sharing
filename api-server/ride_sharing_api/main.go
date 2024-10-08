package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	dbFile := "./database.db"
	err := createDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	fmt.Println("Hello World")
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
