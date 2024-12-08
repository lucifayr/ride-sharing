package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"strconv"
	"time"
)

func rideHandlers(h *http.ServeMux) {
	h.HandleFunc("POST /rides", handle(createRide).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/many", handle(getManyRides).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/by-id/{id}", handle(getRideById).with(bearerAuth(false)).build())
}

type createRideParams struct {
	LocationFrom   *string    `json:"locationFrom" validate:"required"`
	LocationTo     *string    `json:"locationTo" validate:"required"`
	TackingPlaceAt *time.Time `json:"tackingPlaceAt" validate:"required"`
	Driver         *string    `json:"driver" validate:"required"`
	TransportLimit *int64     `json:"transportLimit" validate:"required"`
}

type rideEventData struct {
	RideId         string        `json:"rideId"`
	RideEventId    string        `json:"rideEventId"`
	LocationFrom   string        `json:"locationFrom"`
	LocationTo     string        `json:"locationTo"`
	TackingPlaceAt time.Time     `json:"tackingPlaceAt"`
	Status         string        `json:"status"`
	CreatedBy      string        `json:"createdBy"`
	CreatedByEmail string        `json:"createdByEmail"`
	DriverId       string        `json:"driverId"`
	DriverEmail    string        `json:"driverEmail"`
	TransportLimit int64         `json:"transportLimit"`
	Schedule       *rideSchedule `json:"schedule"`
}

type rideSchedule struct {
	Unit     string    `json:"unit"`
	Interval int64     `json:"interval"`
	Weekdays *[]string `json:"weekdays"`
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

func getManyRides(w http.ResponseWriter, r *http.Request) {
	var offset int64 = 0
	offsetStr := r.FormValue("offset")
	if parsed, err := strconv.ParseInt(offsetStr, 10, 64); err == nil && parsed > 0 {
		offset = parsed
	}

	rides, err := state.queries.RidesGetMany(r.Context(), offset)
	if err != nil {
		log.Println("Error: Failed to get rides.", "error:", err)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to get rides.")
		return
	}

	var resp []byte
	if len(rides) == 0 {
		resp = []byte("[]")
	} else {
		resp, err = json.Marshal(rides)
	}

	assert.Nil(err, "Failed to serialize rides.")
	w.WriteHeader(200)
	w.Write(resp)
}

func getRideById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	ride, err := state.queries.RidesGetById(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No ride with the provide 'id' exists.")
		return
	}

	if err != nil {
		log.Println("Error: Failed to get ride.", "error:", err)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to get rides.")
		return
	}

	var resp []byte
	resp, err = json.Marshal(ride)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(200)
	w.Write(resp)
}
