package data

import (
	"context"
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"gorm.io/gorm"
)

// GormRelationshipLoader, GORM kullanarak ilişkileri yükleyen bir implementasyondur.
//
// Bu yapı, RelationshipLoader interface'ini implement eder ve raw SQL kullanarak
// ilişkileri yükler. Circular dependency sorununu çözmek için struct field'larına
// bağımlı olmadan çalışır.
//
// # Özellikler
//
// - **Raw SQL**: GORM Preload yerine raw SQL kullanır
// - **Batch Loading**: N+1 problemini önlemek için batch loading yapar
// - **Tip Bağımsız**: Struct field'larına bağımlı değildir
// - **Performans**: Optimize edilmiş sorgular ile hızlı yükleme
//
// # Kullanım Örneği
//
//	loader := NewGormRelationshipLoader(db)
//	err := loader.EagerLoad(ctx, items, relationshipField)
type GormRelationshipLoader struct {
	db *gorm.DB
}

// NewGormRelationshipLoader, yeni bir GormRelationshipLoader instance'ı oluşturur.
//
// # Parametreler
//
// - **db**: GORM veritabanı bağlantısı (*gorm.DB)
//
// # Döndürür
//
// - GormRelationshipLoader pointer'ı
//
// # Kullanım Örneği
//
//	loader := NewGormRelationshipLoader(db)
func NewGormRelationshipLoader(db *gorm.DB) *GormRelationshipLoader {
	return &GormRelationshipLoader{db: db}
}

// EagerLoad, birden fazla kayıt için ilişkileri batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların ilişkilerini
// tek seferde yükler. İlişki tipine göre uygun stratejiyi seçer.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (context.Context)
// - **items**: İlişkileri yüklenecek kayıt listesi ([]interface{})
// - **field**: İlişki field tanımı (fields.RelationshipField)
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Desteklenen İlişki Tipleri
//
// - **belongsTo**: BelongsTo ilişkisi
// - **hasMany**: HasMany ilişkisi
// - **hasOne**: HasOne ilişkisi
//
// # Kullanım Örneği
//
//	items := []interface{}{&user1, &user2, &user3}
//	err := loader.EagerLoad(ctx, items, authorField)
func (l *GormRelationshipLoader) EagerLoad(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
	if len(items) == 0 {
		return nil
	}

	// İlişki tipine göre yönlendir
	switch field.GetRelationshipType() {
	case "belongsTo":
		return l.eagerLoadBelongsTo(ctx, items, field)
	case "hasMany":
		return l.eagerLoadHasMany(ctx, items, field)
	case "hasOne":
		return l.eagerLoadHasOne(ctx, items, field)
	case "belongsToMany":
		return l.eagerLoadBelongsToMany(ctx, items, field)
	default:
		return fmt.Errorf("unsupported relationship type: %s", field.GetRelationshipType())
	}
}

// LazyLoad, tek bir kayıt için ilişkiyi yükler.
//
// Bu metod, ihtiyaç anında tek bir kaydın ilişkisini yükler.
// Eager loading'e göre daha az performanslıdır ancak bellek tasarrufu sağlar.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (context.Context)
// - **item**: İlişkisi yüklenecek kayıt (interface{})
// - **field**: İlişki field tanımı (fields.RelationshipField)
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi
// - error: Hata durumunda hata mesajı
//
// # Desteklenen İlişki Tipleri
//
// - **belongsTo**: BelongsTo ilişkisi
// - **hasMany**: HasMany ilişkisi
// - **hasOne**: HasOne ilişkisi
//
// # Kullanım Örneği
//
//	relatedData, err := loader.LazyLoad(ctx, &user, authorField)
func (l *GormRelationshipLoader) LazyLoad(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// İlişki tipine göre yönlendir
	switch field.GetRelationshipType() {
	case "belongsTo":
		return l.lazyLoadBelongsTo(ctx, item, field)
	case "hasMany":
		return l.lazyLoadHasMany(ctx, item, field)
	case "hasOne":
		return l.lazyLoadHasOne(ctx, item, field)
	case "belongsToMany":
		return l.lazyLoadBelongsToMany(ctx, item, field)
	default:
		return nil, fmt.Errorf("unsupported relationship type: %s", field.GetRelationshipType())
	}
}

// LoadWithConstraints, constraint'ler ile ilişkiyi yükler.
//
// Bu metod, ek WHERE koşulları ile ilişkiyi yükler.
// Filtrelenmiş ilişki yükleme için kullanılır.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (context.Context)
// - **item**: İlişkisi yüklenecek kayıt (interface{})
// - **field**: İlişki field tanımı (fields.RelationshipField)
// - **constraints**: Ek WHERE koşulları (map[string]interface{})
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	constraints := map[string]interface{}{
//	    "status": "active",
//	    "published": true,
//	}
//	relatedData, err := loader.LoadWithConstraints(ctx, &user, postsField, constraints)
func (l *GormRelationshipLoader) LoadWithConstraints(ctx context.Context, item interface{}, field fields.RelationshipField, constraints map[string]interface{}) (interface{}, error) {
	// Şu an için LazyLoad ile aynı davranış
	// Gelecekte constraint'leri WHERE koşulu olarak ekleyebiliriz
	return l.LazyLoad(ctx, item, field)
}
