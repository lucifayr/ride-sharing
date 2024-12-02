package rest

import (
	"encoding/json"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/sqlc"
)

func userHandlers(h *http.ServeMux) {
	h.HandleFunc("GET /users/me", handle(getUserMe).with(bearerAuth(false)).build())
}

func getUserMe(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")
	bytes, err := json.Marshal(user)
	assert.True(err == nil, "Failed to serialize user struct.", "user:", user, "error:", err)

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}
