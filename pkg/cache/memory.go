package cache

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type Memory struct{}

var c *cache.Cache

const DefaultCacheTime = 15

func init() {
	c = cache.New(
		time.Minute*DefaultCacheTime,
		time.Minute*(DefaultCacheTime+5),
	)
}

func (m *Memory) Driver() string {
	return "memory"
}

func (m *Memory) Set(key string, val any) (bool, error) {
	c.SetDefault(key, val)

	return true, nil
}

func (m *Memory) SetFor(key string, val any, exp time.Duration) (bool, error) {
	c.Set(key, val, exp)

	return true, nil
}

func (m *Memory) Get(key string) (val any, err error) {
	val, cached := c.Get(key)
	if !cached {
		return nil, fmt.Errorf("key %s has no value", key)
	}

	return val, nil
}

func (m *Memory) Items() map[string]interface{} {
	items := make(map[string]interface{})
	for k, v := range c.Items() {
		items[k] = v.Object
	}

	return items
}
