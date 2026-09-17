package memory

import (
	"fmt"
	"sync"
)

type Cache struct {
	cache sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) cacheKey(name string) string {
	return "cache:" + name
}

func (c *Cache) Get(name string) (*string, error) {
	cacheKey := c.cacheKey(name)

	cache, exists := c.cache.Load(cacheKey)
	if !exists {
		return nil, fmt.Errorf("cache for %s not exists", cacheKey)
	}

	value, ok := cache.(string)
	if !ok {
		return nil, fmt.Errorf("cache value type assertion failed for %s", cacheKey)
	}

	return &value, nil
}

func (c *Cache) Delete(name string) {
	cacheKey := c.cacheKey(name)

	_, exists := c.cache.Load(cacheKey)
	if !exists {
		return
	}

	c.cache.Delete(cacheKey)
}
