# Panel.go

**Panel.go**, modern Go teknolojileri ile geliştirilmiş, tam yığınlı (full-stack) web uygulamasıdır. Fiber web framework'ü, PostgreSQL veritabanı ve HTMX tabanlı interaktif kullanıcı arayüzü ile yüksek performanslı ve güvenli bir yönetim paneli sunar.

## 🚀 Özellikler

- **🔐 Güvenli Kimlik Doğrulama**: JWT tabanlı oturum yönetimi ve çok faktörlü kimlik doğrulama
- **👤 Kullanıcı Yönetimi**: Kayıt, giriş, profil yönetimi ve avatar yükleme
- **📧 Bildirim Sistemi**: E-posta ve SMS entegrasyonları ile gelişmiş bildirim altyapısı
- **🎨 Modern UI/UX**: Tailwind CSS v4 ve HTMX ile responsive, interaktif arayüz
- **🏗️ Clean Architecture**: Katmanlı mimari ile sürdürülebilir ve test edilebilir kod yapısı
- **⚡ Yüksek Performans**: Fiber v2 ile optimize edilmiş HTTP sunucusu
- **🗄️ PostgreSQL Entegrasyonu**: Ent ORM ile tip güvenli veritabanı işlemleri
- **🔄 Hot Reload**: Geliştirme sırasında anında yeniden yükleme desteği
- **📱 Mobil Uyumlu**: Responsive tasarım ile tüm cihazlarda mükemmel görünüm

## 🏗️ Teknoloji Altyapısı

### Backend
- **Go 1.25.0**: Ana programlama dili
- **Fiber v2**: Yüksek performanslı web framework
- **Ent ORM**: Tip güvenli veritabanı işlemleri
- **PostgreSQL**: Güçlü ve ölçeklenebilir veritabanı

### Frontend
- **Templ**: Go tabanlı tip güvenli şablon motoru
- **HTMX v2**: Minimal JavaScript ile interaktif arayüzler
- **Tailwind CSS v4**: Utility-first CSS framework

### Geliştirme Araçları
- **Air**: Hot reload ve canlı geliştirme
- **Docker Compose**: Konteynerleştirilmiş veritabanı
- **Go Modules**: Bağımlılık yönetimi

## 📋 Gereksinimler

- **Go**: 1.25.0 veya üzeri
- **PostgreSQL**: 12+ sürüm
- **Node.js**: 16+ sürüm (web assets için)
- **Docker**: Veritabanı konteyneri için
- **Make**: Build komutları için

## ⚡ Hızlı Başlangıç

### 1. Projeyi Klonlayın
```bash
git clone <repository-url>
cd panel.go
```

### 2. Ortam Değişkenlerini Ayarlayın
`.env` dosyası oluşturun ve aşağıdaki değişkenleri tanımlayın:
```env
# Veritabanı Konfigürasyonu
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=panel_db
BLUEPRINT_DB_USERNAME=panel_user
BLUEPRINT_DB_PASSWORD=your_secure_password

# Uygulama Konfigürasyonu
PORT=8080

# E-posta/SMS API Anahtarları (opsiyonel)
SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USER=your_email
SMTP_PASS=your_password

# SMS Gateway (opsiyonel)
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
```

### 3. Veritabanını Başlatın
```bash
# Docker Compose ile PostgreSQL konteynerini başlatın
make docker-run
```

### 4. Uygulamayı Derleyin ve Çalıştırın
```bash
# Tüm bağımlılıkları yükleyin ve uygulamayı derleyin
make build

# Uygulamayı çalıştırın
make run
```

### 5. Tarayıcıda Açın
Uygulama `http://localhost:8080` adresinde çalışacaktır.

## 🛠️ Geliştirme Komutları

### Temel Komutlar
```bash
# Tüm işlemleri gerçekleştir (derleme + test)
make all

# Uygulamayı derle
make build

# Uygulamayı çalıştır
make run

# Temizlik
make clean
```

### Veritabanı İşlemleri
```bash
# PostgreSQL konteynerini başlat
make docker-run

# PostgreSQL konteynerini durdur
make docker-down
```

### Test İşlemleri
```bash
# Tüm testleri çalıştır
make test

# Sadece entegrasyon testlerini çalıştır
make itest
```

### Geliştirme Modu
```bash
# Hot reload ile geliştirme (Air)
make watch

# Templ CLI'yi yükle
make templ-install

# Web bağımlılıklarını yükle
make web-install
```

## 📁 Proje Yapısı

```
panel.go/
├── cmd/
│   ├── api/                 # Ana uygulama giriş noktası
│   │   └── main.go
│   ├── cli/                 # CLI araçları
│   ├── entgo/               # Ent kod üreteci
│   └── web/                 # Web bileşenleri ve şablonlar
│       ├── assets/          # Statik dosyalar (JS, CSS)
│       ├── deps/            # Şablon bileşenleri
│       ├── templates/       # E-posta şablonları
│       └── *.templ          # Templ şablonları
├── internal/
│   ├── constants/           # Uygulama sabitleri
│   ├── entities/            # Domain modelleri
│   ├── errors/              # Hata tanımları
│   ├── handler/             # HTTP işleyicileri
│   ├── infrastructure/      # Harici servis entegrasyonları
│   │   ├── email/           # E-posta servisi
│   │   ├── notification/    # Bildirim sistemi
│   │   └── sms/             # SMS servisi
│   ├── interfaces/          # Arayüz tanımları
│   ├── middleware/          # HTTP middleware'ler
│   ├── repository/          # Veri erişim katmanı
│   ├── resource/            # Veri dönüştürücüleri
│   ├── server/              # Sunucu konfigürasyonu
│   └── service/             # İş mantığı servisleri
├── shared/                  # Paylaşılan yardımcılar
│   ├── encrypt/             # Şifreleme yardımcıları
│   ├── uuid/                # UUID işlemleri
│   └── validate/            # Doğrulama yardımcıları
├── docker-compose.yml       # Docker konfigürasyonu
├── Makefile                 # Build komutları
├── go.mod                   # Go modülleri
└── README.md
```

## 🏛️ Mimari Tasarım

### Clean Architecture Yaklaşımı

Panel.go, **Clean Architecture** prensiplerine göre tasarlanmıştır:

#### 1. **Sunum Katmanı** (`cmd/web/`)
- HTTP istek/cevap yönetimi
- Templ şablonları ile sunucu taraflı render
- HTMX ile progressive enhancement

#### 2. **Uygulama Katmanı** (`internal/server/`)
- İstek yönlendirme ve middleware orkestrasyonu
- HTTP bağlam yönetimi
- Zincir sorumluluğu (Chain of Responsibility) pattern

#### 3. **İş Mantığı Katmanı** (`internal/service/`)
- Temel iş kuralları ve kullanım senaryoları
- İşlem yönetimi ve validasyon
- Servis katmanı pattern

#### 4. **Veri Erişim Katmanı** (`internal/repository/`)
- Veritabanı işlemleri ve sorgu oluşturma
- Repository pattern implementasyonu
- Veri eşleme işlemleri

#### 5. **Domain Katmanı** (`internal/entities/`)
- Domain varlıkları ve değer nesneleri
- İş kuralları ve kısıtlamalar
- Domain event'leri

## 🔒 Güvenlik Özellikleri

- **Giriş Doğrulama**: Güvenli kimlik doğrulama akışları
- **Şifre Yönetimi**: Argon2 tabanlı şifre hash'leme
- **Oturum Güvenliği**: JWT token yönetimi
- **CSRF Koruması**: Cross-site request forgery önleme
- **XSS Koruması**: Template auto-escaping
- **Rate Limiting**: İstek sınırlaması
- **Güvenlik Başlıkları**: HSTS, CSP, X-Frame-Options

## 📊 Performans Optimizasyonları

- **Connection Pooling**: Veritabanı bağlantı havuzu
- **Template Caching**: Şablon ön derleme
- **Static Asset Optimization**: CSS/JS sıkıştırma
- **Gzip Compression**: Yanıt sıkıştırma
- **Prepared Statements**: SQL performans optimizasyonu

## 🧪 Test Stratejisi

- **Unit Testler**: İşlevsellik testi
- **Integration Testler**: Sistem entegrasyonu testi
- **End-to-End Testler**: Kullanıcı akışı testi
- **Performance Testler**: Yük ve performans testi

## 🚀 Dağıtım

### Production Ortamı
```bash
# Production build
make build

# Docker ile dağıtım (opsiyonel)
docker build -t panel.go .
docker run -p 8080:8080 panel.go
```

### Environment Variables (Production)
```env
PORT=8080
DB_HOST=your_production_db_host
DB_PORT=5432
DB_DATABASE=panel_prod
DB_USERNAME=prod_user
DB_PASSWORD=secure_prod_password
```

## 🤝 Katkıda Bulunma

1. Fork edin
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Commit edin (`git commit -m 'Add amazing feature'`)
4. Push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

## 📝 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakınız.

## 👥 Destek

Sorularınız ve geri bildirimleriniz için:

- **Issues**: GitHub Issues sayfasını kullanın
- **Discussions**: Genel tartışmalar için
- **Documentation**: Detaylı dokümantasyon için `docs/` klasörüne bakın

## 🔄 Güncellemeler

### v1.0.0
- İlk kararlı sürüm
- Temel kullanıcı yönetimi
- Güvenlik özellikleri
- E-posta/SMS entegrasyonları

---

**Panel.go** ile modern, güvenli ve ölçeklenebilir web uygulamaları geliştirin! 🚀
