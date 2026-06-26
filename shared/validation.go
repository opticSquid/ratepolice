package shared

func isAlgoValid(backend Backend, algo Algorithm) bool {
	switch backend {
	case InMemory:
		switch algo {
		case FixedWindowCounter:
			return true
		default:
			return false
		}
	// For now as Redis backend is not implemented we return false
	case Redis:
		return false
	default:
		return false
	}
}
func isBackendValid(backend Backend) bool {
	switch backend {
	case InMemory:
		return true
	// For now as Redis backend is not implemented we return false
	case Redis:
		return false
	default:
		return false
	}
}

// IsValidConfig validates the given configuration
func IsValidConfig(cfg Config) bool {
	if isBackendValid(cfg.Backend) {
		return isAlgoValid(cfg.Backend, cfg.Algorithm)
	}
	return false
}
