# API Key Yönetimi

Bu doküman, Panel.go içinde API key tabanlı erişimin nasıl yönetileceğini anlatır.

## Genel Bakış

Panel.go’da API key yönetimi OpenAI/better-auth tarzı managed key lifecycle ile yapılır.
API key veritabanında hash olarak saklanır ve oluşturulurken raw key yalnızca bir kez gösterilir.

## Nereden Yönetilir?

Varsayılan sayfa:

- Sayfa slug: `api-settings`
- UI yolu: `/api-settings`

Bu sayfa otomatik kayıt edilir ve sadece `admin` role sahip kullanıcılar erişebilir.

## Sayfadaki Alanlar

- `/api/openapi.json`
- `/api/docs`

Bu endpointler public erişime açıktır (session gerektirmez).

## Managed API Key Endpointleri

OpenAI/better-auth benzeri çoklu key lifecycle için:

- `GET /api/api-keys` -> key listesini döner
- `POST /api/api-keys` -> yeni key üretir
- `DELETE /api/api-keys/:id` -> key revoke eder

### Create Request

```json
{
  "name": "CI Key",
  "expires_at": "2026-12-31T23:59:59Z"
}
```

`expires_at` opsiyoneldir ve RFC3339 formatındadır.

### Create Response

```json
{
  "data": {
    "id": 1,
    "name": "CI Key",
    "prefix": "pnl_xxxxxxxx",
    "status": "active",
    "created_at": "2026-02-15T12:00:00Z"
  },
  "key": "pnl_very_secret_generated_value"
}
```

`key` değeri yalnızca create response sırasında bir kez döner.
Veritabanında raw key değil, yalnızca SHA-256 hash saklanır.

## Yetki Kuralları

- Managed key yönetim endpointleri (`/api/api-keys*`) için **admin session** gerekir.
- API key ile authenticate olmuş istekler bu yönetim endpointlerine erişemez (`403`).
- API key doğrulanan istekler resource API’lerine erişebilir.

## Örnek Kullanım (cURL)

### Geçerli key ile resource listesi

```bash
curl -H "X-API-Key: pnl_very_secret_generated_value" \
  http://localhost:8080/api/resource/verifications
```

### Geçersiz key

```bash
curl -H "X-API-Key: wrong-key" \
  http://localhost:8080/api/resource/verifications
```

Beklenen cevap: `401 Unauthorized`

## Auth Davranışı

Bir istek geçerli API key ile doğrulanırsa:

- Session middleware isteği kabul eder
- Context içine sentetik bir kullanıcı (`Role: admin`) yerleştirilir
- Session gerektiren API endpoint'lerine erişim sağlanır

Bu davranış server-to-server entegrasyonlar için tasarlanmıştır.

## Güvenlik Notları

- API key'leri güçlü, uzun ve rastgele üretin
- Key'leri düzenli rotate edin (yenisini oluştur, eskisini revoke et)
- Üretimde HTTPS zorunlu kullanın
- API key erişimi olan servisleri ağ seviyesinde sınırlayın (allowlist/VPN)

## İlgili Dosyalar

- `pkg/panel/default_api_page.go`
- `pkg/panel/api_key_management.go`
- `pkg/middleware/api_key.go`
- `pkg/openapi/spec.go`
- `pkg/handler/auth/handler.go`
- `pkg/handler/openapi_handler.go`
