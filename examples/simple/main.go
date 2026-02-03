package main

import (
	"log"

	"github.com/ferdiunal/panel.go/pkg/panel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cfg := panel.Config{
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Features: panel.FeatureConfig{
			Register:       true,
			ForgotPassword: false,
		},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "development", // Forces usage of embedded assets
		// Storage: panel.StorageConfig{
		// 	Path: "./storage/public",
		// 	URL:  "/storage",
		// },
	}

	app := panel.New(cfg)

	// You can register custom resources here
	// app.RegisterResource(MyResource)

	log.Println("Starting panel on http://localhost:8080")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
