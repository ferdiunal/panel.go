package main

import (
	"log"

	"github.com/ferdiunal/panel.go/examples/simple/blog"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&blog.Author{}, &blog.Profile{}, &blog.Post{}, &blog.Tag{}, &blog.Comment{})

	cfg := panel.Config{
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Features: panel.FeatureConfig{
			Register:       true,
			ForgotPassword: true,
		},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "8787",
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
			blog.NewAuthorResource(),
			blog.NewProfileResource(),
			blog.NewPostResource(),
			blog.NewTagResource(),
			blog.NewCommentResource(),
		},
	}

	app := panel.New(cfg)

	// You can register custom resources here
	// app.RegisterResource(MyResource)

	log.Println("Starting panel on http://localhost:8787")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
