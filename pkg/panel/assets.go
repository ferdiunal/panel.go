// Package panel, Panel uygulamasının temel bileşenlerini ve varlıklarını yönetir.
//
// Bu paket, Panel UI'nin statik dosyalarını (HTML, CSS, JavaScript vb.) yönetmek için
// Go'nun embed özelliğini kullanır. Dosyalar derleme zamanında binary'ye gömülür,
// böylece runtime'da harici dosya sistemine bağımlılık olmaz.
//
// # Temel Özellikler
//
// - **Gömülü Dosya Sistemi**: UI dosyaları binary'ye derleme zamanında gömülür
// - **Esnek Dosya Sistemi Seçimi**: Geliştirme ve üretim ortamları için farklı dosya sistemi kaynakları
// - **Performans**: Gömülü dosyalar disk I/O'dan kaçınır, daha hızlı yükleme sağlar
//
// # Kullanım Senaryoları
//
// 1. **Üretim Ortamı**: Gömülü dosya sistemi kullanılarak bağımsız bir binary dağıtılır
// 2. **Geliştirme Ortamı**: İsteğe bağlı olarak os.DirFS kullanılarak dosyalar disk'ten yüklenir
// 3. **Docker Dağıtımı**: Tüm UI dosyaları container image'ında bulunur
//
// # Örnek Kullanım
//
//	// Üretim ortamında gömülü dosyaları kullan
//	fs, err := GetFileSystem(true)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Geliştirme ortamında disk'ten dosyaları yükle
//	fs, err := GetFileSystem(false)
//	if err != nil {
//		// Fallback olarak os.DirFS kullan
//		fs = os.DirFS("./pkg/panel/ui")
//	}
package panel

import (
	"embed"
	"io/fs"
)

// embedFS, Go'nun embed özelliği kullanılarak ui/ dizinindeki tüm dosyaları içerir.
//
// # Detaylar
//
// - **Derleme Zamanı Gömülmesi**: //go:embed direktifi, derleme sırasında ui/ dizinindeki
//   tüm dosyaları (HTML, CSS, JavaScript, resimler vb.) binary'ye gömülmesini sağlar
// - **Özyinelemeli Gömülme**: ui/* deseni, ui/ dizini ve tüm alt dizinlerindeki dosyaları içerir
// - **Boyut Etkisi**: Gömülü dosyalar binary boyutunu artırır, ancak dağıtım kolaylığı sağlar
//
// # Avantajlar
//
// - Harici dosya bağımlılığı yok
// - Dağıtım sırasında dosya kaybı riski yok
// - Daha hızlı uygulama başlatma
// - Tek bir binary dosyası ile dağıtım
//
// # Dezavantajlar
//
// - Binary boyutu artar
// - Dosya güncellemeleri için yeniden derleme gerekir
// - Geliştirme sırasında sıcak yenileme (hot reload) zor olabilir
//
// # Önemli Notlar
//
// - Dosyalar sadece okunabilir erişime sahiptir
// - Gömülü dosya sistemi thread-safe'dir
// - Dosya izinleri korunmaz, tüm dosyalar okunabilir olur
//
//go:embed ui/*
var embedFS embed.FS

// GetFileSystem, Panel UI dosyalarına erişmek için uygun bir dosya sistemi döndürür.
//
// # Parametreler
//
// - `useEmbed` (bool): Eğer true ise, gömülü dosya sistemini kullanır.
//   Eğer false ise, nil döndürür ve çağıran taraf os.DirFS gibi alternatif bir
//   dosya sistemi kullanmalıdır.
//
// # Dönüş Değerleri
//
// - `fs.FS`: Dosya sistemi arayüzü. Gömülü dosyalara erişim sağlar.
// - `error`: İşlem sırasında oluşan hata. Başarılı durumda nil döndürülür.
//
// # Kullanım Senaryoları
//
// ## Senaryo 1: Üretim Ortamında Gömülü Dosyaları Kullan
//
//	fs, err := GetFileSystem(true)
//	if err != nil {
//		log.Fatalf("Dosya sistemi alınamadı: %v", err)
//	}
//	// fs'i HTTP sunucusunda kullan
//	http.Handle("/ui/", http.FileServer(http.FS(fs)))
//
// ## Senaryo 2: Geliştirme Ortamında Disk'ten Dosyaları Yükle
//
//	fs, err := GetFileSystem(false)
//	if err != nil || fs == nil {
//		// Fallback: Disk'ten dosyaları yükle
//		fs = os.DirFS("./pkg/panel/ui")
//	}
//	// fs'i HTTP sunucusunda kullan
//	http.Handle("/ui/", http.FileServer(http.FS(fs)))
//
// ## Senaryo 3: Ortam Değişkenine Göre Karar Ver
//
//	useEmbedded := os.Getenv("USE_EMBEDDED_UI") == "true"
//	fs, err := GetFileSystem(useEmbedded)
//	if err != nil || fs == nil {
//		fs = os.DirFS("./pkg/panel/ui")
//	}
//
// # Detaylı Açıklama
//
// Bu fonksiyon, Panel uygulamasının UI dosyalarını yönetmek için iki farklı strateji sunar:
//
// 1. **Gömülü Dosya Sistemi (useEmbed=true)**:
//    - Binary'ye gömülü ui/ dizinini kullanır
//    - fs.Sub() kullanarak "ui" alt dizinine erişim sağlar
//    - Üretim ortamında tercih edilir
//    - Harici dosya bağımlılığı yoktur
//
// 2. **Disk Dosya Sistemi (useEmbed=false)**:
//    - nil döndürür
//    - Çağıran taraf os.DirFS() kullanarak disk'ten dosyaları yükler
//    - Geliştirme ortamında tercih edilir
//    - Sıcak yenileme (hot reload) için uygun
//
// # Avantajlar
//
// - **Esnek Dağıtım**: Üretim ve geliştirme ortamları için farklı stratejiler
// - **Performans**: Gömülü dosyalar disk I/O'dan kaçınır
// - **Güvenlik**: Dosyalar binary'ye gömülü olduğu için değiştirilmesi zor
// - **Basitlik**: Tek bir binary dosyası ile dağıtım
//
// # Dezavantajlar
//
// - **Binary Boyutu**: Gömülü dosyalar binary boyutunu artırır
// - **Geliştirme Karmaşıklığı**: Fallback mekanizması gerekebilir
// - **Dosya Güncellemeleri**: Üretim ortamında dosya güncellemeleri için yeniden derleme gerekir
//
// # Önemli Notlar
//
// - fs.Sub() fonksiyonu, embedFS'in "ui" alt dizinine erişim sağlar
// - Döndürülen fs.FS arayüzü thread-safe'dir
// - Dosya izinleri korunmaz, tüm dosyalar okunabilir olur
// - nil döndürüldüğünde, çağıran taraf mutlaka fallback mekanizması uygulamalıdır
//
// # Hata Yönetimi
//
// - fs.Sub() başarısız olursa, hata döndürülür
// - useEmbed=false olduğunda, hata asla döndürülmez (nil, nil)
// - Çağıran taraf, nil fs'i kontrol etmeli ve fallback uygulamalıdır
//
// # Performans Özellikleri
//
// - Gömülü dosya sistemi: O(1) erişim süresi (bellekten)
// - Disk dosya sistemi: O(n) erişim süresi (disk I/O'ya bağlı)
// - fs.Sub() işlemi: Çok hızlı, sadece referans oluşturur
//
// # Güvenlik Özellikleri
//
// - Gömülü dosyalar runtime'da değiştirilemez
// - Dosya izinleri korunmaz, tüm dosyalar okunabilir olur
// - Disk dosya sistemi kullanıldığında, işletim sistemi izinleri uygulanır
//
// # Uyumluluk
//
// - Go 1.16+: embed paketi gereklidir
// - Tüm işletim sistemleri desteklenir (Windows, macOS, Linux)
// - Cross-platform derleme desteklenir
//
// # İlgili Fonksiyonlar
//
// - fs.Sub(): Dosya sisteminin alt dizinine erişim sağlar
// - os.DirFS(): Disk'ten dosyaları yükler
// - http.FileServer(): HTTP sunucusunda dosyaları sunar
//
func GetFileSystem(useEmbed bool) (fs.FS, error) {
	if useEmbed {
		// embedFS'in "ui" alt dizinine erişim sağla.
		// fs.Sub() fonksiyonu, embedFS'in "ui" dizinini kök olarak ayarlar.
		// Böylece http.FileServer() kullanıldığında, dosyalar doğru yoldan sunulur.
		//
		// Örnek:
		// - embedFS'te: ui/index.html
		// - fs.Sub() sonrası: index.html
		return fs.Sub(embedFS, "ui")
	}
	// Geliştirme ortamında, çağıran taraf os.DirFS() kullanarak disk'ten dosyaları yükler.
	// Bu, sıcak yenileme (hot reload) ve hızlı geliştirme için uygun.
	return nil, nil
}
