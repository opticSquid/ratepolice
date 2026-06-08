package redis

import (
	"encoding/json"
	"net/http"

	"github.com/opticSquid/ratepolice/shared"
)

func NewRedisRateLimiter(cfg shared.Config) *RedisRateLimiter {
	return &RedisRateLimiter{
		algorithm:          cfg.Algorithm,
		allowedRequests:    cfg.MaxAllowedRequests,
		keyFunc:            cfg.KeyFunc,
		timeWindow:         cfg.TimeWindow,
		coolDownDur:        cfg.CoolDownDur,
		coolDownMultiplier: cfg.CoolDownMultiplier,
		redisConn:          cfg.RedisConn,
	}
}

func (cfg *RedisRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Ratelimiting process in Redis
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		_ = json.NewEncoder(w).Encode(shared.ResponseHeaders{
			Type:   "Configuration Error",
			Title:  http.StatusText(http.StatusNotImplemented),
			Detail: "Redis backend for rate limiter not implemented yet",
		})
		// return
	})
}

func (cfg *RedisRateLimiter) Stop() {
}
