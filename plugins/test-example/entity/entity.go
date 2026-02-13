package entity

import (
	"time"
)

// Organization - Organizasyon bilgilerini tutan model
// HasMany: Addresses, Products, Shipments
// HasOne: BillingInfo
type Organization struct {
	ID          uint64       `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	Name        string       `gorm:"not null;column:name;varchar(255)" json:"name"`
	Email       string       `gorm:"not null;column:email;varchar(255)" json:"email"`
	Phone       string       `gorm:"not null;column:phone;varchar(255)" json:"phone"`
	Addresses   []Address    `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"addresses,omitempty"`
	BillingInfo *BillingInfo `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"billingInfo,omitempty"`
	Products    []Product    `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"products,omitempty"`
	Shipments   []Shipment   `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipments,omitempty"`
	CreatedAt   time.Time    `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// BillingInfo - Fatura bilgilerini tutan model (HasOne örneği)
// BelongsTo: Organization
type BillingInfo struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;unique;column:organization_id;bigint;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	TaxNumber      string        `gorm:"not null;column:tax_number;varchar(255)" json:"taxNumber"`
	TaxOffice      string        `gorm:"not null;column:tax_office;varchar(255)" json:"taxOffice"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Address - Adres bilgilerini tutan model
// BelongsTo: Organization
type Address struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	Address        string        `gorm:"not null;column:address;text" json:"address"`
	City           string        `gorm:"not null;column:city;varchar(255);index" json:"city"`
	Country        string        `gorm:"not null;column:country;varchar(255);index" json:"country"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Product - Ürün bilgilerini tutan model
// BelongsTo: Organization
// BelongsToMany: Categories
// HasMany: ShipmentRows
type Product struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255);index" json:"name"`
	Categories     []Category    `gorm:"many2many:product_categories;" json:"categories,omitempty"`
	ShipmentRows   []ShipmentRow `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"shipmentRows,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Category - Kategori bilgilerini tutan model (BelongsToMany örneği)
// BelongsToMany: Products
type Category struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	Name      string    `gorm:"not null;column:name;varchar(255);index" json:"name"`
	Products  []Product `gorm:"many2many:product_categories;" json:"products,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Shipment - Gönderi bilgilerini tutan model
// BelongsTo: Organization
// HasMany: ShipmentRows
type Shipment struct {
	ID             uint64        `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	OrganizationID uint64        `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"`
	Name           string        `gorm:"not null;column:name;varchar(255)" json:"name"`
	ShipmentRows   []ShipmentRow `gorm:"foreignKey:ShipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipmentRows,omitempty"`
	CreatedAt      time.Time     `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// ShipmentRow - Gönderi satırı bilgilerini tutan model
// BelongsTo: Shipment, Product
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

// Comment - Yorum bilgilerini tutan model (MorphTo örneği)
// MorphTo: Commentable (Product/Shipment)
type Comment struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	CommentableID   uint64    `gorm:"not null;column:commentable_id;bigint;index:idx_commentable" json:"commentableId"`
	CommentableType string    `gorm:"not null;column:commentable_type;varchar(255);index:idx_commentable" json:"commentableType"`
	Content         string    `gorm:"not null;column:content;text" json:"content"`
	CreatedAt       time.Time `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}

// Tag - Etiket bilgilerini tutan model (MorphToMany örneği)
// MorphToMany: Taggable (Product/Shipment)
type Tag struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"`
	Name      string    `gorm:"not null;unique;column:name;varchar(255);index" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"`
}
