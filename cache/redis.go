package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache redis cache结构体
type redisCache struct {
	client *redis.Client
	op     Options
}

func newRedisCache(opts ...Option) *redisCache {
	cache := &redisCache{}
	for _, o := range opts {
		o(&cache.op)
	}

	cache.client = redis.NewClient(&redis.Options{
		Addr:         cache.op.Endpoint,
		Password:     cache.op.Password,
		DB:           cache.op.Db,
		PoolSize:     cache.op.PoolSize,
		MinIdleConns: cache.op.MinIdle,
	})

	return cache
}

func (p *redisCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	return p.client.Set(ctx, key, val, expiration).Err()
}

func (p *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := p.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (p *redisCache) GetDel(ctx context.Context, key string) (interface{}, error) {
	val, err := p.client.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (p *redisCache) Scan(ctx context.Context, key string, val interface{}) error {
	if val == nil {
		return nil
	}
	return p.client.Get(ctx, key).Scan(val)
}

func (p *redisCache) Delete(ctx context.Context, keys ...string) error {
	return p.client.Del(ctx, keys...).Err()
}

func (p *redisCache) Options() Options {
	return p.op
}

func (p *redisCache) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}

func (p *redisCache) Close(ctx context.Context) error {
	return p.client.Close()
}
