package ratepolice

import (
	"errors"

	"github.com/opticSquid/ratepolice/inmemory"
	"github.com/opticSquid/ratepolice/redis"
	"github.com/opticSquid/ratepolice/shared"
)

func NewRateLimiter(cfg shared.Config) (shared.RatePolice, error) {
	if shared.IsValidConfig(cfg) {
		switch cfg.Backend {
		case shared.InMemory:
			return inmemory.NewInMemoryRateLimiter(cfg), nil
		case shared.Redis:
			return redis.NewRedisRateLimiter(cfg), nil
		}
	}
	return nil, errors.New("Invalid Configuration")
}
