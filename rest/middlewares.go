package rest

import "net/http"

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(wr, req)
	})
}
