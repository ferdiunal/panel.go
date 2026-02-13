# Security

Panel.go, production ortamında güvenli bir uygulama çalıştırmanız için kapsamlı güvenlik özellikleri sunar.

## Özellikler

- **Auth Hardening**: Güçlendirilmiş kimlik doğrulama
- **Brute Force Protection**: Otomatik hesap kilitleme
- **Password Policy**: Güçlü şifre gereksinimleri
- **Email Validation**: RFC-uyumlu email doğrulama
- **Transaction-based Registration**: Atomik kullanıcı kaydı
- **Role-based Access Control**: Admin ve user rolleri
- **Rate Limiting**: Login endpoint'i için rate limiting
- **Audit Logging**: Tüm değişikliklerin loglanması

## Auth Hardening

Panel.go'nun auth servisi, güvenlik best practice'lerini uygular.

### Email Validation

Email adresleri RFC 5322 standardına göre doğrulanır:

```go
// Geçerli email formatları
"user@example.com"          // ✅ Geçerli
"user.name@example.com"     // ✅ Geçerli
"user+tag@example.com"      // ✅ Geçerli

// Geçersiz email formatları
"invalid"                   // ❌ Geçersiz
"@example.com"              // ❌ Geçersiz
"user@"                     // ❌ Geçersiz
```

Email adresleri otomatik olarak normalize edilir:
- Boşluklar temizlenir
- Küçük harfe çevrilir

```go
// Tüm bu email'ler aynı kullanıcıya işaret eder
"  USER@EXAMPLE.COM  "  // → "user@example.com"
"User@Example.Com"      // → "user@example.com"
"user@example.com"      // → "user@example.com"
```

### Password Policy

Şifre gereksinimleri:

- **Minimum uzunluk**: 8 karakter
- **Whitespace kontrolü**: Boşluk, tab, newline karakterleri yasak
- **Karakter çeşitliliği**: Önerilir ama zorunlu değil

```go
// Geçerli şifreler
"Password1"              // ✅ 9 karakter
"MySecurePass123"        // ✅ 16 karakter
"P@ssw0rd!"              // ✅ Özel karakterler
"VeryLongPassword123"    // ✅ Uzun şifre

// Geçersiz şifreler
"Pass1"                  // ❌ Çok kısa (5 karakter)
"Pass word1"             // ❌ Boşluk içeriyor
"Pass\tword1"            // ❌ Tab içeriyor
"Pass\nword1"            // ❌ Newline içeriyor
```

### Password Hashing

Şifreler bcrypt ile hash'lenir:

- **Development/Test**: `bcrypt.MinCost` (4) - Hızlı testler için
- **Production**: `bcrypt.DefaultCost` (10) - Güvenli hash için

```go
// Otomatik olarak environment'a göre seçilir
func resolvePasswordHashCost() int {
    if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
        return bcrypt.MinCost  // Test ortamı
    }
    return bcrypt.DefaultCost  // Production ortamı
}
```

## Brute Force Protection

Panel.go, brute force saldırılarına karşı otomatik koruma sağlar.

### Nasıl Çalışır?

1. **Failed Attempt Tracking**: Her başarısız login denemesi kaydedilir
2. **Attempt Window**: 15 dakikalık zaman penceresi
3. **Max Attempts**: 5 başarısız deneme
4. **Lockout Duration**: 15 dakika kilitleme

### Lockout Mekanizması

```
Attempt 1: ❌ Wrong password → Count: 1
Attempt 2: ❌ Wrong password → Count: 2
Attempt 3: ❌ Wrong password → Count: 3
Attempt 4: ❌ Wrong password → Count: 4
Attempt 5: ❌ Wrong password → Count: 5 → 🔒 LOCKED for 15 minutes

Attempt 6: ❌ Correct password → 🚫 ErrTooManyAttempts
...
After 15 minutes: ✅ Lockout expires, can try again
```

### Attempt Key

Lockout, email + IP kombinasyonuna göre yapılır:

```go
attemptKey := email + "|" + ip

// Örnekler:
"user@example.com|192.168.1.100"
"admin@example.com|10.0.0.1"
```

Bu sayede:
- Aynı kullanıcı farklı IP'lerden deneyebilir
- Farklı kullanıcılar aynı IP'den etkilenmez

### Lockout Temizleme

Başarılı login sonrası lockout otomatik temizlenir:

```go
// Başarılı login
sess, err := service.LoginEmail(ctx, email, password, ip, userAgent)
if err == nil {
    // Lockout temizlendi ✅
}
```

### Test Örneği

```go
func TestLoginBruteForceLockout(t *testing.T) {
    service := newTestService(t)
    ctx := context.Background()

    // Kullanıcı oluştur
    service.RegisterEmail(ctx, "User", "user@example.com", "Password1")

    // 5 başarısız deneme
    for i := 0; i < 5; i++ {
        service.LoginEmail(ctx, "user@example.com", "WrongPass", "127.0.0.1", "test")
    }

    // 6. deneme - doğru şifre bile olsa reddedilir
    _, err := service.LoginEmail(ctx, "user@example.com", "Password1", "127.0.0.1", "test")
    if err != ErrTooManyAttempts {
        t.Fatal("Expected lockout")
    }
}
```

## Transaction-based Registration

Kullanıcı kaydı atomik bir transaction içinde yapılır. Bu, veri tutarlılığını garanti eder.

### Neden Transaction?

Kullanıcı kaydı 2 adımdan oluşur:
1. User kaydı oluştur
2. Account (credential) kaydı oluştur

Transaction olmadan:
```
1. User oluşturuldu ✅
2. Account oluşturma başarısız ❌
→ Orphan user kaydı (şifresiz kullanıcı) 💥
```

Transaction ile:
```
BEGIN TRANSACTION
1. User oluşturuldu ✅
2. Account oluşturma başarısız ❌
ROLLBACK
→ Hiçbir kayıt oluşturulmadı ✅
```

### Implementation

```go
func (s *Service) RegisterEmail(ctx context.Context, name, email, password string) (*user.User, error) {
    var createdUser *user.User

    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Transaction içinde repository'ler oluştur
        txUserRepo := orm.NewUserRepository(tx)
        txAccountRepo := orm.NewAccountRepository(tx)

        // 1. User oluştur
        u := &user.User{...}
        if err := txUserRepo.CreateUser(ctx, u); err != nil {
            return err  // Rollback
        }

        // 2. Account oluştur
        acc := &account.Account{...}
        if err := txAccountRepo.Create(ctx, acc); err != nil {
            return err  // Rollback
        }

        createdUser = u
        return nil  // Commit
    })

    return createdUser, err
}
```

## Role-based Access Control

Panel.go, role-based access control (RBAC) destekler.

### Roller

```go
const (
    RoleAdmin = "admin"  // Tam yetki
    RoleUser  = "user"   // Sınırlı yetki
)
```

### İlk Kullanıcı Admin Olur

Güvenlik için, ilk kayıt olan kullanıcı otomatik olarak admin olur:

```go
role := user.RoleUser
var totalUsers int64
if err := tx.Model(&user.User{}).Count(&totalUsers).Error; err == nil && totalUsers == 0 {
    role = user.RoleAdmin  // İlk kullanıcı admin
}
```

Bu sayede:
- Uygulama ilk kurulumda admin hesabı oluşturulur
- Sonraki kullanıcılar normal user olarak kaydolur

### Policy Örneği

```go
type UserPolicy struct{}

func (p UserPolicy) ViewAny(ctx *appContext.Context) bool {
    authUser := ctx.User()
    return authUser != nil && authUser.Role == domainUser.RoleAdmin
}

func (p UserPolicy) View(ctx *appContext.Context, model interface{}) bool {
    authUser := ctx.User()
    if authUser == nil {
        return false
    }

    // Admin her şeyi görebilir
    if authUser.Role == domainUser.RoleAdmin {
        return true
    }

    // User sadece kendi kaydını görebilir
    userModel := model.(*domainUser.User)
    return userModel.ID == authUser.ID
}
```

## Rate Limiting

Login endpoint'i için rate limiting uygulanır.

### Konfigürasyon

```go
authLoginLimiter := limiter.New(limiter.Config{
    Max:        10,                    // Maksimum 10 request
    Expiration: time.Minute,           // 1 dakika içinde
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()                  // IP bazlı
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
            "error": "too many login requests",
        })
    },
})

authRoutes.Post("/sign-in/email", authLoginLimiter, context.Wrap(authH.LoginEmail))
```

### Davranış

```
Request 1-10: ✅ İzin verilir
Request 11+:  🚫 429 Too Many Requests
After 1 min:  ✅ Counter sıfırlanır
```

### IP Bazlı Limiting

Rate limiting IP adresine göre yapılır:
- Aynı IP'den 1 dakikada maksimum 10 login denemesi
- Farklı IP'ler birbirini etkilemez

## Session Management

### Session Oluşturma

Login başarılı olduğunda session oluşturulur:

```go
sessionId, _ := uuid.NewV7()
sessionToken, _ := uuid.NewV7()

sess := &session.Session{
    ID:        sessionId.String(),
    UserID:    u.ID,
    Token:     sessionToken.String(),
    User:      u,
    ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 gün
    IPAddress: ip,
    UserAgent: userAgent,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

### Session Validation

Session token ile kullanıcı doğrulanır:

```go
sess, err := service.ValidateSession(ctx, token)
if err != nil {
    return nil, err
}

if sess.ExpiresAt.Before(time.Now()) {
    return nil, errors.New("session expired")
}
```

### Session Expiration

Session'lar 7 gün sonra otomatik olarak expire olur. Expired session'lar geçersizdir.

### Logout

Logout işlemi session'ı siler:

```go
err := service.Logout(ctx, token)
// Session veritabanından silindi
```

## Error Handling

Auth servisi, güvenlik için generic error mesajları döndürür.

### Login Errors

```go
// Kullanıcı bulunamadı veya şifre yanlış
ErrInvalidCredentials = errors.New("invalid credentials")

// Çok fazla başarısız deneme
ErrTooManyAttempts = errors.New("too many failed login attempts")
```

**Neden generic?**

Spesifik error mesajları saldırganlara bilgi verir:
- ❌ "User not found" → Email'in sistemde olup olmadığını öğrenir
- ❌ "Wrong password" → Email'in geçerli olduğunu öğrenir
- ✅ "Invalid credentials" → Hiçbir bilgi vermez

### Registration Errors

```go
ErrEmailAlreadyExists = errors.New("email already exists")
ErrInvalidEmail       = errors.New("invalid email")
ErrWeakPassword       = errors.New("password does not meet policy requirements")
ErrInvalidName        = errors.New("invalid name")
```

Registration'da daha spesifik error'lar döndürülür çünkü kullanıcı kendi bilgilerini giriyor.

### HTTP Status Codes

```go
// Registration
400 Bad Request       → Invalid email, weak password, invalid name
409 Conflict          → Email already exists
500 Internal Server   → Database error

// Login
401 Unauthorized      → Invalid credentials
429 Too Many Requests → Rate limit exceeded, too many attempts
500 Internal Server   → Database error
```

## Security Best Practices

### 1. Environment Variables

Hassas bilgileri environment variable'larda saklayın:

```env
# .env
COOKIE_ENCRYPTION_KEY=<32-byte-base64-key>
DATABASE_DSN=postgres://user:pass@localhost/db
```

**Asla commit etmeyin:**
```gitignore
.env
*.key
*.pem
```

### 2. HTTPS Kullanın

Production'da mutlaka HTTPS kullanın:

```go
// Reverse proxy (nginx, caddy) ile HTTPS
// veya
app.ListenTLS(":443", "cert.pem", "key.pem")
```

### 3. CSRF Protection

Production'da CSRF protection aktif:

```go
if config.Environment == "production" {
    app.Use(csrf.New())
}
```

### 4. Helmet Middleware

Security header'ları otomatik eklenir:

```go
app.Use(helmet.New(helmet.Config{
    CrossOriginResourcePolicy: "cross-origin",
}))
```

### 5. Database Injection Prevention

GORM parametreli query'ler kullanır:

```go
// ✅ Güvenli
db.Where("email = ?", email).First(&user)

// ❌ Güvensiz (kullanmayın)
db.Where("email = '" + email + "'").First(&user)
```

### 6. Password Storage

Şifreler asla plain text saklanmaz:

```go
// ✅ Güvenli - bcrypt hash
hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ❌ Güvensiz - plain text (asla kullanmayın)
user.Password = password
```

### 7. Session Token Security

Session token'lar UUID v7 ile oluşturulur:

```go
// ✅ Güvenli - UUID v7 (time-ordered, random)
sessionToken, _ := uuid.NewV7()

// ❌ Güvensiz - tahmin edilebilir
sessionToken := strconv.Itoa(time.Now().Unix())
```

### 8. Audit Logging

Tüm değişiklikler otomatik loglanır:

```go
// Audit middleware otomatik çalışır
api.Use(context.Wrap(obs.AuditMiddleware(db)))
```

### 9. Rate Limiting

Kritik endpoint'lerde rate limiting kullanın:

```go
limiter := limiter.New(limiter.Config{
    Max:        10,
    Expiration: time.Minute,
})
app.Post("/api/auth/sign-in/email", limiter, handler)
```

### 10. Regular Updates

Bağımlılıkları düzenli güncelleyin:

```bash
go get -u ./...
go mod tidy
```

## Security Checklist

Production'a geçmeden önce kontrol edin:

- [ ] HTTPS aktif
- [ ] Environment variable'lar güvenli saklanıyor
- [ ] `.env` dosyası `.gitignore`'da
- [ ] CSRF protection aktif
- [ ] Rate limiting yapılandırıldı
- [ ] Audit logging aktif
- [ ] Password policy uygulanıyor
- [ ] Brute force protection çalışıyor
- [ ] Session expiration ayarlandı
- [ ] Database backup stratejisi var
- [ ] Monitoring ve alerting kuruldu
- [ ] Security header'lar aktif
- [ ] Bağımlılıklar güncel

## Vulnerability Scanning

### govulncheck

Go vulnerability database'ini kontrol edin:

```bash
# Makefile ile
make vuln

# Veya direkt
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

### CI/CD Integration

GitHub Actions workflow'da otomatik vulnerability scanning:

```yaml
- name: Install govulncheck
  run: go install golang.org/x/vuln/cmd/govulncheck@latest

- name: Run govulncheck
  run: govulncheck ./...
```

## Incident Response

### Şüpheli Aktivite Tespiti

Audit log'ları kullanarak şüpheli aktiviteleri tespit edin:

```go
// Çok sayıda başarısız login
var logs []audit.Log
db.Where("action = ? AND status_code = ? AND created_at > ?",
    "auth:sign-in", 401, time.Now().Add(-1*time.Hour)).
   Find(&logs)

// Aynı IP'den çok sayıda farklı kullanıcı denemesi
var logs []audit.Log
db.Where("ip_address = ? AND action = ? AND created_at > ?",
    suspiciousIP, "auth:sign-in", time.Now().Add(-1*time.Hour)).
   Find(&logs)
```

### Hesap Kilitleme

Şüpheli hesapları manuel olarak kilitleyin:

```go
// User entity'ye disabled field ekleyin
type User struct {
    // ...
    Disabled bool `json:"disabled" gorm:"default:false"`
}

// Login'de kontrol edin
if user.Disabled {
    return nil, errors.New("account disabled")
}
```

### Session İptali

Tüm session'ları iptal edin:

```go
// Belirli bir kullanıcının tüm session'larını sil
db.Where("user_id = ?", userID).Delete(&session.Session{})

// Tüm session'ları sil (acil durum)
db.Exec("TRUNCATE TABLE sessions")
```

## İlgili Dökümanlar

- [Observability](OBSERVABILITY.md) - Monitoring ve logging
- [Authentication](Authentication.md) - Auth sistemi
- [Authorization](Authorization.md) - Yetkilendirme
- [API Reference](API-Reference.md) - API dokümantasyonu
