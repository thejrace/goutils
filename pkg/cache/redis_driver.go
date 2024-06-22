package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	Client *redis.Client
}

func (d *RedisDriver) Init(options ...interface{}) error {
	d.Client = options[0].(*redis.Client)

	return nil
}

func (d *RedisDriver) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := d.Client.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("Cache key %s not found", key)
	} else if err != nil {
		return nil, err // Unexpected error
	} else {
		return val, nil
	}
}

func (d *RedisDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return d.Client.Set(ctx, key, value, ttl).Err()
}

func (d *RedisDriver) Delete(ctx context.Context, key string) error {
	return d.Client.Del(ctx, key).Err()
}

func (d *RedisDriver) Inc(ctx context.Context, key string, incBy int64) error {
	return d.Client.IncrBy(ctx, key, incBy).Err()
}

func (d *RedisDriver) Flush(ctx context.Context) error {
	d.Client.FlushDB(ctx)

	return nil
}
