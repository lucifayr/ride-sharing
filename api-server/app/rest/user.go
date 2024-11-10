package rest

import (
	"log"
	"net/http"
	"ride_sharing_api/app/simulator"
	"ride_sharing_api/app/sqlc"
)

func userHandlers(h simulator.HTTPMux) {
	h.HandleFunc("GET /users", handle(getUserById).with(bearerAuth(false)).build())
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")
	log.Println(user)
}
