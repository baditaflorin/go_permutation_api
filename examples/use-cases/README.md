# Use Cases and Real-World Examples

This directory contains concrete, runnable examples for common use cases of the Go Permutation API.

## Use Case Index

| # | Use Case | Technique | Elements |
|---|----------|-----------|----------|
| 1 | [Password combination testing](#1-password-combination-testing) | CLI | Characters |
| 2 | [Route optimization (TSP)](#2-route-optimization) | API POST | City names |
| 3 | [A/B test variant generation](#3-ab-test-variant-generation) | API GET | Feature flags |
| 4 | [DNA sequence analysis](#4-dna-sequence-analysis) | CLI + CSV | ACGT bases |
| 5 | [Menu meal planning](#5-meal-planning-permutations) | WebSocket | Meal items |
| 6 | [Workflow step ordering](#6-workflow-step-ordering) | API + pagination | Steps |

---

## 1. Password Combination Testing

**Scenario:** A security team needs to enumerate all orderings of known password components to test their lockout policies.

### CLI

```bash
# Generate all orderings of 3 components
./bin/permutation-api --format=json Password 2024 Company! | jq '.[].
 | join("")' | head -10
```

### API

```bash
curl -s -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["Password", "2024", "Company!"]' | jq 'length'
# Output: 6
```

### Expected Output

```json
[
  ["Password", "2024", "Company!"],
  ["Password", "Company!", "2024"],
  ["2024", "Password", "Company!"],
  ["2024", "Company!", "Password"],
  ["Company!", "Password", "2024"],
  ["Company!", "2024", "Password"]
]
```

---

## 2. Route Optimization

**Scenario:** A delivery company wants to evaluate all possible orderings of 5 stops to find the optimal route.

### API Request

```bash
curl -s -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '["Warehouse", "Stop-A", "Stop-B", "Stop-C", "Stop-D"]' \
  | jq 'length'
# Output: 120 (5! routes to evaluate)
```

### Integration with a Distance Matrix

```python
import requests
import itertools

stops = ["Warehouse", "Stop-A", "Stop-B", "Stop-C", "Stop-D"]

# Get all permutations
resp = requests.post("http://localhost:8080/", json=stops)
routes = resp.json()

# Evaluate each route with your distance function
def total_distance(route, distance_matrix):
    return sum(distance_matrix[route[i]][route[i+1]] for i in range(len(route)-1))

# Find optimal
best_route = min(routes, key=lambda r: total_distance(r, your_distance_matrix))
print(f"Optimal route: {' → '.join(best_route)}")
```

---

## 3. A/B Test Variant Generation

**Scenario:** A product team wants to test all orderings of 4 UI features to find the best user flow.

### Get All Feature Orderings

```bash
curl -s "http://localhost:8080/?elements=signup,onboarding,dashboard,tutorial" \
  | jq '.[] | join(" → ")' | head -5
```

### Output

```
"dashboard → onboarding → signup → tutorial"
"dashboard → onboarding → tutorial → signup"
"dashboard → signup → onboarding → tutorial"
...
```

### Automated Variant Assignment

```javascript
const response = await fetch('http://localhost:8080/?elements=signup,onboarding,dashboard,tutorial');
const variants = await response.json(); // 24 variants

// Assign user to variant based on user ID modulo
function assignVariant(userId) {
  return variants[userId % variants.length];
}
```

---

## 4. DNA Sequence Analysis

**Scenario:** A bioinformatics researcher needs all possible orderings of 4 DNA base pairs.

### Via CSV File

```bash
# Create bases CSV
echo "base" > /tmp/bases.csv
echo "A" >> /tmp/bases.csv
echo "C" >> /tmp/bases.csv
echo "G" >> /tmp/bases.csv
echo "T" >> /tmp/bases.csv

# Generate all orderings (24 sequences)
./bin/permutation-api --csv=/tmp/bases.csv --column=0 --format=json \
  | jq '.[] | join("")'
```

### Output

```
"ACGT"
"ACTG"
"AGCT"
"AGTC"
"ATCG"
"ATGC"
...
```

### Filter for Palindromes

```bash
./bin/permutation-api --csv=/tmp/bases.csv --column=0 --format=json \
  | jq '.[] | select(join("") == (join("") | explode | reverse | implode)) | join("")'
```

---

## 5. Meal Planning Permutations

**Scenario:** A nutritionist wants to generate all weekly meal orderings for variety scoring.

### WebSocket Real-Time Streaming

```html
<!DOCTYPE html>
<html>
<script>
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  ws.send(JSON.stringify({
    action: 'start',
    elements: ['Salad', 'Pasta', 'Steak', 'Fish'],
    chunk_size: 6
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.type === 'chunk') {
    msg.data.forEach(meal => {
      console.log('Meal plan:', meal.join(' → '));
    });
  }
  if (msg.type === 'done') {
    console.log(`Generated ${msg.total} meal plans in ${msg.elapsed_ms}ms`);
    ws.close();
  }
};
</script>
</html>
```

### Stop After Finding a Preferred Combination

```javascript
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.type === 'chunk') {
    const preferred = msg.data.find(plan => plan[0] === 'Salad' && plan[4] === 'Steak');
    if (preferred) {
      ws.send(JSON.stringify({ action: 'stop' }));
      console.log('Found preferred plan:', preferred);
    }
  }
};
```

---

## 6. Workflow Step Ordering

**Scenario:** A DevOps team wants to evaluate all orderings of 6 deployment steps to find which order minimises rollback risk.

### API with Pagination (Large Result Sets)

```bash
# 6 elements = 720 permutations — use pagination
curl -s "http://localhost:8080/?elements=build,test,lint,deploy,notify,cleanup&page=1&per_page=50" \
  | jq '.data | length'
# Output: 50

# Get metadata
curl -s "http://localhost:8080/?elements=build,test,lint,deploy,notify,cleanup&page=1&per_page=50" \
  | jq '.meta'
```

### Paginated Output

```json
{
  "version": "1",
  "request_id": "a1b2c3d4e5f6",
  "data": [
    ["build", "cleanup", "deploy", "lint", "notify", "test"],
    ...
  ],
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 720,
    "total_pages": 15
  }
}
```

### Python Iterator Over All Pages

```python
import requests

def all_permutations(elements, base_url="http://localhost:8080/", per_page=100):
    page = 1
    while True:
        resp = requests.get(base_url, params={
            "elements": ",".join(elements),
            "page": page,
            "per_page": per_page
        })
        body = resp.json()
        yield from body["data"]
        if page >= body["meta"]["total_pages"]:
            break
        page += 1

steps = ["build", "test", "lint", "deploy", "notify", "cleanup"]
for plan in all_permutations(steps):
    score = evaluate_deployment_risk(plan)  # your function
    print(f"Risk {score:.2f}: {' → '.join(plan)}")
```

---

## Database Integration Example

Load permutation elements directly from a PostgreSQL table:

```bash
# Set up
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_TABLE=products
export DB_COLUMN=sku

# Run the API server — it will query DB for elements
./bin/permutation-api --serve
```

**SQL schema:**

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(255)
);

INSERT INTO products (sku, name) VALUES
    ('SKU-A', 'Widget Alpha'),
    ('SKU-B', 'Widget Beta'),
    ('SKU-C', 'Widget Gamma');
```

**Result:** The API will generate permutations of `['SKU-A', 'SKU-B', 'SKU-C']`.

---

## Performance Reference

| Elements | Permutations | API response time | Memory |
|----------|--------------|------------------|--------|
| 3 | 6 | < 1ms | ~2 KB |
| 5 | 120 | < 5ms | ~10 KB |
| 7 | 5,040 | < 50ms | ~400 KB |
| 9 | 362,880 | ~300ms | ~28 MB |
| 11 | 39,916,800 | ~2.5s | streaming |
| 12 | 479,001,600 | ~35s | streaming |

For inputs > 10 elements, use the **WebSocket endpoint** to stream results and cancel early.
