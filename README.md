# Go Permutation API

A production-ready, highly configurable Go application for generating lexicographical permutations. Supports CLI, HTTP API, web-based GUI configuration, and multiple data sources (CSV, TSV, database).

## Features

- **Multiple Modes**:
  - CLI mode for command-line permutation generation
  - HTTP API server for REST-based access
  - Web GUI for runtime configuration management
  - CSV/TSV file input support
  - Database integration (PostgreSQL, MySQL, SQLite)

- **Configurable**:
  - Environment variables for all settings
  - JSON configuration files
  - Runtime GUI-based configuration
  - No hardcoded values

- **Production-Ready**:
  - Graceful shutdown
  - Comprehensive input validation
  - Memory-efficient streaming
  - Request logging and optional CORS
  - Docker and docker-compose support

- **Well-Architected**:
  - Clean separation of concerns
  - DRY and SOLID principles
  - Comprehensive test coverage
  - Small, focused packages

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [CLI Mode](#cli-mode)
  - [API Server](#api-server)
  - [GUI Configuration](#gui-configuration)
  - [Data Sources](#data-sources)
- [Configuration](#configuration)
- [Docker](#docker)
- [Interactive Documentation & Examples](#interactive-documentation--examples)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Project Structure](#project-structure)

## Installation

### Prerequisites

- Go 1.23 or higher
- Make (optional, for using Makefile)
- Docker and docker-compose (optional, for containerized deployment)

### From Source

```bash
git clone https://github.com/baditaflorin/go_permutation_api.git
cd go_permutation_api
make build
```

Or manually:

```bash
go build -o bin/permutation-api ./cmd/server
```

### Using Docker

```bash
docker-compose up -d
```

## Quick Start

### CLI Mode

Generate permutations from command-line arguments:

```bash
./bin/permutation-api a b c
```

### API Server

Start the HTTP API server:

```bash
./bin/permutation-api --serve
```

Then access at `http://localhost:8080`

### GUI Configuration

Start the web-based configuration interface:

```bash
./bin/permutation-api --gui
```

Then open `http://localhost:3000` in your browser

## Usage

### CLI Mode

**Basic usage:**
```bash
./bin/permutation-api [options] element1 element2 element3 ...
```

**Options:**
- `--quiet`: Suppress output of individual permutations
- `--csv=<file>`: Load elements from CSV file
- `--tsv=<file>`: Load elements from TSV file
- `--column=<n>`: Column index for CSV/TSV (0-based, default: 0)

**Examples:**

```bash
# Generate permutations from arguments
./bin/permutation-api a b c

# Quiet mode (show only statistics)
./bin/permutation-api --quiet 1 2 3 4

# From CSV file
./bin/permutation-api --csv=data.csv --column=0

# From TSV file
./bin/permutation-api --tsv=data.tsv --column=1
```

### API Server

**Start the server:**
```bash
./bin/permutation-api --serve
```

**Endpoints:**

**GET /**: Generate permutations from query parameter
```bash
curl "http://localhost:8080/?elements=a,b,c"
```

**POST /**: Generate permutations from JSON body
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["a", "b", "c"]'
```

**GET /health**: Health check endpoint
```bash
curl http://localhost:8080/health
```

**Response format:**
```json
[
  ["a", "b", "c"],
  ["a", "c", "b"],
  ["b", "a", "c"],
  ["b", "c", "a"],
  ["c", "a", "b"],
  ["c", "b", "a"]
]
```

### GUI Configuration

Start the GUI server:

```bash
./bin/permutation-api --gui
```

Access the configuration interface at `http://localhost:3000`

Features:
- View and edit server configuration
- Modify application settings
- Configure database connections
- Save configuration to file
- Download configuration as JSON

### Data Sources

#### CSV Files

```bash
# Create a CSV file
echo "name,age,city" > data.csv
echo "Alice,30,NYC" >> data.csv
echo "Bob,25,LA" >> data.csv

# Generate permutations from first column
./bin/permutation-api --csv=data.csv --column=0
```

#### TSV Files

```bash
# Create a TSV file
printf "name\tage\tcity\n" > data.tsv
printf "Alice\t30\tNYC\n" >> data.tsv

# Generate permutations
./bin/permutation-api --tsv=data.tsv --column=0
```

#### Database

Configure database connection via environment variables:

```bash
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_DATABASE=permutation_db
export DB_TABLE=elements
export DB_COLUMN=value

# Application will use database automatically
```

## Configuration

### Environment Variables

Create a `.env` file from the example:

```bash
cp .env.example .env
```

**Server Configuration:**
- `SERVER_PORT`: API server port (default: 8080)
- `SERVER_HOST`: Server host (default: localhost)
- `GUI_PORT`: GUI server port (default: 3000)
- `SHUTDOWN_TIMEOUT`: Graceful shutdown timeout in seconds (default: 5)

**Application Settings:**
- `MAX_ELEMENTS`: Maximum number of elements allowed (default: 12)
- `MEMORY_STATS_FREQ`: Memory statistics update frequency (default: 1000)
- `QUIET`: Suppress CLI output (default: false)
- `ENABLE_CORS`: Enable CORS headers (default: true)

**Database Configuration:**
- `DB_DRIVER`: Database driver (postgres, mysql, sqlite3)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USERNAME`: Database username
- `DB_PASSWORD`: Database password
- `DB_DATABASE`: Database name
- `DB_TABLE`: Table name to query
- `DB_COLUMN`: Column name containing elements
- `DB_SSL_MODE`: SSL mode for PostgreSQL (default: disable)

### Configuration File

Use a JSON configuration file:

```bash
export CONFIG_FILE=configs/example.config.json
./bin/permutation-api --serve
```

See `configs/example.config.json` for the format.

## Docker

### Build Image

```bash
make docker-build
```

Or:

```bash
docker build -t permutation-api .
```

### Run with Docker Compose

Start all services (API, GUI, database):

```bash
make docker-up
```

Or:

```bash
docker-compose up -d
```

Services:
- API: http://localhost:8080
- GUI: http://localhost:3000
- PostgreSQL: localhost:5432

### Stop Services

```bash
make docker-down
```

## Interactive Documentation & Examples

A complete interactive documentation site is available in the `docs/` directory with:

- **Landing Page** (`docs/index.html`): Quick-start guide with live demo widget
- **Interactive Examples** (`docs/examples/index.html`): 7 runnable examples you can execute directly in your browser
  - Basic GET/POST permutations
  - Password combination testing
  - Route optimization (TSP)
  - A/B test variant generation
  - DNA sequence analysis
  - WebSocket streaming with progress
  - WebSocket with real-time cancellation
- **OpenAPI Spec** (`docs/openapi/index.html`): Full Swagger UI viewer
- **Architecture Decisions** (`docs/adr/index.html`): Six Thinking Hats ADR index

### View Documentation Locally

Serve the docs with any static file server:

```bash
# Python
cd docs && python3 -m http.server 8000

# Node.js
npx http-server docs -p 8000

# Go
cd docs && go run -mod=mod github.com/shurcooL/goexec@latest 'http.ListenAndServe(":8000", http.FileServer(http.Dir(".")))'
```

Then open `http://localhost:8000` in your browser.

**Note:** The interactive examples require a running API server. Start one with:
```bash
docker run -p 8080:8080 -e SERVER_HOST=0.0.0.0 ghcr.io/baditaflorin/go_permutation_api --serve
```

### GitHub Pages

To enable the hosted version at `https://baditaflorin.github.io/go_permutation_api/`:

1. Go to **Settings → Pages** in your GitHub repository
2. Under **Source**, select **GitHub Actions**
3. Merge this branch to `main`

The `.github/workflows/pages.yml` workflow will deploy automatically.

## API Documentation

### Endpoints

#### `GET /`
Generate permutations from comma-separated query parameter.

**Parameters:**
- `elements` (required): Comma-separated list of elements

**Example:**
```bash
curl "http://localhost:8080/?elements=x,y,z"
```

#### `POST /`
Generate permutations from JSON array in request body.

**Body:**
```json
["x", "y", "z"]
```

**Example:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["x", "y", "z"]'
```

#### `GET /health`
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

### Error Responses

**400 Bad Request:**
- Too many elements (exceeds MAX_ELEMENTS)
- Empty or whitespace-only elements
- Missing required parameters

**405 Method Not Allowed:**
- Unsupported HTTP method

## Development

### Makefile Commands

```bash
make help          # Display all available commands
make build         # Build the binary
make test          # Run all tests
make coverage      # Generate coverage report
make run           # Run in CLI mode
make run-api       # Run API server
make run-gui       # Run GUI server
make fmt           # Format code
make vet           # Run go vet
make clean         # Remove build artifacts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./internal/permutation/... -v
```

### Code Structure

The project follows Go best practices with clear separation of concerns:

- `cmd/server/`: Application entry point
- `internal/`: Private application code
  - `api/`: HTTP API server and handlers
  - `config/`: Configuration management
  - `datasource/`: CSV, TSV, and database data sources
  - `gui/`: Web-based configuration interface
  - `permutation/`: Core permutation algorithm
- `pkg/`: Public reusable packages
  - `validator/`: Input validation

## Project Structure

```
go_permutation_api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── server.go            # HTTP server setup
│   │   ├── handlers.go          # Request handlers
│   │   ├── middleware.go        # HTTP middleware
│   │   └── api_test.go          # API tests
│   ├── config/
│   │   ├── config.go            # Configuration structs
│   │   └── config_test.go       # Config tests
│   ├── datasource/
│   │   ├── datasource.go        # DataSource interface
│   │   ├── csv.go               # CSV reader
│   │   ├── tsv.go               # TSV reader
│   │   ├── database.go          # Database connector
│   │   └── datasource_test.go   # DataSource tests
│   ├── gui/
│   │   ├── server.go            # GUI server
│   │   └── handlers.go          # GUI handlers
│   └── permutation/
│       ├── generator.go         # Permutation algorithm
│       └── generator_test.go    # Algorithm tests
├── pkg/
│   └── validator/
│       ├── validator.go         # Input validation
│       └── validator_test.go    # Validation tests
├── configs/
│   └── example.config.json      # Example configuration
├── .env.example                 # Example environment variables
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## Performance

The application uses an efficient in-place permutation algorithm:

- **Time Complexity**: O(n!) - generates all n! permutations
- **Space Complexity**: O(n) - uses minimal extra memory
- **Streaming**: Permutations are generated and sent one at a time
- **Memory Tracking**: CLI mode tracks and reports memory usage

**Benchmarks:**
- ~2.5 seconds for 11 elements (39,916,800 permutations)
- Memory-efficient: minimal allocations per permutation

## Limitations

- Maximum elements configurable via `MAX_ELEMENTS` (default: 12)
- Factorial growth: 12! = 479,001,600 permutations
- For very large inputs, consider streaming to file or database

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure all tests pass: `make test`
5. Format your code: `make fmt`
6. Submit a pull request

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Submit a pull request
- Contact: [Your contact information]

## Acknowledgments

- Built with Go standard library (no external dependencies for core functionality)
- Database drivers: lib/pq, go-sql-driver/mysql, mattn/go-sqlite3
- Inspired by efficient permutation algorithms

---

**Made with ❤️ in Go**
