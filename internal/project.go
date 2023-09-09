package internal

import (
	"fmt"
	"projects-feed/pkg/project"
	"sort"
	"strings"
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

func GetProjects(vendor string, page int, tag string) (p []project.Project, err error) {
	key := fmt.Sprintf("vendor%spage%dtag%s", vendor, page, tag)

	cached, exists := memcache.Get(key)
	if exists {
		return cached.([]project.Project), nil
	}

	var vendors []string
	if vendor != "" {
		names := strings.Split(vendor, ".")
		if len(names) >= 1 {
			vendors = append(vendors, names[0])
		} else {
			vendors = append(vendors, vendor)
		}
	} else {
		vendors = []string{"sribu", "projects"}
	}

	var wg sync.WaitGroup
	for _, v := range vendors {
		wg.Add(1)
		go func(name string) {
			vendor := project.New(name)
			if vendor == nil {
				wg.Done()

				return
			}

			latest, err := vendor.GetProjects(page, tag)
			if err == nil {
				p = append(p, latest...)
			}

			wg.Done()
		}(v)
	}

	wg.Wait()

	// sort datetime from different vendor
	sort.Slice(p, func(i, j int) bool {
		return p[j].PublishedAt.Before(p[i].PublishedAt)
	})

	memcache.SetDefault(key, p)

	return
}
