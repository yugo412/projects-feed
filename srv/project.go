package srv

import (
	"fmt"
	"projects-feed/pkg/cache"
	"projects-feed/pkg/project"
	"sort"
	"strings"
	"sync"

	"github.com/gookit/slog"
)

func GetDetail(URL, vendor string) (p project.Project, err error) {
	vendors := strings.Split(vendor, ".")
	p, err = project.New(vendors[0]).GetDetail(URL)

	return
}

func GetProjects(vendor string, page int, tag string) (p []project.Project, err error) {
	key := fmt.Sprintf("projects_vendor%spage%dtag%s", vendor, page, tag)

	c, err := cache.New("memory")
	if err != nil {
		slog.Errorf("failed to initialize cache: %v", err)
	}

	if val, err := c.Get(key); err == nil && c != nil {
		return val.([]project.Project), nil
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

	if ok, _ := c.Set(key, p); !ok {
		slog.Errorf("Failed to set cache for \"%s\".", key)
	}

	return
}
