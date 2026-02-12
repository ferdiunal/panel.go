package entity

import (
	"testing"
	"time"
)

// TestOrganizationCreation - Organization entity oluşturma testi
func TestOrganizationCreation(t *testing.T) {
	org := &Organization{
		ID:    1,
		Name:  "Test Organization",
		Email: "test@example.com",
		Phone: "+90 555 123 4567",
	}

	if org.Name != "Test Organization" {
		t.Errorf("Expected name 'Test Organization', got '%s'", org.Name)
	}

	if org.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", org.Email)
	}

	if org.Phone != "+90 555 123 4567" {
		t.Errorf("Expected phone '+90 555 123 4567', got '%s'", org.Phone)
	}
}

// TestAddressTypeConstants - AddressType enum değerleri testi
func TestAddressTypeConstants(t *testing.T) {
	if AddressTypeSender != 0 {
		t.Errorf("Expected AddressTypeSender to be 0, got %d", AddressTypeSender)
	}

	if AddressTypeReceiver != 1 {
		t.Errorf("Expected AddressTypeReceiver to be 1, got %d", AddressTypeReceiver)
	}

	if AddressTypeInvoice != 2 {
		t.Errorf("Expected AddressTypeInvoice to be 2, got %d", AddressTypeInvoice)
	}
}

// TestAddressCreation - Address entity oluşturma testi
func TestAddressCreation(t *testing.T) {
	addr := &Address{
		ID:             1,
		OrganizationID: 1,
		Name:           "Merkez Ofis",
		Address:        "Atatürk Cad. No:123",
		City:           "İstanbul",
		State:          "İstanbul",
		ZipCode:        "34000",
		Country:        "Türkiye",
		Type:           AddressTypeSender,
	}

	if addr.Name != "Merkez Ofis" {
		t.Errorf("Expected name 'Merkez Ofis', got '%s'", addr.Name)
	}

	if addr.City != "İstanbul" {
		t.Errorf("Expected city 'İstanbul', got '%s'", addr.City)
	}

	if addr.Type != AddressTypeSender {
		t.Errorf("Expected type AddressTypeSender, got %d", addr.Type)
	}
}

// TestCargoCompanyCreation - CargoCompany entity oluşturma testi
func TestCargoCompanyCreation(t *testing.T) {
	company := &CargoCompany{
		ID:   1,
		Name: "Aras Kargo",
	}

	if company.Name != "Aras Kargo" {
		t.Errorf("Expected name 'Aras Kargo', got '%s'", company.Name)
	}
}

// TestPriceListCreation - PriceList entity oluşturma testi
func TestPriceListCreation(t *testing.T) {
	priceList := &PriceList{
		ID:             1,
		CargoCompanyID: 1,
		Name:           "Standart Fiyat Listesi",
	}

	if priceList.Name != "Standart Fiyat Listesi" {
		t.Errorf("Expected name 'Standart Fiyat Listesi', got '%s'", priceList.Name)
	}

	if priceList.CargoCompanyID != 1 {
		t.Errorf("Expected CargoCompanyID 1, got %d", priceList.CargoCompanyID)
	}
}

// TestCommissionCreation - Commission entity oluşturma testi
func TestCommissionCreation(t *testing.T) {
	commission := &Commission{
		ID:             1,
		OrganizationID: 1,
		PriceListID:    1,
		Name:           "Standart Komisyon",
		Commission:     15,
	}

	if commission.Name != "Standart Komisyon" {
		t.Errorf("Expected name 'Standart Komisyon', got '%s'", commission.Name)
	}

	if commission.Commission != 15 {
		t.Errorf("Expected commission 15, got %d", commission.Commission)
	}
}

// TestPriceCreation - Price entity oluşturma testi
func TestPriceCreation(t *testing.T) {
	price := &Price{
		ID:            1,
		PriceListID:   1,
		Price:         5000,
		DesiThreshold: 10,
	}

	if price.Price != 5000 {
		t.Errorf("Expected price 5000, got %d", price.Price)
	}

	if price.DesiThreshold != 10 {
		t.Errorf("Expected DesiThreshold 10, got %d", price.DesiThreshold)
	}
}

// TestShipmentCreation - Shipment entity oluşturma testi
func TestShipmentCreation(t *testing.T) {
	shipment := &Shipment{
		ID:                1,
		OrganizationID:    1,
		SenderAddressID:   1,
		ReceiverAddressID: 2,
		Name:              "Sipariş #12345",
	}

	if shipment.Name != "Sipariş #12345" {
		t.Errorf("Expected name 'Sipariş #12345', got '%s'", shipment.Name)
	}

	if shipment.SenderAddressID != 1 {
		t.Errorf("Expected SenderAddressID 1, got %d", shipment.SenderAddressID)
	}

	if shipment.ReceiverAddressID != 2 {
		t.Errorf("Expected ReceiverAddressID 2, got %d", shipment.ReceiverAddressID)
	}
}

// TestShipmentRowCreation - ShipmentRow entity oluşturma testi
func TestShipmentRowCreation(t *testing.T) {
	row := &ShipmentRow{
		ID:         1,
		ShipmentID: 1,
		ProductID:  1,
		Quantity:   5,
	}

	if row.Quantity != 5 {
		t.Errorf("Expected quantity 5, got %d", row.Quantity)
	}

	if row.ShipmentID != 1 {
		t.Errorf("Expected ShipmentID 1, got %d", row.ShipmentID)
	}

	if row.ProductID != 1 {
		t.Errorf("Expected ProductID 1, got %d", row.ProductID)
	}
}

// TestProductCreation - Product entity oluşturma testi
func TestProductCreation(t *testing.T) {
	product := &Product{
		ID:             1,
		OrganizationID: 1,
		Name:           "Laptop",
	}

	if product.Name != "Laptop" {
		t.Errorf("Expected name 'Laptop', got '%s'", product.Name)
	}

	if product.OrganizationID != 1 {
		t.Errorf("Expected OrganizationID 1, got %d", product.OrganizationID)
	}
}

// TestEntityRelationships - Entity ilişkileri testi
func TestEntityRelationships(t *testing.T) {
	// Organization ile Address ilişkisi
	org := &Organization{
		ID:   1,
		Name: "Test Org",
		Addresses: []Address{
			{ID: 1, Name: "Address 1"},
			{ID: 2, Name: "Address 2"},
		},
	}

	if len(org.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(org.Addresses))
	}

	// Shipment ile ShipmentRow ilişkisi
	shipment := &Shipment{
		ID:   1,
		Name: "Test Shipment",
		ShipmentRows: []ShipmentRow{
			{ID: 1, Quantity: 5},
			{ID: 2, Quantity: 10},
		},
	}

	if len(shipment.ShipmentRows) != 2 {
		t.Errorf("Expected 2 shipment rows, got %d", len(shipment.ShipmentRows))
	}

	totalQuantity := 0
	for _, row := range shipment.ShipmentRows {
		totalQuantity += row.Quantity
	}

	if totalQuantity != 15 {
		t.Errorf("Expected total quantity 15, got %d", totalQuantity)
	}
}

// TestEntityTimestamps - Entity timestamp alanları testi
func TestEntityTimestamps(t *testing.T) {
	now := time.Now()

	org := &Organization{
		ID:        1,
		Name:      "Test Org",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if org.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if org.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if !org.CreatedAt.Equal(org.UpdatedAt) {
		t.Error("Expected CreatedAt and UpdatedAt to be equal")
	}
}

// TestEntityPointerRelationships - Entity pointer ilişkileri testi
func TestEntityPointerRelationships(t *testing.T) {
	// ShipmentRow ile Product ilişkisi
	product := &Product{
		ID:   1,
		Name: "Test Product",
	}

	row := &ShipmentRow{
		ID:        1,
		ProductID: 1,
		Product:   product,
		Quantity:  5,
	}

	if row.Product == nil {
		t.Error("Expected Product to be set")
	}

	if row.Product.Name != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", row.Product.Name)
	}

	// Address ile Organization ilişkisi
	org := &Organization{
		ID:   1,
		Name: "Test Org",
	}

	addr := &Address{
		ID:             1,
		OrganizationID: 1,
		Organization:   org,
		Name:           "Test Address",
	}

	if addr.Organization == nil {
		t.Error("Expected Organization to be set")
	}

	if addr.Organization.Name != "Test Org" {
		t.Errorf("Expected organization name 'Test Org', got '%s'", addr.Organization.Name)
	}
}
