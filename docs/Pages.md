# Sayfalar (Pages)

Sayfalar (Pages), Panel SDK'da kaynaklar (resources) dışında kalan özel görünümleri ve gösterge panellerini (dashboards) tanımlamanızı sağlar. Genellikle widget'ları bir araya getirmek veya özel raporlama ekranları oluşturmak için kullanılır.

## Sayfa Tanımlama

Bir sayfa oluşturmak için `page.Page` arayüzünü implemente eden bir struct tanımlamanız gerekir.

```go
import (
    "panel.go/internal/page"
    "panel.go/internal/widget"
)

type Dashboard struct {
    page.Base // Temel implementasyon için gömülü struct
}

// Slug, sayfanın URL'deki kimliğidir (örn: /api/pages/dashboard)
func (d *Dashboard) Slug() string {
    return "dashboard"
}

// Title, menüde ve sayfa başlığında görünür
func (d *Dashboard) Title() string {
    return "Genel Bakış"
}

// Widgets, sayfada gösterilecek bileşenleri tanımlar
func (d *Dashboard) Widgets() []widget.Widget {
    return []widget.Widget{
        widget.NewCountWidget("Toplam Kullanıcı", &user.User{}),
        widget.NewTrendWidget("Kayıt Trendi", &user.User{}, "created_at"),
    }
}
```

## Sayfa Kaydetme (Registration)

Oluşturduğunuz sayfayı uygulamanıza tanıtmak için `RegisterPage` metodunu kullanmalısınız. Bu işlem genellikle `main.go` içerisinde yapılır.

```go
func main() {
    // ... uygulama yapılandırması ...
    
    p := panel.New(config)
    
    // Sayfayı kaydet
    p.RegisterPage(&Dashboard{})
    
    p.Start()
}
```

## API Kullanımı

Kayıtlı sayfalara aşağıdaki endpointler üzerinden erişilebilir:

- `GET /api/pages`: Tüm kayıtlı sayfaların listesini döner (Menü oluşturmak için).
- `GET /api/pages/:slug`: Belirtilen sayfanın detaylarını ve widget verilerini döner.

### Örnek Yanıt (Detay)

```json
{
  "slug": "dashboard",
  "title": "Genel Bakış",
  "meta": {
    "widgets": [
      {
        "component": "value-metric",
        "title": "Toplam Kullanıcı",
        "value": 150,
        "width": "1/3"
      },
      // ... diğer widget'lar
    ]
  }
}
```
