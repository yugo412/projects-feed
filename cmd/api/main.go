package main

import (
	"fmt"
	"net/http"
	"os"
	"projects-feed/cmd/api/router"
	"projects-feed/cron"

	"github.com/gookit/slog"
)

func init() {
	cron.Run()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		os.Setenv("PORT", port)
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
