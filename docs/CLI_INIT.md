# Panel CLI - `init` Komutu

`panel init`, yeni bir Panel.go projesi için başlangıç iskeletini üretir.

Bu dokümanın amacı, komutun ne ürettiğini ve komut sonrası hangi adımlarla geliştirmeye devam edeceğinizi netleştirmektir.

## Komut Özeti

### İnteraktif kullanım (önerilen)

```bash
panel init
```

### Flag ile kullanım

```bash
panel init -d sqlite
panel init -d postgres
panel init -d mysql
```

## `panel init` Ne Yapar?

- Başlangıç proje dosyalarını oluşturur
- Seçilen veritabanına göre örnek bağlantı konfigürasyonu yazar
- Kod üretim şablonlarını `.panel/stubs/` altına kopyalar
- Yardımcı skill dosyalarını `.claude/skills/` altına kopyalar
- `.env` içine güvenli `COOKIE_ENCRYPTION_KEY` üretir

## Veritabanı Seçenekleri

### SQLite (default)

```env
DATABASE_DRIVER=sqlite
DATABASE_DSN=proje-adi.db
```

### PostgreSQL

```env
DATABASE_DRIVER=postgres
DATABASE_DSN=host=localhost user=postgres password=postgres dbname=proje-adi port=5432 sslmode=disable TimeZone=UTC
```

### MySQL

```env
DATABASE_DRIVER=mysql
DATABASE_DSN=user:password@tcp(localhost:3306)/proje-adi?charset=utf8mb4&parseTime=True&loc=Local
```

## Oluşturulan Dosyalar (Başlangıçta Kritik Olanlar)

### 1) `main.go`

Panel uygulamasının giriş noktasıdır. Veritabanı bağlantısı ve `panel.New(...)` burada başlar.

### 2) `.env`

Veritabanı bağlantısı, host/port ve environment burada tutulur.

### 3) `go.mod`

Go module tanımı ve bağımlılık yönetimi.

### 4) `.panel/stubs/`

`make:model`, `make:resource`, `make:page` gibi komutlar için şablonlar.

### 5) `.claude/skills/`

Kod üretimini hızlandıran yardımcı skill içerikleri.

## Init Sonrası Net Akış (Önerilen)

### 1) Model oluştur

```bash
panel make:model Post
```

### 2) Resource oluştur

```bash
panel make:resource Post
```

### 3) `main.go` içinde bağla ve çalıştır

- Veritabanı bağlantısını doğrula
- Resource'u panele kaydet
- Uygulamayı başlat

```bash
go mod tidy
go run main.go
```

## Tamamlandıktan Sonra Hangi Dokümanı Okumalıyım?

Detaylı teknik kurulum akışı için doğrudan şuraya geçin:

- [Başlarken](Getting-Started)

## Sık Sorunlar

### `panel: command not found`

CLI binary PATH içinde değildir.

```bash
go install github.com/ferdiunal/panel.go/cmd/panel@latest
```

### Veritabanı bağlantı hatası

`.env` içindeki `DATABASE_DSN` değerini gerçek bağlantı bilgilerinizle güncelleyin.

### Resource oluştu ama panelde görünmüyor

`main.go` içinde resource import/register akışını ve slug çakışmalarını kontrol edin.
