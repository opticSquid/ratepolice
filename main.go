package ratepolice

import (
	"errors"
	"net/http"
	"time"
)

func NewRateLimiter(cfg Config) (*RateLimiter, error) {
	if isValidConfig(cfg) {
		return &RateLimiter{
			data:            make(map[string][]time.Time),
			allowedRequests: cfg.AllowedRequests,
			keyFunc:         cfg.KeyFunc,
			timeWindow:      cfg.TimeWindow,
		}, nil
	}
	return nil, errors.New("Invalid Configuration")
}

func (cfg *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Ratelimiting process
		// handing over to next middleware
		next.ServeHTTP(w, r)
	})

}
