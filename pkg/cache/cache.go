package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Cache interface {
	Driver() string
	Set(string, any) (bool, error)
	SetFor(string, any, time.Duration) (bool, error)
	Get(string) (any, error)
	Items() map[string]interface{}
}

var (
	once sync.Once

	memoryInstance *Memory
)

func New(driver string) (Cache, error) {
	driver = strings.ToLower(driver)

	switch driver {
	case "memory":
		// memory cache is always on singleton
		once.Do(func() {
			memoryInstance = &Memory{}
		})

		return memoryInstance, nil
	}

	return nil, fmt.Errorf("driver %s is not supported", driver)
}
