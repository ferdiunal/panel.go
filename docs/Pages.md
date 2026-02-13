# Sayfalar (Pages)

Sayfalar (Pages), Panel SDK'da kaynaklar (resources) dışında kalan özel görünümleri ve gösterge panellerini (dashboards) tanımlamanızı sağlar. Genellikle widget'ları bir araya getirmek veya özel raporlama ekranları oluşturmak için kullanılır.

## Sayfa Tanımlama

Bir sayfa oluşturmak için `page.Page` arayüzünü implemente eden bir struct tanımlamanız veya `page.Base` yapısını kullanmanız gerekir.

```go
import (
    "github.com/ferdiunal/panel.go/pkg/page"
    "github.com/ferdiunal/panel.go/pkg/widget"
)

type Dashboard struct {
    page.Base // Temel implementasyon
}

// Slug, sayfanın URL'deki kimliğidir (örn: /api/pages/dashboard)
func (d *Dashboard) Slug() string {
    return "dashboard"
}

// Title, menüde ve sayfa başlığında görünür
func (d *Dashboard) Title() string {
    return "Genel Bakış"
}

// Icon (Lucide React ikon ismi)
func (d *Dashboard) Icon() string {
    return "layout-dashboard"
}

// Cards, sayfada gösterilecek bileşenleri tanımlar
func (d *Dashboard) Cards() []widget.Card {
    return []widget.Card{
        widget.NewCountWidget("Toplam Kullanıcı", &user.User{}),
        widget.NewTrendWidget("Kayıt Trendi", &user.User{}, "created_at"),
    }
}
```

## Sayfa Kaydetme (Registration)

Sayfaları `main.go` içerisinde yapılandırabilirsiniz:

```go
func main() {
    cfg := panel.Config{
        // Varsayılan Dashboard sayfasını değiştirme
        DashboardPage: &Dashboard{},
        
        // Settings sayfasına özel alanlar ekleme
        SettingsPage: &page.Settings{
            Elements: []fields.Element{
                 fields.Text("Support Email", "support_email"),
            },
        },
    }
    
    app := panel.New(cfg)
    
    // Özel bir sayfayı manuel kaydetme
    app.RegisterPage(&CustomPage{})
    
    app.Start()
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
    "cards": [
      {
        "component": "value-metric",
        "title": "Toplam Kullanıcı",
        "data": {
            "value": 150
        },
        "width": "1/3"
      },
      // ... diğer widget'lar
    ]
  }
}
```
