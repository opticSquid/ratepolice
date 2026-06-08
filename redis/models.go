package redis

import (
	"net/http"
	"time"

	"github.com/opticSquid/ratepolice/shared"
)

type RedisRateLimiter struct {
	algorithm          shared.Algorithm
	allowedRequests    int64
	keyFunc            func(*http.Request) string
	timeWindow         time.Duration
	coolDownDur        time.Duration
	coolDownMultiplier int
	redisConn          string
}
