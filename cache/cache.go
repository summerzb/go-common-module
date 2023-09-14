package cache

import (
	"context"
	"time"
)

// Cache 定义cache驱动接口
type Cache interface {
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	GetDel(ctx context.Context, key string) (interface{}, error)
	Scan(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, keys ...string) error
	Options() Options
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

// NewCache returns a new cache.
func NewCache(opts ...Option) Cache {
	return newRedisCache(opts...)
}
