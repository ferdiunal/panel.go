# Settings Page - Sistem Ayarları Sayfası

## Genel Bakış

Settings Page, Panel.go projesinde sistem genelinde ayarları yönetmek için kullanılan özel bir sayfa türüdür. Veritabanında key-value çiftleri olarak saklanan ayarları dinamik form alanları aracılığıyla kullanıcı arayüzünde sunar.

### Temel Özellikler

- **Dinamik Form Alanları**: Farklı field türleri ile esnek ayar yönetimi
- **Veritabanı Entegrasyonu**: Ayarlar veritabanında kalıcı olarak saklanır
- **Navigasyon Kontrolü**: Menüde gösterilip gizlenebilir
- **Grup Desteği**: Ayarları kategorize etmek için grup sistemi
- **Otomatik Kaydetme**: Form submit edildiğinde otomatik olarak veritabanına kaydedilir

---

## Mimari

### Backend (Go)

```
pkg/
├── page/
│   └── settings.go              # Settings Page struct ve metodları
├── domain/
│   └── setting/
│       └── entity.go            # Setting domain modeli
└── panel/
    ├── app.go                   # Settings Page kaydı
    └── config.go                # Settings Page yapılandırması
```

**Ana Bileşenler:**
- `Settings`: Ana sayfa struct'ı (Base'i embed eder)
- `Setting`: Domain modeli (key-value çiftleri)
- `Config.SettingsPage`: Opsiyonel yapılandırma

### Frontend (React/TypeScript)

Settings Page, panel'in standart sayfa render sistemi tarafından otomatik olarak işlenir. Özel bir frontend component'i gerekmez.

---

## Backend Kullanımı

### 1. Varsayılan Settings Page

Panel başlatıldığında otomatik olarak varsayılan bir Settings Page oluşturulur:

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
)

func main() {
    config := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server: panel.ServerConfig{Host: "localhost", Port: "8080"},
        Environment: "development",
    }

    p := panel.New(config)
    p.Start()
}
```

**Varsayılan Field'lar:**
- Site Name (text)
- Site URL (text)
- Site Description (textarea)
- Contact Email (email)
- Contact Phone (tel)
- Contact Address (textarea)
- Register Enable (switch)
- Forgot Password Enable (switch)
- Maintenance Mode (switch)
- Debug Mode (switch)

### 2. Özel Settings Page Oluşturma

Kendi Settings Page'inizi oluşturmak için:

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/page"
    "github.com/ferdiunal/panel.go/pkg/panel"
)

func main() {
    // Özel Settings Page oluştur
    customSettings := &page.Settings{
        Elements: []fields.Element{
            // Temel Ayarlar
            fields.Text("Site Name", "site_name").
                Label("Site Adı").
                Placeholder("Site adını girin").
                Required().
                Default("Panel.go"),

            fields.Text("Site URL", "site_url").
                Label("Site URL").
                Placeholder("https://example.com").
                Required(),

            fields.Textarea("Site Description", "site_description").
                Label("Site Açıklaması").
                Placeholder("Site açıklamasını girin").
                Rows(3),

            // İletişim Bilgileri
            fields.Email("Contact Email", "contact_email").
                Label("İletişim E-posta").
                Placeholder("contact@example.com"),

            fields.Tel("Contact Phone", "contact_phone").
                Label("İletişim Telefon").
                Placeholder("+90 555 123 4567"),

            // Özellikler
            fields.Switch("Register Enable", "register_enable").
                Label("Kullanıcı Kaydı").
                HelpText("Yeni kullanıcıların kayıt olmasına izin ver").
                Default(true),

            fields.Switch("Maintenance Mode", "maintenance_mode").
                Label("Bakım Modu").
                HelpText("Siteyi bakım moduna al").
                Default(false),

            // Gelişmiş Ayarlar
            fields.Number("Session Timeout", "session_timeout").
                Label("Oturum Zaman Aşımı (dakika)").
                Min(5).
                Max(1440).
                Default(60),

            fields.Select("Default Language", "default_language").
                Label("Varsayılan Dil").
                Options(map[string]string{
                    "en": "English",
                    "tr": "Türkçe",
                    "de": "Deutsch",
                }).
                Default("en"),

            fields.Select("Date Format", "date_format").
                Label("Tarih Formatı").
                Options(map[string]string{
                    "DD/MM/YYYY": "Gün/Ay/Yıl",
                    "MM/DD/YYYY": "Ay/Gün/Yıl",
                    "YYYY-MM-DD": "Yıl-Ay-Gün",
                }).
                Default("DD/MM/YYYY"),
        },
        HideInNavigation: false, // Menüde göster
    }

    config := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server: panel.ServerConfig{Host: "localhost", Port: "8080"},
        SettingsPage: customSettings, // Özel Settings Page kullan
    }

    p := panel.New(config)
    p.Start()
}
```

### 3. Gruplandırılmış Ayarlar

Ayarları gruplara ayırmak için field'ların sırasını ve label'larını kullanabilirsiniz:

```go
customSettings := &page.Settings{
    Elements: []fields.Element{
        // === GENEL AYARLAR ===
        fields.Text("Site Name", "site_name").
            Label("Site Adı").
            Required(),

        fields.Text("Site URL", "site_url").
            Label("Site URL").
            Required(),

        // === İLETİŞİM BİLGİLERİ ===
        fields.Email("Contact Email", "contact_email").
            Label("İletişim E-posta"),

        fields.Tel("Contact Phone", "contact_phone").
            Label("İletişim Telefon"),

        // === ÖZELLİKLER ===
        fields.Switch("Register Enable", "register_enable").
            Label("Kullanıcı Kaydı").
            Default(true),

        fields.Switch("Maintenance Mode", "maintenance_mode").
            Label("Bakım Modu").
            Default(false),
    },
}
```

### 4. Conditional Visibility (Koşullu Görünürlük)

Belirli kullanıcılara veya koşullara göre field'ları göstermek için:

```go
customSettings := &page.Settings{
    Elements: []fields.Element{
        fields.Text("Site Name", "site_name").
            Label("Site Adı").
            Required(),

        // Sadece admin kullanıcılar görebilir
        fields.Switch("Debug Mode", "debug_mode").
            Label("Debug Modu").
            CanSee(func(ctx *core.ResourceContext) bool {
                user := ctx.User
                return user.IsAdmin()
            }),

        // Sadece development ortamında göster
        fields.Switch("SQL Logging", "sql_logging").
            Label("SQL Loglama").
            CanSee(func(ctx *core.ResourceContext) bool {
                return os.Getenv("ENVIRONMENT") == "development"
            }),
    },
}
```

### 5. Navigasyonda Gizleme

Settings Page'i navigasyon menüsünde gizlemek için:

```go
customSettings := &page.Settings{
    Elements: []fields.Element{
        // Field'lar...
    },
    HideInNavigation: true, // Menüde gizle
}
```

---

## Veritabanı Yapısı

Settings verileri `settings` tablosunda saklanır:

```sql
CREATE TABLE settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    key VARCHAR(255) UNIQUE NOT NULL,
    value TEXT,
    type VARCHAR(50),
    group VARCHAR(100),
    label VARCHAR(255),
    help TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    INDEX idx_key (key),
    INDEX idx_group (group)
);
```

**Alan Açıklamaları:**
- `key`: Ayarın benzersiz anahtarı (örn: "site_name")
- `value`: Ayarın değeri (string formatında)
- `type`: Değer tipi ("string", "integer", "boolean", "json")
- `group`: Ayar grubu (örn: "general", "email", "security")
- `label`: Kullanıcı dostu etiket
- `help`: Yardım metni

---

## Ayarları Okuma ve Yazma

### Backend'de Ayarları Okuma

```go
package main

import (
    "context"
    "github.com/ferdiunal/panel.go/pkg/domain/setting"
    "gorm.io/gorm"
)

func getSettings(db *gorm.DB) {
    ctx := context.Background()

    // Tek bir ayarı oku
    var siteName setting.Setting
    db.Where("key = ?", "site_name").First(&siteName)
    println("Site Name:", siteName.Value)

    // Grup bazında ayarları oku
    var generalSettings []setting.Setting
    db.Where("group = ?", "general").Find(&generalSettings)

    for _, s := range generalSettings {
        println(s.Key, ":", s.Value)
    }
}
```

### Backend'de Ayarları Yazma

```go
package main

import (
    "context"
    "github.com/ferdiunal/panel.go/pkg/domain/setting"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

func saveSettings(db *gorm.DB) {
    ctx := context.Background()

    // Tek bir ayarı kaydet (upsert)
    s := setting.Setting{
        Key:   "site_name",
        Value: "My Panel",
        Type:  "string",
        Group: "general",
        Label: "Site Name",
    }

    db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "key"}},
        DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
    }).Create(&s)

    // Birden fazla ayarı kaydet
    settings := []setting.Setting{
        {Key: "site_name", Value: "My Panel", Type: "string", Group: "general"},
        {Key: "site_url", Value: "https://example.com", Type: "string", Group: "general"},
        {Key: "register_enable", Value: "true", Type: "boolean", Group: "features"},
    }

    for _, setting := range settings {
        db.Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "key"}},
            DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
        }).Create(&setting)
    }
}
```

---

## Field Türleri ve Kullanımları

### Text Field

```go
fields.Text("Site Name", "site_name").
    Label("Site Adı").
    Placeholder("Site adını girin").
    Required().
    Default("Panel.go").
    HelpText("Sitenizin görüntülenecek adı")
```

### Textarea Field

```go
fields.Textarea("Site Description", "site_description").
    Label("Site Açıklaması").
    Placeholder("Site açıklamasını girin").
    Rows(5).
    MaxLength(500)
```

### Email Field

```go
fields.Email("Contact Email", "contact_email").
    Label("İletişim E-posta").
    Placeholder("contact@example.com").
    Required()
```

### Tel Field

```go
fields.Tel("Contact Phone", "contact_phone").
    Label("İletişim Telefon").
    Placeholder("+90 555 123 4567")
```

### Switch Field

```go
fields.Switch("Maintenance Mode", "maintenance_mode").
    Label("Bakım Modu").
    HelpText("Siteyi bakım moduna al").
    Default(false)
```

### Number Field

```go
fields.Number("Session Timeout", "session_timeout").
    Label("Oturum Zaman Aşımı (dakika)").
    Min(5).
    Max(1440).
    Default(60).
    Step(5)
```

### Select Field

```go
fields.Select("Default Language", "default_language").
    Label("Varsayılan Dil").
    Placeholder("Dil seçin").
    Options(map[string]string{
        "en": "English",
        "tr": "Türkçe",
        "de": "Deutsch",
    }).
    Default("en")
```

### Image Field

```go
fields.Image("Site Logo", "site_logo").
    Label("Site Logosu").
    HelpText("Site logosunu yükleyin (PNG, JPG, SVG)").
    StoreAs(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
        // Dosya yükleme mantığı
        storageUrl := "/storage/"
        storagePath := "./storage/public"
        ext := filepath.Ext(file.Filename)
        filename := fmt.Sprintf("logo_%d%s", time.Now().UnixNano(), ext)
        localPath := filepath.Join(storagePath, filename)

        _ = os.MkdirAll(storagePath, 0755)

        if err := c.SaveFile(file, localPath); err != nil {
            return "", err
        }
        return fmt.Sprintf("%s/%s", storageUrl, filename), nil
    })
```

---

## API Endpoints

Settings Page otomatik olarak aşağıdaki API endpoint'lerini oluşturur:

### GET /api/pages/settings

Settings sayfasının yapılandırmasını ve field'larını döndürür.

**Response:**
```json
{
  "slug": "settings",
  "title": "Settings",
  "description": "Sistem ayarlarını yönetin",
  "group": "System",
  "fields": [
    {
      "key": "site_name",
      "name": "Site Name",
      "view": "text-field",
      "label": "Site Adı",
      "placeholder": "Site adını girin",
      "required": true,
      "default": "Panel.go"
    }
  ]
}
```

### GET /api/pages/settings/data

Mevcut ayar değerlerini döndürür.

**Response:**
```json
{
  "site_name": "My Panel",
  "site_url": "https://example.com",
  "register_enable": true,
  "maintenance_mode": false
}
```

### POST /api/pages/settings

Ayarları kaydeder.

**Request:**
```json
{
  "site_name": "My Panel",
  "site_url": "https://example.com",
  "register_enable": true,
  "maintenance_mode": false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Ayarlar başarıyla kaydedildi"
}
```

---

## Best Practices

### 1. Key Naming Convention

Ayar anahtarları için tutarlı bir isimlendirme kullanın:

```go
// ✅ İyi: snake_case ve açıklayıcı
"site_name"
"contact_email"
"register_enable"
"session_timeout"

// ❌ Kötü: tutarsız ve belirsiz
"siteName"
"email"
"reg"
"timeout"
```

### 2. Default Değerler

Her field için mantıklı default değerler belirleyin:

```go
// ✅ İyi: Güvenli default değerler
fields.Switch("Register Enable", "register_enable").
    Default(false) // Varsayılan olarak kapalı (güvenli)

fields.Number("Session Timeout", "session_timeout").
    Default(60) // Makul bir değer

// ❌ Kötü: Güvensiz veya belirsiz default değerler
fields.Switch("Debug Mode", "debug_mode").
    Default(true) // Production'da tehlikeli

fields.Number("Max Upload Size", "max_upload_size").
    Default(0) // Belirsiz
```

### 3. Validation

Kritik ayarlar için validation ekleyin:

```go
fields.Number("Session Timeout", "session_timeout").
    Label("Oturum Zaman Aşımı (dakika)").
    Min(5).      // Minimum 5 dakika
    Max(1440).   // Maximum 24 saat
    Required()   // Zorunlu alan
```

### 4. Help Text

Karmaşık ayarlar için yardım metni ekleyin:

```go
fields.Switch("Maintenance Mode", "maintenance_mode").
    Label("Bakım Modu").
    HelpText("Aktif edildiğinde site sadece admin kullanıcılar tarafından erişilebilir olur")
```

### 5. Gruplandırma

İlgili ayarları gruplandırın:

```go
// Genel Ayarlar
fields.Text("Site Name", "site_name"),
fields.Text("Site URL", "site_url"),

// İletişim Bilgileri
fields.Email("Contact Email", "contact_email"),
fields.Tel("Contact Phone", "contact_phone"),

// Özellikler
fields.Switch("Register Enable", "register_enable"),
fields.Switch("Maintenance Mode", "maintenance_mode"),
```

---

## Örnekler

### Örnek 1: E-posta Ayarları

```go
emailSettings := &page.Settings{
    Elements: []fields.Element{
        fields.Text("SMTP Host", "smtp_host").
            Label("SMTP Sunucusu").
            Placeholder("smtp.gmail.com").
            Required(),

        fields.Number("SMTP Port", "smtp_port").
            Label("SMTP Port").
            Default(587).
            Required(),

        fields.Text("SMTP Username", "smtp_username").
            Label("SMTP Kullanıcı Adı").
            Required(),

        fields.Password("SMTP Password", "smtp_password").
            Label("SMTP Şifre").
            Required(),

        fields.Email("From Email", "from_email").
            Label("Gönderen E-posta").
            Required(),

        fields.Text("From Name", "from_name").
            Label("Gönderen Adı").
            Default("Panel.go"),

        fields.Switch("Use TLS", "smtp_use_tls").
            Label("TLS Kullan").
            Default(true),
    },
}
```

### Örnek 2: Güvenlik Ayarları

```go
securitySettings := &page.Settings{
    Elements: []fields.Element{
        fields.Number("Password Min Length", "password_min_length").
            Label("Minimum Şifre Uzunluğu").
            Min(6).
            Max(32).
            Default(8),

        fields.Switch("Require Uppercase", "password_require_uppercase").
            Label("Büyük Harf Zorunlu").
            Default(true),

        fields.Switch("Require Lowercase", "password_require_lowercase").
            Label("Küçük Harf Zorunlu").
            Default(true),

        fields.Switch("Require Number", "password_require_number").
            Label("Rakam Zorunlu").
            Default(true),

        fields.Switch("Require Special Char", "password_require_special").
            Label("Özel Karakter Zorunlu").
            Default(false),

        fields.Number("Max Login Attempts", "max_login_attempts").
            Label("Maksimum Giriş Denemesi").
            Min(3).
            Max(10).
            Default(5),

        fields.Number("Lockout Duration", "lockout_duration").
            Label("Kilitleme Süresi (dakika)").
            Min(5).
            Max(60).
            Default(15),
    },
}
```

### Örnek 3: Görünüm Ayarları

```go
appearanceSettings := &page.Settings{
    Elements: []fields.Element{
        fields.Select("Default Theme", "default_theme").
            Label("Varsayılan Tema").
            Options(map[string]string{
                "light": "Açık Tema",
                "dark":  "Koyu Tema",
                "auto":  "Otomatik",
            }).
            Default("auto"),

        fields.Select("Primary Color", "primary_color").
            Label("Ana Renk").
            Options(map[string]string{
                "blue":   "Mavi",
                "green":  "Yeşil",
                "purple": "Mor",
                "red":    "Kırmızı",
            }).
            Default("blue"),

        fields.Number("Items Per Page", "items_per_page").
            Label("Sayfa Başına Öğe").
            Min(10).
            Max(100).
            Default(25).
            Step(5),

        fields.Switch("Show Sidebar", "show_sidebar").
            Label("Kenar Çubuğunu Göster").
            Default(true),

        fields.Switch("Compact Mode", "compact_mode").
            Label("Kompakt Mod").
            Default(false),
    },
}
```

---

## Troubleshooting

### Ayarlar Kaydedilmiyor

**Sorun**: Form submit edildiğinde ayarlar veritabanına kaydedilmiyor.

**Çözüm**:
1. Veritabanı bağlantısını kontrol edin
2. `settings` tablosunun oluşturulduğundan emin olun
3. Field key'lerinin doğru olduğundan emin olun
4. Backend log'larını kontrol edin

### Field Görünmüyor

**Sorun**: Eklediğim field Settings sayfasında görünmüyor.

**Çözüm**:
1. Field'ın `CanSee` metodunu kontrol edin
2. Field'ın `OnlyOnForm`, `OnlyOnList` gibi visibility metodlarını kontrol edin
3. Frontend build'i yeniden yapın: `cd web && npm run build`

### Validation Hataları

**Sorun**: Form submit edildiğinde validation hataları alıyorum.

**Çözüm**:
1. Required field'ların doldurulduğundan emin olun
2. Min/Max değerlerinin doğru olduğundan emin olun
3. Email, Tel gibi field'ların formatının doğru olduğundan emin olun

---

## Lisans

Bu özellik Panel.go projesinin bir parçasıdır ve aynı lisans altında dağıtılır.
