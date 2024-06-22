package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisDriver(t *testing.T) {
	t.Run("it sets key", func(t *testing.T) {
		key := "meditopia"
		expectedVal := "test123"

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Set(ctx, key, expectedVal, time.Second*10)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		val, err := driver.Get(ctx, key)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		if val == nil {
			t.Errorf("Get failed, expected %s got nil", expectedVal)
		}

		if val.(string) != expectedVal {
			t.Errorf("Get failed, expected %s, got %s", expectedVal, val.(string))
		}

		driver.Flush(ctx)
	})

	t.Run("it returns nil if key is not found", func(t *testing.T) {
		key := "meditopia"

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		val, err := driver.Get(ctx, key)

		if err == nil {
			t.Errorf("Expecting error got nil")
		}

		if val != nil {
			t.Errorf("Get failed, expected nil got %q", val)
		}

		driver.Flush(ctx)
	})

	t.Run("it sets with ttl", func(t *testing.T) {
		key := "meditopia"

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Set(ctx, key, "expire", time.Millisecond*100)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		time.Sleep(time.Millisecond * 300)

		val, err := driver.Get(ctx, key)

		if err == nil {
			t.Errorf("Expecting error got nil")
		}

		if val != nil {
			t.Errorf("Get failed, expected nil got %q", val)
		}

		driver.Flush(ctx)
	})

	t.Run("it increments existing value", func(t *testing.T) {
		key := "meditopia"
		existingVal := 10
		expectedVal := "11" // Redis returns values as strings

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Set(ctx, key, existingVal, time.Second*10)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		err = driver.Inc(ctx, key, 1)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		val, _ := driver.Get(ctx, key)

		if val.(string) != expectedVal {
			t.Errorf("Get failed, expected %s, got %s", expectedVal, val.(string))
		}

		driver.Flush(ctx)
	})

	t.Run("it increments non existing value", func(t *testing.T) {
		key := "meditopia"
		expectedVal := "1" // Redis returns values as strings

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Inc(ctx, key, 1)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		val, _ := driver.Get(ctx, key)

		if val.(string) != expectedVal {
			t.Errorf("Get failed, expected %s, got %s", expectedVal, val.(string))
		}

		driver.Flush(ctx)
	})

	t.Run("it deletes existing key from cache", func(t *testing.T) {
		key := "meditopia"

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Set(ctx, key, "med", time.Second*10)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		err = driver.Delete(ctx, key)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		val, _ := driver.Get(ctx, key)

		if val != nil {
			t.Errorf("Get failed, expected nil, got %s", val.(string))
		}

		driver.Flush(ctx)
	})

	t.Run("it attempts to delete nonexisting key from cache without errors", func(t *testing.T) {
		key := "meditopia"

		driver := GetRedisDriver()
		err := driver.Init(getClient())

		if err != nil {
			t.Errorf("Unexpected connection error occurred %s", err.Error())
		}

		ctx := context.Background()

		err = driver.Delete(ctx, key)

		if err != nil {
			t.Errorf("Unexpected error occurred %s", err.Error())
		}

		val, _ := driver.Get(ctx, key)

		if val != nil {
			t.Errorf("Get failed, expected nil, got %s", val.(string))
		}

		driver.Flush(ctx)
	})
}

func getClient() *redis.Client {
	url := "redis://localhost:6379/0?protocol=3"
	opts, err := redis.ParseURL(url)

	if err != nil {
		panic(err)
	}

	return redis.NewClient(opts)
}
