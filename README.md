# Go Permutation API

A lightweight Go-based API and CLI for generating all permutations of an array. Designed to be fast, flexible, and production-ready.

## Features

- **CLI Mode**: Run directly from the terminal to compute permutations.
- **HTTP API**: Expose permutations as a REST API with GET and POST endpoints.
- **Functional Style**: Pure Go implementation of permutations.
- **Production-Ready**: Includes graceful shutdown, input validation, and structured logging.

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/baditaflorin/go_permutation_api.git
cd go_permutation_api
go build -o go_permutation_api
```

Usage
CLI Mode
Generate permutations directly from the command line:

```bash
./go_permutation_api a b c
```
Output:

```
1: a, b, c
2: a, c, b
3: b, a, c
4: b, c, a
5: c, a, b
6: c, b, a
```


HTTP API
Run the API server:

```shell
./go_permutation_api serve [port]
```

Endpoints
GET: /permutations?elements=a,b,c

```
curl "http://localhost:8080/permutations?elements=a,b,c"
```
Response:

```
[["a","b","c"],["a","c","b"],["b","a","c"],["b","c","a"],["c","a","b"],["c","b","a"]]
```
POST: /permutations
```
curl -X POST -H "Content-Type: application/json" -d '["a","b","c"]' http://localhost:8080/permutations
```
Response:
```
[["a","b","c"],["a","c","b"],["b","a","c"],["b","c","a"],["c","a","b"],["c","b","a"]]
```
Development
Run tests locally:
```
go test ./... -v
```

Contributing
Contributions are welcome! Please open issues or submit pull requests.

License
This project is licensed under the MIT License. See the LICENSE file for details.
