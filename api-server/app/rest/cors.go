package rest

import (
	"net/http"
)

func WithCors(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
		w.Header()["Access-Control-Allow-Methods"] = []string{"GET", "POST", "OPTIONS"}
		w.Header()["Access-Control-Allow-Headers"] = []string{"*"}

		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
		} else {
			mux.ServeHTTP(w, r)
		}
	})
}
