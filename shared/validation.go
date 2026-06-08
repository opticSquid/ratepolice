package shared

func isAlgoValid(algo Algorithm) bool {
	switch algo {
	case FixedWindowCounter, SlidingWindow, SlidingWindowCounter, TokenBucket, LeakyBucket:
		return true
	default:
		return false
	}
}

func IsValidConfig(cfg Config) bool {
	return isAlgoValid(cfg.Algorithm)
}
