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
)

func userHandlers(h *http.ServeMux) {
	h.HandleFunc("GET /users/me", handle(getUserMe).with(bearerAuth(false)).build())
	h.HandleFunc("GET /users/by-id/{id}", handle(getUserById).with(bearerAuth(false)).build())
	h.HandleFunc("POST /users/by-id/{id}/ban-status", handle(setUserBanStatus).with(bearerAuth(false)).build())
}

func getUserMe(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")
	if user.IsBlocked {
		w.Header().Add("Content-Type", "application/json")
		httpWriteErr(w, http.StatusForbidden, "Your account has been blocked by an admin.")
		return
	}

	bytes, err := json.Marshal(user)
	assert.True(err == nil, "Failed to serialize user struct.", "user:", user, "error:", err)

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	getMiddlewareData[sqlc.User](r, "user")

	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	user, err := state.queries.UsersGetById(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		httpWriteErr(w, http.StatusNotFound, "No user exists with 'id'.")
		return
	}

	var resp []byte
	resp, err = json.Marshal(user)
	assert.Nil(err, "Failed to serialize group.")
	w.WriteHeader(200)
	w.Write(resp)
}

type setBanStatusParams struct {
	IsBanned *bool `json:"isBanned" validate:"required"`
}

func setUserBanStatus(w http.ResponseWriter, r *http.Request) {
	getMiddlewareData[sqlc.User](r, "user")

	id := r.PathValue("id")
	if id == "" {
		httpWriteErr(w, http.StatusBadRequest, "Must provide 'id' path parameter.")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var params setBanStatusParams
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

	args := sqlc.UsersSetBlockedParams{
		IsBlocked: *params.IsBanned,
		ID:        id,
	}

	err = state.queries.UsersSetBlocked(r.Context(), args)
	assert.Nil(err)
	w.WriteHeader(200)
}
