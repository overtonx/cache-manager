package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheManager struct {
	rdb *redis.Client
	ns  string
}

func NewRedisCacheManager(namespace string, rdb *redis.Client) (*RedisCacheManager, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is empty")
	}

	return &RedisCacheManager{
		rdb: rdb,
		ns:  namespace,
	}, nil
}

func (cm *RedisCacheManager) Get(ctx context.Context, key Key) (string, error) {
	res, err := cm.rdb.Get(ctx, cm.keyString(key)).Result()
	slog.Debug("get result", "res", res, "err", err)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrKeyNotExists
		}

		return "", fmt.Errorf("get error: %w", err)
	}

	return res, nil
}

func (cm *RedisCacheManager) Set(ctx context.Context, key Key, value string, ttl time.Duration) error {
	res := cm.rdb.Set(ctx, cm.keyString(key), value, ttl)
	slog.Debug("set result", "res", res)

	return fmt.Errorf("set error: %w", res.Err())
}

func (cm *RedisCacheManager) Del(ctx context.Context, key Key) error {
	cnt, err := cm.rdb.Del(ctx, cm.keyString(key)).Result()
	slog.Debug("del result", "deleted count", cnt, "err", err)

	if err != nil {
		return fmt.Errorf("del error: %w", err)
	}

	slog.Debug("del success", "keyString", key)

	return nil
}

func (cm *RedisCacheManager) keyString(k Key) string {
	return newKeyWrapper(cm.ns, k).String()
}
