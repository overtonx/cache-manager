package cache

import "errors"

var (
	ErrKeyNotExists      = errors.New("keyString does not exist")
	ErrCacheNotAvailable = errors.New("cache is not available")
	ErrCacheReadTimeout  = errors.New("read timeout")
)
