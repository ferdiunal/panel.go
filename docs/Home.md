# Panel SDK Dokümantasyonu

**Panel SDK**'ya hoş geldiniz. Go ekosisteminin gücünü ve tip sistemini sonuna kadar kullanan, modern, tip güvenli (type-safe) ve akıcı (fluent) bir UI aracıdır. Bu SDK, güçlü yönetim panelleri ve dashboard'ları kolaylıkla oluşturmanız için tasarlanmıştır.

## Felsefemiz

Panel SDK birkaç temel prensip üzerine kuruludur:

1.  **Tip Güvenliği (Type Safety)**: Alan tanımlarından veri sorgularına kadar her şey tip güvenlidir. Sihirli string'lere veya tahmin oyunlarına yer yok.
2.  **Akıcı API (Fluent API)**: Kaynaklarınızı ve alanlarınızı, doğal bir dil gibi okunan temiz, zincirleme (chainable) bir API kullanarak tanımlayın.
3.  **Performans**: Eager Loading (`With`) gibi yerleşik optimizasyonlar ve verimli JSON serileştirme sayesinde uygulamalarınız hızlı çalışır.
4.  **Basitlik**: Sihir yerine netliği tercih ediyoruz. Açık ilişkiler, anlaşılır isimlendirme ve öngörülebilir davranışlar.

## Başlarken

Panel SDK'yı kullanmak için genellikle **Kaynaklar (Resources)** tanımlarsınız. Bir Kaynak, bir veri modeline (örneğin bir Kullanıcı veya Ürün) karşılık gelir ve bunun nasıl görüntüleneceğini ve etkileşime girileceğini tanımlar.

### Örnek Kaynak (Resource)

```go
type UserResource struct{}

func (u *UserResource) Model() interface{} {
    return &User{}
}

func (u *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("İsim", "Name").Sortable().Required(),
        fields.Email("E-posta", "Email"),
        fields.Link("Şirket", "Company"), // İlişki
    }
}
```

## Dokümantasyon

### Başlangıç

-   **[Başlarken (Getting Started)](Getting-Started)** - Panel SDK'yı kurun ve ilk kaynağınızı oluşturun
-   **[Kaynaklar (Resources)](Resources)** - Kaynakları tanımlayın ve yapılandırın

### Temel Kavramlar

-   **[Alanlar (Fields)](Fields)** - Tüm alan türleri ve seçenekleri
-   **[İlişkiler (Relationships)](Relationships)** - BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo
-   **[Yetkilendirme (Authorization)](Authorization)** - Erişim kontrolü ve politikalar

### İleri Seviye

-   **[Gelişmiş Kullanım (Advanced Usage)](Advanced-Usage)** - Özel alanlar, middleware, hooks
-   **[API Referansı (API Reference)](API-Reference)** - Tüm metodlar ve parametreler
-   **[Lensler (Lenses)](Lenses)** - Özel raporlar ve görünümler
-   **[Sayfalar (Pages)](Pages)** - Özel gösterge panelleri
-   **[Ayarlar (Settings)](Settings)** - Uygulama ayarları
-   **[Widgets](Widgets)** - Gösterge paneli widget'ları
