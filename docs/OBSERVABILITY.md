# Observability

Panel.go, production ortamında uygulamanızı izlemek ve debug etmek için kapsamlı observability özellikleri sunar.

## Özellikler

- **Metrics**: Prometheus-uyumlu HTTP metrics
- **Health Checks**: Liveness ve readiness probe'ları
- **Audit Logging**: Tüm API değişikliklerinin otomatik loglanması
- **Request Logging**: Structured JSON logging

## Metrics

Panel.go, Prometheus formatında HTTP metrics sağlar.

### Endpoint

```
GET /metrics
```

### Metrikler

**panel_requests_total** (counter)
- Toplam HTTP request sayısı

**panel_request_errors_total** (counter)
- Hata dönen request sayısı (status >= 400)

**panel_request_duration_seconds_avg** (gauge)
- Ortalama request süresi (saniye)

**panel_request_status_total{status="200"}** (counter)
- HTTP status code'a göre request sayısı

**panel_request_route_total{method="GET",path="/api/resource/users"}** (counter)
- Route ve method'a göre request sayısı

### Prometheus Konfigürasyonu

```yaml
scrape_configs:
  - job_name: 'panel-go'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Örnek Metrik Çıktısı

```
# TYPE panel_requests_total counter
panel_requests_total 1523
# TYPE panel_request_errors_total counter
panel_request_errors_total 42
# TYPE panel_request_duration_seconds_avg gauge
panel_request_duration_seconds_avg 0.125000
panel_request_status_total{status="200"} 1450
panel_request_status_total{status="404"} 31
panel_request_status_total{status="500"} 11
panel_request_route_total{method="GET",path="/api/resource/users"} 234
panel_request_route_total{method="POST",path="/api/auth/sign-in/email"} 89
```

## Health Checks

Panel.go, Kubernetes ve diğer orchestration sistemleri için health check endpoint'leri sağlar.

### Liveness Probe

Uygulamanın çalışıp çalışmadığını kontrol eder.

```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "time": "2026-02-13T15:00:00Z"
}
```

**Status Code:** 200 OK

### Readiness Probe

Uygulamanın trafiği kabul etmeye hazır olup olmadığını kontrol eder. Veritabanı bağlantısını test eder.

```
GET /ready
```

**Success Response:**
```json
{
  "status": "ready"
}
```

**Status Code:** 200 OK

**Failure Response:**
```json
{
  "status": "not_ready",
  "error": "database connection failed"
}
```

**Status Code:** 503 Service Unavailable

### Kubernetes Örneği

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: panel-go
spec:
  containers:
  - name: panel-go
    image: panel-go:latest
    ports:
    - containerPort: 8080
    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 10
      periodSeconds: 30
    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
```

## Audit Logging

Panel.go, tüm API değişikliklerini otomatik olarak loglar. Bu, güvenlik denetimi ve compliance için kritiktir.

### Audit Log Yapısı

Audit log'lar `audit_logs` tablosunda saklanır:

```go
type Log struct {
    ID         string                 // UUID
    UserID     string                 // İşlemi yapan kullanıcı
    SessionID  string                 // Session ID
    Action     string                 // create, update, delete
    Resource   string                 // resource:users, page:settings, auth:sign-in
    ResourceID string                 // İşlem yapılan kaydın ID'si
    Method     string                 // HTTP method (POST, PUT, DELETE)
    Path       string                 // Request path
    StatusCode int                    // HTTP status code
    IPAddress  string                 // Client IP
    UserAgent  string                 // User agent
    RequestID  string                 // X-Request-ID header
    Metadata   map[string]interface{} // Ek bilgiler
    CreatedAt  time.Time              // İşlem zamanı
}
```

### Hangi İşlemler Loglanır?

Audit logging sadece değişiklik yapan işlemleri loglar:

- ✅ **POST** - Create işlemleri
- ✅ **PUT** - Update işlemleri
- ✅ **PATCH** - Partial update işlemleri
- ✅ **DELETE** - Delete işlemleri
- ❌ **GET** - Read işlemleri (loglanmaz)
- ❌ **HEAD** - Head işlemleri (loglanmaz)
- ❌ **OPTIONS** - Options işlemleri (loglanmaz)

### Audit Log Sorgulama

```go
// Belirli bir kullanıcının tüm işlemlerini getir
var logs []audit.Log
db.Where("user_id = ?", userID).
   Order("created_at DESC").
   Limit(100).
   Find(&logs)

// Belirli bir resource'un değişiklik geçmişi
var logs []audit.Log
db.Where("resource = ? AND resource_id = ?", "resource:users", "123").
   Order("created_at DESC").
   Find(&logs)

// Son 24 saatteki tüm delete işlemleri
var logs []audit.Log
db.Where("action = ? AND created_at > ?", "delete", time.Now().Add(-24*time.Hour)).
   Order("created_at DESC").
   Find(&logs)
```

### Audit Log Retention

Audit log'lar sürekli büyüyeceği için retention policy uygulamanız önerilir:

```go
// 90 günden eski audit log'ları sil
db.Where("created_at < ?", time.Now().Add(-90*24*time.Hour)).
   Delete(&audit.Log{})
```

Bu işlemi cron job veya scheduled task olarak çalıştırabilirsiniz.

## Request Logging

Panel.go, tüm HTTP request'leri structured JSON formatında loglar.

### Log Formatı

```json
{
  "time": "2026-02-13T15:00:00Z",
  "level": "INFO",
  "msg": "http_request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/resource/users",
  "status": 200,
  "duration_ms": 125,
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0..."
}
```

### Log Seviyesi

Request logging varsayılan olarak `INFO` seviyesindedir. Hata durumlarında otomatik olarak `ERROR` seviyesine yükseltilmez, ancak status code'dan hata olup olmadığını anlayabilirsiniz.

### Log Çıktısı

Log'lar `stdout`'a JSON formatında yazılır. Production ortamında bu log'ları bir log aggregation sistemine (ELK, Loki, CloudWatch, vb.) yönlendirebilirsiniz.

### Docker Örneği

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o panel-go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/panel-go .
# Log'lar stdout'a yazılır, Docker bunları otomatik toplar
CMD ["./panel-go"]
```

### Log Filtreleme

Belirli path'leri loglamak istemiyorsanız, middleware'i özelleştirebilirsiniz:

```go
// internal/observability/logging.go dosyasını düzenleyin
func RequestLoggerMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Health check endpoint'lerini loglama
        if c.Path() == "/health" || c.Path() == "/ready" {
            return c.Next()
        }

        start := time.Now()
        err := c.Next()
        duration := time.Since(start)

        requestLogger.Info("http_request",
            "request_id", c.GetRespHeader(fiber.HeaderXRequestID),
            "method", c.Method(),
            "path", c.Path(),
            "status", c.Response().StatusCode(),
            "duration_ms", duration.Milliseconds(),
            "ip", c.IP(),
            "user_agent", c.Get(fiber.HeaderUserAgent),
        )
        return err
    }
}
```

## Monitoring Dashboard Örnekleri

### Grafana Dashboard

Prometheus metrics'lerini Grafana'da görselleştirebilirsiniz:

**Request Rate:**
```promql
rate(panel_requests_total[5m])
```

**Error Rate:**
```promql
rate(panel_request_errors_total[5m]) / rate(panel_requests_total[5m])
```

**Average Response Time:**
```promql
panel_request_duration_seconds_avg
```

**Top 10 Slowest Endpoints:**
```promql
topk(10, sum by (path) (rate(panel_request_route_total[5m])))
```

### Alerting

Prometheus Alertmanager ile alert'ler oluşturabilirsiniz:

```yaml
groups:
  - name: panel-go
    rules:
      - alert: HighErrorRate
        expr: rate(panel_request_errors_total[5m]) / rate(panel_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }}"

      - alert: SlowResponseTime
        expr: panel_request_duration_seconds_avg > 1.0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Slow response time detected"
          description: "Average response time is {{ $value }}s"
```

## Best Practices

1. **Metrics**: Prometheus'u 15-30 saniye aralıklarla scrape edin
2. **Health Checks**: Liveness probe'u 30 saniye, readiness probe'u 10 saniye aralıklarla kontrol edin
3. **Audit Logs**: 90-180 gün retention policy uygulayın
4. **Request Logs**: Log aggregation sistemi kullanın (ELK, Loki, CloudWatch)
5. **Alerting**: Error rate ve response time için alert'ler oluşturun
6. **Dashboard**: Grafana veya benzeri bir tool ile metrics'leri görselleştirin

## Troubleshooting

### Metrics endpoint'i çalışmıyor

**Sorun:** `/metrics` endpoint'i 404 döndürüyor

**Çözüm:** Panel.go'nun en son versiyonunu kullandığınızdan emin olun. Metrics middleware'i otomatik olarak eklenir.

### Audit log'lar oluşmuyor

**Sorun:** `audit_logs` tablosunda kayıt yok

**Çözüm:**
1. Veritabanı migration'ının çalıştığından emin olun
2. Sadece POST, PUT, PATCH, DELETE işlemleri loglanır
3. `/api/*` path'lerinde çalışır, diğer path'lerde çalışmaz

### Request log'ları görünmüyor

**Sorun:** Console'da request log'ları görünmüyor

**Çözüm:** Log'lar JSON formatında `stdout`'a yazılır. Terminal'inizin JSON log'ları desteklediğinden emin olun veya `jq` gibi bir tool kullanın:

```bash
./panel-go | jq
```

## İlgili Dökümanlar

- [Security](SECURITY.md) - Güvenlik özellikleri
- [API Reference](API-Reference.md) - API dokümantasyonu
- [Deployment](Deployment.md) - Production deployment
