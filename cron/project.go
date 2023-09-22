package cron

import (
	"projects-feed/pkg/project"
	"sync"
)

func FetchProjects() {
	vendors := [2]string{"projects", "sribu"}

	var projects []project.Project
	var wg sync.WaitGroup
	for _, v := range vendors {
		wg.Add(1)
		go func(name string) {
			vendor := project.New(name)
			if vendor == nil {
				wg.Done()

				return
			}

			if latest, err := vendor.GetProjects(1, ""); err == nil {
				projects = append(projects, latest...)
			}

			wg.Done()
		}(v)
	}

	wg.Wait()

	return
}
