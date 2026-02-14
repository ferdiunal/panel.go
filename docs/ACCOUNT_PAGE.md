# Account Page - Kullanıcı Hesap Ayarları Sayfası

## Genel Bakış

Account Page, Panel.go projesinde kullanıcıların kendi hesap ayarlarını yönetmek için kullanılan özel bir sayfa türüdür. Her kullanıcı sadece kendi profil bilgilerini, şifresini, bildirim tercihlerini ve görünüm ayarlarını düzenleyebilir.

### Temel Özellikler

- **Profil Yönetimi**: Ad, email, profil resmi güncelleme
- **Şifre Değiştirme**: Güvenli şifre değiştirme işlemi
- **Bildirim Tercihleri**: Email ve SMS bildirim ayarları
- **Görünüm Ayarları**: Dil ve tema tercihleri
- **Güvenlik**: Her kullanıcı sadece kendi bilgilerini görebilir ve düzenleyebilir
- **Otomatik Kaydetme**: Form submit edildiğinde otomatik olarak veritabanına kaydedilir

---

## Mimari

### Backend (Go)

```
pkg/
├── page/
│   └── account.go               # Account Page struct ve metodları
├── domain/
│   └── user/
│       └── entity.go            # User domain modeli
└── panel/
    ├── app.go                   # Account Page kaydı
    └── config.go                # Account Page yapılandırması
```

**Ana Bileşenler:**
- `Account`: Ana sayfa struct'ı (Base'i embed eder)
- `User`: Domain modeli (kullanıcı bilgileri)
- `Config.Pages`: Sayfaları `Pages` dizisi ile kaydedin

### Frontend (React/TypeScript)

Account Page, panel'in standart sayfa render sistemi tarafından otomatik olarak işlenir. Özel bir frontend component'i gerekmez.

---

## Backend Kullanımı

### 1. Varsayılan Account Page

Panel başlatıldığında varsayılan sayfa **otomatik oluşturulmaz** (SDK modu). `panel init` komutu ile proje oluşturulduğunda `internal/pages/account.go` dosyası otomatik oluşturulur.

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/panel"
    "github.com/ferdiunal/panel.go/pkg/page"
    "your-module/internal/pages"
)

func main() {
    config := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server: panel.ServerConfig{Host: "localhost", Port: "8080"},
        Environment: "development",
        Pages: []page.Page{
            pages.NewAccount(), // internal/pages/account.go
        },
    }

    p := panel.New(config)
    p.Start()
}
```

**Varsayılan Field'lar:**
- Name (text) - Tam ad
- Email (email) - Email adresi
- Image (image) - Profil resmi
- Current Password (password) - Mevcut şifre
- New Password (password) - Yeni şifre
- Confirm Password (password) - Şifre tekrar
- Email Notifications (switch) - Email bildirimleri
- SMS Notifications (switch) - SMS bildirimleri
- Language (select) - Dil seçimi
- Theme (select) - Tema seçimi

### 2. Özel Account Page Oluşturma

Kendi Account Page'inizi oluşturmak için:

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/page"
    "github.com/ferdiunal/panel.go/pkg/panel"
)

func main() {
    // Özel Account Page oluştur
    customAccount := &page.Account{
        Elements: []fields.Element{
            // === PROFİL BİLGİLERİ ===
            fields.Text("Name", "name").
                Label("Tam Adınız").
                Placeholder("Adınızı ve soyadınızı girin").
                Required().
                HelpText("Profilinizde görüntülenecek adınız"),

            fields.Email("Email", "email").
                Label("Email Adresi").
                Placeholder("email@example.com").
                Required().
                HelpText("Giriş yapmak için kullanacağınız email adresi"),

            fields.Image("Image", "image").
                Label("Profil Resmi").
                HelpText("Profil resminizi yükleyin (PNG, JPG, max 2MB)"),

            fields.Tel("Phone", "phone").
                Label("Telefon Numarası").
                Placeholder("+90 555 123 4567").
                HelpText("İletişim için telefon numaranız"),

            fields.Textarea("Bio", "bio").
                Label("Biyografi").
                Placeholder("Kendiniz hakkında kısa bir açıklama").
                Rows(3).
                MaxLength(500),

            // === ŞİFRE DEĞİŞTİRME ===
            fields.Password("Current Password", "current_password").
                Label("Mevcut Şifre").
                Placeholder("Mevcut şifrenizi girin").
                HelpText("Şifre değiştirmek için mevcut şifrenizi girmelisiniz"),

            fields.Password("New Password", "new_password").
                Label("Yeni Şifre").
                Placeholder("Yeni şifrenizi girin").
                HelpText("Minimum 8 karakter, en az 1 büyük harf, 1 küçük harf ve 1 rakam"),

            fields.Password("Confirm Password", "confirm_password").
                Label("Şifre Tekrar").
                Placeholder("Yeni şifrenizi tekrar girin").
                HelpText("Yeni şifrenizi doğrulamak için tekrar girin"),

            // === BİLDİRİM TERCİHLERİ ===
            fields.Switch("Email Notifications", "email_notifications").
                Label("Email Bildirimleri").
                HelpText("Önemli güncellemeler için email bildirimleri alın").
                Default(true),

            fields.Switch("SMS Notifications", "sms_notifications").
                Label("SMS Bildirimleri").
                HelpText("Acil durumlar için SMS bildirimleri alın").
                Default(false),

            fields.Switch("Push Notifications", "push_notifications").
                Label("Push Bildirimleri").
                HelpText("Tarayıcı push bildirimleri alın").
                Default(true),

            // === GÖRÜNÜM AYARLARI ===
            fields.Select("Language", "language").
                Label("Dil").
                Placeholder("Dil seçin").
                Options(map[string]string{
                    "en": "English",
                    "tr": "Türkçe",
                    "de": "Deutsch",
                    "fr": "Français",
                }).
                Default("en").
                HelpText("Arayüz dili"),

            fields.Select("Theme", "theme").
                Label("Tema").
                Placeholder("Tema seçin").
                Options(map[string]string{
                    "light": "Açık Tema",
                    "dark":  "Koyu Tema",
                    "auto":  "Otomatik (Sistem)",
                }).
                Default("auto").
                HelpText("Arayüz teması"),

            fields.Select("Timezone", "timezone").
                Label("Saat Dilimi").
                Placeholder("Saat dilimi seçin").
                Options(map[string]string{
                    "Europe/Istanbul": "İstanbul (GMT+3)",
                    "Europe/London":   "Londra (GMT+0)",
                    "America/New_York": "New York (GMT-5)",
                }).
                Default("Europe/Istanbul"),
        },
        HideInNavigation: false, // Menüde göster
    }

    config := panel.Config{
        Database: panel.DatabaseConfig{Instance: db},
        Server: panel.ServerConfig{Host: "localhost", Port: "8080"},
        Pages: []page.Page{
            customAccount, // Özel Account Page kullan
        },
    }

    p := panel.New(config)
    p.Start()
}
```

### 3. Şifre Değiştirme Mantığı

Account Page'de şifre değiştirme işlemi özel olarak ele alınır:

```go
// pkg/page/account.go içinde Save metodu

func (p *Account) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
    // Oturumdaki kullanıcı ID'sini al
    userID := c.Ctx.Locals("user_id")
    if userID == nil {
        return fmt.Errorf("kullanıcı oturumu bulunamadı")
    }

    // Kullanıcıyı veritabanından al
    var u user.User
    if err := db.First(&u, userID).Error; err != nil {
        return fmt.Errorf("kullanıcı bulunamadı: %w", err)
    }

    // Şifre değiştirme kontrolü
    if newPassword, ok := data["new_password"].(string); ok && newPassword != "" {
        // Mevcut şifre kontrolü
        if currentPassword, ok := data["current_password"].(string); ok && currentPassword != "" {
            // Mevcut şifreyi doğrula
            if !verifyPassword(u.Password, currentPassword) {
                return fmt.Errorf("mevcut şifre yanlış")
            }

            // Yeni şifreyi hash'le
            hashedPassword, err := hashPassword(newPassword)
            if err != nil {
                return fmt.Errorf("şifre hash'lenemedi: %w", err)
            }

            data["password"] = hashedPassword
            delete(data, "new_password")
            delete(data, "current_password")
            delete(data, "confirm_password")
        } else {
            return fmt.Errorf("şifre değiştirmek için mevcut şifrenizi girmelisiniz")
        }
    }

    // Kullanıcı bilgilerini güncelle
    if err := db.Model(&u).Updates(data).Error; err != nil {
        return fmt.Errorf("hesap ayarları güncellenemedi: %w", err)
    }

    return nil
}
```

### 4. Profil Resmi Yükleme

Profil resmi yükleme işlemi için özel bir StoreAs fonksiyonu kullanılır:

```go
fields.Image("Image", "image").
    Label("Profil Resmi").
    StoreAs(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
        // Dosya boyutu kontrolü (max 2MB)
        if file.Size > 2*1024*1024 {
            return "", fmt.Errorf("dosya boyutu 2MB'dan büyük olamaz")
        }

        // Dosya tipi kontrolü
        allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
        contentType := file.Header.Get("Content-Type")
        if !contains(allowedTypes, contentType) {
            return "", fmt.Errorf("sadece JPG, PNG ve GIF dosyaları yüklenebilir")
        }

        // Dosya adı oluştur
        storageUrl := "/storage/avatars/"
        storagePath := "./storage/public/avatars"
        ext := filepath.Ext(file.Filename)
        filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
        localPath := filepath.Join(storagePath, filename)

        // Dizin oluştur
        _ = os.MkdirAll(storagePath, 0755)

        // Dosyayı kaydet
        if err := c.SaveFile(file, localPath); err != nil {
            return "", fmt.Errorf("dosya kaydedilemedi: %w", err)
        }

        return fmt.Sprintf("%s%s", storageUrl, filename), nil
    })
```

### 5. Navigasyonda Gizleme

Account Page'i navigasyon menüsünde gizlemek için:

```go
customAccount := &page.Account{
    Elements: []fields.Element{
        // Field'lar...
    },
    HideInNavigation: true, // Menüde gizle
}
```

---

## API Endpoints

Account Page otomatik olarak aşağıdaki API endpoint'lerini oluşturur:

### GET /api/pages/account

Account sayfasının yapılandırmasını ve field'larını döndürür.

**Response:**
```json
{
  "slug": "account",
  "title": "Account",
  "description": "Hesap ayarlarınızı yönetin",
  "group": "User",
  "fields": [
    {
      "key": "name",
      "name": "Name",
      "view": "text-field",
      "label": "Tam Adınız",
      "placeholder": "Adınızı ve soyadınızı girin",
      "required": true
    }
  ]
}
```

### GET /api/pages/account/data

Oturumdaki kullanıcının mevcut bilgilerini döndürür.

**Response:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "image": "/storage/avatars/1234567890.jpg",
  "phone": "+90 555 123 4567",
  "email_notifications": true,
  "sms_notifications": false,
  "language": "en",
  "theme": "auto"
}
```

### POST /api/pages/account

Kullanıcı bilgilerini günceller.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+90 555 123 4567",
  "current_password": "old_password",
  "new_password": "new_password",
  "confirm_password": "new_password",
  "email_notifications": true,
  "sms_notifications": false,
  "language": "en",
  "theme": "dark"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Hesap ayarları başarıyla güncellendi"
}
```

---

## Güvenlik

### 1. Kullanıcı Yetkilendirmesi

Account Page, sadece oturumdaki kullanıcının kendi bilgilerini gösterir ve düzenler:

```go
func (p *Account) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
    // Oturumdaki kullanıcı ID'sini al
    userID := c.Ctx.Locals("user_id")
    if userID == nil {
        return fmt.Errorf("kullanıcı oturumu bulunamadı")
    }

    // Sadece oturumdaki kullanıcının bilgilerini güncelle
    var u user.User
    if err := db.First(&u, userID).Error; err != nil {
        return fmt.Errorf("kullanıcı bulunamadı: %w", err)
    }

    // Güncelleme işlemi...
}
```

### 2. Şifre Doğrulama

Şifre değiştirme işleminde mevcut şifre doğrulaması yapılır:

```go
// Mevcut şifre kontrolü
if currentPassword, ok := data["current_password"].(string); ok && currentPassword != "" {
    // Mevcut şifreyi doğrula
    if !verifyPassword(u.Password, currentPassword) {
        return fmt.Errorf("mevcut şifre yanlış")
    }
}
```

### 3. Şifre Hash'leme

Yeni şifreler güvenli bir şekilde hash'lenir:

```go
// Bcrypt ile şifre hash'leme
func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// Şifre doğrulama
func verifyPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

### 4. Dosya Yükleme Güvenliği

Profil resmi yükleme işleminde güvenlik kontrolleri yapılır:

```go
// Dosya boyutu kontrolü
if file.Size > 2*1024*1024 {
    return "", fmt.Errorf("dosya boyutu 2MB'dan büyük olamaz")
}

// Dosya tipi kontrolü
allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
contentType := file.Header.Get("Content-Type")
if !contains(allowedTypes, contentType) {
    return "", fmt.Errorf("sadece JPG, PNG ve GIF dosyaları yüklenebilir")
}
```

---

## Best Practices

### 1. Şifre Politikası

Güçlü şifre politikası uygulayın:

```go
fields.Password("New Password", "new_password").
    Label("Yeni Şifre").
    HelpText("Minimum 8 karakter, en az 1 büyük harf, 1 küçük harf ve 1 rakam").
    MinLength(8).
    Pattern("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).+$")
```

### 2. Email Doğrulama

Email değişikliklerinde doğrulama yapın:

```go
// Email değişikliği kontrolü
if newEmail, ok := data["email"].(string); ok && newEmail != u.Email {
    // Email doğrulama kodu gönder
    sendEmailVerification(newEmail)
    
    // Email'i pending olarak işaretle
    data["email_verified"] = false
    data["email_verification_token"] = generateToken()
}
```

### 3. Profil Resmi Optimizasyonu

Yüklenen resimleri optimize edin:

```go
// Resim boyutlandırma ve optimizasyon
func optimizeImage(filePath string) error {
    // Resmi aç
    img, err := imaging.Open(filePath)
    if err != nil {
        return err
    }

    // Boyutlandır (max 500x500)
    img = imaging.Fit(img, 500, 500, imaging.Lanczos)

    // Kaydet (JPEG, kalite 85)
    return imaging.Save(img, filePath, imaging.JPEGQuality(85))
}
```

### 4. Audit Log

Önemli değişiklikleri logla:

```go
// Hesap değişikliklerini logla
func logAccountChange(userID uint, field string, oldValue, newValue interface{}) {
    log := AuditLog{
        UserID:    userID,
        Action:    "account_update",
        Field:     field,
        OldValue:  fmt.Sprintf("%v", oldValue),
        NewValue:  fmt.Sprintf("%v", newValue),
        Timestamp: time.Now(),
    }
    db.Create(&log)
}
```

### 5. Rate Limiting

Şifre değiştirme işlemlerinde rate limiting uygulayın:

```go
// Rate limiting middleware
func rateLimitPasswordChange(c *fiber.Ctx) error {
    userID := c.Locals("user_id")
    
    // Son 1 saatte kaç kez şifre değiştirildi?
    var count int64
    db.Model(&PasswordChangeLog{}).
        Where("user_id = ? AND created_at > ?", userID, time.Now().Add(-1*time.Hour)).
        Count(&count)
    
    if count >= 3 {
        return c.Status(429).JSON(fiber.Map{
            "error": "Çok fazla şifre değiştirme denemesi. Lütfen 1 saat sonra tekrar deneyin.",
        })
    }
    
    return c.Next()
}
```

---

## Örnekler

### Örnek 1: Basit Account Page

```go
simpleAccount := &page.Account{
    Elements: []fields.Element{
        fields.Text("Name", "name").
            Label("Ad Soyad").
            Required(),

        fields.Email("Email", "email").
            Label("Email").
            Required(),

        fields.Password("Current Password", "current_password").
            Label("Mevcut Şifre"),

        fields.Password("New Password", "new_password").
            Label("Yeni Şifre"),
    },
}
```

### Örnek 2: Gelişmiş Account Page

```go
advancedAccount := &page.Account{
    Elements: []fields.Element{
        // Profil Bilgileri
        fields.Text("Name", "name").Label("Ad Soyad").Required(),
        fields.Email("Email", "email").Label("Email").Required(),
        fields.Image("Image", "image").Label("Profil Resmi"),
        fields.Tel("Phone", "phone").Label("Telefon"),
        fields.Textarea("Bio", "bio").Label("Biyografi").Rows(3),

        // Şifre Değiştirme
        fields.Password("Current Password", "current_password").Label("Mevcut Şifre"),
        fields.Password("New Password", "new_password").Label("Yeni Şifre"),
        fields.Password("Confirm Password", "confirm_password").Label("Şifre Tekrar"),

        // Bildirim Tercihleri
        fields.Switch("Email Notifications", "email_notifications").Label("Email Bildirimleri").Default(true),
        fields.Switch("SMS Notifications", "sms_notifications").Label("SMS Bildirimleri").Default(false),

        // Görünüm Ayarları
        fields.Select("Language", "language").Label("Dil").Options(map[string]string{
            "en": "English",
            "tr": "Türkçe",
        }).Default("en"),
        fields.Select("Theme", "theme").Label("Tema").Options(map[string]string{
            "light": "Açık",
            "dark":  "Koyu",
            "auto":  "Otomatik",
        }).Default("auto"),
    },
}
```

---

## Troubleshooting

### Şifre Değişmiyor

**Sorun**: Yeni şifre girildiğinde şifre değişmiyor.

**Çözüm**:
1. Mevcut şifrenin doğru girildiğinden emin olun
2. Yeni şifrenin confirm password ile eşleştiğinden emin olun
3. Şifre politikasına uygun olduğundan emin olun
4. Backend log'larını kontrol edin

### Profil Resmi Yüklenmiyor

**Sorun**: Profil resmi yüklendiğinde hata alıyorum.

**Çözüm**:
1. Dosya boyutunun limitlere uygun olduğundan emin olun
2. Dosya tipinin desteklendiğinden emin olun (JPG, PNG, GIF)
3. Storage dizininin yazma izinlerini kontrol edin
4. Disk alanının yeterli olduğundan emin olun

### Email Değişmiyor

**Sorun**: Email adresi değiştirilemiyor.

**Çözüm**:
1. Yeni email adresinin benzersiz olduğundan emin olun
2. Email formatının doğru olduğundan emin olun
3. Email doğrulama sisteminin aktif olup olmadığını kontrol edin

---

## Lisans

Bu özellik Panel.go projesinin bir parçasıdır ve aynı lisans altında dağıtılır.
