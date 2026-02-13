// Package widget, panel.go uygulamasının dashboard widget'larını içerir.
// Bu paket, çeşitli metrik ve veri görüntüleme widget'larının tanımlanması ve yönetilmesini sağlar.
package widget

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// Bu yapı, dashboard'da sayısal değerleri (metrik) görüntülemek için kullanılan bir widget'ı temsil eder.
// Value widget'ı, veritabanından dinamik olarak veri çekerek, bu verileri dashboard'da
// "value-metric" bileşeni olarak gösterir. Genellikle toplam kullanıcı sayısı, sipariş sayısı,
// gelir gibi önemli metrikleri göstermek için kullanılır.
//
// Kullanım Senaryoları:
// - Toplam ürün sayısını göstermek
// - Aktif kullanıcı sayısını göstermek
// - Aylık gelir toplamını göstermek
// - Beklemede olan siparişlerin sayısını göstermek
//
// Örnek Kullanım:
//   widget := &Value{
//       Title: "Toplam Ürünler",
//       QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
//           var count int64
//           db.Model(&Product{}).Count(&count)
//           return count, nil
//       },
//   }
type Value struct {
	// Title, widget'ın başlığını ve adını belirtir.
	// Dashboard'da bu başlık, metrik değerinin üstünde gösterilir.
	// Örnek: "Toplam Ürünler", "Aktif Kullanıcılar", "Aylık Gelir"
	Title string

	// QueryFunc, veritabanından veri çekmek için kullanılan özel bir fonksiyondur.
	// Bu fonksiyon, panel context'i ve GORM veritabanı bağlantısını alarak,
	// int64 türünde bir değer döndürür. Hata durumunda error döndürür.
	//
	// Parametreler:
	// - ctx: *context.Context - Panel context'i, kullanıcı bilgisi ve diğer bağlamsal veriler içerir
	// - db: *gorm.DB - GORM veritabanı bağlantısı, sorguları çalıştırmak için kullanılır
	//
	// Dönüş Değerleri:
	// - int64 - Sorgu sonucu elde edilen sayısal değer
	// - error - Sorgu sırasında oluşan hata (başarılı ise nil)
	//
	// Örnek:
	//   QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
	//       var total int64
	//       if err := db.Model(&User{}).Count(&total).Error; err != nil {
	//           return 0, err
	//       }
	//       return total, nil
	//   }
	QueryFunc func(ctx *context.Context, db *gorm.DB) (int64, error)
}

// Bu metod, widget'ın adını (başlığını) döndürür.
// Dashboard sisteminde widget'ı tanımlamak ve referans almak için kullanılır.
//
// Dönüş Değeri:
// - string - Widget'ın başlığı (Title alanının değeri)
//
// Kullanım Senaryosu:
// Widget'ın adını almak için kullanılır, örneğin logging veya debug amaçlı.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   fmt.Println(widget.Name()) // Çıktı: "Toplam Ürünler"
func (w *Value) Name() string {
	return w.Title
}

// Bu metod, widget'ın frontend'te kullanılacak bileşen adını döndürür.
// Vue.js veya React gibi frontend framework'ünde, bu bileşen adı kullanılarak
// widget'ın görsel temsili render edilir.
//
// Dönüş Değeri:
// - string - Bileşen adı: "value-metric"
//
// Önemli Not:
// Bu değer sabit olup, tüm Value widget'ları için aynıdır.
// Frontend'te "value-metric" adında bir bileşen tanımlanmış olmalıdır.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   fmt.Println(widget.Component()) // Çıktı: "value-metric"
func (w *Value) Component() string {
	return "value-metric"
}

// Bu metod, widget'ın dashboard grid sisteminde kaplayacağı genişliği belirtir.
// Dashboard, responsive grid sistemi kullanarak widget'ları düzenler.
// "1/3" değeri, widget'ın dashboard genişliğinin 1/3'ünü kaplayacağını gösterir.
//
// Dönüş Değeri:
// - string - Genişlik oranı: "1/3" (dashboard genişliğinin üçte biri)
//
// Önemli Notlar:
// - Bu değer Tailwind CSS veya benzer CSS framework'ü ile uyumludur
// - Responsive tasarımda, farklı ekran boyutlarında farklı genişlikler kullanılabilir
// - Sabit değer olup, tüm Value widget'ları için aynıdır
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   fmt.Println(widget.Width()) // Çıktı: "1/3"
func (w *Value) Width() string {
	return "1/3"
}

// Bu metod, widget'ın türünü belirtir.
// Widget sistemi, farklı widget türlerini ayırt etmek için bu metodu kullanır.
// CardTypeValue, bu widget'ın bir değer/metrik widget'ı olduğunu gösterir.
//
// Dönüş Değeri:
// - CardType - Widget türü: CardTypeValue
//
// Kullanım Senaryosu:
// Widget'ın türüne göre farklı işlemler yapmak için kullanılır.
// Örneğin, widget'ı serialize ederken veya render ederken türü kontrol etmek.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   if widget.GetType() == CardTypeValue {
//       fmt.Println("Bu bir değer widget'ıdır")
//   }
func (w *Value) GetType() CardType {
	return CardTypeValue
}

// Bu metod, widget'ın veri çekerek sonucu döndürür.
// Dashboard'da widget'ı render etmek için gerekli olan verileri hazırlar.
// QueryFunc'ı çalıştırarak veritabanından veri alır ve bunu
// frontend'e gönderilebilecek bir harita (map) formatında döndürür.
//
// Parametreler:
// - ctx: *context.Context - Panel context'i, kullanıcı bilgisi ve diğer bağlamsal veriler içerir
// - db: *gorm.DB - GORM veritabanı bağlantısı, sorguları çalıştırmak için kullanılır
//
// Dönüş Değerleri:
// - interface{} - Aşağıdaki yapıda bir harita:
//   {
//       "value": int64,    // QueryFunc'tan dönen sayısal değer
//       "title": string,   // Widget'ın başlığı
//   }
// - error - Sorgu sırasında oluşan hata (başarılı ise nil)
//
// Hata Yönetimi:
// QueryFunc'ta hata oluşursa, bu metod nil ve error döndürür.
// Hata durumunda HandleError() metodu çağrılmalıdır.
//
// Kullanım Senaryosu:
// Dashboard API'si, widget'ı render etmek için bu metodu çağırır.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   data, err := widget.Resolve(ctx, db)
//   if err != nil {
//       errorData := widget.HandleError(err)
//       // Hata verilerini frontend'e gönder
//   } else {
//       // data'yı frontend'e gönder
//   }
func (w *Value) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	val, err := w.QueryFunc(ctx, db)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"value": val,
		"title": w.Title,
	}, nil
}

// Bu metod, widget'ı render ederken hata oluştuğunda çağrılır.
// Hata bilgisini frontend'e gönderilebilecek bir format'ta hazırlar.
// Kullanıcıya hata mesajını göstermek için gerekli verileri sağlar.
//
// Parametreler:
// - err: error - Oluşan hata nesnesi
//
// Dönüş Değeri:
// - map[string]interface{} - Aşağıdaki yapıda bir harita:
//   {
//       "error": string,           // Hata mesajı (err.Error())
//       "title": string,           // Widget'ın başlığı
//       "type": CardType,          // Widget türü (CardTypeValue)
//   }
//
// Kullanım Senaryosu:
// Resolve() metodu hata döndürdüğünde, bu metod çağrılarak
// hata bilgisi frontend'e gönderilir.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   data, err := widget.Resolve(ctx, db)
//   if err != nil {
//       errorResponse := widget.HandleError(err)
//       // errorResponse'u JSON olarak frontend'e gönder
//   }
func (w *Value) HandleError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
		"title": w.Title,
		"type":  CardTypeValue,
	}
}

// Bu metod, widget'ın meta verilerini döndürür.
// Widget'ın yapılandırması ve özellikleri hakkında bilgi sağlar.
// Dashboard sisteminde widget'ı tanımlamak ve yönetmek için kullanılır.
//
// Dönüş Değeri:
// - map[string]interface{} - Aşağıdaki yapıda bir harita:
//   {
//       "name": string,           // Widget'ın adı (Title)
//       "component": string,      // Bileşen adı ("value-metric")
//       "width": string,          // Genişlik oranı ("1/3")
//       "type": CardType,         // Widget türü (CardTypeValue)
//   }
//
// Kullanım Senaryosu:
// Dashboard konfigürasyonunu almak veya widget'ı tanımlamak için kullanılır.
// Frontend'te widget'ın özelliklerini göstermek için de kullanılabilir.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   metadata := widget.GetMetadata()
//   fmt.Printf("Widget Adı: %s\n", metadata["name"])
//   fmt.Printf("Bileşen: %s\n", metadata["component"])
func (w *Value) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":      w.Title,
		"component": "value-metric",
		"width":     "1/3",
		"type":      CardTypeValue,
	}
}

// Bu metod, widget'ı JSON formatında serialize eder.
// API yanıtlarında widget'ı göndermek için kullanılır.
// Frontend'te widget'ı render etmek için gerekli tüm bilgileri içerir.
//
// Dönüş Değeri:
// - map[string]interface{} - Aşağıdaki yapıda bir harita:
//   {
//       "component": string,      // Bileşen adı ("value-metric")
//       "title": string,          // Widget'ın başlığı
//       "width": string,          // Genişlik oranı ("1/3")
//       "type": CardType,         // Widget türü (CardTypeValue)
//   }
//
// Önemli Not:
// Bu metod, widget'ın yapılandırmasını döndürür, veri değerini değil.
// Veri değeri Resolve() metodu tarafından döndürülür.
//
// Kullanım Senaryosu:
// Dashboard konfigürasyonunu frontend'e göndermek için kullanılır.
// Widget'ın statik özelliklerini JSON olarak serialize etmek için.
//
// Örnek:
//   widget := NewCountWidget("Toplam Ürünler", &Product{})
//   jsonData := widget.JsonSerialize()
//   // jsonData'yı JSON olarak frontend'e gönder
func (w *Value) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": "value-metric",
		"title":     w.Title,
		"width":     "1/3",
		"type":      CardTypeValue,
	}
}

// ============================================================================
// YARDIMCI FONKSİYONLAR (Helpers)
// ============================================================================

// Bu fonksiyon, veritabanında bir modelin toplam kayıt sayısını gösteren
// bir Value widget'ı oluşturur. Yaygın olarak kullanılan bir yardımcı fonksiyondur.
// Örneğin, toplam ürün sayısı, toplam kullanıcı sayısı gibi metrikleri
// hızlı bir şekilde dashboard'a eklemek için kullanılır.
//
// Parametreler:
// - title: string - Widget'ın başlığı (örn: "Toplam Ürünler")
// - model: interface{} - GORM modeli (örn: &Product{}, &User{})
//   Bu parametre, GORM'a hangi tablodan veri çekeceğini söyler.
//
// Dönüş Değeri:
// - *Value - Yapılandırılmış Value widget pointer'ı
//   Bu widget, çağrıldığında veritabanından modelin toplam sayısını döndürür.
//
// Kullanım Senaryoları:
// - Toplam ürün sayısını göstermek
// - Toplam kullanıcı sayısını göstermek
// - Toplam sipariş sayısını göstermek
// - Herhangi bir modelin toplam kayıt sayısını göstermek
//
// Önemli Notlar:
// - QueryFunc, GORM'un Count() metodunu kullanarak toplam sayıyı alır
// - Hata durumunda (örn: tablo bulunamadı), error döndürülür
// - Model parametresi, GORM tarafından tanınan bir struct olmalıdır
//
// Örnek Kullanım:
//   // Toplam ürün sayısını gösteren widget oluştur
//   productWidget := NewCountWidget("Toplam Ürünler", &Product{})
//
//   // Toplam kullanıcı sayısını gösteren widget oluştur
//   userWidget := NewCountWidget("Aktif Kullanıcılar", &User{})
//
//   // Widget'ı dashboard'a ekle
//   dashboard.AddWidget(productWidget)
//
//   // Widget'ın verilerini al
//   data, err := productWidget.Resolve(ctx, db)
//   if err != nil {
//       errorData := productWidget.HandleError(err)
//       // Hata verilerini frontend'e gönder
//   }
func NewCountWidget(title string, model interface{}) *Value {
	return &Value{
		Title: title,
		QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
			var total int64
			if err := db.Model(model).Count(&total).Error; err != nil {
				return 0, err
			}
			return total, nil
		},
	}
}
