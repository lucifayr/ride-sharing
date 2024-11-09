package rest

import (
	"net/http"
	"ride_sharing_api/app/simulator"
	sqlc "ride_sharing_api/app/sqlc"
)

type restState struct {
	queries *sqlc.Queries
}

var state = &restState{}

const (
	AUTH_PROVIDER_GOOGLE    = "google"
	AUTH_PROVIDER_MICROSOFT = "microsoft"
)

func NewRESTApi(queries *sqlc.Queries) http.Handler {
	state.queries = queries

	mux := simulator.S.NewHttpServerMux()
	authHandlersGoogle(mux)

	return mux
}
