// Bu paket, UUID (Universally Unique Identifier) oluşturma işlemlerini yönetir.
// UUID'ler, sistemde benzersiz tanımlayıcılar olarak kullanılır ve veri tabanında
// kayıtları, işlemleri veya diğer varlıkları tanımlamak için gereklidir.
//
// Paket, UUID v7 (zaman tabanlı, sıralanabilir UUID) oluşturmayı tercih eder,
// ancak v7 oluşturma başarısız olursa UUID v4 (rastgele UUID) ile geri döner.
// Bu yaklaşım, modern UUID standartlarının avantajlarından yararlanırken
// uyumluluk sağlar.
package uuid

import (
	"github.com/google/uuid"
)

// Bu fonksiyon, sistem içinde benzersiz tanımlayıcı oluşturmak için kullanılır.
//
// # Açıklama
// NewUUID, yeni bir UUID oluşturur. Öncelikle UUID v7 (zaman tabanlı, sıralanabilir)
// oluşturmaya çalışır. UUID v7, zaman bilgisini içerdiği için veri tabanında
// daha iyi performans sağlar ve kayıtları kronolojik sırada tutar.
// Eğer UUID v7 oluşturma başarısız olursa, UUID v4 (rastgele UUID) ile geri döner.
//
// # Kullanım Senaryoları
// - Yeni kullanıcı kaydı oluştururken benzersiz ID atama
// - Veri tabanında yeni kayıt oluştururken birincil anahtar olarak kullanma
// - API isteklerinde işlem takibi için benzersiz istek ID'si oluşturma
// - Dış sistemlerle entegrasyon sırasında benzersiz referans ID'si oluşturma
// - Loglama ve denetim izleri için benzersiz olay ID'si oluşturma
//
// # Parametreler
// Hiçbir parametre almaz.
//
// # Dönüş Değeri
// uuid.UUID: Oluşturulan UUID değeri. UUID v7 başarılı olursa v7 döner,
// aksi takdirde UUID v4 döner. Her zaman geçerli bir UUID döner.
//
// # Kullanım Örneği
// ```go
// // Basit kullanım
// newID := uuid.NewUUID()
// fmt.Println("Oluşturulan UUID:", newID.String())
//
// // Veri tabanı kaydında kullanma
// user := User{
//     ID:   uuid.NewUUID(),
//     Name: "Ahmet Yılmaz",
// }
// db.Create(&user)
//
// // API yanıtında kullanma
// response := map[string]interface{}{
//     "id":      uuid.NewUUID().String(),
//     "message": "Kayıt başarıyla oluşturuldu",
// }
// ```
//
// # Önemli Notlar
// - UUID v7, UUID v4'ten daha yeni bir standarttır ve daha iyi performans sağlar
// - UUID v7 zaman bilgisini içerdiği için sıralanabilir ve indeksleme için iyidir
// - Geri dönüş mekanizması (fallback), UUID v7 kütüphanesinin başarısız olması
//   durumunda sistem çökmesini önler
// - Döndürülen UUID her zaman geçerli ve benzersizdir
// - UUID'ler ağ üzerinde güvenli bir şekilde iletilmek için String() metoduyla
//   string'e dönüştürülebilir
//
// # Uyarılar
// - UUID'ler rastgele veya zaman tabanlı olsa da, kriptografik olarak güvenli
//   değildir. Güvenlik açısından kritik işlemler için ek güvenlik önlemleri alınmalıdır
// - UUID v7 oluşturma başarısız olursa UUID v4 döner, bu durumda sıralanabilirlik
//   kaybedilir ancak benzersizlik korunur
// - Çok yüksek çağrı hızında (saniyede milyonlarca çağrı) UUID çakışması teorik
//   olarak mümkün olsa da, pratikte ihmal edilebilir düzeydedir
func NewUUID() uuid.UUID {
	// UUID v7 oluşturmaya çalış. UUID v7, zaman bilgisini içeren modern UUID standardıdır.
	// Başarılı olursa, sıralanabilir ve performanslı bir UUID döner.
	id, err := uuid.NewV7()
	if err != nil {
		// UUID v7 oluşturma başarısız olursa, UUID v4 ile geri dön.
		// UUID v4 rastgele bir UUID'dir ve her zaman başarıyla oluşturulabilir.
		// Bu geri dönüş mekanizması, sistem stabilitesini sağlar.
		return uuid.New()
	}
	// Başarıyla oluşturulan UUID v7'yi döndür
	return id
}
