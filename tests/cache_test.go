package tests

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testbackend/internal/cache"
	"testbackend/internal/weather"
	"testing"
)

func TestCache(t *testing.T) {
	t.Parallel()

	testCache := cache.NewMapCache()

	t.Run("correctly stored value", func(t *testing.T) {
		t.Parallel()
		key := "some key"
		value := "some value"

		err := testCache.Set(key, value)
		assert.NoError(t, err)
		storedValue, err := testCache.Get(key)
		assert.NoError(t, err)

		assert.Equal(t, value, storedValue)
	})

	t.Run("no data races", func(t *testing.T) {
		t.Parallel()

		parallelFactor := 100_000
		emulateLoad(t, testCache, parallelFactor)
	})
}

func TestCacheInterface(t *testing.T) {
	c := cache.NewMapCache()
	key := "01-01-2022"
	value := weather.Weather{
		Description: weather.DescriptionModel{
			Main:        "A",
			Description: "B",
		},
		Metrics: weather.MetricsModel{
			TempMin:   1,
			TempMax:   2,
			Pressure:  3,
			Humidity:  4,
			MoonPhase: "C",
		},
		Name: "D",
		Icon: "E",
	}
	err := c.Set(key, value)
	assert.NoError(t, err)

	_, err = c.Get(key)
	assert.NoError(t, err)
	_, err = c.Get("not exist")
	assert.ErrorIs(t, err, cache.ErrorNotFound)
	
}

func emulateLoad(t *testing.T, c cache.Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("%d-key", i)
		value := fmt.Sprintf("%d-value", i)

		wg.Add(1)
		go func(k string) {
			err := c.Set(k, value)
			assert.NoError(t, err)
			wg.Done()
		}(key)

		wg.Add(1)
		go func(k, v string) {
			storedValue, err := c.Get(k)

			if !errors.Is(err, cache.ErrorNotFound) {
				assert.Equal(t, v, storedValue)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(k string) {
			err := c.Delete(k)
			assert.NoError(t, err)
			wg.Done()
		}(key)
	}
}
