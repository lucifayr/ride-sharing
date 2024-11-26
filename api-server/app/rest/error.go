package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ride_sharing_api/app/assert"
)

type restErrors struct {
	Errors []restError `json:"errors"`
}

type restError struct {
	Title   string  `json:"title"`
	Details *string `json:"details"`
}

func httpWriteErr(w http.ResponseWriter, status int, title string, details ...string) {
	var detail *string = nil
	if len(details) > 0 {
		str := fmt.Sprint(details)
		detail = &str
	}

	bytes, err := json.Marshal(restErrors{
		Errors: []restError{
			{
				Title:   title,
				Details: detail,
			},
		},
	})
	assert.Nil(err, "Failed to serialize error message.")

	http.Error(w, string(bytes), status)
}
