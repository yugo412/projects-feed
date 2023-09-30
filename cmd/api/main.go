package main

import (
	"fmt"
	"net/http"
	"os"
	"projects-feed/cmd/api/router"
	"projects-feed/cron"

	"github.com/gookit/slog"
)

// Run every job that registered in cron/schedulers
func init() {
	cron.Run()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		err := os.Setenv("PORT", port)
		if err != nil {
			slog.Errorf("Failed to set port %s: %v", port, err)
		}
	}

	env := os.Getenv("ENV")
	if env != "production" {
		slog.Infof("Running [%s] in port: %s", env, port)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router.RegisterRoutes())
	if err != nil {
		panic(err)
	}
}
