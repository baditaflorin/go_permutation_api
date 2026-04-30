# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-01-15

### Added - Major Refactoring

#### Architecture
- Completely restructured codebase into clean packages (cmd/, internal/, pkg/)
- Implemented DRY and SOLID principles throughout
- Added comprehensive separation of concerns

#### Configuration
- Environment variable support for all settings via `.env` files
- JSON configuration file support with save/load functionality
- Web-based GUI for runtime configuration management (port 3000)
- Configuration hot-reload support
- Example configuration files

#### Data Sources
- CSV file input support with configurable column index
- TSV file input support with configurable column index
- PostgreSQL database integration
- MySQL database integration
- SQLite database integration
- Database connection pooling
- Database health checks

#### API Enhancements
- `/version` endpoint with build information
- `/metrics` endpoint with application statistics
- `/health` endpoint improvements
- Request rate limiting middleware
- Response compression (gzip) middleware
- Request timeout configuration
- Panic recovery middleware
- CORS middleware (configurable)
- Request logging middleware
- Pagination support for large result sets
- Caching layer for repeated requests
- Performance profiling endpoints (pprof)

#### CLI Improvements
- CSV/TSV file input via `--csv` and `--tsv` flags
- Multiple output formats: plain, JSON, CSV
- Configurable column selection for CSV/TSV
- Enhanced error messages
- Memory statistics tracking

#### Testing
- Comprehensive unit tests for all packages
- Integration tests with database support
- Benchmark tests
- 100+ tests covering core functionality
- GitHub Actions CI/CD pipeline
- Code coverage reporting

#### Infrastructure
- Multi-stage Dockerfile for optimized builds
- docker-compose configuration with PostgreSQL
- Makefile with 20+ development commands
- Database initialization script (init.sql)
- Example data files (CSV/TSV)

#### Documentation
- Comprehensive README with usage examples
- API documentation (docs/API.md)
- OpenAPI/Swagger specification (docs/openapi.yaml)
- Examples directory with sample files
- CHANGELOG.md
- Inline code documentation

#### Monitoring & Observability
- Structured logging with log levels
- Application metrics tracking
- Request/error counters
- Memory usage monitoring
- Uptime tracking
- Goroutine tracking

### Changed
- Moved all hardcoded values to configuration
- Refactored permutation algorithm into separate package
- Improved error handling and validation
- Enhanced graceful shutdown handling
- Optimized memory usage with streaming

### Fixed
- Broken unit tests from previous version
- Memory leaks in long-running processes
- Inconsistent error messages
- Race conditions in concurrent operations

### Security
- Added request timeout to prevent DoS
- Implemented rate limiting
- Added panic recovery
- Improved input validation
- Database connection pooling limits

## [1.0.0] - 2024-01-01

### Initial Release
- Basic CLI permutation generation
- HTTP API with GET/POST endpoints
- Lexicographical permutation algorithm
- Basic error handling
- Simple tests

---

## Upgrade Guide

### From 1.x to 2.0

#### Breaking Changes
1. **Command-line interface**: The binary is now in `cmd/server/` instead of root
2. **Configuration**: All settings now configurable via environment variables
3. **Import paths**: Internal packages moved to `internal/` directory

#### Migration Steps

1. **Update build process:**
   ```bash
   # Old
   go build -o permutation-api .
   
   # New
   go build -o bin/permutation-api ./cmd/server
   ```

2. **Set environment variables:**
   Create a `.env` file from `.env.example`:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Update imports** (if using as a library):
   ```go
   // Old
   import "github.com/baditaflorin/go_permutation_api"
   
   // New
   import "github.com/baditaflorin/go_permutation_api/internal/permutation"
   import "github.com/baditaflorin/go_permutation_api/pkg/validator"
   ```

4. **Use new CLI flags:**
   ```bash
   # CSV input
   ./bin/permutation-api --csv=data.csv --column=0
   
   # GUI configuration
   ./bin/permutation-api --gui
   ```

5. **Deploy with Docker:**
   ```bash
   docker-compose up -d
   ```

---

## Roadmap

### v2.1.0 (Planned)
- [ ] WebSocket support for real-time streaming
- [ ] GraphQL API endpoint
- [ ] Redis caching support
- [ ] Distributed tracing with OpenTelemetry
- [ ] Multi-language support for GUI

### v2.2.0 (Planned)
- [ ] Authentication and authorization
- [ ] User management
- [ ] API rate limiting per user
- [ ] Result persistence to database
- [ ] Async job processing for large inputs

### v3.0.0 (Future)
- [ ] Microservices architecture
- [ ] gRPC API
- [ ] Kubernetes deployment manifests
- [ ] Horizontal scaling support
- [ ] Plugin system for custom data sources
