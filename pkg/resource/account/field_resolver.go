// Bu paket, Account kaynağı için alan çözücü (field resolver) işlevselliğini sağlar.
// Account nesnelerinin yönetim panelinde gösterilecek alanlarını tanımlar ve yapılandırır.
package account

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// Bu yapı, Account kaynağının alanlarını çözmek ve tanımlamak için kullanılır.
//
// AccountFieldResolver, panel uygulamasında Account nesnelerinin hangi alanlarının
// gösterileceğini, nasıl görüntüleneceğini ve hangi özelliklere sahip olacağını belirler.
//
// Kullanım Senaryoları:
// - Yönetim panelinde Account verilerini görüntülemek
// - Account alanlarının salt okunur veya düzenlenebilir olmasını kontrol etmek
// - Kullanıcı-Account ilişkisini yönetmek
// - Kimlik doğrulama sağlayıcı bilgilerini göstermek
//
// Örnek Kullanım:
//   resolver := &AccountFieldResolver{}
//   fields := resolver.ResolveFields(ctx)
//   // fields artık Account'un tüm alanlarını içerir
//
// Önemli Notlar:
// - Tüm alanlar varsayılan olarak salt okunur (ReadOnly) olarak yapılandırılmıştır
// - Hassas bilgiler (Access Token, Refresh Token) yorum satırı olarak depolanmıştır
// - User ilişkisi "users" tablosuna ve "user" alanına bağlıdır
type AccountFieldResolver struct{}

// Bu metod, Account kaynağının tüm alanlarını çözer ve döner.
//
// ResolveFields, yönetim panelinde Account nesnelerinin hangi alanlarının
// gösterileceğini belirleyen bir alan listesi oluşturur. Her alan, belirli
// görüntüleme ve etkileşim özellikleriyle yapılandırılmıştır.
//
// Parametreler:
//   - ctx (*context.Context): İstek bağlamı, kimlik doğrulama ve yetkilendirme
//     bilgilerini içerir. Gelecekte alan görünürlüğünü kontrol etmek için
//     kullanılabilir.
//
// Dönüş Değeri:
//   - []fields.Element: Account kaynağının tüm alanlarını içeren bir dilim.
//     Her element, bir alanın adı, veritabanı sütunu, türü ve özelliklerini
//     tanımlar.
//
// Döndürülen Alanlar:
//   1. ID: Hesabın benzersiz tanımlayıcısı (salt okunur)
//   2. Provider ID: Kimlik doğrulama sağlayıcısının hesap kimliği (salt okunur)
//   3. User: Hesabın ait olduğu kullanıcıya ilişki (salt okunur)
//   4. Provider: Kimlik doğrulama sağlayıcısının adı (salt okunur)
//   5. Created At: Hesabın oluşturulma tarihi ve saati (salt okunur)
//   6. Updated At: Hesabın son güncellenme tarihi ve saati (salt okunur)
//
// Kullanım Senaryoları:
// - Yönetim panelinde Account listesini görüntülemek
// - Account detay sayfasını oluşturmak
// - Kullanıcı-Account ilişkisini yönetmek
// - Kimlik doğrulama sağlayıcı bilgilerini göstermek
//
// Örnek Kullanım:
//   resolver := &AccountFieldResolver{}
//   ctx := &context.Context{} // Gerçek bağlam ile değiştirin
//   fields := resolver.ResolveFields(ctx)
//   for _, field := range fields {
//       fmt.Println(field.GetLabel()) // Alan adını yazdır
//   }
//
// Önemli Notlar:
// - Tüm alanlar salt okunur olarak yapılandırılmıştır (ReadOnly)
// - Hassas bilgiler (Access Token, Refresh Token) güvenlik nedeniyle yorum satırı olarak bırakılmıştır
// - User ilişkisi "users" tablosuna bağlıdır ve "user" alanı üzerinden erişilir
// - Tarih alanları (Created At, Updated At) otomatik olarak sistem tarafından yönetilir
// - Gelecekte, bağlam (ctx) parametresi kullanıcı izinlerine göre alanları filtrelemek için kullanılabilir
func (r *AccountFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		// Bu alan, Account kaynağının benzersiz tanımlayıcısını temsil eder.
		// Veritabanında otomatik olarak oluşturulan bir birincil anahtardır.
		// Salt okunur olarak yapılandırılmıştır çünkü sistem tarafından yönetilir.
		fields.ID("ID").ReadOnly(),

		// Bu alan, kimlik doğrulama sağlayıcısı tarafından atanan hesap kimliğini gösterir.
		// Örneğin: Google OAuth için Google User ID, GitHub için GitHub User ID vb.
		// Salt okunur olarak yapılandırılmıştır çünkü sağlayıcı tarafından belirlenir.
		// Veritabanı sütunu: providerId
		fields.Text("Provider", "providerId").ReadOnly(),

		// Bu alan yorum satırı olarak bırakılmıştır.
		// Gelecekte Account ID alanını göstermek için kullanılabilir.
		// fields.Text("Account ID", "accountId").ReadOnly(),

		// Bu alan, Account'un ait olduğu User (Kullanıcı) ile ilişkisini temsil eder.
		// Yönetim panelinde kullanıcıya hızlı erişim sağlar.
		//
		// İlişki Yapılandırması:
		// - Etiket: "User" (yönetim panelinde gösterilecek ad)
		// - Tablo: "users" (ilişkili kaynağın tablosu)
		// - Alan: "user" (Account struct'ında ilişkiyi temsil eden alan adı)
		//
		// Kullanım Senaryosu:
		// - Yönetim panelinde Account detaylarını görüntülerken, ilişkili User'ı göstermek
		// - Kullanıcı-Account ilişkisini yönetmek
		//
		// Salt okunur olarak yapılandırılmıştır çünkü Account oluşturulduktan sonra
		// kullanıcı değiştirilmemelidir.
		fields.Link("User", "users", "user").
			ReadOnly(),

		// Bu alan, Account'un ait olduğu kimlik doğrulama sağlayıcısının adını gösterir.
		// Örneğin: "google", "github", "microsoft" vb.
		// Salt okunur olarak yapılandırılmıştır çünkü Account oluşturulduktan sonra
		// sağlayıcı değiştirilmemelidir.
		// Veritabanı sütunu: provider
		fields.Text("Provider", "provider").
			ReadOnly(),

		// Bu alan yorum satırı olarak bırakılmıştır.
		// Güvenlik nedeniyle Access Token yönetim panelinde gösterilmemektedir.
		// Access Token'lar hassas bilgilerdir ve:
		// - Veritabanında şifrelenmiş olarak saklanmalıdır
		// - Yönetim panelinde gösterilmemelidir
		// - Sadece gerekli olduğunda sunucu tarafında kullanılmalıdır
		// fields.Text("Access Token", "accessToken").ReadOnly(),

		// Bu alan yorum satırı olarak bırakılmıştır.
		// Güvenlik nedeniyle Refresh Token yönetim panelinde gösterilmemektedir.
		// Refresh Token'lar çok hassas bilgilerdir ve:
		// - Veritabanında şifrelenmiş olarak saklanmalıdır
		// - Yönetim panelinde hiçbir zaman gösterilmemelidir
		// - Sadece token yenileme işlemi sırasında sunucu tarafında kullanılmalıdır
		// - Süresi dolmuş token'lar düzenli olarak temizlenmelidir
		// fields.Text("Refresh Token", "refreshToken").ReadOnly(),

		// Bu alan, Account'un oluşturulma tarihi ve saatini gösterir.
		// Sistem tarafından otomatik olarak ayarlanır ve değiştirilemez.
		// Salt okunur olarak yapılandırılmıştır.
		// Veritabanı sütunu: createdAt
		// Tarih Formatı: RFC3339 (ISO 8601)
		fields.DateTime("Created At", "createdAt").ReadOnly(),

		// Bu alan, Account'un son güncellenme tarihi ve saatini gösterir.
		// Sistem tarafından otomatik olarak güncellenir ve manuel olarak değiştirilemez.
		// Salt okunur olarak yapılandırılmıştır.
		// Veritabanı sütunu: updatedAt
		// Tarih Formatı: RFC3339 (ISO 8601)
		// Kullanım Senaryosu:
		// - Account bilgilerinin ne zaman son kez güncellendiğini görmek
		// - Denetim (audit) ve izleme (logging) amaçları için
		fields.DateTime("Updated At", "updatedAt").ReadOnly(),
	}
}
