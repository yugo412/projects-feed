package main

import (
	"fmt"
	"net/http"
	"os"
	"projects-feed/cmd/api/router"
	"projects-feed/cron"
	"projects-feed/models"

	"github.com/gookit/slog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Init database connection, by default it uses SQLite
func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	env := os.Getenv("ENV")

	// auto migrate for non prod env
	if env != "production" {
		err = db.AutoMigrate(
			&models.Author{},
			&models.Vendor{},
			&models.Project{},
		)
		if err != nil {
			slog.Errorf("Failed to auto-migrate database: %v", err)
		}
	}
}

// Run every job that registered in cron/schedulers
func init() {
	cron.Run(cron.Cron{
		DB: db,
	})
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
