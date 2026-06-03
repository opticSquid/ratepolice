package ratepolice

import "net/http"

func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Ratelimiting process
		// handing over to next middleware
		next.ServeHTTP(w, r)
	})
}
