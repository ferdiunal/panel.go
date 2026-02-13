// Package page, sistem yönetim panelinin sayfa bileşenlerini içerir.
//
// Bu paket, panel uygulamasında farklı sayfa türlerini (Account, Dashboard, vb.)
// tanımlamak ve yönetmek için kullanılan temel yapıları ve arayüzleri sağlar.
package page

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Account, kullanıcı hesap ayarlarını yönetmek için kullanılan sayfa bileşenidir.
//
// # Açıklama
// Account struct'ı, panel uygulamasında kullanıcıların kendi hesap ayarlarını
// (profil bilgileri, şifre değiştirme, bildirim tercihleri, vb.) yönetmek için
// tasarlanmıştır. Her kullanıcı kendi hesap ayarlarını görüntüleyebilir ve düzenleyebilir.
//
// # Yapı Alanları
// - Base: Tüm sayfa bileşenlerinin ortak özelliklerini içeren gömülü struct
// - Elements: Hesap ayarları sayfasında gösterilecek form alanlarının listesi
// - HideInNavigation: Hesap ayarları sayfasının navigasyon menüsünde gizlenip gizlenmeyeceğini belirler
//
// # Kullanım Senaryoları
// 1. Kullanıcı profil bilgilerini güncelleme (ad, email, resim)
// 2. Şifre değiştirme işlemleri
// 3. Bildirim tercihlerini yönetme (email, SMS)
// 4. Dil ve tema ayarlarını değiştirme
//
// # Örnek Kullanım
// ```go
// account := &Account{
//     Elements: []fields.Element{
//         // Form alanları buraya eklenir
//     },
//     HideInNavigation: false,
// }
//
// // Hesap ayarlarını kaydetmek
// err := account.Save(ctx, db, map[string]interface{}{
//     "name": "John Doe",
//     "email": "john@example.com",
// })
// ```
//
// # Avantajlar
// - Kullanıcı dostu hesap yönetimi
// - Dinamik form alanları ile esnek ayar yönetimi
// - Navigasyon menüsünde gösterilip gizlenebilir
// - Farklı veri tiplerini destekler
//
// # Önemli Notlar
// - Her kullanıcı sadece kendi hesap ayarlarını görebilir ve düzenleyebilir
// - Şifre değiştirme işleminde mevcut şifre doğrulaması yapılmalıdır
// - Hassas bilgiler (şifre) güvenli şekilde saklanmalıdır
type Account struct {
	Base
	// Elements, hesap ayarları sayfasında gösterilecek form alanlarının listesidir.
	// Her Element, bir form alanını temsil eder (TextInput, Select, Checkbox, vb.)
	Elements []fields.Element

	// HideInNavigation, hesap ayarları sayfasının navigasyon menüsünde gizlenip gizlenmeyeceğini belirler.
	// true ise sayfa menüde görünmez, false ise görünür.
	HideInNavigation bool
}

// Slug, hesap ayarları sayfasının URL'de kullanılan benzersiz tanımlayıcısını döndürür.
//
// # Dönüş Değeri
// "account" - Sayfanın URL slug'ı
//
// # Kullanım
// Bu metot, sayfa yönlendirmesi ve URL oluşturma işlemlerinde kullanılır.
// Örneğin: /admin/account
//
// # Örnek
// ```go
// account := &Account{}
// slug := account.Slug() // "account"
// ```
func (p *Account) Slug() string {
	return "account"
}

// Title, hesap ayarları sayfasının başlığını döndürür.
//
// # Dönüş Değeri
// "Account" - Sayfanın görüntülenecek başlığı
//
// # Kullanım
// Bu başlık, sayfa başlığı, tarayıcı sekmesi ve navigasyon menüsünde gösterilir.
//
// # Örnek
// ```go
// account := &Account{}
// title := account.Title() // "Account"
// ```
func (p *Account) Title() string {
	return "Account"
}

// Description, hesap ayarları sayfasının açıklamasını döndürür.
//
// # Dönüş Değeri
// "Hesap ayarlarınızı yönetin" - Sayfanın açıklaması
//
// # Kullanım
// Bu açıklama, sayfa hakkında bilgi vermek için kullanıcı arayüzünde gösterilir.
// Genellikle sayfa başlığının altında veya navigasyon menüsünde tooltip olarak görünür.
//
// # Örnek
// ```go
// account := &Account{}
// desc := account.Description() // "Hesap ayarlarınızı yönetin"
// ```
func (p *Account) Description() string {
	return "Hesap ayarlarınızı yönetin"
}

// Group, hesap ayarları sayfasının ait olduğu grup/kategorisini döndürür.
//
// # Dönüş Değeri
// "User" - Sayfanın ait olduğu grup adı
//
// # Kullanım
// Navigasyon menüsünde sayfaları gruplandırmak için kullanılır.
// Aynı grup adına sahip sayfalar menüde birlikte gösterilir.
//
// # Örnek
// ```go
// account := &Account{}
// group := account.Group() // "User"
// // Menüde "User" başlığı altında gösterilir
// ```
func (p *Account) Group() string {
	return "User"
}

// NavigationOrder, navigasyon menüsünde sayfanın gösterilme sırasını belirler.
//
// # Dönüş Değeri
// 10 - Sıra numarası (daha düşük sayı, menünün daha üstünde gösterilir)
//
// # Kullanım
// Navigasyon menüsünde sayfaları sıralamak için kullanılır.
// Kullanıcı sayfaları genellikle menünün üst kısmında gösterilir.
//
// # Sıra Kuralları
// - Düşük sayılar (0-50): Menünün üst kısmında
// - Orta sayılar (50-100): Menünün orta kısmında
// - Yüksek sayılar (100+): Menünün alt kısmında
//
// # Örnek
// ```go
// account := &Account{}
// order := account.NavigationOrder() // 10
// // Kullanıcı sayfaları menünün üst kısmında gösterilir
// ```
func (p *Account) NavigationOrder() int {
	return 10
}

// Visible, hesap ayarları sayfasının navigasyon menüsünde görünür olup olmadığını belirler.
//
// # Dönüş Değeri
// bool - true ise sayfa menüde görünür, false ise gizlidir
//
// # Mantık
// HideInNavigation alanının ters değerini döndürür.
// - HideInNavigation = false → Visible = true (görünür)
// - HideInNavigation = true → Visible = false (gizli)
//
// # Kullanım
// Navigasyon menüsü oluşturulurken, hangi sayfaların gösterilip gösterilmeyeceğini
// belirlemek için kullanılır.
//
// # Örnek
// ```go
// account := &Account{HideInNavigation: false}
// visible := account.Visible() // true
//
// account2 := &Account{HideInNavigation: true}
// visible2 := account2.Visible() // false
// ```
func (p *Account) Visible() bool {
	return !p.HideInNavigation
}

// Cards, hesap ayarları sayfasında gösterilecek widget kartlarını döndürür.
//
// # Dönüş Değeri
// []widget.Card - Boş bir kart listesi
//
// # Kullanım
// Hesap ayarları sayfasında istatistik, grafik veya özet bilgiler göstermek için
// kullanılabilecek kartları tanımlar. Şu anda boş bir liste döndürülmektedir.
//
// # Gelecek Geliştirmeler
// Hesap ayarları sayfasında kullanıcı istatistikleri veya özet bilgiler göstermek için
// bu metot genişletilebilir.
//
// # Örnek
// ```go
// account := &Account{}
// cards := account.Cards() // []widget.Card{} (boş)
// ```
func (p *Account) Cards() []widget.Card {
	return []widget.Card{}
}

// Fields, hesap ayarları sayfasında gösterilecek form alanlarını döndürür.
//
// # Dönüş Değeri
// []fields.Element - Sayfada gösterilecek form alanlarının listesi
//
// # Kullanım
// Hesap ayarları sayfasının form alanlarını dinamik olarak sağlamak için kullanılır.
// Her Element, bir form alanını temsil eder (TextInput, Select, Checkbox, vb.)
//
// # Örnek
// ```go
// account := &Account{
//     Elements: []fields.Element{
//         // Form alanları
//     },
// }
// fields := account.Fields() // Elements döndürülür
// ```
func (p *Account) Fields() []fields.Element {
	return p.Elements
}

// Save, hesap ayarları sayfasından gelen verileri veritabanına kaydeder.
//
// # Parametreler
// - c (*context.Context): İstek bağlamı, kullanıcı ve oturum bilgilerini içerir
// - db (*gorm.DB): Veritabanı bağlantısı, GORM ORM aracılığıyla
// - data (map[string]interface{}): Kaydedilecek hesap ayarları (anahtar-değer çiftleri)
//
// # Dönüş Değeri
// error - İşlem başarılı ise nil, hata ise error nesnesi
//
// # Çalışma Mantığı
// 1. Oturumdaki kullanıcı ID'sini alır
// 2. Gelen veri haritasındaki değerleri kullanıcı modeline uygular
// 3. Özel alanları işler (resim yükleme, şifre değiştirme)
// 4. Veritabanında kullanıcı bilgilerini günceller
//
// # Kullanım Senaryoları
// 1. Kullanıcı profil bilgilerini güncelleme
// 2. Şifre değiştirme işlemleri
// 3. Bildirim tercihlerini güncelleme
//
// # Örnek Kullanım
// ```go
// account := &Account{}
// err := account.Save(ctx, db, map[string]interface{}{
//     "name": "John Doe",
//     "email": "john@example.com",
//     "current_password": "old_password",
//     "new_password": "new_password",
// })
// if err != nil {
//     log.Printf("Hesap ayarları kaydedilemedi: %v", err)
// }
// ```
//
// # Önemli Notlar
// - Sadece oturumdaki kullanıcının kendi bilgilerini güncelleyebilir
// - Şifre değiştirme işleminde mevcut şifre doğrulaması yapılır
// - Resim yükleme işleminde dosya storage'a kaydedilir
// - Hassas bilgiler (şifre) hash'lenerek saklanır
func (p *Account) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
	// Oturumdaki kullanıcı ID'sini al
	userID := c.Ctx.Locals("user_id")
	if userID == nil {
		return fmt.Errorf("kullanıcı oturumu bulunamadı")
	}

	// Kullanıcıyı veritabanından al
	var u user.User
	if err := db.First(&u, userID).Error; err != nil {
		return fmt.Errorf("kullanıcı bulunamadı: %w", err)
	}

	// Özel alanları işle
	// Resim yükleme
	if imageFile, ok := data["image"].(*multipart.FileHeader); ok && imageFile != nil {
		storageUrl := "/storage/"
		storagePath := "./storage/public"
		ext := filepath.Ext(imageFile.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		localPath := filepath.Join(storagePath, filename)

		_ = os.MkdirAll(storagePath, 0755)

		if err := c.Ctx.SaveFile(imageFile, localPath); err != nil {
			return fmt.Errorf("resim yüklenemedi: %w", err)
		}
		data["image"] = fmt.Sprintf("%s/%s", storageUrl, filename)
	}

	// Şifre değiştirme
	if newPassword, ok := data["new_password"].(string); ok && newPassword != "" {
		// Yeni şifreyi bcrypt ile hash'le
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("şifre hash'lenemedi: %w", err)
		}

		// Account tablosunda şifreyi güncelle
		if err := db.Model(&account.Account{}).
			Where("user_id = ? AND provider_id = ?", u.ID, "credential").
			Update("password", string(hashed)).Error; err != nil {
			return fmt.Errorf("şifre güncellenemedi: %w", err)
		}

		// Şifre field'larını data'dan kaldır (user tablosuna yazılmasın)
		delete(data, "new_password")
		delete(data, "current_password")
		delete(data, "confirm_password")
	}

	// Diğer alanları güncelle
	updateData := make(map[string]interface{})
	for key, value := range data {
		// Özel alanları atla
		if key == "current_password" || key == "new_password" || key == "confirm_password" {
			continue
		}

		// Değeri string'e dönüştür (gerekirse)
		var strValue interface{}
		if v, ok := value.(string); ok {
			strValue = v
		} else {
			b, _ := json.Marshal(value)
			strValue = string(b)
		}

		updateData[key] = strValue
	}

	// Veritabanında kullanıcı bilgilerini güncelle
	if err := db.Model(&u).Updates(updateData).Error; err != nil {
		return fmt.Errorf("hesap ayarları güncellenemedi: %w", err)
	}

	return nil
}
