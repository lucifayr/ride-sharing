package rest

import (
	"fmt"
	"net/http"
	"ride_sharing_api/app/simulator"
	sqlc "ride_sharing_api/app/sqlc"
	"slices"
	"strings"
)

type restState struct {
	queries *sqlc.Queries
}

var state = &restState{}

func NewRESTApi(queries *sqlc.Queries) http.Handler {
	state.queries = queries

	mux := simulator.S.NewHttpServerMux()
	authHandlersGoogle(mux)

	return mux
}

func withAllowedMethods(handler func(w http.ResponseWriter, r *http.Request), methods ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(methods, r.Method) {
			handler(w, r)
		} else {
			http.Error(w, fmt.Sprintf("Invalid method. Allowed are [%s]", strings.Join(methods, ",")), http.StatusMethodNotAllowed)
		}
	}
}
