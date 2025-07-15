package cache

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

const cacheNamespace = "svc"

var rdb *redis.Client

type SomeKey struct {
	k string
}

func (k SomeKey) String() string {
	return k.k
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("redis", "8.0.3", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		ctx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()

		return rdb.Ping(ctx).Err()
	}); err != nil {
		log.Fatalf("Could not connect to redis: %s", err)
	}

	// as of go1.15 testing.M returns the exit code of m.Run(), so it is safe to use defer here
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func newRedisTestCacheManager() *RedisCacheManager {
	manager, _ := NewRedisCacheManager(cacheNamespace, rdb)

	return manager
}

func TestCacheManager_Get(t *testing.T) {
	manager := newRedisTestCacheManager()

	ctx := context.Background()
	key := SomeKey{k: "exists_key"}
	want := "value"

	err := manager.Set(ctx, key, want, time.Minute)
	require.NoError(t, err, "error on set cache value")

	got, err := rdb.Get(ctx, wrapperKeyValue(key)).Result()
	require.NoError(t, err, "error on get value from redis")
	require.Equal(t, want, got, "value from redis is not equal to value from cache")
}

func TestCacheManager_GetNotExists(t *testing.T) {
	manager := newRedisTestCacheManager()

	ctx := context.Background()
	key := SomeKey{k: "not_exists_key"}

	_, err := manager.Get(ctx, key)
	require.ErrorIs(t, err, ErrKeyNotExists, "error is not ErrNotExists")
}

func TestCacheManager_GetSuccess(t *testing.T) {
	manager := newRedisTestCacheManager()

	ctx := context.Background()
	key := SomeKey{k: "exists_key"}
	want := "value"

	rdb.Set(ctx, wrapperKeyValue(key), want, time.Minute)

	got, err := manager.Get(ctx, key)
	require.NoError(t, err, "on success get error")
	require.Equal(t, got, want)
}

func TestCacheManager_DelSuccess(t *testing.T) {
	manager := newRedisTestCacheManager()

	ctx := context.Background()
	key := SomeKey{k: "del_key"}

	err := manager.Del(ctx, key)
	require.NoError(t, err, "on success get error")
}

func wrapperKeyValue(key Key) string {
	return newKeyWrapper(cacheNamespace, key).String()
}
