package rest

import (
	"context"
	"database/sql"
	"net/http"
	sqlc "ride_sharing_api/app/sqlc"
	"sync"
	"time"
)

var state *apiState

type apiState struct {
	queries     *sqlc.Queries
	mutex       sync.Mutex
	oauthStates map[string]time.Time
	getDBTx     func(ctx context.Context) (*sql.Tx, error)
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

func NewRESTApi(db *sql.DB) http.Handler {
	state = &apiState{oauthStates: make(map[string]time.Time), queries: sqlc.New(db), getDBTx: func(ctx context.Context) (*sql.Tx, error) {
		return db.BeginTx(ctx, &sql.TxOptions{})
	}}

	mux := http.NewServeMux()

	authHandlers(mux)
	authHandlersGoogle(mux)
	userHandlers(mux)
	rideHandlers(mux)
	groupHandlers(mux)
	groupMessageHandlers(mux)

	return WithCors(mux)
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

func bearerAuth(ignoreExpired bool) func(w http.ResponseWriter, r *http.Request) (bool, *middlewareData) {
	return func(w http.ResponseWriter, r *http.Request) (bool, *middlewareData) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing header 'Authorization'.", http.StatusBadRequest)
			return true, nil
		}

		tokens, err := decodeAccessToken([]byte(token))
		if err != nil || (!ignoreExpired && time.Now().After(tokens.ExpiresAt)) {
			http.Error(w, "Invalid access token in 'Authorization' header.", http.StatusUnauthorized)
			return true, nil
		}

		user, err := state.queries.UsersGetById(r.Context(), *tokens.Id)
		if err != nil {
			http.Error(w, "Invalid access token in 'Authorization' header.", http.StatusUnauthorized)
			return true, nil
		}

		if !user.AccessToken.Valid || user.AccessToken.String != token {
			http.Error(w, "Invalid access token in 'Authorization' header.", http.StatusUnauthorized)
			return true, nil
		}

		return false, &middlewareData{key: "user", value: user}
	}
}
