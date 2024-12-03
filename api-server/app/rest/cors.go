package rest

import (
	"net/http"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/utils"
)

func WithCors(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["Access-Control-Allow-Origin"] = []string{utils.GetEnvRequired(common.ENV_WEB_APP_URL)}
		w.Header()["Access-Control-Allow-Methods"] = []string{"GET", "POST", "OPTIONS"}
		w.Header()["Access-Control-Allow-Headers"] = []string{"*"}

		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
		} else {
			mux.ServeHTTP(w, r)
		}
	})
}
