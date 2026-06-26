package ratepolice

import (
	"errors"

	"github.com/opticSquid/ratepolice/inmemory"
	"github.com/opticSquid/ratepolice/redis"
	"github.com/opticSquid/ratepolice/shared"
)

type Config = shared.Config

const (
	InMemory             = shared.InMemory
	Redis                = shared.Redis
	FixedWindowCounter   = shared.FixedWindowCounter
	SlidingWindow        = shared.SlidingWindow
	SlidingWindowCounter = shared.SlidingWindowCounter
	TokenBucket          = shared.TokenBucket
	LeakyBucket          = shared.LeakyBucket
)

// NewRateLimiter creates a new rate limiter instance based on the given configuration
func NewRateLimiter(cfg shared.Config) (shared.RatePolice, error) {
	if shared.IsValidConfig(cfg) {
		switch cfg.Backend {
		case shared.InMemory:
			return inmemory.NewInMemoryRateLimiter(cfg), nil
		case shared.Redis:
			limiter, err := redis.NewRedisRateLimiter(cfg)
			if err != nil {
				return nil, err
			}
			return limiter, nil
		}
	}
	return nil, errors.New("Invalid Configuration")
}
