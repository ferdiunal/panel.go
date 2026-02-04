package main

import (
	"log"

	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceSetting "github.com/ferdiunal/panel.go/pkg/resource/setting"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
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
		Permissions: panel.PermissionConfig{
			Path: "examples/simple/permissions.toml",
		},
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
			resourceSession.NewSessionResource(),
			resourceVerification.NewVerificationResource(),
			resourceSetting.NewSettingResource(),
		},
	}

	app := panel.New(cfg)

	// You can register custom resources here
	// app.RegisterResource(MyResource)

	log.Println("Starting panel on http://localhost:8080")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
