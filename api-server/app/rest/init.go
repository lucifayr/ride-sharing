package rest

import (
	"context"
	"net/http"
	"ride_sharing_api/app/simulator"
	sqlc "ride_sharing_api/app/sqlc"
)

var state = &apiState{}

type apiState struct {
	queries *sqlc.Queries
}

const middlewareKey = "middleware"

type handleFuncBuilder struct {
	handler    func(w http.ResponseWriter, r *http.Request)
	middleware [](func(w http.ResponseWriter, r *http.Request) (bool, *middlewareData))
	data       map[string]any
}

type middlewareData struct {
	key   string
	value any
}

func NewRESTApi(queries *sqlc.Queries) http.Handler {
	state.queries = queries

	mux := simulator.S.NewHttpServerMux()
	authHandlersGoogle(mux)
	userHandlers(mux)

	return mux
}

func handle(handler func(w http.ResponseWriter, r *http.Request)) *handleFuncBuilder {
	return &handleFuncBuilder{handler: handler, middleware: make([]func(w http.ResponseWriter, r *http.Request) (bool, *middlewareData), 0), data: make(map[string]any)}
}

func getMiddlewareData[T any](r *http.Request, key string) T {
	return r.Context().Value(middlewareKey).(map[string]any)[key].(T)
}

func (b *handleFuncBuilder) with(middleware func(w http.ResponseWriter, r *http.Request) (bool, *middlewareData)) *handleFuncBuilder {
	b.middleware = append(b.middleware, middleware)
	return b
}

func (b *handleFuncBuilder) build() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range b.middleware {
			stop, data := handler(w, r)
			if stop {
				return
			}

			if data != nil {
				b.data[data.key] = data.value
			}
		}

		ctx := context.WithValue(r.Context(), middlewareKey, b.data)
		b.handler(w, r.WithContext(ctx))
	}
}

func bearerAuth(w http.ResponseWriter, r *http.Request) (bool, *middlewareData) {
	token := r.Header.Get("Authentication")
	if token == "" {
		http.Error(w, "Missing header 'Authentication'.", http.StatusBadRequest)
		return true, nil
	}

	tokens, err := decodeAccessToken([]byte(token))
	if err != nil {
		http.Error(w, "Invalid access token in 'Authentication' header.", http.StatusBadRequest)
		return true, nil
	}

	// TODO: check if expired

	user, err := state.queries.UsersGetById(r.Context(), *tokens.Id)
	if err != nil {
		http.Error(w, "Failed to get user with matching token.", http.StatusNotFound)
		return true, nil
	}

	if !user.AccessToken.Valid || user.AccessToken.String != token {
		http.Error(w, "Token is invalid.", http.StatusUnauthorized)
		return true, nil
	}

	return false, &middlewareData{key: "user", value: user}
}
