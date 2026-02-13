// Bu paket, oturum (session) kaynağı için alan çözümleme işlevselliğini sağlar.
// Session yönetimi, kullanıcı oturumlarının oluşturulması, güncellenmesi ve görüntülenmesi
// için gerekli olan tüm alanları tanımlar.
package session

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// Bu yapı, Session kaynağının alanlarını çözmek ve tanımlamak için kullanılır.
// SessionFieldResolver, panel uygulamasında oturum verilerinin nasıl görüntüleneceğini
// ve etkileşime gireceğini belirler.
//
// Kullanım Senaryosu:
// - Admin panelinde oturum listesini görüntülemek
// - Oturum detaylarını incelemek
// - Oturum bilgilerini yönetmek
//
// Örnek:
//   resolver := &SessionFieldResolver{}
//   fields := resolver.ResolveFields(ctx)
//   // fields, oturum alanlarının tam listesini içerir
type SessionFieldResolver struct{}

// Bu metod, Session kaynağı için tüm kullanılabilir alanları çözer ve döndürür.
// Panel uygulamasında oturum verilerinin nasıl sunulacağını tanımlar.
//
// Parametreler:
//   - ctx (*context.Context): İstek bağlamı, kullanıcı bilgisi ve diğer bağlamsal
//     verileri içerir. Alan çözümleme sırasında yetkilendirme ve filtreleme için
//     kullanılabilir.
//
// Dönüş Değeri:
//   - []core.Element: Oturum kaynağı için tanımlanan tüm alanların listesi.
//     Her alan, panel UI'de nasıl görüntüleneceğini belirten yapılandırmalar içerir.
//
// Döndürülen Alanlar:
//   1. ID: Oturum benzersiz tanımlayıcısı (salt okunur, sadece detay sayfasında)
//   2. Token: Oturum token'ı (salt okunur, liste ve detay sayfalarında)
//   3. User: Oturuma ait kullanıcı ilişkisi (liste ve detay sayfalarında)
//   4. IP Address: İstemci IP adresi (liste ve detay sayfalarında)
//   5. User Agent: İstemci tarayıcı bilgisi (sadece detay sayfasında)
//   6. Expires At: Oturum sona erme tarihi (liste ve detay sayfalarında)
//   7. Created At: Oturum oluşturma tarihi (salt okunur, liste ve detay sayfalarında)
//   8. Updated At: Oturum güncelleme tarihi (salt okunur, liste ve detay sayfalarında)
//
// Kullanım Örneği:
//   resolver := &SessionFieldResolver{}
//   ctx := &context.Context{} // İstek bağlamı
//   elements := resolver.ResolveFields(ctx)
//   for _, element := range elements {
//       // Her alan, panel UI'de işlenir
//       fmt.Println(element)
//   }
//
// Önemli Notlar:
//   - Tüm alanlar method chaining kullanılarak yapılandırılır
//   - ReadOnly() ile işaretlenen alanlar düzenlenemez
//   - OnList() alanı liste görünümünde gösterir
//   - OnDetail() alanı detay görünümünde gösterir
//   - OnlyOnDetail() alanı sadece detay görünümünde gösterir
//   - Link() alanı ilişkili kaynağa bağlantı oluşturur
//
// Uyarılar:
//   - Token alanı hassas bilgi içerebilir, erişim kontrolü önemlidir
//   - User Agent ve IP Address, gizlilik düzenlemeleri açısından dikkat gerektirir
func (r *SessionFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		// ID Alanı: Oturum benzersiz tanımlayıcısı
		// - Salt okunur (ReadOnly): Sistem tarafından otomatik olarak atanır, değiştirilemez
		// - Sadece detay sayfasında görüntülenir (OnlyOnDetail)
		// - Veritabanında birincil anahtar olarak kullanılır
		fields.ID("ID").ReadOnly().OnlyOnDetail(),

		// Token Alanı: Oturum kimlik doğrulama token'ı
		// - Salt okunur (ReadOnly): Güvenlik nedeniyle değiştirilemez
		// - Liste ve detay sayfalarında görüntülenir (OnList, OnDetail)
		// - Kullanıcı oturumunu tanımlamak için kullanılır
		// - Hassas bilgi: Erişim kontrolü ve loglama önemlidir
		fields.Text("Token", "token").ReadOnly().OnList().OnDetail(),

		// User İlişkisi: Oturuma ait kullanıcı
		// - Link alanı, ilişkili kaynağa bağlantı oluşturur
		// - Parametreler: ("Görüntü Adı", "Tablo Adı", "Alan Adı")
		// - "users" tablosundaki "user" alanı ile ilişkilidir
		// - Liste ve detay sayfalarında görüntülenir
		// - Tıklanarak ilişkili kullanıcı detaylarına gidilebilir
		fields.Link("User", "users", "user").OnList().OnDetail(),

		// IP Address Alanı: İstemci IP adresi
		// - Oturumun oluşturulduğu istemcinin IP adresini içerir
		// - Liste ve detay sayfalarında görüntülenir
		// - Güvenlik denetimi ve oturum izleme için kullanılır
		// - Gizlilik: IP adresleri kişisel veri olarak kabul edilebilir
		fields.Text("IP Address", "ipAddress").OnList().OnDetail(),

		// User Agent Alanı: İstemci tarayıcı ve işletim sistemi bilgisi
		// - Oturumun oluşturulduğu tarayıcı ve cihaz bilgisini içerir
		// - Sadece detay sayfasında görüntülenir (OnDetail)
		// - Oturum güvenliği ve cihaz tanımlama için kullanılır
		// - Örnek: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
		fields.Text("User Agent", "userAgent").OnDetail(),

		// Expires At Alanı: Oturum sona erme tarihi ve saati
		// - DateTime alanı, tarih ve saat bilgisini içerir
		// - Liste ve detay sayfalarında görüntülenir
		// - Oturumun ne zaman geçersiz hale geleceğini gösterir
		// - Otomatik oturum kapatma için kullanılır
		fields.DateTime("Expires At", "expiresAt").OnList().OnDetail(),

		// Created At Alanı: Oturum oluşturma tarihi ve saati
		// - Salt okunur (ReadOnly): Sistem tarafından otomatik olarak atanır
		// - DateTime alanı, tarih ve saat bilgisini içerir
		// - Liste ve detay sayfalarında görüntülenir
		// - Oturumun ne zaman başlatıldığını gösterir
		// - Denetim ve güvenlik izleme için önemlidir
		fields.DateTime("Created At", "createdAt").ReadOnly().OnList().OnDetail(),

		// Updated At Alanı: Oturum son güncelleme tarihi ve saati
		// - Salt okunur (ReadOnly): Sistem tarafından otomatik olarak güncellenir
		// - DateTime alanı, tarih ve saat bilgisini içerir
		// - Liste ve detay sayfalarında görüntülenir
		// - Oturumun son etkinlik zamanını gösterir
		// - Oturum aktivitesi izleme için kullanılır
		fields.DateTime("Updated At", "updatedAt").ReadOnly().OnList().OnDetail(),
	}
}
