package main

import (
	"fmt"
	"os"
	"testing"

	"cargo.go/entity"
	"cargo.go/resources"
	_ "cargo.go/resources/address"
	_ "cargo.go/resources/cargo_company"
	_ "cargo.go/resources/commissions"
	_ "cargo.go/resources/organization"
	_ "cargo.go/resources/prices"
	_ "cargo.go/resources/products"
	_ "cargo.go/resources/shipment"
	_ "cargo.go/resources/shipment_row"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB - Test veritabanı bağlantısı oluşturur
func setupTestDB(t *testing.T) *gorm.DB {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		t.Skip("Skipping integration test: .env file not found")
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
		t.Skipf("Skipping integration test: database connection failed: %v", err)
	}

	return db
}

// TestResourceRegistry - Resource registry entegrasyon testi
func TestResourceRegistry(t *testing.T) {
	// Resource'lar import sırasında init() fonksiyonları ile kayıt edilir
	// Clear() çağrılmamalı çünkü init() fonksiyonları sadece bir kez çağrılır

	// Kayıtlı resource'ları kontrol et
	slugs := resources.List()
	expectedSlugs := []string{
		"organizations",
		"addresses",
		"cargo-companies",
		"commissions",
		"price-lists",
		"products",
		"shipments",
		"shipment-rows",
	}

	// En az beklenen resource'lar kayıtlı olmalı
	if len(slugs) < len(expectedSlugs) {
		t.Errorf("Expected at least %d registered resources, got %d", len(expectedSlugs), len(slugs))
	}

	// Her bir resource'un kayıtlı olduğunu kontrol et
	for _, expectedSlug := range expectedSlugs {
		res := resources.Get(expectedSlug)
		if res == nil {
			t.Errorf("Expected resource '%s' to be registered", expectedSlug)
		}
	}
}

// TestDatabaseConnection - Veritabanı bağlantı testi
func TestDatabaseConnection(t *testing.T) {
	db := setupTestDB(t)

	// Basit bir sorgu çalıştır
	var result int
	err := db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Errorf("Database query failed: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result 1, got %d", result)
	}
}

// TestOrganizationCRUD - Organization CRUD işlemleri testi
func TestOrganizationCRUD(t *testing.T) {
	db := setupTestDB(t)

	// Test organizasyonu oluştur
	org := &entity.Organization{
		Name:  "Test Organization",
		Email: "test@example.com",
		Phone: "+90 555 123 4567",
	}

	// Create
	result := db.Create(org)
	if result.Error != nil {
		t.Fatalf("Failed to create organization: %v", result.Error)
	}

	if org.ID == 0 {
		t.Error("Expected organization ID to be set after creation")
	}

	// Read
	var fetchedOrg entity.Organization
	result = db.First(&fetchedOrg, org.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch organization: %v", result.Error)
	}

	if fetchedOrg.Name != org.Name {
		t.Errorf("Expected name '%s', got '%s'", org.Name, fetchedOrg.Name)
	}

	// Update
	fetchedOrg.Name = "Updated Organization"
	result = db.Save(&fetchedOrg)
	if result.Error != nil {
		t.Fatalf("Failed to update organization: %v", result.Error)
	}

	// Verify update
	var updatedOrg entity.Organization
	result = db.First(&updatedOrg, org.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch updated organization: %v", result.Error)
	}

	if updatedOrg.Name != "Updated Organization" {
		t.Errorf("Expected name 'Updated Organization', got '%s'", updatedOrg.Name)
	}

	// Delete
	result = db.Delete(&updatedOrg)
	if result.Error != nil {
		t.Fatalf("Failed to delete organization: %v", result.Error)
	}

	// Verify deletion
	var deletedOrg entity.Organization
	result = db.First(&deletedOrg, org.ID)
	if result.Error == nil {
		t.Error("Expected organization to be deleted")
	}
}

// TestRelationshipLoading - İlişki yükleme testi
func TestRelationshipLoading(t *testing.T) {
	db := setupTestDB(t)

	// Test organizasyonu oluştur
	org := &entity.Organization{
		Name:  "Test Organization",
		Email: "test@example.com",
		Phone: "+90 555 123 4567",
	}

	result := db.Create(org)
	if result.Error != nil {
		t.Fatalf("Failed to create organization: %v", result.Error)
	}

	// Test adresi oluştur
	addr := &entity.Address{
		OrganizationID: org.ID,
		Name:           "Test Address",
		Address:        "Test Street",
		City:           "Test City",
		State:          "Test State",
		ZipCode:        "12345",
		Country:        "Test Country",
		Type:           entity.AddressTypeSender,
	}

	result = db.Create(addr)
	if result.Error != nil {
		t.Fatalf("Failed to create address: %v", result.Error)
	}

	// Organization'ı Addresses ile birlikte yükle
	var fetchedOrg entity.Organization
	result = db.Preload("Addresses").First(&fetchedOrg, org.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch organization with addresses: %v", result.Error)
	}

	if len(fetchedOrg.Addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(fetchedOrg.Addresses))
	}

	if fetchedOrg.Addresses[0].Name != "Test Address" {
		t.Errorf("Expected address name 'Test Address', got '%s'", fetchedOrg.Addresses[0].Name)
	}

	// Cleanup
	db.Delete(addr)
	db.Delete(org)
}

// TestShipmentWithRows - Shipment ve ShipmentRow ilişki testi
func TestShipmentWithRows(t *testing.T) {
	db := setupTestDB(t)

	// Test organizasyonu oluştur
	org := &entity.Organization{
		Name:  "Test Organization",
		Email: "test@example.com",
		Phone: "+90 555 123 4567",
	}
	db.Create(org)

	// Test adresleri oluştur
	senderAddr := &entity.Address{
		OrganizationID: org.ID,
		Name:           "Sender Address",
		Address:        "Sender Street",
		City:           "Sender City",
		State:          "Sender State",
		ZipCode:        "12345",
		Country:        "Sender Country",
		Type:           entity.AddressTypeSender,
	}
	db.Create(senderAddr)

	receiverAddr := &entity.Address{
		OrganizationID: org.ID,
		Name:           "Receiver Address",
		Address:        "Receiver Street",
		City:           "Receiver City",
		State:          "Receiver State",
		ZipCode:        "54321",
		Country:        "Receiver Country",
		Type:           entity.AddressTypeReceiver,
	}
	db.Create(receiverAddr)

	// Test ürünü oluştur
	product := &entity.Product{
		OrganizationID: org.ID,
		Name:           "Test Product",
	}
	db.Create(product)

	// Test gönderisi oluştur
	shipment := &entity.Shipment{
		OrganizationID:    org.ID,
		SenderAddressID:   senderAddr.ID,
		ReceiverAddressID: receiverAddr.ID,
		Name:              "Test Shipment",
	}
	db.Create(shipment)

	// Test gönderi satırı oluştur
	row := &entity.ShipmentRow{
		ShipmentID: shipment.ID,
		ProductID:  product.ID,
		Quantity:   5,
	}
	db.Create(row)

	// Shipment'ı tüm ilişkilerle birlikte yükle
	var fetchedShipment entity.Shipment
	result := db.Preload("Organization").
		Preload("SenderAddress").
		Preload("ReceiverAddress").
		Preload("ShipmentRows").
		Preload("ShipmentRows.Product").
		First(&fetchedShipment, shipment.ID)

	if result.Error != nil {
		t.Fatalf("Failed to fetch shipment with relationships: %v", result.Error)
	}

	// İlişkileri kontrol et
	if fetchedShipment.Organization == nil {
		t.Error("Expected Organization to be loaded")
	}

	if fetchedShipment.SenderAddress == nil {
		t.Error("Expected SenderAddress to be loaded")
	}

	if fetchedShipment.ReceiverAddress == nil {
		t.Error("Expected ReceiverAddress to be loaded")
	}

	if len(fetchedShipment.ShipmentRows) != 1 {
		t.Errorf("Expected 1 shipment row, got %d", len(fetchedShipment.ShipmentRows))
	}

	if fetchedShipment.ShipmentRows[0].Product == nil {
		t.Error("Expected Product to be loaded in ShipmentRow")
	}

	if fetchedShipment.ShipmentRows[0].Quantity != 5 {
		t.Errorf("Expected quantity 5, got %d", fetchedShipment.ShipmentRows[0].Quantity)
	}

	// Cleanup
	db.Delete(row)
	db.Delete(shipment)
	db.Delete(product)
	db.Delete(receiverAddr)
	db.Delete(senderAddr)
	db.Delete(org)
}
