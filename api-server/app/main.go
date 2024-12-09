package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"ride_sharing_api/app/common"
	"ride_sharing_api/app/rest"
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

	handler := rest.NewRESTApi(db)
	server := &http.Server{
		Addr:    utils.GetEnvRequired(common.ENV_HOST_ADDR),
		Handler: handler,
	}

	log.Println("Listening on", server.Addr)

	if strings.ToLower(os.Getenv(common.ENV_NO_TLS)) == "true" {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Closing HTTP server.")
		}
	} else {
		err = server.ListenAndServeTLS(utils.GetEnvRequired(common.ENV_TLS_CERT), utils.GetEnvRequired(common.ENV_TLS_KEY))
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Closing HTTP server.")
		}
	}

}
