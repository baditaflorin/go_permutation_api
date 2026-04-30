# API Documentation

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API does not require authentication. Future versions may include API key or OAuth support.

## Endpoints

### Generate Permutations (GET)

Generate permutations from comma-separated query parameters.

**Endpoint:** `GET /`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| elements | string | Yes | Comma-separated list of elements to permute |

**Example Request:**
```bash
curl "http://localhost:8080/?elements=a,b,c"
```

**Example Response:**
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

**Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid input (too many elements, empty elements, missing parameter)
- `405 Method Not Allowed` - Unsupported HTTP method

---

### Generate Permutations (POST)

Generate permutations from JSON array in request body.

**Endpoint:** `POST /`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
["a", "b", "c"]
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["a", "b", "c"]'
```

**Example Response:**
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

**Status Codes:**
- `200 OK` - Success
- `400 Bad Request` - Invalid JSON, too many elements, or empty elements

---

### Health Check

Check if the service is running and healthy.

**Endpoint:** `GET /health`

**Example Request:**
```bash
curl http://localhost:8080/health
```

**Example Response:**
```json
{
  "status": "healthy"
}
```

**Status Codes:**
- `200 OK` - Service is healthy

---

## Error Handling

### Error Response Format

When an error occurs, the API returns a plain text error message with an appropriate HTTP status code.

**Example Error Response:**
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain

too many elements: maximum is 12, got 15
```

### Common Error Messages

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `too many elements: maximum is N, got M` | Request exceeds MAX_ELEMENTS limit | Reduce number of elements or increase MAX_ELEMENTS config |
| `element at index N is empty or contains only whitespace` | Empty or whitespace-only element | Remove empty elements from input |
| `'elements' query parameter is required` | Missing query parameter | Add `?elements=...` to URL |
| `invalid JSON body: ...` | Malformed JSON in POST body | Fix JSON syntax |

---

## Rate Limiting

Currently, there is no rate limiting. For production deployments, consider:
- Using a reverse proxy (nginx, Caddy) with rate limiting
- Implementing application-level rate limiting middleware
- Using API gateway services

---

## CORS

CORS is configurable via the `ENABLE_CORS` environment variable (default: enabled).

When enabled, the API returns these headers:
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

To disable CORS, set:
```bash
export ENABLE_CORS=false
```

---

## Performance Considerations

### Complexity

Permutation generation has factorial complexity:
- Time: O(n!)
- Space: O(n)

### Recommended Limits

| Elements | Permutations | Approx. Time | Use Case |
|----------|--------------|--------------|----------|
| 3 | 6 | < 1ms | Interactive |
| 5 | 120 | < 10ms | Real-time |
| 8 | 40,320 | < 100ms | API requests |
| 10 | 3,628,800 | ~1s | Batch processing |
| 11 | 39,916,800 | ~2.5s | Large batch |
| 12 | 479,001,600 | ~30s | Maximum recommended |

### Optimization Tips

1. **Streaming**: The API streams results, so memory usage remains constant
2. **Caching**: Consider caching results for repeated identical requests
3. **Async Processing**: For large inputs (>10 elements), use background jobs
4. **Database**: For very large datasets, consider storing results in a database

---

## Configuration

All configuration is managed via environment variables or config files. See the main [README.md](../README.md) for details.

### Key Settings

- `MAX_ELEMENTS`: Maximum number of elements (default: 12)
- `SERVER_PORT`: API server port (default: 8080)
- `ENABLE_CORS`: Enable CORS headers (default: true)
- `SHUTDOWN_TIMEOUT`: Graceful shutdown timeout in seconds (default: 5)

---

## Examples

### JavaScript (fetch)

```javascript
// GET request
const response = await fetch('http://localhost:8080/?elements=x,y,z');
const permutations = await response.json();
console.log(permutations);

// POST request
const response = await fetch('http://localhost:8080/', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(['x', 'y', 'z'])
});
const permutations = await response.json();
console.log(permutations);
```

### Python (requests)

```python
import requests

# GET request
response = requests.get('http://localhost:8080/', params={'elements': 'x,y,z'})
permutations = response.json()
print(permutations)

# POST request
response = requests.post(
    'http://localhost:8080/',
    json=['x', 'y', 'z']
)
permutations = response.json()
print(permutations)
```

### cURL

```bash
# GET request
curl "http://localhost:8080/?elements=x,y,z"

# POST request
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["x", "y", "z"]'

# Health check
curl http://localhost:8080/health
```

---

## Future Enhancements

Planned features for future versions:

- [ ] Authentication (API keys, OAuth)
- [ ] Rate limiting
- [ ] WebSocket support for real-time streaming
- [ ] Pagination for large result sets
- [ ] Result caching with TTL
- [ ] Async job processing for large inputs
- [ ] GraphQL endpoint
- [ ] Metrics and monitoring endpoints
- [ ] Request tracing and correlation IDs

---

## Support

For issues or questions:
- GitHub Issues: https://github.com/baditaflorin/go_permutation_api/issues
- Email: [Your contact]
