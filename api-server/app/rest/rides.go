package rest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"strconv"
	"time"
)

const (
	RIDE_STATUS_DONE     = "done"
	RIDE_STATUS_CANCELED = "canceled"
	RIDE_STATUS_UPCOMING = "upcoming"
)

func rideHandlers(h *http.ServeMux) {
	h.HandleFunc("POST /rides", handle(createRide).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/many", handle(getNextRideById).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/by-id/{id}", handle(getNextRideById).with(bearerAuth(false)).build())
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

// TODO: update to handle new fields
func getManyRides(w http.ResponseWriter, r *http.Request) {
	var offset int64 = 0
	offsetStr := r.FormValue("offset")
	if parsed, err := strconv.ParseInt(offsetStr, 10, 64); err == nil && parsed > 0 {
		offset = parsed
	}

	now := time.Now()
	err := state.queries.RidesMarkPastEventsDone(r.Context(), now.UTC().Format(time.RFC3339))
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, "Failed to update rides.")
		return
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

func getNextRideById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	now := time.Now()
	err := state.queries.RidesMarkPastEventsDone(r.Context(), now.UTC().Format(time.RFC3339))
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, "Failed to update rides.")
		return
	}

	rideLatest, err := state.queries.RidesGetLatest(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No ride with the provide 'id' exists.")
		return
	}

	if rideLatest.Status != RIDE_STATUS_UPCOMING && !rideLatest.RideScheduleID.Valid {
		httpWriteErr(w, http.StatusNotFound, "No next ride exists for the ride with 'id'.")
		return
	}

	if err != nil {
		log.Println("Error: Failed to get ride.", "error:", err)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to get rides.")
		return
	}

	var weekdays *[]string = nil
	if rideLatest.RideScheduleUnit.String == "weekdays" {
		days, err := state.queries.RidesGetScheduleWeekdays(r.Context(), rideLatest.RideScheduleID.String)
		weekdays = &days
	}

	var ride *rideEventData = nil
	if rideLatest.Status != RIDE_STATUS_UPCOMING {
		next, err := nextScheduledTime(now, rideLatest.RideScheduleID.String, rideLatest.RideScheduleUnit.String, rideLatest.RideScheduleInterval.Int64, weekdays)

		argsCreateEvent := sqlc.RidesCreateEventParams{
			RideID:         rideLatest.RideID,
			LocationFrom:   rideLatest.BaseLocationFrom,
			LocationTo:     rideLatest.BaseLocationTo,
			TransportLimit: rideLatest.BaseTransportLimit,
			Driver:         rideLatest.BaseDriver,
			TackingPlaceAt: next.UTC().Format(time.RFC3339),
		}

		err = state.queries.RidesCreateEvent(r.Context(), argsCreateEvent)

		rideNext, err := state.queries.RidesGetLatest(r.Context(), id)
		ride, err = buildRideEventData(rideNext, weekdays)
		return
	} else {
		ride, err = buildRideEventData(rideLatest, weekdays)
	}

	var resp []byte
	resp, err = json.Marshal(ride)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(200)
	w.Write(resp)
}

func buildRideEventData(ride sqlc.RidesGetLatestRow, weekdays *[]string) (*rideEventData, error) {
	tackingPlaceAt, err := time.Parse(time.RFC3339, ride.TackingPlaceAt)
	if err != nil {
		return nil, err
	}

	var schedule *rideSchedule = nil
	if ride.RideScheduleID.Valid {
		schedule = &rideSchedule{
			Unit:     ride.RideScheduleUnit.String,
			Interval: ride.RideScheduleInterval.Int64,
			Weekdays: weekdays,
		}
	}

	rideEvent := rideEventData{
		RideId:         ride.RideID,
		RideEventId:    ride.RideEventID,
		LocationFrom:   ride.LocationFrom,
		LocationTo:     ride.LocationTo,
		TackingPlaceAt: tackingPlaceAt,
		Status:         ride.Status,
		CreatedBy:      ride.CreatedBy,
		CreatedByEmail: ride.CreatedByEmail,
		DriverId:       ride.Driver,
		DriverEmail:    ride.DriverEmail,
		TransportLimit: ride.TransportLimit,
		Schedule:       schedule,
	}

	return &rideEvent, nil
}

func nextScheduledTime(current time.Time, scheduleId string, scheduleUnit string, scheduleInterval int64, weekdays *[]string) (time.Time, error) {
	switch scheduleUnit {
	case "days":
		{
			return current.AddDate(0, 0, int(scheduleInterval)), nil
		}
	case "weeks":
		{
			return current.AddDate(0, 0, int(scheduleInterval)*7), nil
		}
	case "months":
		{
			return current.AddDate(0, int(scheduleInterval), 0), nil
		}
	case "years":
		{
			return current.AddDate(int(scheduleInterval), 0, 0), nil
		}
	case "weekdays":
		{
			return current, fmt.Errorf("TODO")
		}
	default:
		{
			return current, fmt.Errorf("Invalid schedule unit '%s'", scheduleUnit)
		}
	}
}
