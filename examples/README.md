# Examples

This directory contains example data files and usage examples for the Permutation API.

## Data Files

### CSV File (`data.csv`)

Example CSV file with employee data:
- Column 0: name
- Column 1: age
- Column 2: city
- Column 3: department

**Usage:**
```bash
# Generate permutations from names (column 0)
./bin/permutation-api --csv=examples/data.csv --column=0

# Generate permutations from cities (column 2)
./bin/permutation-api --csv=examples/data.csv --column=2
```

### TSV File (`data.tsv`)

Example TSV file with employee data (tab-separated).

**Usage:**
```bash
# Generate permutations from names
./bin/permutation-api --tsv=examples/data.tsv --column=0
```

## CLI Examples

### Basic Usage

```bash
# Simple permutation
./bin/permutation-api a b c

# Quiet mode (statistics only)
./bin/permutation-api --quiet 1 2 3 4 5

# From CSV file
./bin/permutation-api --csv=examples/data.csv --column=0

# From TSV file
./bin/permutation-api --tsv=examples/data.tsv --column=0
```

## API Examples

### cURL Examples

```bash
# Start the API server
./bin/permutation-api --serve

# GET request
curl "http://localhost:8080/?elements=a,b,c"

# POST request
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["a", "b", "c"]'

# Health check
curl http://localhost:8080/health
```

### JavaScript Example

```javascript
// Using fetch API
async function getPermutations(elements) {
  const response = await fetch(
    `http://localhost:8080/?elements=${elements.join(',')}`
  );
  return await response.json();
}

// Usage
const perms = await getPermutations(['x', 'y', 'z']);
console.log(perms);
```

### Python Example

```python
import requests

# GET request
response = requests.get(
    'http://localhost:8080/',
    params={'elements': 'a,b,c'}
)
permutations = response.json()

# POST request
response = requests.post(
    'http://localhost:8080/',
    json=['a', 'b', 'c']
)
permutations = response.json()
```

## GUI Configuration

```bash
# Start the GUI server
./bin/permutation-api --gui

# Open http://localhost:3000 in your browser
```

From the GUI you can:
- Configure server settings
- Set up database connections
- Save/load configuration files
- Modify application parameters

## Docker Examples

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

### Using Docker Directly

```bash
# Build image
docker build -t permutation-api .

# Run API server
docker run -p 8080:8080 permutation-api --serve

# Run GUI server
docker run -p 3000:3000 permutation-api --gui

# Run CLI mode
docker run permutation-api a b c
```

## Database Examples

### PostgreSQL Setup

```bash
# Start PostgreSQL with Docker
docker-compose up -d db

# Configure environment
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_DATABASE=permutation_db
export DB_TABLE=elements
export DB_COLUMN=value

# Application will automatically use database
./bin/permutation-api --serve
```

### MySQL Setup

```bash
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_DATABASE=permutation_db
export DB_TABLE=elements
export DB_COLUMN=value
```

### SQLite Setup

```bash
export DB_DRIVER=sqlite3
export DB_DATABASE=./permutations.db
export DB_TABLE=elements
export DB_COLUMN=value
```

## Makefile Examples

```bash
# Build the application
make build

# Run tests
make test

# Generate coverage report
make coverage

# Run API server
make run-api

# Run GUI server
make run-gui

# Build Docker image
make docker-build

# Start all services
make docker-up

# Clean build artifacts
make clean
```

## Performance Testing

### Benchmarking Different Sizes

```bash
# 3 elements (6 permutations)
time ./bin/permutation-api --quiet a b c

# 5 elements (120 permutations)
time ./bin/permutation-api --quiet 1 2 3 4 5

# 8 elements (40,320 permutations)
time ./bin/permutation-api --quiet 1 2 3 4 5 6 7 8

# 10 elements (3,628,800 permutations)
time ./bin/permutation-api --quiet 1 2 3 4 5 6 7 8 9 10
```

### Load Testing the API

Using `ab` (Apache Bench):

```bash
# Install ab
sudo apt-get install apache2-utils

# Test with 1000 requests, 10 concurrent
ab -n 1000 -c 10 "http://localhost:8080/?elements=a,b,c"
```

Using `wrk`:

```bash
# Install wrk
sudo apt-get install wrk

# Test for 30 seconds with 10 connections
wrk -t10 -c10 -d30s "http://localhost:8080/?elements=a,b,c"
```
