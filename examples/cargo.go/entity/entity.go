package entity

import (
	"time"
)

// Organization - Organizasyon bilgilerini tutan model
// HasMany: Addresses (bir organizasyonun birden fazla adresi olabilir)
// HasMany: Commissions (bir organizasyonun birden fazla komisyonu olabilir)
// HasMany: Shipments (bir organizasyonun birden fazla gönderisi olabilir)
// HasMany: Products (bir organizasyonun birden fazla ürünü olabilir)
//
// Kullanım örneği:
// ```go
// // Organizasyon oluşturma
// org := &entity.Organization{
//     Name:  "Acme Corp",
//     Email: "info@acme.com",
//     Phone: "+90 555 123 4567",
// }
// db.Create(&org)
//
// // İlişkili verileri yükleme (Eager Loading)
// var organization entity.Organization
// db.Preload("Addresses").Preload("Products").First(&organization, 1)
//
// // Has Many ilişki üzerinden veri ekleme
// address := entity.Address{Name: "Merkez Ofis", City: "İstanbul"}
// db.Model(&org).Association("Addresses").Append(&address)
// ```
type Organization struct {
	ID          uint64       `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	Name        string       `gorm:"not null;column:name;varchar(255)" json:"name"`
	Email       string       `gorm:"not null;column:email;varchar(255)" json:"email"`
	Phone       string       `gorm:"not null;column:phone;varchar(255)" json:"phone"`
	Addresses   []Address    `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"addresses,omitempty"`
	Commissions []Commission `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"commissions,omitempty"`
	Shipments   []Shipment   `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipments,omitempty"`
	Products    []Product    `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"products,omitempty"`
	CreatedAt   time.Time    `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// AddressType - Adres tipi enum değerleri
// 0: Gönderici adresi, 1: Alıcı adresi, 2: Fatura adresi
type AddressType int

const (
	AddressTypeSender   AddressType = 0 // Gönderici adresi
	AddressTypeReceiver AddressType = 1 // Alıcı adresi
	AddressTypeInvoice  AddressType = 2 // Fatura adresi
)

// Address - Adres bilgilerini tutan model
// BelongsTo: Organization (bir adres bir organizasyona aittir)
//
// Kullanım örneği:
// ```go
// // Adres oluşturma
// address := &entity.Address{
//     OrganizationID: 1,
//     Name:           "Merkez Ofis",
//     Address:        "Atatürk Cad. No:123",
//     City:           "İstanbul",
//     State:          "İstanbul",
//     ZipCode:        "34000",
//     Country:        "Türkiye",
//     Type:           entity.AddressTypeSender,
// }
// db.Create(&address)
//
// // Organization ile birlikte yükleme
// var address entity.Address
// db.Preload("Organization").First(&address, 1)
// ```
type Address struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index:idx_addresses_org_type,priority:1;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	Address        string        `gorm:"not null;column:address;text" json:"address"`
	City           string        `gorm:"not null;column:city;varchar(255);index" json:"city"`
	State          string        `gorm:"not null;column:state;varchar(255)" json:"state"`
	ZipCode        string        `gorm:"not null;column:zip_code;varchar(255)" json:"zipCode"`
	Country        string        `gorm:"not null;column:country;varchar(255);index" json:"country"`
	Type           AddressType   `gorm:"not null;column:type;int;index:idx_addresses_org_type,priority:2" json:"type"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// CargoCompany - Kargo şirketi bilgilerini tutan model
// HasMany: PriceLists (bir kargo şirketinin birden fazla fiyat listesi olabilir)
//
// Kullanım örneği:
// ```go
// // Kargo şirketi oluşturma
// company := &entity.CargoCompany{Name: "Aras Kargo"}
// db.Create(&company)
//
// // Fiyat listeleri ile birlikte yükleme
// var company entity.CargoCompany
// db.Preload("PriceLists").First(&company, 1)
// ```
type CargoCompany struct {
	ID         uint64      `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	Name       string      `gorm:"not null;column:name;varchar(255)" json:"name"`
	PriceLists []PriceList `gorm:"foreignKey:CargoCompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"priceLists,omitempty"`
	CreatedAt  time.Time   `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// PriceList - Fiyat listesi bilgilerini tutan model
// BelongsTo: CargoCompany (bir fiyat listesi bir kargo şirketine aittir)
// HasMany: Prices (bir fiyat listesinin birden fazla fiyatı olabilir)
// HasMany: Commissions (bir fiyat listesinin birden fazla komisyonu olabilir)
//
// Kullanım örneği:
// ```go
// // Fiyat listesi oluşturma
// priceList := &entity.PriceList{
//     CargoCompanyID: 1,
//     Name:           "Standart Fiyat Listesi",
// }
// db.Create(&priceList)
//
// // İlişkili verilerle birlikte yükleme
// var priceList entity.PriceList
// db.Preload("CargoCompany").Preload("Prices").Preload("Commissions").First(&priceList, 1)
// ```
type PriceList struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	CargoCompanyID uint64        `gorm:"not null;column:cargo_company_id;bigint;index" json:"cargoCompanyId"`
	CargoCompany   *CargoCompany `gorm:"foreignKey:CargoCompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cargoCompany,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	Prices         []Price       `gorm:"foreignKey:PriceListID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"prices,omitempty"`
	Commissions    []Commission  `gorm:"foreignKey:PriceListID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"commissions,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Commission - Komisyon bilgilerini tutan model
// BelongsTo: Organization (bir komisyon bir organizasyona aittir)
// BelongsTo: PriceList (bir komisyon bir fiyat listesine aittir)
//
// Kullanım örneği:
// ```go
// // Komisyon oluşturma
// commission := &entity.Commission{
//     OrganizationID: 1,
//     PriceListID:    1,
//     Name:           "Standart Komisyon",
//     Commission:     15, // %15 komisyon
// }
// db.Create(&commission)
//
// // İlişkili verilerle birlikte yükleme
// var commission entity.Commission
// db.Preload("Organization").Preload("PriceList").First(&commission, 1)
// ```
type Commission struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	PriceListID    uint64        `gorm:"not null;column:price_list_id;bigint;index" json:"priceListId"`
	PriceList      *PriceList    `gorm:"foreignKey:PriceListID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"priceList,omitempty"`
	Commission     uint64        `gorm:"not null;column:commission;bigint" json:"commission"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Price - Fiyat bilgilerini tutan model
// BelongsTo: PriceList (bir fiyat bir fiyat listesine aittir)
//
// Kullanım örneği:
// ```go
// // Fiyat oluşturma
// price := &entity.Price{
//     PriceListID:   1,
//     Price:         5000, // 50.00 TL (kuruş cinsinden)
//     DesiThreshold: 10,   // 10 desi
// }
// db.Create(&price)
//
// // PriceList ile birlikte yükleme
// var price entity.Price
// db.Preload("PriceList").First(&price, 1)
// ```
type Price struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	PriceListID   uint64     `gorm:"not null;column:price_list_id;bigint;index" json:"priceListId"`
	PriceList     *PriceList `gorm:"foreignKey:PriceListID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"priceList,omitempty"`
	Price         uint64     `gorm:"not null;column:price;bigint" json:"price"`
	DesiThreshold uint64     `gorm:"not null;column:desi_threshold;bigint" json:"desiThreshold"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Shipment - Gönderi bilgilerini tutan model
// BelongsTo: Organization (bir gönderi bir organizasyona aittir)
// BelongsTo: SenderAddress (gönderici adresi)
// BelongsTo: ReceiverAddress (alıcı adresi)
// HasMany: ShipmentRows (bir gönderinin birden fazla satırı olabilir)
//
// Kullanım örneği:
// ```go
// // Gönderi oluşturma
// shipment := &entity.Shipment{
//     OrganizationID:    1,
//     SenderAddressID:   1,
//     ReceiverAddressID: 2,
//     Name:              "Sipariş #12345",
// }
// db.Create(&shipment)
//
// // İlişkili verilerle birlikte yükleme
// var shipment entity.Shipment
// db.Preload("Organization").
//    Preload("SenderAddress").
//    Preload("ReceiverAddress").
//    Preload("ShipmentRows").
//    First(&shipment, 1)
// ```
type Shipment struct {
	ID                uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID    uint64        `gorm:"not null;column:organization_id;bigint;index:idx_shipments_org_created,priority:1;index" json:"organizationId"`
	Organization      *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	SenderAddressID   uint64        `gorm:"not null;column:sender_address_id;bigint;index" json:"senderAddressId"`
	SenderAddress     *Address      `gorm:"foreignKey:SenderAddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"senderAddress,omitempty"`
	ReceiverAddressID uint64        `gorm:"not null;column:receiver_address_id;bigint;index" json:"receiverAddressId"`
	ReceiverAddress   *Address      `gorm:"foreignKey:ReceiverAddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"receiverAddress,omitempty"`
	Name              string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	ShipmentRows      []ShipmentRow `gorm:"foreignKey:ShipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipmentRows,omitempty"`
	CreatedAt         time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz;index:idx_shipments_org_created,priority:2" json:"createdAt"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// ShipmentRow - Gönderi satırı bilgilerini tutan model
// BelongsTo: Shipment (bir gönderi satırı bir gönderiye aittir)
// BelongsTo: Product (bir gönderi satırı bir ürüne aittir)
//
// Kullanım örneği:
// ```go
// // Gönderi satırı oluşturma
// row := &entity.ShipmentRow{
//     ShipmentID: 1,
//     ProductID:  1,
//     Quantity:   5,
// }
// db.Create(&row)
//
// // İlişkili verilerle birlikte yükleme
// var row entity.ShipmentRow
// db.Preload("Shipment").Preload("Product").First(&row, 1)
// ```
type ShipmentRow struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	ShipmentID uint64    `gorm:"not null;column:shipment_id;bigint;index" json:"shipmentId"`
	Shipment   *Shipment `gorm:"foreignKey:ShipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipment,omitempty"`
	ProductID  uint64    `gorm:"not null;column:product_id;bigint;index" json:"productId"`
	Product    *Product  `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product,omitempty"`
	Quantity   int       `gorm:"not null;column:quantity;int" json:"quantity"`
	CreatedAt  time.Time `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Product - Ürün bilgilerini tutan model
// BelongsTo: Organization (bir ürün bir organizasyona aittir)
// HasMany: ShipmentRows (bir ürünün birden fazla gönderi satırı olabilir)
//
// Kullanım örneği:
// ```go
// // Ürün oluşturma
// product := &entity.Product{
//     OrganizationID: 1,
//     Name:           "Laptop",
// }
// db.Create(&product)
//
// // İlişkili verilerle birlikte yükleme
// var product entity.Product
// db.Preload("Organization").Preload("ShipmentRows").First(&product, 1)
// ```
type Product struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index:idx_products_org_created,priority:1;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255);index" json:"name"`
	ShipmentRows   []ShipmentRow `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"shipmentRows,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz;index:idx_products_org_created,priority:2" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}
