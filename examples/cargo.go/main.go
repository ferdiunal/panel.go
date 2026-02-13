package main

import (
	"fmt"
	"log"
	"os"

	"cargo.go/resources/address"
	"cargo.go/resources/cargo_company"
	"cargo.go/resources/commissions"
	"cargo.go/resources/organization"
	"cargo.go/resources/prices"
	"cargo.go/resources/products"
	"cargo.go/resources/shipment"
	"cargo.go/resources/shipment_row"
	"github.com/ferdiunal/panel.go/pkg/migration"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	password := os.Getenv("DB_PASSWORD")
	var dsn string
	if password == "" {
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"))
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			password,
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"))
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	mg := migration.NewMigrationGenerator(db)
	mg.RegisterResources(
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		cargo_company.NewCargoCompanyResource(),
		organization.NewOrganizationResource(),
		address.NewAddressResource(),
		commissions.NewCommissionResource(),
		prices.NewPriceListResource(),
		prices.NewPriceResource(),
		products.NewProductResource(),
		shipment.NewShipmentResource(),
		shipment_row.NewShipmentRowResource(),
	)
	if err := mg.AutoMigrate(); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	cfg := panel.Config{
		EncryptionCookie: encryptcookie.Config{
			Key:    os.Getenv("COOKIE_ENCRYPTION_KEY"),
			Except: []string{"csrf_token"},
		},
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
			Path: "permissions.toml",
		},
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
			resourceSession.NewSessionResource(),
			resourceVerification.NewVerificationResource(),
			cargo_company.NewCargoCompanyResource(),
			organization.NewOrganizationResource(),
			address.NewAddressResource(),
			commissions.NewCommissionResource(),
			prices.NewPriceListResource(),
			prices.NewPriceResource(),
			products.NewProductResource(),
			shipment.NewShipmentResource(),
			shipment.NewShipmentRowResource(),
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
