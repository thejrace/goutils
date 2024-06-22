package cache

import (
	"context"
	"time"
)

type CacheDriverContract interface {
	Init(options ...interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Inc(ctx context.Context, key string, incBy int64) error
	Flush(ctx context.Context) error
}

func GetRedisDriver() CacheDriverContract {
	return &RedisDriver{}
}

func GetInMemoryDriver() CacheDriverContract {
	return &InMemoryDriver{}
}
