# Roadmap

## Phase 0: Project Foundation

### Goal

Define the project shape, API philosophy, and compatibility rules.

### Tasks

- Create repository
- Add README
- Add license
- Add Go module
- Define package naming
- Define semantic versioning policy
- Add contribution guidelines
- Add API stability policy
- Set up CI
- Add linting and formatting checks

### Deliverables

- GitHub repository
- Initial README
- go.mod
- CI workflow
- Project structure
- Initial design notes

### Success Criteria

- Repository is usable
- Contributors understand package goals
- Compatibility expectations are documented

## Phase 1: MVP Release

### Goal

Ship a simple, useful rate limiter middleware quickly.

### Tasks

- net/http middleware
- Fixed window rate limiting
- In-memory storage
- Basic configuration
- Default key extraction by client IP
- HTTP 429 Too Many Requests
- Basic rate limit headers
- Unit tests
- Basic examples

### Deliverables

- Algorithm: Fixed Window
- Backend: In-memory
- Middleware: net/http
- Dependencies: none
- Basic Headers
  - X-RateLimit-Limit
  - X-RateLimit-Remaining
  - X-RateLimit-Reset
  - Retry-After

### Success Criteria

- A user can install the package and protect an HTTP endpoint in under 5 minutes
- No external dependencies
- Works with net/http
- API surface remains small

## Phase 2: API Hardening

### Goal

Improve API design before declaring stability.

### Tasks

- Review exported types
- Minimize public API
- Add godoc comments
- Add examples for public APIs
- Ensure errors and results are extensible
- Add option-based configuration if needed
- Add validation for invalid configs
- Decide whether to keep Config struct, functional options, or both

### Deliverables

- Documented API
- Improved examples
- API review notes
- Compatibility policy
- Clear deprecation policy

### Success Criteria

- Public API feels idiomatic
- Internal implementation can change without breaking users
- Package is ready for broader adoption

## Phase 3: Algorithm Abstraction

### Goal

Prepare the package to support multiple rate limiting algorithms.

### Algorithms to Support

1. Fixed window
2. Token bucket
3. Sliding window
4. Leaky bucket <!-- optional -->

### Tasks

1. Introduce internal algorithm abstraction
2. Avoid making algorithm interfaces too broad too early
3. Add token bucket implementation
4. Add algorithm selection through configuration
5. Add tests comparing behavior of algorithms
6. Add benchmarks

### Deliverables

1. Algorithm abstraction
2. Token bucket support
3. Sliding window design proposal
4. Algorithm-specific tests
5. Benchmarks

### Success Criteria

- Users can choose an algorithm
- Existing users are not broken
- Fixed window remains the default behavior unless explicitly changed

## Phase 4: Storage Abstraction

### Goal

Prepare for multiple backends while keeping the core lightweight.

### Storage Backends

- Initial
  - In Memory
- Future
  - Redis
  - Memcached

### Storage Design Principles

- Core package should not import Redis clients
- Redis support should live in an optional package
- Storage interfaces should be minimal
- Storage should work with context.Context

### Deliverables

- Storage interface
- Memory store implementation
- Storage conformance tests
- Redis design document
- Backend extension guide

### Success Criteria

- Core package still has no external dependencies
- Custom storage can be implemented by users
- Existing in-memory usage remains unchanged

## Phase 5: Redis Backend

### Goal

Add distributed rate limiting support.

### Scope

- Redis-backed counters
- Atomic operations
- Expiry handling
- Context support
- Redis-specific tests
- Optional dependency

### Deliverables

- Redis store package
- Redis integration tests
- Failure behavior documentation
- Example distributed setup

### Success Criteria

- Multiple application instances can share rate limit state
- Redis package is optional
- Core package remains dependency-light

## Phase 6: Cooldown and Penalty Features

### Goal

Support temporary penalties for repeated violations.

### Deliverables

- Cooldown support
- Retry-After behavior
- Tests for cooldown expiry
- Documentation

### Success Criteria

- Cooldown behavior is predictable
- Existing users are unaffected unless cooldown is configured
- Headers reflect blocked duration properly

## Phase 7: Advanced Middleware Features

### Goal

Improve real-world usability.

### Features

- Custom response handler
- Custom key extractor
- Route-specific limits
- Method-specific limits
- Per-user limits
- IP-based limits
- Header-based limits
- Skip function for health checks
- Trust proxy configuration
- Structured result object

### Deliverables

- Custom key support
- Skip behavior
- Custom limit response
- Route-specific examples
- Middleware customization docs

### Success Criteria

- Middleware supports common production use cases
- API remains simple for basic users
- Advanced users can customize behavior without forking

## Phase 8: Observability and Instrumentation

### Goal

Make the package easier to operate in production.

### Features

Hooks for allow / deny decisions
Optional logging interface
Metrics callbacks
Expose result metadata
Debug-friendly behavior

### Avoid

Do not hard-depend on:

- Prometheus
- OpenTelemetry
- Zap
- Zerolog
- Logrus
Instead, expose hooks.

### Deliverables

- Hook system
- Metrics example
- Logging example
- Observability guide

### Success Criteria

- Users can integrate with their own observability stack
- Core package remains dependency-free
- No forced logging or metrics implementation

## Phase 9: v1.0.0 Stable Release

### Goal

Publish the first stable version.

### Deliverables

- v1.0.0 tag
- Release notes
- Migration notes
- Stable API documentation
- Compatibility promise

### Success Criteria

- Users can confidently depend on the package
- Future enhancements can be added without breaking existing users