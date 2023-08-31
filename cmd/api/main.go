package main

import (
	"fmt"
	"log"
	"projects-rss/cmd/api/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

import (
	"net/http"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	c := chi.NewRouter()

	c.Get("/projects", handlers.GetProjects)
	c.Get("/projects/{type}", handlers.GetProjectsFeed)
	c.Get("/projects/go", handlers.RedirectProject)

	if viper.GetString("ENV") != "production" {
		log.Println("Running in port:", viper.GetInt("APP_PORT"))
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("APP_PORT")), c)
	if err != nil {
		panic(err)
	}
}
