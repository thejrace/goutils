package cache

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryDriver(t *testing.T) {
	t.Run("it sets key", func(t *testing.T) {
		key := "meditopia"
		expectedVal := "inmemory"

		driver := GetInMemoryDriver()
		_ = driver.Init()

		ctx := context.Background()

		_ = driver.Set(ctx, key, expectedVal, time.Millisecond)

		val, _ := driver.Get(ctx, key)

		if val == nil {
			t.Errorf("Get failed, expected %s got nil", expectedVal)
		}

		if val.(string) != expectedVal {
			t.Errorf("Get failed, expected %s, got %s", expectedVal, val.(string))
		}

		driver.Flush(ctx)
	})

	t.Run("it increments value if exist", func(t *testing.T) {
		key := "meditopia"
		existingVal := 10
		expectedVal := 11

		driver := GetInMemoryDriver()
		_ = driver.Init()

		ctx := context.Background()

		_ = driver.Set(ctx, key, existingVal, time.Millisecond)

		_ = driver.Inc(ctx, key, 1)

		val, _ := driver.Get(ctx, key)

		if val == nil {
			t.Errorf("Get failed, expected %d got nil", expectedVal)
		}

		if val.(int) != expectedVal {
			t.Errorf("Get failed, expected %d, got %d", expectedVal, val.(int))
		}

		driver.Flush(ctx)
	})

	t.Run("it increments value as 1 if it does not exist", func(t *testing.T) {
		key := "meditopia"
		expectedVal := 1

		driver := GetInMemoryDriver()
		_ = driver.Init()

		ctx := context.Background()

		_ = driver.Inc(ctx, key, 1)

		val, _ := driver.Get(ctx, key)

		if val == nil {
			t.Errorf("Get failed, expected %d got nil", expectedVal)
		}

		if val.(int) != expectedVal {
			t.Errorf("Get failed, expected %d, got %d", expectedVal, val.(int))
		}

		driver.Flush(ctx)
	})

	t.Run("it deletes val from cache", func(t *testing.T) {
		key := "meditopia"
		cacheVal := 1

		driver := GetInMemoryDriver()
		_ = driver.Init()

		ctx := context.Background()

		_ = driver.Set(ctx, key, cacheVal, time.Millisecond)

		_ = driver.Delete(ctx, key)

		val, _ := driver.Get(ctx, key)

		if val != nil {
			t.Errorf("Delete failed, expected nil got %d", cacheVal)
		}

		driver.Flush(ctx)
	})

	t.Run("it deletes key from cache without errors if key does not exist", func(t *testing.T) {
		key := "meditopia"

		driver := GetInMemoryDriver()
		_ = driver.Init()

		ctx := context.Background()

		err := driver.Delete(ctx, key)

		if err != nil {
			t.Errorf("Delete failed, got error %s", err.Error())
		}

		driver.Flush(ctx)
	})
}
