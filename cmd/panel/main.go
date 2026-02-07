// Bu paket, Panel CLI uygulamasının ana giriş noktasıdır.
//
// Panel, Go tabanlı bir kod oluşturma aracıdır (code generator) ve aşağıdaki
// komutları destekler:
//   - make:resource: Yeni bir resource (kaynak) oluşturur
//   - make:page: Yeni bir sayfa oluşturur
//   - make:model: Yeni bir model (veri modeli) oluşturur
//
// Tüm komutlar, gömülü stub dosyalarından şablonlar kullanarak dosyalar oluşturur.
package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Bu değişken, stubs klasöründeki tüm .stub dosyalarını gömülü dosya sistemi
// olarak içerir. Go'nun embed özelliği sayesinde, bu dosyalar derleme zamanında
// ikili dosyaya dahil edilir ve çalışma zamanında erişilebilir hale gelir.
//
// # Kullanım Senaryosu
//
// Stub dosyaları, yeni kaynaklar, sayfalar ve modeller oluştururken şablon
// olarak kullanılır. Bu sayede, tutarlı ve standartlaştırılmış kod yapısı
// sağlanır.
//
// # Önemli Notlar
//
//   - Stub dosyaları, Go template sözdizimini kullanır
//   - Dosyalar, stubs/ klasöründe *.stub uzantısıyla saklanır
//   - Gömülü dosyalar, derleme zamanında sabitlenir ve değiştirilemez
//
//go:embed stubs/*.stub
var stubsFS embed.FS

// Bu fonksiyon, Panel CLI uygulamasının ana giriş noktasıdır.
//
// Komut satırı argümanlarını işler ve ilgili komut fonksiyonlarını çağırır.
// Uygulama, en az 3 argüman gerektirir: program adı, komut ve kaynak adı.
//
// # Desteklenen Komutlar
//
//   - make:resource <name>: Yeni bir resource oluşturur
//   - make:page <name>: Yeni bir sayfa oluşturur
//   - make:model <name>: Yeni bir model oluşturur
//
// # Kullanım Örnekleri
//
//     panel make:resource blog
//     panel make:page dashboard
//     panel make:model post
//
// # Parametreler
//
//   - os.Args[0]: Program adı (panel)
//   - os.Args[1]: Komut (make:resource, make:page, make:model)
//   - os.Args[2]: Kaynak/Sayfa/Model adı
//
// # Hata Yönetimi
//
//   - Eğer 3'ten az argüman sağlanırsa, kullanım bilgisi gösterilir
//   - Bilinmeyen komut girilirse, "Unknown command" mesajı gösterilir
//
// # Önemli Notlar
//
//   - Komut adları case-sensitive'dir (make:resource, make:page, make:model)
//   - Kaynak adları küçük harfle yazılmalıdır (örn: blog, post, user)
//   - Oluşturulan dosyalar, proje yapısına göre otomatik olarak yerleştirilir
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: panel make:resource <name>")
		return
	}

	command := os.Args[1]
	name := os.Args[2]

	if command == "make:resource" {
		makeResource(name)
	} else if command == "make:page" {
		makePage(name)
	} else if command == "make:model" {
		makeModel(name)
	} else {
		fmt.Println("Unknown command")
	}
}

// Bu fonksiyon, yeni bir resource (kaynak) oluşturur.
//
// Resource, Panel uygulamasında veri yönetimi için temel yapı taşıdır. Bu fonksiyon,
// verilen addan hareketle resource için gerekli tüm dosyaları otomatik olarak oluşturur.
// Oluşturulan dosyalar: resource tanımı, policy (yetkilendirme) ve repository (veri erişimi).
//
// # Oluşturulan Dosyalar
//
//   - <name>_resource.go: Resource tanımı ve konfigürasyonu
//   - <name>_policy.go: Yetkilendirme kuralları
//   - <name>_repository.go: Veri erişim katmanı
//
// # Parametreler
//
//   - name: Resource adı (küçük harfle, tekil form)
//     Örnek: "blog", "post", "user"
//
// # Kullanım Örnekleri
//
//     makeResource("blog")
//     // Oluşturur:
//     // - internal/resource/blog/blog_resource.go
//     // - internal/resource/blog/blog_policy.go
//     // - internal/resource/blog/blog_repository.go
//
//     makeResource("post")
//     // Oluşturur:
//     // - internal/resource/post/post_resource.go
//     // - internal/resource/post/post_policy.go
//     // - internal/resource/post/post_repository.go
//
// # İsim Dönüşümleri
//
//   - Giriş: "blog"
//   - ResourceName: "Blog" (başlık harfi büyük)
//   - PackageName: "blog" (küçük harf)
//   - Identifier: "blogs" (çoğul form)
//   - Label: "Blogs" (başlık harfi büyük, çoğul)
//   - ModelName: "Blog" (model adı)
//
// # Dizin Yapısı
//
//   internal/resource/<packageName>/
//   ├── <packageName>_resource.go
//   ├── <packageName>_policy.go
//   └── <packageName>_repository.go
//
// # Önemli Notlar
//
//   - Dizin otomatik olarak oluşturulur (0755 izinleriyle)
//   - Stub dosyaları, gömülü dosya sisteminden okunur
//   - Şablonlar, sağlanan veri haritası kullanılarak işlenir
//   - Varsayılan grup: "Content", varsayılan ikon: "circle"
//   - Model adının önceden var olduğu veya oluşturulacağı varsayılır
//
// # Hata Yönetimi
//
//   - Dizin oluşturma başarısız olursa, hata mesajı gösterilir
//   - Stub dosyası okunamıyorsa, hata mesajı gösterilir
//   - Şablon işleme başarısız olursa, hata mesajı gösterilir
func makeResource(name string) {
	// İsim normalizasyonu
	// Örn: "blog" -> "Blog"
	caser := cases.Title(language.English)
	resourceName := caser.String(name)        // Blog
	packageName := strings.ToLower(name)      // blog
	identifier := strings.ToLower(name) + "s" // blogs
	label := resourceName + "s"               // Blogs
	modelName := resourceName                 // Blog (Model'in var olduğu veya oluşturulacağı varsayılır)

	// Dizin: internal/resource/<name>
	// SDK tüketicileri, kaynakları pkg/resource veya internal/resource içine koyabilir.
	// Şimdilik internal/resource varsayıyoruz veya yapılandırılabilir hale getirebiliriz.
	// internal/resource'a varsayılan olarak ayarlamak, önceki davranışla eşleşir.
	dir := filepath.Join("internal", "resource", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Şablonlar için veri
	data := map[string]string{
		"PackageName":  packageName,
		"ResourceName": resourceName,
		"ModelName":    modelName,
		"Slug":         identifier,
		"Title":        label,
		"Label":        label,
		"Identifier":   identifier,
		"Group":        "Content", // Varsayılan grup
		"Icon":         "circle",  // Varsayılan ikon
	}

	// İşlenecek stub'lar
	stubs := map[string]string{
		"resource.stub":   filepath.Join(dir, fmt.Sprintf("%s_resource.go", packageName)),
		"policy.stub":     filepath.Join(dir, fmt.Sprintf("%s_policy.go", packageName)),
		"repository.stub": filepath.Join(dir, fmt.Sprintf("%s_repository.go", packageName)),
	}

	for stub, target := range stubs {
		createFileFromStub(stub, target, data)
	}

	fmt.Printf("Resource %s generated successfully in %s\n", resourceName, dir)
}

// Bu fonksiyon, yeni bir sayfa (page) oluşturur.
//
// Sayfa, Panel uygulamasında kullanıcı arayüzü bileşenlerini temsil eder. Bu fonksiyon,
// verilen addan hareketle sayfa için gerekli dosyaları otomatik olarak oluşturur.
// Oluşturulan dosya, sayfa tanımı ve konfigürasyonunu içerir.
//
// # Oluşturulan Dosyalar
//
//   - <name>.go: Sayfa tanımı ve konfigürasyonu
//
// # Parametreler
//
//   - name: Sayfa adı (küçük harfle)
//     Örnek: "dashboard", "settings", "profile"
//
// # Kullanım Örnekleri
//
//     makePage("dashboard")
//     // Oluşturur:
//     // - internal/page/dashboard.go
//
//     makePage("settings")
//     // Oluşturur:
//     // - internal/page/settings.go
//
// # İsim Dönüşümleri
//
//   - Giriş: "dashboard"
//   - PageName: "Dashboard" (başlık harfi büyük)
//   - PackageName: "page" (sabit, tüm sayfalar aynı pakette)
//   - Slug: "dashboard" (küçük harf)
//   - Title: "Dashboard" (başlık harfi büyük)
//
// # Dizin Yapısı
//
//   internal/page/
//   └── <packageName>.go
//
// # Önemli Notlar
//
//   - Tüm sayfalar, "internal/page" dizininde saklanır
//   - Paket adı her zaman "page" olarak ayarlanır
//   - Varsayılan grup: "System", varsayılan ikon: "circle"
//   - Dizin otomatik olarak oluşturulur (0755 izinleriyle)
//   - Stub dosyası, gömülü dosya sisteminden okunur
//
// # Hata Yönetimi
//
//   - Dizin oluşturma başarısız olursa, hata mesajı gösterilir
//   - Stub dosyası okunamıyorsa, hata mesajı gösterilir
//   - Şablon işleme başarısız olursa, hata mesajı gösterilir
func makePage(name string) {
	// İsim normalizasyonu
	// Örn: "dashboard" -> "Dashboard"
	caser := cases.Title(language.English)
	pageName := caser.String(name)       // Dashboard
	packageName := strings.ToLower(name) // dashboard
	slug := strings.ToLower(name)        // dashboard
	title := pageName                    // Dashboard

	dir := filepath.Join("internal", "page")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, fmt.Sprintf("%s.go", packageName))

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": "page", // Mevcut yapıyla eşleşmesi için 'page' paketini kullanıyoruz
		"PageName":    pageName,
		"Slug":        slug,
		"Title":       title,
		"Group":       "System",
		"Icon":        "circle",
	}

	createFileFromStub("page.stub", targetPath, data)
	fmt.Printf("Page %s generated successfully at %s\n", pageName, targetPath)
}

// Bu fonksiyon, yeni bir model (veri modeli) oluşturur.
//
// Model, Panel uygulamasında veri yapısını temsil eder. Bu fonksiyon,
// verilen addan hareketle model için gerekli dosyaları otomatik olarak oluşturur.
// Oluşturulan dosya, model tanımı ve veri yapısını içerir.
//
// # Oluşturulan Dosyalar
//
//   - entity.go: Model tanımı ve veri yapısı
//
// # Parametreler
//
//   - name: Model adı (küçük harfle, tekil form)
//     Örnek: "blog", "post", "user"
//
// # Kullanım Örnekleri
//
//     makeModel("blog")
//     // Oluşturur:
//     // - internal/domain/blog/entity.go
//
//     makeModel("post")
//     // Oluşturur:
//     // - internal/domain/post/entity.go
//
// # İsim Dönüşümleri
//
//   - Giriş: "blog"
//   - ModelName: "Blog" (başlık harfi büyük)
//   - PackageName: "blog" (küçük harf)
//
// # Dizin Yapısı
//
//   internal/domain/<packageName>/
//   └── entity.go
//
// # Önemli Notlar
//
//   - Her model, "internal/domain/<packageName>" dizininde saklanır
//   - Dosya adı her zaman "entity.go" olarak ayarlanır
//   - Dizin otomatik olarak oluşturulur (0755 izinleriyle)
//   - Stub dosyası, gömülü dosya sisteminden okunur
//   - Model, resource ve page tarafından referans alınabilir
//
// # Hata Yönetimi
//
//   - Dizin oluşturma başarısız olursa, hata mesajı gösterilir
//   - Stub dosyası okunamıyorsa, hata mesajı gösterilir
//   - Şablon işleme başarısız olursa, hata mesajı gösterilir
func makeModel(name string) {
	// İsim normalizasyonu
	// Örn: "blog" -> "Blog"
	caser := cases.Title(language.English)
	modelName := caser.String(name)      // Blog
	packageName := strings.ToLower(name) // blog

	// Dizin: internal/domain/<name>
	dir := filepath.Join("internal", "domain", packageName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	targetPath := filepath.Join(dir, "entity.go")

	// Şablonlar için veri
	data := map[string]string{
		"PackageName": packageName,
		"ModelName":   modelName,
	}

	createFileFromStub("model.stub", targetPath, data)
	fmt.Printf("Model %s generated successfully at %s\n", modelName, targetPath)
}

// Bu fonksiyon, stub dosyasından şablon işleyerek yeni bir dosya oluşturur.
//
// Bu fonksiyon, Panel'in kod oluşturma sisteminin temel bileşenidir. Gömülü dosya
// sisteminden stub dosyasını okur, Go template sözdizimini kullanarak işler ve
// sonucu hedef dosyaya yazar. Tüm resource, page ve model oluşturma işlemleri
// bu fonksiyonu kullanır.
//
// # Parametreler
//
//   - stubName: Stub dosyasının adı (örn: "resource.stub", "page.stub", "model.stub")
//     Eğer "stubs/" ön eki yoksa otomatik olarak eklenir.
//   - targetPath: Oluşturulacak dosyanın tam yolu (örn: "internal/resource/blog/blog_resource.go")
//   - data: Şablon değişkenlerini içeren harita
//     Örn: map[string]string{"PackageName": "blog", "ResourceName": "Blog"}
//
// # Kullanım Örnekleri
//
//     data := map[string]string{
//         "PackageName":  "blog",
//         "ResourceName": "Blog",
//         "ModelName":    "Blog",
//     }
//     createFileFromStub("resource.stub", "internal/resource/blog/blog_resource.go", data)
//
//     data := map[string]string{
//         "PackageName": "page",
//         "PageName":    "Dashboard",
//         "Slug":        "dashboard",
//     }
//     createFileFromStub("page.stub", "internal/page/dashboard.go", data)
//
// # İşlem Adımları
//
//   1. Stub dosyasının yolunu kontrol eder ve gerekirse "stubs/" ön ekini ekler
//   2. Gömülü dosya sisteminden stub dosyasını okur
//   3. Dosya içeriğini Go template olarak ayrıştırır
//   4. Hedef dosyayı oluşturur
//   5. Şablonu sağlanan veri haritası ile işler ve dosyaya yazar
//
// # Şablon Sözdizimi
//
// Stub dosyaları, Go template sözdizimini kullanır:
//
//     package {{.PackageName}}
//
//     type {{.ResourceName}} struct {
//         // ...
//     }
//
// # Hata Yönetimi
//
//   - Stub dosyası okunamıyorsa: "Error reading stub" mesajı gösterilir
//   - Şablon ayrıştırılamıyorsa: "Error parsing template" mesajı gösterilir
//   - Hedef dosya oluşturulamıyorsa: "Error creating file" mesajı gösterilir
//   - Şablon işleme başarısız olursa: "Error executing template" mesajı gösterilir
//
// # Önemli Notlar
//
//   - Stub dosyaları, gömülü dosya sisteminden okunur (derleme zamanında sabitlenir)
//   - Hedef dosya zaten varsa, üzerine yazılır
//   - Dosya oluşturulduktan sonra başarı mesajı gösterilir
//   - Şablon değişkenleri, sağlanan veri haritasından alınır
//   - Stub dosyası yolu, "stubs/" ile başlamıyorsa otomatik olarak eklenir
//
// # Döndürür
//
//   - Hata durumunda: Hata mesajı yazdırılır, fonksiyon döner
//   - Başarı durumunda: Dosya oluşturulur ve başarı mesajı yazdırılır
func createFileFromStub(stubName, targetPath string, data map[string]string) {
	// Stub dosyasını gömülü dosya sisteminden oku
	// Stub'lar, main.go'ya göre stubs/ klasöründe yer alır
	// Eğer dosyalar stubs/*.stub içindeyse, "stubs/stubName" olarak okuruz

	// Stub adının zaten ön ek içerip içermediğini kontrol et
	path := stubName
	if !strings.HasPrefix(path, "stubs/") {
		path = filepath.Join("stubs", stubName)
	}

	content, err := stubsFS.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading stub %s: %v\n", path, err)
		return
	}

	// Şablonu işle
	tmpl, err := template.New(stubName).Parse(string(content))
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", stubName, err)
		return
	}

	// Dosya oluştur
	f, err := os.Create(targetPath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", targetPath, err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("Error executing template %s: %v\n", stubName, err)
	}
	fmt.Printf("Created: %s\n", targetPath)
}
