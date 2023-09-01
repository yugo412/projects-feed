package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projects-feed/cmd/api/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		os.Setenv("PORT", port)
	}

	env := os.Getenv("ENV")
	if env != "production" {
		log.Printf("Running [%s] in port: %s", env, port)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router.RegisterRoutes())
	if err != nil {
		panic(err)
	}
}
