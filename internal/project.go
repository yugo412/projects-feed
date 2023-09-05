package internal

import (
	"fmt"
	"projects-feed/pkg/projects"
	"sort"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const CacheExpiration = 15 // in minutes

var memcache *cache.Cache

func init() {
	memcache = cache.New(
		time.Minute*CacheExpiration,
		time.Minute*(CacheExpiration+5),
	)
}

func GetProjects(page int, tag string) (p []projects.Project, err error) {
	key := fmt.Sprintf("page%dtag%s", page, tag)

	cached, exists := memcache.Get(key)
	if exists {
		return cached.([]projects.Project), nil
	}

	vendors := []string{"projects"}

	var wg sync.WaitGroup
	for _, v := range vendors {
		wg.Add(1)
		go func() {
			vendor := projects.New(v)
			latest, err := vendor.GetProjects(page, tag)
			if err == nil {
				p = append(p, latest...)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	// sort datetime from different vendor
	sort.Slice(p, func(i, j int) bool {
		return p[j].PublishedAt.Before(p[i].PublishedAt)
	})

	memcache.SetDefault(key, p)

	return
}
