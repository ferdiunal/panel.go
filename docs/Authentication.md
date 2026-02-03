# Authentication System Documentation

## Genel Bakış
Panel auth sistemi, modüler ve genişletilebilir bir yapıya sahiptir. `better-auth` kütüphanesinden esinlenerek tasarlanmıştır ve güvenlik (Bcrypt, UUID v7) standartlarına uyar.

## Sıkça Sorulan Sorular

### 1. Kullanıcı kendi auth handle'ını ekleyebilir mi?
**Evet.** Auth sistemi modülerdir. Kullanıcılar şu yöntemlerle sisteme müdahale edebilir:
- **Middleware**: Fiber middleware'leri ile route koruması yapabilirler.
- **Provider Pattern**: `Account` entity'si `ProviderID` alanına sahiptir. Bu, sadece "credential" (e-posta/şifre) değil, gelecekte `google`, `github` gibi farklı sağlayıcıların da eklenmesine olanak tanır.
- **Service Wrapping**: `AuthService` struct'ı, `Repository` arayüzlerine bağlıdır. Kullanıcılar bu repoları mocklayabilir veya kendi implementasyonlarını enjekte edebilirler.

### 2. Plugin mantığı var mı? OTP ile giriş, OAuth desteği?
**Evet, yapı buna uygundur.** 
- **OAuth**: `Account` tablosundaki `ProviderID`, `AccessToken`, `RefreshToken` alanları OAuth entegrasyonu için hazırdır. Yeni bir route (örn: `/auth/sign-in/google`) ekleyerek ve `AuthService` içinde ilgili provider doğrulamasını yaparak sisteme entegre edilebilir.
- **OTP/2FA**: `Verification` domain'i bu amaçla oluşturulmuştur. OTP kodları `token` olarak saklanabilir ve doğrulama sonrası oturum açılabilir.

### 3. Register, Şifremi unuttum var mı? E-posta doğrulaması dahil.
**Temel yapı hazırdır.**
- **Register**: `/api/auth/sign-up/email` endpointi mevcuttur.
- **E-posta Doğrulama**: `User` entity'sinde `EmailVerified` alanı vardır. `Verification` domain'i kullanılarak kayıt sonrası token oluşturulup e-posta ile gönderilebilir.
- **Şifremi Unuttum**: Henüz API endpointi eklenmemiştir ancak `Verification` domain'i üzerinden şifre sıfırlama token'ı oluşturulup `User` şifresini güncelleme akışı kolayca eklenebilir.

### 4. UUID v7 kullanımı
**Evet.** Tüm sistem (User, Session, Account, Verification) artık **UUID v7** standardını kullanmaktadır. Bu, zaman bazlı sıralanabilirlik ve benzersizlik sağlar.

## Mimari

### Domainler
- **User**: Kullanıcı temel bilgileri (`id`, `name`, `email`, `emailVerified`).
- **Session**: Oturum yönetimi (`token`, `expiresAt`, `userAgent`).
- **Account**: Giriş yöntemleri (`providerId`, `password` hash).
- **Verification**: Geçici tokenlar (OTP, Email Verify, Password Reset).

### API Endpointleri
- `POST /api/auth/sign-in/email`: E-posta/Şifre ile giriş.
- `POST /api/auth/sign-up/email`: Yeni üye kaydı.
- `POST /api/auth/sign-out`: Çıkış yap.
- `GET /api/auth/session`: Mevcut oturum bilgisini getir.
