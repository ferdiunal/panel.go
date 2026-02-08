# Circuit Breaker ve i18n Test Uygulaması

Bu örnek uygulama, Panel.go'da Circuit Breaker ve i18n middleware'lerinin nasıl çalıştığını gösterir.

## Çalıştırma

```bash
cd examples/circuit-breaker-i18n
go run main.go
```

## Test Endpoint'leri

### 1. Basit Çeviri Testi

**Türkçe:**
```bash
curl http://localhost:8080/api/test/welcome?lang=tr
```

**İngilizce:**
```bash
curl http://localhost:8080/api/test/welcome?lang=en
```

**Yanıt:**
```json
{
  "message": "Hoş geldiniz",
  "lang": "tr"
}
```

### 2. Template ile Çeviri Testi

**Türkçe:**
```bash
curl http://localhost:8080/api/test/welcome/Ahmet?lang=tr
```

**İngilizce:**
```bash
curl http://localhost:8080/api/test/welcome/John?lang=en
```

**Yanıt:**
```json
{
  "message": "Hoş geldiniz, Ahmet",
  "name": "Ahmet",
  "lang": "tr"
}
```

### 3. Circuit Breaker Testi - Hata Simülasyonu

Bu endpoint, ilk 5 istekte hata döndürür ve circuit breaker'ı tetikler.

**İlk 5 istek (hata):**
```bash
curl http://localhost:8080/api/test/error
```

**Yanıt (1-5. istek):**
```json
{
  "error": "Simulated error",
  "count": 1,
  "message": "Bu hata circuit breaker'ı tetiklemek için simüle edildi"
}
```

**6. istek (circuit breaker açık):**
```json
{
  "error": "Service temporarily unavailable",
  "message": "The service is experiencing high failure rates. Please try again later.",
  "code": "CIRCUIT_BREAKER_OPEN"
}
```

**Timeout sonrası (half-open):**
```json
{
  "error": "Service is recovering",
  "message": "The service is testing recovery. Please wait.",
  "code": "CIRCUIT_BREAKER_HALF_OPEN"
}
```

**Kurtarma sonrası (başarılı):**
```json
{
  "success": true,
  "count": 6,
  "message": "Servis kurtarıldı"
}
```

### 4. Başarılı İstek Testi

```bash
curl http://localhost:8080/api/test/success?lang=tr
```

**Yanıt:**
```json
{
  "success": true,
  "message": "Kayıt başarıyla oluşturuldu",
  "lang": "tr"
}
```

### 5. Tüm Çevirileri Listele

```bash
curl http://localhost:8080/api/test/translations?lang=tr
```

**Yanıt:**
```json
{
  "lang": "tr",
  "translations": {
    "welcome": "Hoş geldiniz",
    "error.notFound": "Kayıt bulunamadı",
    "error.unauthorized": "Bu işlem için yetkiniz yok",
    "error.serverError": "Sunucu hatası oluştu",
    "circuitBreaker.open": "Servis geçici olarak kullanılamıyor...",
    "success.created": "Kayıt başarıyla oluşturuldu",
    "button.save": "Kaydet",
    "navigation.dashboard": "Kontrol Paneli"
  }
}
```

## Circuit Breaker Durumları

Circuit breaker üç durum arasında geçiş yapar:

1. **Closed (Kapalı)**: Normal çalışma, istekler geçer
2. **Open (Açık)**: Devre açık, istekler reddedilir (503)
3. **Half-Open (Yarı Açık)**: Test modu, sınırlı istek geçer (429)

## Dil Değiştirme

Dil, şu sırayla belirlenir:

1. **Query parametresi**: `?lang=tr` veya `?lang=en`
2. **Accept-Language header**: `Accept-Language: tr-TR,tr;q=0.9`
3. **DefaultLanguage**: Varsayılan dil (tr)

## Monitoring

Uygulama çalışırken console'da şu log'ları göreceksiniz:

```
❌ Hata simülasyonu: 1/5
❌ Hata simülasyonu: 2/5
❌ Hata simülasyonu: 3/5
[CIRCUIT_BREAKER] State changed: CLOSED -> OPEN
✅ Başarılı yanıt: 6
[CIRCUIT_BREAKER] State changed: OPEN -> HALF_OPEN
[CIRCUIT_BREAKER] State changed: HALF_OPEN -> CLOSED
```

## Notlar

- Circuit breaker yapılandırması test için düşük değerlerle ayarlanmıştır
- Üretim ortamında daha yüksek threshold ve timeout değerleri kullanın
- Dil dosyaları `locales/` dizininde bulunmalıdır
