# Sayfalar (Pages) - Legacy Teknik Akış

Sayfalar (Pages), Panel SDK'da kaynaklar (resources) dışında kalan özel görünümleri ve gösterge panellerini (dashboards) tanımlamanızı sağlar. Genellikle widget'ları bir araya getirmek veya özel raporlama ekranları oluşturmak için kullanılır.

## Bu Doküman Ne Zaman Okunmalı?

Önerilen sıra:
1. [Başlarken](Getting-Started)
2. [Kaynaklar (Resource)](Resources)
3. [Widget'lar (Cards)](Widgets)
4. Bu doküman (`Pages`)

## Hızlı Page Akışı

1. Page struct'ını oluştur ve `Slug`, `Title`, `Icon`, `Cards` metotlarını tanımla.
2. `Config.Pages` içinde page'i kaydet.
3. Gerekirse `app.RegisterPage(...)` ile manuel ekleme yap.
4. `GET /api/pages` ve `GET /api/pages/:slug` endpoint'leri ile doğrula.

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

Sayfaları `main.go` içerisinde `Config.Pages` ile yapılandırabilirsiniz:

```go
func main() {
    cfg := panel.Config{
        // Sayfaları Pages dizisi ile kaydedin
        Pages: []page.Page{
            &Dashboard{},
            &page.Settings{
                Elements: []fields.Element{
                     fields.Text("Support Email", "support_email"),
                },
            },
            &page.Account{
                Elements: []fields.Element{
                    fields.Text("Name", "name").Required(),
                    fields.Email("Email", "email").Required(),
                },
            },
        },
    }
    
    app := panel.New(cfg)
    
    // Özel bir sayfayı manuel kaydetme
    app.RegisterPage(&CustomPage{})
    
    app.Start()
}
```

> **Not:** SDK varsayılan sayfalar oluşturmaz. `panel init` komutu çalıştırıldığında
> `internal/pages/` dizininde Dashboard, Settings ve Account dosyaları otomatik oluşturulur.

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

## Sık Hata Kontrolü (Pages)

- Sayfa listede yok: `Config.Pages` veya `RegisterPage(...)` akışını kontrol edin.
- Endpoint 404: `Slug()` değeri ile çağrılan URL'nin birebir eşleştiğini doğrulayın.
- Kartlar boş: `Cards()` metodunun gerçek veri döndürdüğünü ve widget sorgularını kontrol edin.
- Menüde ikon yanlış: `Icon()` içinde geçerli ikon adının verildiğinden emin olun.

## Sonraki Adım

- Sistem ayar sayfası için: [Ayarlar (Settings)](Settings)
- Kullanıcı hesabı ve profil için: [Account Page](ACCOUNT_PAGE)
- Dashboard metrikleri için: [Widget'lar (Cards)](Widgets)
