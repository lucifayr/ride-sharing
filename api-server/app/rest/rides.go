package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
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
	TransportLimit *int64     `json:"transportLimit" validate:"required"`
}

func createRide(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var params createRideParams
	err = json.Unmarshal(data, &params)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(params)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	tackingPlaceAt := params.TackingPlaceAt.UTC().Format(time.RFC3339)
	args := sqlc.RidesCreateParams{
		LocationFrom:   *params.LocationFrom,
		LocationTo:     *params.LocationTo,
		TackingPlaceAt: tackingPlaceAt,
		Driver:         *params.Driver,
		CreatedBy:      user.ID,
		TransportLimit: *params.TransportLimit,
	}

	ride, err := state.queries.RidesCreate(r.Context(), args)
	if err != nil {
		log.Println("Error: Failed to create ride.", "error:", err, "args:", args)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to create ride. This might be due to invalid data or because of an internal server error.")
		return
	}

	resp, err := json.Marshal(ride)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(201)
	w.Write(resp)
}
