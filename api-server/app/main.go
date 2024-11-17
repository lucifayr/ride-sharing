package main

import (
	"log"
	"net/http"

	"ride_sharing_api/app/common"
	"ride_sharing_api/app/rest"
	sqlc "ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbFile := utils.GetEnvRequired(common.ENV_DB_NAME)
	err := utils.CreateDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	db, err := utils.InitDb(dbFile)
	if err != nil {
		log.Fatalln("Failed to initialize database.", dbFile, err)
	}

	handler := rest.NewRESTApi(sqlc.New(db))
	server := &http.Server{
		Addr:    utils.GetEnvRequired(common.ENV_HOST_ADDR),
		Handler: handler,
	}

	log.Println("Listening on", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Closing HTTP server.")
	}
}
