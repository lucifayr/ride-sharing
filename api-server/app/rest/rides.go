package rest

import (
	"context"
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
	h.HandleFunc("POST /rides/update", handle(updateRide).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/many", handle(getManyRides).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/by-id/{id}", handle(getEventById).with(bearerAuth(false)).build())
	h.HandleFunc("GET /rides/upcoming/by-id/{id}", handle(getUpcomingById).with(bearerAuth(false)).build())
}

type RideEventData struct {
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

type createRideParams struct {
	LocationFrom   *string       `json:"locationFrom" validate:"required"`
	LocationTo     *string       `json:"locationTo" validate:"required"`
	TackingPlaceAt *time.Time    `json:"tackingPlaceAt" validate:"required"`
	Driver         *string       `json:"driver" validate:"required"`
	TransportLimit *int64        `json:"transportLimit" validate:"required"`
	Schedule       *rideSchedule `json:"schedule"`
}

type updateRideParams struct {
	RideEventId *string       `json:"rideEventId" validate:"required"`
	Schedule    *rideSchedule `json:"schedule"`
	Status      *string       `json:"status"`
}

type rideSchedule struct {
	Unit     *string   `json:"unit" validate:"required"`
	Interval *int64    `json:"interval" validate:"required"`
	Weekdays *[]string `json:"weekdays"`
}

func updateRide(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var updateParams updateRideParams
	err = json.Unmarshal(data, &updateParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(updateParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	tx, err := state.getDBTx(r.Context())
	assert.Nil(err)

	queriesTx := state.queries.WithTx(tx)
	err = markPastRideEventsAsDoneAndCreateScheduled(queriesTx, r.Context())
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	event, err := queriesTx.RidesGetEvent(r.Context(), *updateParams.RideEventId)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No ride event exists for the event with 'id'.")
		return
	}

	assert.Nil(err)

	if event.CreatedBy != user.ID {
		httpWriteErr(w, http.StatusBadRequest, "You are not the owner of this ride event.")
		return
	}

	// update schedule
	if updateParams.Schedule != nil {
		err := queriesTx.RidesDropScheduleWeekdays(r.Context(), event.RideScheduleID.String)
		assert.Nil(err)
		queriesTx.RidesDropSchedule(r.Context(), event.RideScheduleID.String)
		assert.Nil(err)

		argsCreateSchedule := sqlc.RidesCreateScheduleParams{
			RideID:           event.RideID,
			ScheduleInterval: *updateParams.Schedule.Interval,
			Unit:             *updateParams.Schedule.Unit,
		}

		scheduleId, err := queriesTx.RidesCreateSchedule(r.Context(), argsCreateSchedule)
		assert.Nil(err)

		if *updateParams.Schedule.Unit == "weekdays" {
			if updateParams.Schedule.Weekdays == nil || len(*updateParams.Schedule.Weekdays) == 0 {
				httpWriteErr(w, http.StatusBadRequest, "Invalid schedule. Field 'unit' is 'weekdays' but the field 'weekdays' is empty or undefined.")
				return
			}

			for _, day := range *updateParams.Schedule.Weekdays {
				_, err = weekdayToInt(day)
				if err != nil {
					httpWriteErr(w, http.StatusBadRequest, "Invalid schedule weekday. Only lowercase standard english weekday names are allowed.")
					return
				}

				argsCreateScheduleWeekday := sqlc.RidesCreateScheduleWeekdayParams{
					RideScheduleID: scheduleId,
					Weekday:        day,
				}

				err := queriesTx.RidesCreateScheduleWeekday(r.Context(), argsCreateScheduleWeekday)
				assert.Nil(err)
			}
		}
	}

	if updateParams.Status != nil {
		argsUpdateEventStatus := sqlc.RidesUpdateEventStatusParams{
			Status: *updateParams.Status,
			ID:     event.RideEventID,
		}
		err = queriesTx.RidesUpdateEventStatus(r.Context(), argsUpdateEventStatus)
		assert.Nil(err)
	}

	err = tx.Commit()
	assert.Nil(err)
}

func createRide(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var createParams createRideParams
	err = json.Unmarshal(data, &createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	tx, err := state.getDBTx(r.Context())
	assert.Nil(err)

	queriesTx := state.queries.WithTx(tx)

	tackingPlaceAt := createParams.TackingPlaceAt.UTC().Format(time.RFC3339)
	argsCreateBase := sqlc.RidesCreateParams{
		LocationFrom:   *createParams.LocationFrom,
		LocationTo:     *createParams.LocationTo,
		TackingPlaceAt: tackingPlaceAt,
		Driver:         *createParams.Driver,
		CreatedBy:      user.ID,
		TransportLimit: *createParams.TransportLimit,
	}

	rideId, err := queriesTx.RidesCreate(r.Context(), argsCreateBase)
	if err != nil {
		log.Println("Error: Failed to create ride.", "error:", err, "args:", argsCreateBase)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to create ride. This might be due to invalid data or because of an internal server error.")
		return
	}

	if createParams.Schedule != nil {
		argsCreateSchedule := sqlc.RidesCreateScheduleParams{
			RideID:           rideId,
			ScheduleInterval: *createParams.Schedule.Interval,
			Unit:             *createParams.Schedule.Unit,
		}

		scheduleId, err := queriesTx.RidesCreateSchedule(r.Context(), argsCreateSchedule)
		assert.Nil(err)

		if *createParams.Schedule.Unit == "weekdays" {
			if createParams.Schedule.Weekdays == nil || len(*createParams.Schedule.Weekdays) == 0 {
				httpWriteErr(w, http.StatusBadRequest, "Invalid schedule. Field 'unit' is 'weekdays' but the field 'weekdays' is empty or undefined.")
				return
			}

			for _, day := range *createParams.Schedule.Weekdays {
				_, err = weekdayToInt(day)
				if err != nil {
					httpWriteErr(w, http.StatusBadRequest, "Invalid schedule weekday. Only lowercase standard english weekday names are allowed.")
					return
				}

				argsCreateScheduleWeekday := sqlc.RidesCreateScheduleWeekdayParams{
					RideScheduleID: scheduleId,
					Weekday:        day,
				}

				err := queriesTx.RidesCreateScheduleWeekday(r.Context(), argsCreateScheduleWeekday)
				assert.Nil(err)
			}
		}
	}

	rideLatest, err := queriesTx.RidesGetLatest(r.Context(), rideId)
	assert.Nil(err)

	err = tx.Commit()
	assert.Nil(err)

	resp, err := json.Marshal(rideLatest)
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

	tx, err := state.getDBTx(r.Context())
	assert.Nil(err)

	queriesTx := state.queries.WithTx(tx)
	err = markPastRideEventsAsDoneAndCreateScheduled(queriesTx, r.Context())
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tx.Commit()
	assert.Nil(err)

	rows, err := state.queries.RidesGetMany(r.Context(), offset)
	if err != nil {
		log.Println("Error: Failed to get rides.", "error:", err)
		httpWriteErr(w, http.StatusInternalServerError, "Failed to get rides.")
		return
	}

	rides := make([]*RideEventData, len(rows))
	for idx, row := range rows {
		var weekdays *[]string = nil
		if row.RideScheduleID.Valid && row.RideScheduleUnit.String == "weekdays" {
			days, err := state.queries.RidesGetScheduleWeekdays(r.Context(), row.RideScheduleID.String)
			assert.Nil(err)
			weekdays = &days
		}

		event, err := buildRideEventData(rideRow(row), weekdays)
		assert.Nil(err)

		rides[idx] = event
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

func getEventById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	tx, err := state.getDBTx(r.Context())
	assert.Nil(err)

	queriesTx := state.queries.WithTx(tx)
	err = markPastRideEventsAsDoneAndCreateScheduled(queriesTx, r.Context())
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tx.Commit()
	assert.Nil(err)

	event, err := state.queries.RidesGetEvent(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No ride event exists for the event with 'id'.")
		return
	}

	var weekdays *[]string = nil
	if event.RideScheduleUnit.String == "weekdays" {
		days, err := state.queries.RidesGetScheduleWeekdays(r.Context(), event.RideScheduleID.String)
		assert.Nil(err)
		weekdays = &days
	}
	ride, err := buildRideEventData(eventToRideRow(event), weekdays)
	assert.Nil(err)

	var resp []byte
	resp, err = json.Marshal(ride)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(200)
	w.Write(resp)
}

func getUpcomingById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	tx, err := state.getDBTx(r.Context())
	assert.Nil(err)

	queriesTx := state.queries.WithTx(tx)
	err = markPastRideEventsAsDoneAndCreateScheduled(queriesTx, r.Context())
	if err != nil {
		httpWriteErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tx.Commit()
	assert.Nil(err)

	rideLatest, err := state.queries.RidesGetLatest(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No next ride exists for the ride with 'id'.")
		return
	}

	var weekdays *[]string = nil
	if rideLatest.RideScheduleUnit.String == "weekdays" {
		days, err := state.queries.RidesGetScheduleWeekdays(r.Context(), rideLatest.RideScheduleID.String)
		assert.Nil(err)
		weekdays = &days
	}
	ride, err := buildRideEventData(latestToRideRow(rideLatest), weekdays)
	assert.Nil(err)

	var resp []byte
	resp, err = json.Marshal(ride)
	assert.Nil(err, "Failed to serialize ride.")
	w.WriteHeader(200)
	w.Write(resp)
}

type rideRow struct {
	RideID               string
	RideEventID          string
	LocationFrom         string
	LocationTo           string
	TackingPlaceAt       string
	CreatedBy            string
	TransportLimit       int64
	Driver               string
	Status               string
	DriverEmail          string
	CreatedByEmail       string
	RideScheduleID       sql.NullString
	RideScheduleUnit     sql.NullString
	RideScheduleInterval sql.NullInt64
}

func eventToRideRow(row sqlc.RidesGetEventRow) rideRow {
	return rideRow{
		RideID:               row.RideID,
		RideEventID:          row.RideEventID,
		LocationFrom:         row.LocationFrom,
		LocationTo:           row.LocationTo,
		TackingPlaceAt:       row.TackingPlaceAt,
		CreatedBy:            row.CreatedBy,
		TransportLimit:       row.TransportLimit,
		Driver:               row.Driver,
		Status:               row.Status,
		DriverEmail:          row.DriverEmail,
		CreatedByEmail:       row.CreatedByEmail,
		RideScheduleID:       row.RideScheduleID,
		RideScheduleUnit:     row.RideScheduleUnit,
		RideScheduleInterval: row.RideScheduleInterval,
	}
}

func latestToRideRow(row sqlc.RidesGetLatestRow) rideRow {
	return rideRow{
		RideID:               row.RideID,
		RideEventID:          row.RideEventID,
		LocationFrom:         row.LocationFrom,
		LocationTo:           row.LocationTo,
		TackingPlaceAt:       row.TackingPlaceAt,
		CreatedBy:            row.CreatedBy,
		TransportLimit:       row.TransportLimit,
		Driver:               row.Driver,
		Status:               row.Status,
		DriverEmail:          row.DriverEmail,
		CreatedByEmail:       row.CreatedByEmail,
		RideScheduleID:       row.RideScheduleID,
		RideScheduleUnit:     row.RideScheduleUnit,
		RideScheduleInterval: row.RideScheduleInterval,
	}
}

func buildRideEventData(ride rideRow, weekdays *[]string) (*RideEventData, error) {
	tackingPlaceAt, err := time.Parse(time.RFC3339, ride.TackingPlaceAt)
	if err != nil {
		return nil, err
	}

	var schedule *rideSchedule = nil
	if ride.RideScheduleID.Valid {
		schedule = &rideSchedule{
			Unit:     &ride.RideScheduleUnit.String,
			Interval: &ride.RideScheduleInterval.Int64,
			Weekdays: weekdays,
		}
	}

	rideEvent := RideEventData{
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

func nextScheduledTime(current time.Time, scheduleUnit string, scheduleInterval int64, weekdays *[]string) (time.Time, error) {
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
			if weekdays == nil || len(*weekdays) == 0 {
				return current, errors.New(fmt.Sprint("Invalid weekdays schedule, received weekdays ", weekdays))
			}

			weekdayInts := make([]int, len(*weekdays))
			for idx, day := range *weekdays {
				dayInt, err := weekdayToInt(day)
				if err != nil {
					return current, err
				}

				weekdayInts[idx] = dayInt
			}

			minDayDist := 8
			nowDayInt := int(current.Weekday())
			for _, dayInt := range weekdayInts {
				var dist int

				if nowDayInt < dayInt {
					dist = dayInt - nowDayInt
				} else {
					dist = 7 - (nowDayInt - dayInt)
				}

				if dist < minDayDist {
					minDayDist = dist
				}
			}

			assert.True(minDayDist <= 7, "Invalid next weekday. Distance to now day", minDayDist)

			return current.AddDate(0, 0, minDayDist), nil
		}
	default:
		{
			return current, fmt.Errorf("Invalid schedule unit '%s'", scheduleUnit)
		}
	}
}

func weekdayToInt(day string) (int, error) {
	switch day {
	case "sunday":
		{
			return 0, nil
		}
	case "monday":
		{
			return 1, nil
		}
	case "tuesday":
		{
			return 2, nil
		}
	case "wednesday":
		{
			return 3, nil
		}
	case "thursday":
		{
			return 4, nil
		}
	case "friday":
		{
			return 5, nil
		}
	case "saturday":
		{
			return 6, nil
		}
	default:
		{
			return -1, fmt.Errorf("Invalid weekday '%s'", day)
		}
	}
}

func markPastRideEventsAsDoneAndCreateScheduled(queriesTx *sqlc.Queries, ctx context.Context) error {
	now := time.Now()
	updated, err := queriesTx.RidesMarkPastEventsDone(ctx, now.UTC().Format(time.RFC3339))
	if err != nil {
		return errors.New("Failed to update rides.")
	}

	for _, id := range updated {
		event, err := queriesTx.RidesGetEvent(ctx, id)
		assert.Nil(err)

		schedule, err := queriesTx.RidesGetSchedule(ctx, event.RideID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		assert.Nil(err)

		var weekdays *[]string = nil
		if event.RideScheduleUnit.String == "weekdays" {
			days, err := queriesTx.RidesGetScheduleWeekdays(ctx, schedule.ID)
			assert.Nil(err)
			weekdays = &days
		}
		next, err := nextScheduledTime(now, schedule.Unit, schedule.ScheduleInterval, weekdays)
		assert.Nil(err)

		argsCreateEvent := sqlc.RidesCreateEventParams{
			RideID:         event.RideID,
			LocationFrom:   event.LocationFrom,
			LocationTo:     event.LocationTo,
			TransportLimit: event.TransportLimit,
			Driver:         event.Driver,
			TackingPlaceAt: next.UTC().Format(time.RFC3339),
		}

		err = queriesTx.RidesCreateEvent(ctx, argsCreateEvent)
		assert.Nil(err)
	}

	return nil
}
