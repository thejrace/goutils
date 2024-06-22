package cache

import (
	"context"
	"fmt"
	"time"
)

var cache map[string]interface{} = map[string]interface{}{}

type InMemoryDriver struct {
}

func (d *InMemoryDriver) Init(options ...interface{}) error {
	return nil
}

func (d *InMemoryDriver) Get(ctx context.Context, key string) (interface{}, error) {
	val, ok := cache[key]

	if !ok {
		return nil, fmt.Errorf("Cache key %s not found", key)
	}

	return val, nil
}

func (d *InMemoryDriver) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	cache[key] = value

	return nil
}

func (d *InMemoryDriver) Inc(ctx context.Context, key string, incBy int64) error {
	val, _ := d.Get(ctx, key)
	result := incBy

	if val != nil {
		result = val.(int64) + incBy
	}

	cache[key] = result

	return nil
}

func (d *InMemoryDriver) Delete(ctx context.Context, key string) error {
	delete(cache, key)

	return nil
}

func (d *InMemoryDriver) Flush(ctx context.Context) error {
	cache = map[string]interface{}{}

	return nil
}
