package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"time"
)

func rideHandlers(h *http.ServeMux) {
	h.HandleFunc("POST /rides", handle(createRide).with(bearerAuth(false)).build())
}

type createRideParams struct {
	LocationFrom   *string    `json:"locationFrom" validate:"required"`
	LocationTo     *string    `json:"locationTo" validate:"required"`
	TackingPlaceAt *time.Time `json:"tackingPlaceAt" validate:"required"`
	Driver         *string    `json:"driver" validate:"required"`
}

func createRide(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	var params createRideParams
	err = json.Unmarshal(data, &params)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		http.Error(w, "Invalid JSON in request body.", http.StatusBadRequest)
		return
	}

	err = utils.Validate.Struct(params)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		http.Error(w, "Invalid JSON in request body. "+err.Error(), http.StatusBadRequest)
		return
	}

	args := sqlc.RidesCreateParams{LocationFrom: *params.LocationFrom, LocationTo: *params.LocationTo, TackingPlaceAt: params.TackingPlaceAt.Format(time.RFC3339), Driver: *params.Driver, CreateBy: user.ID}
	ride, err := state.queries.RidesCreate(r.Context(), args)
	if err != nil {
		log.Println("Error: Failed to create ride.", "error:", err, "args:", args)
		http.Error(w, "Failed to create ride. This might be due to invalid data or because of an internal server error.", http.StatusInternalServerError)
		return
	}

	log.Println("Created ride", ride)
}
