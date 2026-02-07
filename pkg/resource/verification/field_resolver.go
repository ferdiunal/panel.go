// Bu paket, doğrulama (verification) kaynağı için alan çözümleme işlevselliğini sağlar.
// Doğrulama varlıklarının UI panelinde nasıl gösterileceğini tanımlar.
package verification

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// Bu yapı, Verification kaynağının alanlarını çözmek ve yapılandırmak için kullanılır.
// VerificationFieldResolver, panel UI'sinde doğrulama verilerinin nasıl sunulacağını belirler.
//
// Kullanım Senaryosu:
// - Doğrulama kayıtlarının liste, detay ve form görünümlerinde gösterilecek alanları tanımlar
// - Her alan için okuma-yazma izinleri, görünüm türü ve görüntülenme konumlarını belirtir
// - Doğrulama token'ları, son kullanma tarihleri ve kimlik bilgilerini yönetir
//
// Örnek Kullanım:
//   resolver := &VerificationFieldResolver{}
//   fields := resolver.ResolveFields(ctx)
//   // fields, UI'da gösterilecek tüm alanları içerir
//
// Önemli Notlar:
// - Bu yapı, panel.go framework'ün alan çözümleme sisteminin bir parçasıdır
// - Boş bir yapı olarak tasarlanmıştır (receiver olarak kullanılır)
// - Tüm alanlar, fields.Schema türü kullanılarak tanımlanır
type VerificationFieldResolver struct{}

// Bu metod, Verification kaynağı için tüm UI alanlarını çözer ve döner.
// ResolveFields, panel uygulamasında doğrulama verilerinin nasıl gösterileceğini belirler.
//
// Parametreler:
//   - ctx (*context.Context): İstek bağlamı, kullanıcı bilgisi ve diğer bağlamsal verileri içerir
//
// Dönüş Değeri:
//   - []core.Element: Yapılandırılmış alan öğelerinin dilimi
//
// Alanlar Açıklaması:
//
// 1. ID Alanı:
//    - Anahtar: "id"
//    - Görünüm Türü: text (metin)
//    - Özellikler: ReadOnly() - Sadece okunabilir
//    - Görüntülenme: OnlyOnDetail() - Yalnızca detay sayfasında gösterilir
//    - Açıklama: Doğrulama kaydının benzersiz tanımlayıcısı
//
// 2. Identifier Alanı:
//    - Anahtar: "identifier"
//    - Görünüm Türü: text (metin)
//    - Özellikler: Düzenlenebilir
//    - Görüntülenme: OnList() + OnDetail() + OnForm() - Liste, detay ve form sayfalarında gösterilir
//    - Açıklama: Doğrulama için kullanılan tanımlayıcı (e-posta, telefon vb.)
//
// 3. Token Alanı:
//    - Anahtar: "token"
//    - Görünüm Türü: text (metin)
//    - Özellikler: ReadOnly() - Sadece okunabilir
//    - Görüntülenme: OnList() + OnDetail() - Liste ve detay sayfalarında gösterilir
//    - Açıklama: Doğrulama token'ı (hassas veri, düzenlenemez)
//
// 4. Expires At Alanı:
//    - Anahtar: "expires_at"
//    - Görünüm Türü: datetime (tarih-saat)
//    - Özellikler: Düzenlenebilir
//    - Görüntülenme: OnList() + OnDetail() + OnForm() - Liste, detay ve form sayfalarında gösterilir
//    - Açıklama: Doğrulama token'ının son kullanma tarihi ve saati
//
// 5. Created At Alanı:
//    - Anahtar: "created_at"
//    - Görünüm Türü: datetime (tarih-saat)
//    - Özellikler: ReadOnly() - Sadece okunabilir
//    - Görüntülenme: OnList() + OnDetail() - Liste ve detay sayfalarında gösterilir
//    - Açıklama: Doğrulama kaydının oluşturulma tarihi ve saati (sistem tarafından otomatik)
//
// 6. Updated At Alanı:
//    - Anahtar: "updated_at"
//    - Görünüm Türü: datetime (tarih-saat)
//    - Özellikler: ReadOnly() - Sadece okunabilir
//    - Görüntülenme: OnList() + OnDetail() - Liste ve detay sayfalarında gösterilir
//    - Açıklama: Doğrulama kaydının son güncellenme tarihi ve saati (sistem tarafından otomatik)
//
// Kullanım Örneği:
//   resolver := &VerificationFieldResolver{}
//   ctx := &context.Context{} // Bağlam oluştur
//   fields := resolver.ResolveFields(ctx)
//   for _, field := range fields {
//       fmt.Println(field) // Her alanı işle
//   }
//
// Önemli Notlar:
// - Tüm alanlar, fields.Schema yapısı kullanılarak tanımlanır
// - Method chaining kullanılarak alanların görüntülenme konumları ve özellikleri belirtilir
// - ReadOnly() alanlar, panel UI'sında düzenlenemez (salt okunur)
// - OnList(), OnDetail(), OnForm() metodları, alanın hangi sayfada gösterileceğini belirler
// - Tarih-saat alanları (datetime), otomatik olarak uygun format ile gösterilir
// - Token alanı hassas veri olduğu için ReadOnly olarak işaretlenmiştir
// - Props (properties) haritası, ek UI konfigürasyonları için kullanılabilir
//
// Döndürür:
// - Yapılandırılmış core.Element öğelerinin dilimi
func (r *VerificationFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		// ID Alanı: Doğrulama kaydının benzersiz tanımlayıcısı
		// Sadece okunabilir ve yalnızca detay sayfasında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		// Identifier Alanı: Doğrulama için kullanılan tanımlayıcı
		// Düzenlenebilir ve liste, detay, form sayfalarında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "identifier",
			Name:  "Identifier",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		// Token Alanı: Doğrulama token'ı (hassas veri)
		// Sadece okunabilir ve liste, detay sayfalarında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "token",
			Name:  "Token",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		// Expires At Alanı: Token'ın son kullanma tarihi ve saati
		// Düzenlenebilir ve liste, detay, form sayfalarında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "expires_at",
			Name:  "Expires At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		// Created At Alanı: Kaydın oluşturulma tarihi ve saati
		// Sadece okunabilir (sistem tarafından otomatik) ve liste, detay sayfalarında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "created_at",
			Name:  "Created At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		// Updated At Alanı: Kaydın son güncellenme tarihi ve saati
		// Sadece okunabilir (sistem tarafından otomatik) ve liste, detay sayfalarında gösterilir
		// Döndürür: - Yapılandırılmış fields.Schema pointer'ı
		(&fields.Schema{
			Key:   "updated_at",
			Name:  "Updated At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}
