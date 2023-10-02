package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// test existing cache
	driver := "memory"
	c, err := New(driver)
	if err != nil {
		t.Errorf("Failed to initialize cache: %v.", err)
	}

	if c != nil && c.Driver() != driver {
		t.Errorf("Expected driver name is %s, %s returned.", driver, c.Driver())
	}

	// test invalid cache
	driver = "memcached"
	_, err = New(driver)
	if err == nil {
		t.Errorf("Expected error for invalid driver, got nil.")

		assert.Equal(t, fmt.Sprintf("driver %s is not supported", driver), err.Error())
	}
}

func TestSetCache(t *testing.T) {
	driver := "memory"
	c, err := New(driver)
	if err != nil {
		t.Errorf("Failed to initialize cache: %v.", err)
	}

	// store the cache first
	key := "key"
	_, err = c.Set(key, "value")
	if err != nil {
		t.Errorf("Failed to set cache %s: %v", key, err)
	}

	// fetch from previous cache
	val, err := c.Get(key)
	if err != nil {
		t.Errorf("Failed to get cache %s: %v", key, err)
	}

	if val != nil && val.(string) != "value" {
		t.Errorf("Expected value is %s, got %s", "value", val)
	}
}

func TestSetExpirationCache(t *testing.T) {
	driver := "memory"
	c, err := New(driver)
	if err != nil {
		t.Errorf("Failed to initialize cache: %v.", err)
	}

	_, err = c.SetFor("key", "value", 1*time.Second)
	if err != nil {
		t.Errorf("Failed to set expiration cache: %v.", err)
	}

	time.Sleep(2 * time.Second)

	_, err = c.Get("key")
	if err == nil {
		t.Errorf("Failed to get expiration cache: %v.", err)
	}

	assert.Equal(t, "key \"key\" has no value", err.Error())
}
