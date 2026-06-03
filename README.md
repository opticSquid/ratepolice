# RatePolice

## Project Vision

Build a stable, lightweight, idiomatic Go rate limiter middleware package that works seamlessly with the Go standard library and can evolve over time to support multiple algorithms, storage backends, cooldown strategies, and advanced production use cases.

### Priorities of this Package

- Shipping a useful version quickly
- Maintaining a stable public API
- Keeping 0 to minimal dependencies
- Supporting net/http from day one
- Designing internals for future extensibility
- Preserving backwards compatibility across non-major releases

## Core Goals

### Immediate Goal

Release a minimal but production-usable rate limiter middleware for Go’s net/http.

### Long-Term Goal

Evolve the package into a flexible rate limiting toolkit supporting:

- Multiple rate limiting algorithms
- Multiple storage backends
- In-memory and distributed usage
- Cooldown / penalty behavior
- Custom key extraction
- Observability hooks
- Stable API across versions

## Design Principles

1. Standard Library First
2. Framework specific adapters can come later, but the core package should remain independent.
3. Minimal Dependencies
4. Stable Public API, initial API is intentionally small.
5. Do not expose any implementation that does not need to be exposed

## Proposed Package Scope

### In Scope

- HTTP middleware for net/http
- In-memory rate limiting
- Configurable limits
- Custom key generation
- Response headers
- Extensible algorithm design
- Extensible storage design
- Optional Redis backend later
- Cooldown behavior later
- Tests, benchmarks, and examples

### Out of Scope Initially

- Distributed rate limiting
- Redis
- SQL storage
- Framework-specific middleware
- Highly advanced adaptive throttling
- Complex quota billing systems
