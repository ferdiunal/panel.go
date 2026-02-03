# Alanlar (Fields)

Alanlar, Panel kaynaklarınızın yapı taşlarıdır. Hangi verilerin gösterileceğini ve nasıl giriş yapılacağını tanımlarlar. SDK, basit metin girişlerinden karmaşık ilişkilere kadar zengin bir alan tipi seti sunar.

## Standart Alanlar

Bu alanlar standart veri tiplerini işler.

-   **`fields.ID()`**: Birincil anahtarları (Primary Key) otomatik olarak işler.
-   **`fields.Text("Etiket", "key")`**: Standart metin girişi.
-   **`fields.Password("Etiket")`**: Parola girişi (listelerde varsayılan olarak gizlidir).
-   **`fields.Number("Etiket")`**: Sayısal giriş.
-   **`fields.Email("Etiket")`**: Doğrulama ipuçları içeren e-posta girişi.
-   **`fields.Tel("Etiket")`**: Telefon numarası girişi.
-   **`fields.Audio("Etiket")`** & **`fields.Video("Etiket")`**: Medya dosyası işleme.
-   **`fields.Date("Doğum Tarihi")`**: Tarih seçici.
-   **`fields.DateTime("Oluşturulma Tarihi")`**: Tarih ve saat seçici.
-   **`fields.File("Avatar")`**: Dosya yükleme işleyicisi.
-   **`fields.KeyValue("Metadata")`**: Anahtar-Değer çifti düzenleyicisi (JSON alanları).

## İlişkiler (Relationships)

İlişkileri daha anlamsal ve UI odaklı hale getirmek için **Yaratıcı İsimlendirme** kullanıyoruz.

### Link (`BelongsTo`)
Bir üst kaynağa bağlantı oluşturur. Mevcut kaynak başka bir kaynağa ait olduğunda bunu kullanın.

```go
fields.Link("Şirket", "company")
```
*UI: Tıklanabilir bir link veya seçim menüsü (select dropdown) olarak render edilir.*

### Collection (`HasMany`)
İlişkili öğelerin bir koleksiyonunu görüntüler. Mevcut kaynağın birçok alt öğesi olduğunda bunu kullanın.

```go
fields.Collection("Yazılar", "posts")
```
*UI: İlişkili öğelerin bir listesi veya tablosu olarak render edilir.*

### Detail (`HasOne`)
Mevcut kaynakla ilişkili tek bir detay kaydını gösterir.

```go
fields.Detail("Profil", "profile")
```
*UI: Gömülü bir detay görünümü veya detaya giden bir link olarak render edilir.*

### Connect (`BelongsToMany`)
Çoktan çoğa (many-to-many) bir bağlantıyı yönetir.

```go
fields.Connect("Roller", "roles")
```
*UI: Çoklu seçim (multi-select) veya etiket yöneticisi olarak render edilir.*

## Polimorfik İlişkiler

Polimorfik ilişkiler için `Poly` ön ekini kullanın:

-   **`fields.PolyLink("Yorumlanabilir")`** (`MorphTo`)
-   **`fields.PolyCollection("Yorumlar")`** (`MorphMany`)
-   **`fields.PolyDetail("Görsel")`** (`MorphOne`)
-   **`fields.PolyConnect("Etiketler")`** (`MorphToMany`)
