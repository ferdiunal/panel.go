package core

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

/// # VisibilityFunc
///
/// Bu fonksiyon tipi, alanların, butonların ve diğer UI elementlerinin görünürlüğünü
/// dinamik olarak kontrol etmek için kullanılır. ResourceContext'e göre bir boolean
/// değer döndürerek elementin gösterilip gösterilmeyeceğine karar verir.
///
/// ## Kullanım Senaryoları
///
/// - **Yetki Bazlı Görünürlük**: Kullanıcının rolüne veya izinlerine göre alanları göster/gizle
/// - **Durum Bazlı Görünürlük**: Kaydın durumuna göre (örn: yayınlanmış/taslak) belirli alanları göster
/// - **Koşullu Alan Gösterimi**: Başka bir alanın değerine bağlı olarak alanları göster/gizle
/// - **Dinamik Form Yapısı**: Kullanıcı etkileşimine göre form yapısını değiştir
/// - **Çoklu Tenant Sistemler**: Tenant'a özel alanları göster/gizle
///
/// ## Parametreler
///
/// - `ctx`: `*ResourceContext` - Mevcut istek bağlamı, kullanıcı bilgileri, kaynak durumu ve
///   diğer bağlamsal verileri içerir
///
/// ## Dönüş Değeri
///
/// - `bool`: Element görünür olmalıysa `true`, gizlenmeli ise `false` döndürür
///
/// ## Kullanım Örnekleri
///
/// ### Örnek 1: Yetki Bazlı Görünürlük
///
/// ```go
/// field := fields.Text("secret_key").
///     Visible(func(ctx *core.ResourceContext) bool {
///         // Sadece admin kullanıcılar bu alanı görebilir
///         return ctx.User != nil && ctx.User.Role == "admin"
///     })
/// ```
///
/// ### Örnek 2: Durum Bazlı Görünürlük
///
/// ```go
/// field := fields.Text("published_at").
///     Visible(func(ctx *core.ResourceContext) bool {
///         // Sadece yayınlanmış kayıtlarda göster
///         if ctx.Record == nil {
///             return false
///         }
///         status, _ := ctx.Record["status"].(string)
///         return status == "published"
///     })
/// ```
///
/// ### Örnek 3: Sayfa Tipine Göre Görünürlük
///
/// ```go
/// field := fields.Text("internal_notes").
///     Visible(func(ctx *core.ResourceContext) bool {
///         // Sadece detay ve düzenleme sayfalarında göster
///         return ctx.Page == core.PageDetail || ctx.Page == core.PageEdit
///     })
/// ```
///
/// ### Örnek 4: Koşullu Alan Gösterimi
///
/// ```go
/// field := fields.Text("discount_code").
///     Visible(func(ctx *core.ResourceContext) bool {
///         // Sadece indirim aktif ise göster
///         if ctx.Record == nil {
///             return true // Yeni kayıtlarda varsayılan olarak göster
///         }
///         hasDiscount, _ := ctx.Record["has_discount"].(bool)
///         return hasDiscount
///     })
/// ```
///
/// ## Avantajlar
///
/// - **Esneklik**: Karmaşık görünürlük mantığını kolayca uygulayabilirsiniz
/// - **Güvenlik**: Hassas verileri yetkisiz kullanıcılardan gizleyebilirsiniz
/// - **Kullanıcı Deneyimi**: Kullanıcıya sadece ilgili alanları göstererek form karmaşıklığını azaltır
/// - **Dinamik UI**: Runtime'da UI yapısını değiştirebilirsiniz
/// - **Tip Güvenliği**: Go'nun tip sistemi sayesinde derleme zamanı kontrolü
///
/// ## Dezavantajlar
///
/// - **Performans**: Her render'da fonksiyon çağrılır, karmaşık mantık performansı etkileyebilir
/// - **Test Edilebilirlik**: Closure kullanımı test yazmayı zorlaştırabilir
/// - **Hata Ayıklama**: Dinamik görünürlük mantığı hata ayıklamayı zorlaştırabilir
///
/// ## Önemli Notlar
///
/// ⚠️ **Güvenlik Uyarısı**: Görünürlük kontrolü sadece UI katmanında çalışır. Backend'de
/// mutlaka yetki kontrolü yapılmalıdır. Bir alanı gizlemek, kullanıcının o veriye erişimini
/// engellemez!
///
/// ⚠️ **Performans**: Fonksiyon her render'da çağrılır. Veritabanı sorguları veya ağır
/// hesaplamalar yapmaktan kaçının. Gerekirse sonuçları ResourceContext'te önbellekleyin.
///
/// ⚠️ **Nil Kontrolleri**: `ctx.Record`, `ctx.User` gibi alanlar nil olabilir. Özellikle
/// yeni kayıt oluşturma (create) sayfalarında `ctx.Record` nil olacaktır.
///
/// 💡 **İpucu**: Karmaşık görünürlük mantığını ayrı fonksiyonlara çıkararak kodunuzu
/// daha okunabilir ve test edilebilir hale getirebilirsiniz.
///
/// ## İlgili Dokümantasyon
///
/// - `docs/Fields.md` - Alan yapılandırması ve görünürlük örnekleri
/// - `ResourceContext` - Bağlam yapısı ve kullanılabilir alanlar
///
/// ## Benzer Kavramlar
///
/// - **Authorization Middleware**: Backend seviyesinde yetki kontrolü
/// - **Field Dependencies**: Alanlar arası bağımlılıklar
/// - **Conditional Rendering**: React/Vue'daki koşullu render mantığı
type VisibilityFunc func(ctx *ResourceContext) bool

/// # StorageCallbackFunc
///
/// Bu fonksiyon tipi, dosya yükleme işlemlerinde özel depolama stratejileri uygulamak
/// için kullanılır. Yüklenen dosyayı alır, istenen depolama sistemine kaydeder ve
/// dosyanın erişim yolunu veya URL'ini döndürür.
///
/// ## Kullanım Senaryoları
///
/// - **Yerel Disk Depolama**: Dosyaları sunucu diskine kaydetme
/// - **Cloud Storage**: AWS S3, Google Cloud Storage, Azure Blob Storage gibi bulut servislerine yükleme
/// - **CDN Entegrasyonu**: Dosyaları CDN'e yükleme ve CDN URL'i döndürme
/// - **Görsel İşleme**: Yüklenen görselleri yeniden boyutlandırma, optimize etme
/// - **Virus Tarama**: Dosyaları güvenlik kontrolünden geçirme
/// - **Metadata Ekleme**: Dosyalara özel metadata (EXIF, watermark) ekleme
/// - **Çoklu Depolama**: Dosyayı birden fazla yere (backup) kaydetme
///
/// ## Parametreler
///
/// - `c`: `*fiber.Ctx` - Fiber HTTP context, request bilgileri, headers, kullanıcı bilgisi vb.
///   içerir
/// - `file`: `*multipart.FileHeader` - Yüklenen dosyanın metadata'sı (isim, boyut, MIME type)
///
/// ## Dönüş Değerleri
///
/// - `string`: Dosyanın kaydedildiği yol veya erişim URL'i. Bu değer veritabanına kaydedilir
///   ve daha sonra dosyaya erişim için kullanılır
/// - `error`: İşlem başarısız olursa hata döndürür (disk dolu, izin hatası, network hatası vb.)
///
/// ## Kullanım Örnekleri
///
/// ### Örnek 1: Yerel Disk Depolama
///
/// ```go
/// field := fields.File("avatar").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         // Benzersiz dosya adı oluştur
///         filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
///         uploadPath := filepath.Join("uploads", "avatars", filename)
///
///         // Dosyayı kaydet
///         if err := c.SaveFile(file, uploadPath); err != nil {
///             return "", fmt.Errorf("dosya kaydedilemedi: %w", err)
///         }
///
///         // Public URL döndür
///         return "/uploads/avatars/" + filename, nil
///     })
/// ```
///
/// ### Örnek 2: AWS S3 Yükleme
///
/// ```go
/// field := fields.File("document").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         // Dosyayı aç
///         src, err := file.Open()
///         if err != nil {
///             return "", err
///         }
///         defer src.Close()
///
///         // S3'e yükle
///         key := fmt.Sprintf("documents/%s/%s",
///             time.Now().Format("2006/01/02"),
///             file.Filename)
///
///         _, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
///             Bucket: aws.String("my-bucket"),
///             Key:    aws.String(key),
///             Body:   src,
///             ContentType: aws.String(file.Header.Get("Content-Type")),
///         })
///
///         if err != nil {
///             return "", fmt.Errorf("S3 yükleme hatası: %w", err)
///         }
///
///         // S3 URL döndür
///         return fmt.Sprintf("https://my-bucket.s3.amazonaws.com/%s", key), nil
///     })
/// ```
///
/// ### Örnek 3: Görsel İşleme ve Optimizasyon
///
/// ```go
/// field := fields.Image("product_image").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         // Dosyayı aç
///         src, err := file.Open()
///         if err != nil {
///             return "", err
///         }
///         defer src.Close()
///
///         // Görseli decode et
///         img, _, err := image.Decode(src)
///         if err != nil {
///             return "", fmt.Errorf("görsel decode hatası: %w", err)
///         }
///
///         // Yeniden boyutlandır (max 1200px genişlik)
///         resized := resize.Thumbnail(1200, 1200, img, resize.Lanczos3)
///
///         // Optimize edilmiş dosyayı kaydet
///         filename := fmt.Sprintf("%s.webp", uuid.New().String())
///         outputPath := filepath.Join("uploads", "products", filename)
///
///         out, err := os.Create(outputPath)
///         if err != nil {
///             return "", err
///         }
///         defer out.Close()
///
///         // WebP formatında kaydet (daha küçük boyut)
///         if err := webp.Encode(out, resized, &webp.Options{Quality: 85}); err != nil {
///             return "", err
///         }
///
///         return "/uploads/products/" + filename, nil
///     })
/// ```
///
/// ### Örnek 4: Virus Tarama ve Güvenlik Kontrolü
///
/// ```go
/// field := fields.File("attachment").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         // Dosya boyutu kontrolü (max 10MB)
///         if file.Size > 10*1024*1024 {
///             return "", fmt.Errorf("dosya çok büyük (max 10MB)")
///         }
///
///         // MIME type kontrolü
///         allowedTypes := []string{"application/pdf", "image/jpeg", "image/png"}
///         contentType := file.Header.Get("Content-Type")
///         if !contains(allowedTypes, contentType) {
///             return "", fmt.Errorf("izin verilmeyen dosya tipi: %s", contentType)
///         }
///
///         // Dosyayı geçici konuma kaydet
///         tempPath := filepath.Join(os.TempDir(), file.Filename)
///         if err := c.SaveFile(file, tempPath); err != nil {
///             return "", err
///         }
///         defer os.Remove(tempPath)
///
///         // Virus taraması yap (örnek: ClamAV)
///         if err := scanForVirus(tempPath); err != nil {
///             return "", fmt.Errorf("güvenlik kontrolü başarısız: %w", err)
///         }
///
///         // Güvenli, kalıcı konuma taşı
///         finalPath := filepath.Join("uploads", "safe", file.Filename)
///         if err := os.Rename(tempPath, finalPath); err != nil {
///             return "", err
///         }
///
///         return "/uploads/safe/" + file.Filename, nil
///     })
/// ```
///
/// ### Örnek 5: Çoklu Depolama (Backup)
///
/// ```go
/// field := fields.File("important_doc").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         filename := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)
///
///         // 1. Yerel diske kaydet
///         localPath := filepath.Join("uploads", filename)
///         if err := c.SaveFile(file, localPath); err != nil {
///             return "", fmt.Errorf("yerel kayıt hatası: %w", err)
///         }
///
///         // 2. S3'e backup yükle (async)
///         go func() {
///             if err := uploadToS3(localPath, filename); err != nil {
///                 log.Printf("S3 backup hatası: %v", err)
///             }
///         }()
///
///         // 3. Yerel yolu döndür (primary)
///         return "/uploads/" + filename, nil
///     })
/// ```
///
/// ## Avantajlar
///
/// - **Esneklik**: İstediğiniz depolama stratejisini uygulayabilirsiniz
/// - **Entegrasyon**: Herhangi bir depolama servisi ile entegre olabilirsiniz
/// - **Özelleştirme**: Dosya işleme, optimizasyon, güvenlik kontrolü ekleyebilirsiniz
/// - **Çoklu Backend**: Farklı alanlar için farklı depolama stratejileri kullanabilirsiniz
/// - **Kontrol**: Dosya adlandırma, klasör yapısı üzerinde tam kontrol
///
/// ## Dezavantajlar
///
/// - **Karmaşıklık**: Hata yönetimi, retry logic, cleanup gibi konuları ele almanız gerekir
/// - **Performans**: Senkron işlemler request süresini uzatabilir (async kullanın)
/// - **Güvenlik**: Dosya validasyonu, sanitization sizin sorumluluğunuzdadır
/// - **Bakım**: Depolama servisi değişikliklerinde kod güncellenmeli
///
/// ## Önemli Notlar
///
/// ⚠️ **Güvenlik Kritik**: Dosya yüklemesi ciddi güvenlik riskleri içerir:
/// - Dosya tipini MIME type'a değil, içeriğe bakarak kontrol edin
/// - Dosya boyutunu sınırlayın
/// - Dosya adlarını sanitize edin (path traversal saldırılarına karşı)
/// - Yüklenen dosyaları executable olmayan bir dizine kaydedin
/// - Mümkünse virus taraması yapın
///
/// ⚠️ **Hata Yönetimi**: Hata durumlarında:
/// - Geçici dosyaları temizleyin
/// - Anlamlı hata mesajları döndürün
/// - Kritik hataları loglayın
/// - Partial upload'ları temizleyin
///
/// ⚠️ **Performans**:
/// - Büyük dosyalar için streaming kullanın
/// - Ağır işlemleri (görsel işleme, virus tarama) async yapın
/// - Timeout değerlerini uygun ayarlayın
/// - Progress tracking için webhook/SSE kullanın
///
/// ⚠️ **Dosya Adı**: Döndürdüğünüz string:
/// - Veritabanına kaydedilir
/// - Frontend'de dosyaya erişim için kullanılır
/// - Mutlak URL veya relative path olabilir
/// - CDN URL'i de olabilir
///
/// 💡 **İpucu**: Depolama mantığını ayrı bir service katmanına çıkarın:
/// ```go
/// type StorageService interface {
///     Upload(file *multipart.FileHeader) (string, error)
/// }
///
/// field := fields.File("avatar").
///     Storage(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
///         return storageService.Upload(file)
///     })
/// ```
///
/// 💡 **İpucu**: Farklı ortamlar için farklı storage kullanın:
/// ```go
/// var storage StorageCallbackFunc
/// if os.Getenv("ENV") == "production" {
///     storage = s3Storage
/// } else {
///     storage = localStorage
/// }
/// ```
///
/// ## İlgili Dokümantasyon
///
/// - `docs/Fields.md` - File ve Image field kullanımı
/// - `fiber.Ctx.SaveFile()` - Fiber dosya kaydetme metodu
/// - `multipart.FileHeader` - Go multipart dosya yapısı
///
/// ## Benzer Kavramlar
///
/// - **Laravel Storage**: Laravel'in Storage facade'i
/// - **Multer (Node.js)**: Express için dosya yükleme middleware'i
/// - **CarrierWave (Rails)**: Rails dosya yükleme gem'i
type StorageCallbackFunc func(c *fiber.Ctx, file *multipart.FileHeader) (string, error)
