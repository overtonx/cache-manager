# Cache Manager

A Redis-based cache manager for Go applications with namespace support and comprehensive error handling.

## Features

- Redis-based caching with namespace support
- Context-aware operations
- Comprehensive error handling
- TTL (Time To Live) support
- Structured logging with slog
- Full test coverage

## Installation

```bash
go get github.com/overtonx/cache-manager
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    "github.com/redis/go-redis/v9"
    "github.com/overtonx/cache-manager"
)

// Implement the Key interface for your cache keys
type UserKey struct {
    ID string
}

func (k UserKey) String() string {
    return fmt.Sprintf("user:%s", k.ID)
}

func main() {
    // Initialize Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create cache manager with namespace
    cm, err := cache.NewRedisCacheManager("myapp", rdb)
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    key := UserKey{ID: "123"}
    
    // Set a value with TTL
    err = cm.Set(ctx, key, "user data", 5*time.Minute)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get a value
    value, err := cm.Get(ctx, key)
    if err != nil {
        if errors.Is(err, cache.ErrKeyNotExists) {
            log.Println("Key not found")
        } else {
            log.Fatal(err)
        }
    }
    
    // Delete a key
    err = cm.Del(ctx, key)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Key Interface

All cache keys must implement the `Key` interface:

```go
type Key interface {
    String() string
}
```

Keys are automatically namespaced using the format `{namespace}:{key}`.

## API Reference

### RedisCacheManager

#### `NewRedisCacheManager(namespace string, rdb *redis.Client) (*RedisCacheManager, error)`

Creates a new cache manager with the specified namespace and Redis client.

#### `Get(ctx context.Context, key Key) (string, error)`

Retrieves a value from the cache. Returns `ErrKeyNotExists` if the key doesn't exist.

#### `Set(ctx context.Context, key Key, value string, ttl time.Duration) error`

Sets a value in the cache with the specified TTL.

#### `Del(ctx context.Context, key Key) error`

Deletes a key from the cache.

## Error Handling

The package provides the following error types:

- `ErrKeyNotExists`: Returned when a key doesn't exist in the cache
- `ErrCacheNotAvailable`: Returned when the cache is not available
- `ErrCacheReadTimeout`: Returned when a read operation times out

## Testing

Run the tests with:

```bash
go test -v
```

For integration tests with Redis:

```bash
cd local
docker-compose up -d
cd ..
go test -v
```

## Dependencies

- [redis/go-redis/v9](https://github.com/redis/go-redis) - Redis client for Go
- [stretchr/testify](https://github.com/stretchr/testify) - Testing toolkit
- [ory/dockertest/v3](https://github.com/ory/dockertest) - Integration testing with Docker

## Development

### Local Development

1. Start Redis using Docker Compose:
   ```bash
   cd local
   docker-compose up -d
   ```

2. Run tests:
   ```bash
   go test -v
   ```

### Requirements

- Go 1.24.5 or later
- Redis server for testing

## License

[Add your license information here]

## Contributing

[Add contributing guidelines here]