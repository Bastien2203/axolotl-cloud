package cache

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type TypedCache[T any] struct {
	mu   sync.RWMutex
	data T
	ttl  time.Duration
	ts   time.Time
	sfg  singleflight.Group
	load func(...any) (T, error)
}

func NewTypedCache[T any](
	initial T,
	ttl time.Duration,
	loader func(...any) (T, error),
) *TypedCache[T] {
	return &TypedCache[T]{
		data: initial,
		ttl:  ttl,
		load: loader,
	}
}

func (c *TypedCache[T]) Get() (T, error) {
	c.mu.RLock()
	valid := !c.expired()
	data := c.data
	c.mu.RUnlock()

	if valid {
		return data, nil
	}

	return c.refresh()
}

func (c *TypedCache[T]) refresh() (T, error) {
	v, err, _ := c.sfg.Do("refresh", func() (any, error) {
		newData, err := c.load()
		if err != nil {
			return c.getStale(), err
		}

		c.mu.Lock()
		c.data = newData
		c.ts = time.Now()
		c.mu.Unlock()

		return newData, nil
	})

	return v.(T), err
}

func (c *TypedCache[T]) expired() bool {
	return time.Since(c.ts) > c.ttl
}

func (c *TypedCache[T]) getStale() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}

func (c *TypedCache[T]) Update(updateFn func(T) T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = updateFn(c.data)
	c.ts = time.Now()
}
