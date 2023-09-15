package cache

import (
	"strings"
	"time"
)

type Cache interface {
	Driver() string
	Set(string, any) (bool, error)
	SetFor(string, any, time.Duration) (bool, error)
	Get(string) (any, error)
}

var (
	memoryInstance *Memory
)

func New(driver string) Cache {
	driver = strings.ToLower(driver)

	switch driver {
	case "memory":
		// memory cache is always on singleton
		if memoryInstance == nil {
			memoryInstance = &Memory{}
		}

		return memoryInstance
	}

	return nil
}
