package ratepolice

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func NewRateLimiter(cfg Config) (*RateLimiter, error) {
	if isValidConfig(cfg) {
		return &RateLimiter{
			data:            make(map[string][]time.Time),
			algoritm:        cfg.Algoritm,
			backend:         cfg.Backend,
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
		switch cfg.backend {
		case InMemory:
		case Redis:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			_ = json.NewEncoder(w).Encode(ProblemDetails{
				Type:   "Configuration Error",
				Title:  http.StatusText(http.StatusNotImplemented),
				Detail: "Redis backend for rate limiter not implemented yet",
			})
			return
		}
		// handing over to next middleware
		next.ServeHTTP(w, r)
	})
}
