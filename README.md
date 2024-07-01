# KV Database API

A small key-value database API.

Run the service with:
```bash
make run
```

Coverage report:
```bash
make coverage
```

Thread-safety is insured with a mutex lock on the database.

Tests & Benchmarks are run with go's race detector.

Benchmarks:
```bash
make bench
```

## Usage
### GET ALL KEYS
```
GET {SERVICEADDR}:8080/
```
Returns a list of all keys in the database.

### GET VALUE
```
GET {SERVICEADDR}:8080/{KEY}
```
Returns the value of the key.

### SET VALUE
```
PUT {SERVICEADDR}:8080/{KEY}
```
Sets the value of the key to the request body. Supports both valid JSON and plain text.

**NOTE:** If the key already exists, the value will be overwritten.


### DELETE
```
DELETE {SERVICEADDR}:8080/{KEY}
```
Deletes the key from the database. 
Returns 404 if the key does not exist.
