// # Sayfa Çözümleme ve Navigasyon Sistemi
//
// Bu paket, panel.go uygulamasında sayfaların dinamik olarak yapılandırılması,
// alanlarının ve kartlarının çözümlenmesi, ve navigasyon özelliklerinin yönetilmesi
// için gerekli interface'ler ve struct'ları sağlar.
//
// ## Temel Konseptler
//
// - **Resolver**: Sayfanın bileşenlerini (alanlar, kartlar) dinamik olarak oluşturan mekanizma
// - **Mixin**: Struct'lara ek işlevsellik eklemek için kullanılan kalıp
// - **OptimizedBase**: Tüm sayfa özelliklerini bir arada sunan temel struct
//
// ## Kullanım Örneği
//
// ```go
// // Özel bir sayfa oluşturma
// type UserPage struct {
//     page.OptimizedBase
// }
//
// // Field resolver'ı ayarlama
// userPage := &UserPage{}
// userPage.SetFieldResolver(&MyFieldResolver{})
// userPage.SetCardResolver(&MyCardResolver{})
// userPage.SetSlug("users")
// userPage.SetTitle("Kullanıcılar")
// userPage.SetIcon("users")
// ```
package page

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// FieldResolver, sayfanın form alanlarını dinamik olarak çözen (resolve) interface'i.
//
// # Açıklama
//
// FieldResolver, bir sayfanın hangi form alanlarını göstereceğini belirlemek için
// kullanılır. Bu interface'i implement eden struct'lar, runtime'da alanları dinamik
// olarak oluşturabilir ve koşullu olarak gösterebilir.
//
// # Kullanım Senaryoları
//
// - Kullanıcı rolüne göre farklı alanlar gösterme
// - Veritabanı verilerine göre dinamik form oluşturma
// - Koşullu alan görünürlüğü (bir alanın değerine göre başka alanları gösterme)
// - Çok dilli form desteği
//
// # Parametreler
//
// - `ctx`: İstek bağlamı, kullanıcı bilgisi, veritabanı bağlantısı vb. içerir
//
// # Dönüş Değeri
//
// Sayfada gösterilecek form alanlarının slice'ı. Boş slice dönüldüğünde hiçbir alan gösterilmez.
//
// # Örnek Implementasyon
//
// ```go
// type UserFieldResolver struct {
//     db *gorm.DB
// }
//
// func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
//     fields := []fields.Element{
//         fields.NewText("name", "Ad Soyad"),
//         fields.NewEmail("email", "E-posta"),
//     }
//
//     // Admin ise ek alanlar ekle
//     if ctx.User.IsAdmin {
//         fields = append(fields, fields.NewSelect("role", "Rol"))
//     }
//
//     return fields
// }
// ```
//
// # Avantajlar
//
// - Dinamik form yapılandırması
// - Koşullu alan gösterimi
// - Yeniden kullanılabilir resolver'lar
// - Test edilebilir yapı
//
// # Önemli Notlar
//
// - Resolver nil olabilir, bu durumda boş slice döndürülür
// - Context nil olabilir, bu durumda varsayılan alanlar döndürülmelidir
// - Alanlar sırasıyla gösterilir, sıra önemlidir
type FieldResolver interface {
	ResolveFields(ctx *context.Context) []fields.Element
}

// CardResolver, sayfanın dashboard kartlarını dinamik olarak çözen (resolve) interface'i.
//
// # Açıklama
//
// CardResolver, bir sayfanın dashboard'unda gösterilecek kartları belirlemek için
// kullanılır. Kartlar istatistikler, grafikler, özet bilgiler vb. gösterebilir.
//
// # Kullanım Senaryoları
//
// - Dashboard istatistiklerini gösterme
// - Gerçek zamanlı veriler içeren kartlar
// - Kullanıcı izinlerine göre farklı kartlar
// - Özel raporlar ve analizler
//
// # Parametreler
//
// - `ctx`: İstek bağlamı, kullanıcı bilgisi, veritabanı bağlantısı vb. içerir
//
// # Dönüş Değeri
//
// Sayfada gösterilecek kartların slice'ı. Boş slice dönüldüğünde hiçbir kart gösterilmez.
//
// # Örnek Implementasyon
//
// ```go
// type DashboardCardResolver struct {
//     db *gorm.DB
// }
//
// func (r *DashboardCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
//     var userCount int64
//     r.db.Model(&User{}).Count(&userCount)
//
//     return []widget.Card{
//         {
//             Title: "Toplam Kullanıcı",
//             Value: fmt.Sprintf("%d", userCount),
//             Icon:  "users",
//         },
//     }
// }
// ```
//
// # Avantajlar
//
// - Dinamik dashboard yapılandırması
// - Gerçek zamanlı veri gösterimi
// - Koşullu kart gösterimi
// - Yeniden kullanılabilir resolver'lar
//
// # Önemli Notlar
//
// - Resolver nil olabilir, bu durumda boş slice döndürülür
// - Kartlar sırasıyla gösterilir, sıra önemlidir
// - Ağır hesaplamalar için caching kullanılması önerilir
type CardResolver interface {
	ResolveCards(ctx *context.Context) []widget.Card
}

// Resolvable, sayfa çözümleme işlevselliğini sağlayan mixin struct'ı.
//
// # Açıklama
//
// Resolvable, bir struct'a form alanları ve dashboard kartları çözümleme
// (resolve) yeteneği eklemek için tasarlanmış bir mixin'dir. Embedding yoluyla
// kullanılarak, herhangi bir struct'a dinamik içerik oluşturma özelliği
// kazandırılabilir.
//
// # Kullanım Senaryoları
//
// - Sayfa struct'larına resolver'lar eklemek
// - Dinamik form ve dashboard yapılandırması
// - Koşullu içerik gösterimi
// - Yeniden kullanılabilir sayfa bileşenleri
//
// # Avantajlar
//
// - Mixin kalıbı sayesinde esnek yapı
// - Nil-safe implementasyon (resolver nil olabilir)
// - Boş slice döndürerek güvenli varsayılan davranış
// - Kolay test edilebilirlik
//
// # Dezavantajlar
//
// - Nil resolver'lar sessizce göz ardı edilir (hata fırlatmaz)
// - Resolver'lar runtime'da değiştirilirse, önceki çağrılar etkilenmez
//
// # Örnek Kullanım
//
// ```go
// type UserPage struct {
//     page.Resolvable
//     db *gorm.DB
// }
//
// userPage := &UserPage{db: db}
// userPage.SetFieldResolver(&UserFieldResolver{db: db})
// userPage.SetCardResolver(&UserCardResolver{db: db})
//
// fields := userPage.ResolveFields(ctx)
// cards := userPage.ResolveCards(ctx)
// ```
//
// # Önemli Notlar
//
// - Resolver'lar nil olabilir, bu durumda boş slice döndürülür
// - Context nil olabilir, resolver'lar bunu kontrol etmelidir
// - Resolver'lar thread-safe olmalıdır
type Resolvable struct {
	fieldResolver FieldResolver
	cardResolver  CardResolver
}

// SetFieldResolver, form alanlarını çözmek için kullanılacak resolver'ı ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın hangi form alanlarını göstereceğini belirlemek için
// kullanılacak FieldResolver'ı ayarlar. Resolver nil olabilir, bu durumda
// ResolveFields boş slice döndürür.
//
// # Parametreler
//
// - `fr`: FieldResolver interface'ini implement eden struct. Nil olabilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// resolvable := &page.Resolvable{}
// resolvable.SetFieldResolver(&MyFieldResolver{})
// ```
//
// # Önemli Notlar
//
// - Resolver nil olarak ayarlanabilir
// - Aynı resolver birden fazla Resolvable'a atanabilir
// - Resolver'ı değiştirmek önceki çağrıları etkilemez
func (r *Resolvable) SetFieldResolver(fr FieldResolver) {
	r.fieldResolver = fr
}

// SetCardResolver, dashboard kartlarını çözmek için kullanılacak resolver'ı ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın dashboard'unda hangi kartları göstereceğini belirlemek için
// kullanılacak CardResolver'ı ayarlar. Resolver nil olabilir, bu durumda
// ResolveCards boş slice döndürür.
//
// # Parametreler
//
// - `cr`: CardResolver interface'ini implement eden struct. Nil olabilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// resolvable := &page.Resolvable{}
// resolvable.SetCardResolver(&MyCardResolver{})
// ```
//
// # Önemli Notlar
//
// - Resolver nil olarak ayarlanabilir
// - Aynı resolver birden fazla Resolvable'a atanabilir
// - Resolver'ı değiştirmek önceki çağrıları etkilemez
func (r *Resolvable) SetCardResolver(cr CardResolver) {
	r.cardResolver = cr
}

// ResolveFields, ayarlanmış FieldResolver'ı kullanarak form alanlarını çözer.
//
// # Açıklama
//
// Bu metod, ayarlanmış FieldResolver'ı çağırarak sayfada gösterilecek form
// alanlarını dinamik olarak oluşturur. Resolver nil ise boş slice döndürür.
//
// # Parametreler
//
// - `ctx`: İstek bağlamı. Resolver'a iletilir, nil olabilir.
//
// # Dönüş Değeri
//
// Sayfada gösterilecek form alanlarının slice'ı. Resolver nil ise boş slice.
//
// # Örnek
//
// ```go
// ctx := &context.Context{User: user}
// fields := resolvable.ResolveFields(ctx)
// for _, field := range fields {
//     fmt.Println(field.Name())
// }
// ```
//
// # Avantajlar
//
// - Nil-safe implementasyon
// - Dinamik alan oluşturma
// - Context'e erişim
//
// # Önemli Notlar
//
// - Resolver nil ise boş slice döndürülür (hata fırlatmaz)
// - Context nil olabilir, resolver'lar bunu kontrol etmelidir
// - Alanlar sırasıyla gösterilir
func (r *Resolvable) ResolveFields(ctx *context.Context) []fields.Element {
	if r.fieldResolver != nil {
		return r.fieldResolver.ResolveFields(ctx)
	}
	return []fields.Element{}
}

// ResolveCards, ayarlanmış CardResolver'ı kullanarak dashboard kartlarını çözer.
//
// # Açıklama
//
// Bu metod, ayarlanmış CardResolver'ı çağırarak sayfanın dashboard'unda
// gösterilecek kartları dinamik olarak oluşturur. Resolver nil ise boş slice döndürür.
//
// # Parametreler
//
// - `ctx`: İstek bağlamı. Resolver'a iletilir, nil olabilir.
//
// # Dönüş Değeri
//
// Sayfada gösterilecek kartların slice'ı. Resolver nil ise boş slice.
//
// # Örnek
//
// ```go
// ctx := &context.Context{User: user}
// cards := resolvable.ResolveCards(ctx)
// for _, card := range cards {
//     fmt.Println(card.Title)
// }
// ```
//
// # Avantajlar
//
// - Nil-safe implementasyon
// - Dinamik kart oluşturma
// - Context'e erişim
//
// # Önemli Notlar
//
// - Resolver nil ise boş slice döndürülür (hata fırlatmaz)
// - Context nil olabilir, resolver'lar bunu kontrol etmelidir
// - Kartlar sırasıyla gösterilir
func (r *Resolvable) ResolveCards(ctx *context.Context) []widget.Card {
	if r.cardResolver != nil {
		return r.cardResolver.ResolveCards(ctx)
	}
	return []widget.Card{}
}

// Navigable, sayfa navigasyon ve menü işlevselliğini sağlayan mixin struct'ı.
//
// # Açıklama
//
// Navigable, bir struct'a navigasyon menüsünde görünmesi için gerekli özellikleri
// eklemek için tasarlanmış bir mixin'dir. Sayfaların menüde nasıl gösterileceğini,
// hangi grupta yer alacağını, sırasını ve görünürlüğünü kontrol eder.
//
// # Kullanım Senaryoları
//
// - Menü öğelerinin sıralanması
// - Sayfaları gruplar halinde organize etme
// - Dinamik menü görünürlüğü (rol bazlı)
// - İkon ve başlık gösterimi
// - Menü hiyerarşisi oluşturma
//
// # Avantajlar
//
// - Mixin kalıbı sayesinde esnek yapı
// - Menü yapılandırması için merkezi yer
// - Kolay sıralama ve gruplama
// - Dinamik görünürlük kontrolü
//
// # Dezavantajlar
//
// - Sıralama için integer kullanımı (string daha esnektir)
// - Grup adı string olduğu için yazım hataları mümkün
//
// # Örnek Kullanım
//
// ```go
// type UserPage struct {
//     page.Navigable
// }
//
// userPage := &UserPage{}
// userPage.SetIcon("users")
// userPage.SetGroup("Management")
// userPage.SetNavigationOrder(1)
// userPage.SetVisible(true)
// ```
//
// # Önemli Notlar
//
// - Sıra küçük sayılar için daha üst konumda gösterilir
// - Aynı grupta sayfalar sıra numarasına göre sıralanır
// - Görünürlük false ise sayfa menüde gösterilmez
// - İkon adı UI framework'ün desteklediği bir ikon olmalıdır
type Navigable struct {
	icon            string
	group           string
	navigationOrder int
	visible         bool
}

// SetIcon, sayfanın menüde gösterilecek ikon adını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın navigasyon menüsünde gösterilecek ikon adını belirler.
// İkon adı, UI framework'ün desteklediği bir ikon olmalıdır (örn: "users", "settings", "dashboard").
//
// # Parametreler
//
// - `icon`: Ikon adı. Boş string olabilir (ikon gösterilmez).
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetIcon("users")
// navigable.SetIcon("settings")
// navigable.SetIcon("") // İkon gösterilmez
// ```
//
// # Önemli Notlar
//
// - İkon adı UI framework'ün desteklediği bir ikon olmalıdır
// - Boş string geçilirse ikon gösterilmez
// - İkon adı case-sensitive olabilir
func (n *Navigable) SetIcon(icon string) {
	n.icon = icon
}

// GetIcon, sayfanın menüde gösterilecek ikon adını döner.
//
// # Açıklama
//
// Bu metod, SetIcon ile ayarlanmış ikon adını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Ikon adı. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetIcon("users")
// icon := navigable.GetIcon() // "users"
// ```
//
// # Önemli Notlar
//
// - Boş string dönüldüğünde ikon gösterilmez
func (n *Navigable) GetIcon() string {
	return n.icon
}

// SetGroup, sayfanın menüde yer alacağı grup adını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın navigasyon menüsünde hangi grup altında gösterileceğini
// belirler. Aynı grup adına sahip sayfalar menüde birlikte gösterilir.
//
// # Parametreler
//
// - `group`: Grup adı. Boş string olabilir (grup gösterilmez).
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetGroup("Management")
// navigable.SetGroup("Settings")
// navigable.SetGroup("") // Grup gösterilmez
// ```
//
// # Kullanım Senaryoları
//
// - Yönetim sayfalarını "Management" grubu altında gösterme
// - Ayarlar sayfalarını "Settings" grubu altında gösterme
// - Raporları "Reports" grubu altında gösterme
//
// # Önemli Notlar
//
// - Grup adı yazım hataları menüde yanlış gruplama oluşturabilir
// - Boş string geçilirse sayfa grup olmadan gösterilir
// - Grup adı case-sensitive olabilir
func (n *Navigable) SetGroup(group string) {
	n.group = group
}

// GetGroup, sayfanın menüde yer alacağı grup adını döner.
//
// # Açıklama
//
// Bu metod, SetGroup ile ayarlanmış grup adını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Grup adı. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetGroup("Management")
// group := navigable.GetGroup() // "Management"
// ```
//
// # Önemli Notlar
//
// - Boş string dönüldüğünde sayfa grup olmadan gösterilir
func (n *Navigable) GetGroup() string {
	return n.group
}

// SetNavigationOrder, sayfanın menüde gösterilme sırasını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın navigasyon menüsünde hangi sırada gösterileceğini belirler.
// Küçük sayılar daha üst konumda gösterilir. Aynı grupta sayfalar bu sıraya göre
// sıralanır.
//
// # Parametreler
//
// - `order`: Sıra numarası. Küçük sayılar daha üst konumda gösterilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetNavigationOrder(1)  // En üstte
// navigable.SetNavigationOrder(10) // Daha aşağıda
// navigable.SetNavigationOrder(100) // En altta
// ```
//
// # Kullanım Senaryoları
//
// - Önemli sayfaları menünün üstüne koymak
// - Sayfaları mantıksal sıraya göre düzenlemek
// - Dinamik sıralama (rol bazlı)
//
// # Önemli Notlar
//
// - Küçük sayılar daha üst konumda gösterilir
// - Aynı sıraya sahip sayfalar alfabetik sıraya göre sıralanabilir
// - Negatif sayılar kullanılabilir (daha üst konumda gösterilir)
// - Varsayılan değer 0'dır
func (n *Navigable) SetNavigationOrder(order int) {
	n.navigationOrder = order
}

// GetNavigationOrder, sayfanın menüde gösterilme sırasını döner.
//
// # Açıklama
//
// Bu metod, SetNavigationOrder ile ayarlanmış sıra numarasını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Sıra numarası. Ayarlanmamışsa 0.
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetNavigationOrder(5)
// order := navigable.GetNavigationOrder() // 5
// ```
//
// # Önemli Notlar
//
// - Küçük sayılar daha üst konumda gösterilir
// - Varsayılan değer 0'dır
func (n *Navigable) GetNavigationOrder() int {
	return n.navigationOrder
}

// SetVisible, sayfanın menüde görünüp görünmeyeceğini ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın navigasyon menüsünde görünüp görünmeyeceğini kontrol eder.
// False olarak ayarlanırsa sayfa menüde gösterilmez.
//
// # Parametreler
//
// - `visible`: true ise sayfa menüde gösterilir, false ise gösterilmez.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetVisible(true)  // Menüde göster
// navigable.SetVisible(false) // Menüde gösterme
// ```
//
// # Kullanım Senaryoları
//
// - Rol bazlı sayfa görünürlüğü
// - Geliştirme aşamasındaki sayfaları gizlemek
// - Dinamik menü yapılandırması
// - Koşullu sayfa erişimi
//
// # Önemli Notlar
//
// - Varsayılan değer false'dur (görünmez)
// - Görünmez sayfalar yine de erişilebilir olabilir (URL'den)
// - Görünürlük kontrol etmek güvenlik sağlamaz, sadece UI'da gizler
func (n *Navigable) SetVisible(visible bool) {
	n.visible = visible
}

// IsVisible, sayfanın menüde görünüp görünmeyeceğini döner.
//
// # Açıklama
//
// Bu metod, SetVisible ile ayarlanmış görünürlük durumunu döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// true ise sayfa menüde gösterilir, false ise gösterilmez.
//
// # Örnek
//
// ```go
// navigable := &page.Navigable{}
// navigable.SetVisible(true)
// visible := navigable.IsVisible() // true
// ```
//
// # Önemli Notlar
//
// - Varsayılan değer false'dur
// - Görünmez sayfalar yine de erişilebilir olabilir
func (n *Navigable) IsVisible() bool {
	return n.visible
}

// OptimizedBase, Page interface'ini implement eden ve tüm sayfa işlevselliğini
// bir arada sunan temel struct'ı.
//
// # Açıklama
//
// OptimizedBase, Resolvable ve Navigable mixin'lerini embed ederek, bir sayfanın
// ihtiyaç duyduğu tüm özellikleri (slug, başlık, açıklama, alanlar, kartlar,
// navigasyon) bir arada sağlar. Bu struct, diğer sayfa struct'larının temel
// olarak kullanılması için optimize edilmiştir.
//
// # Yapı
//
// OptimizedBase aşağıdaki bileşenleri içerir:
// - Resolvable: Alan ve kart çözümleme işlevselliği
// - Navigable: Navigasyon menüsü işlevselliği
// - slug: Sayfanın URL slug'ı
// - title: Sayfanın başlığı
// - description: Sayfanın açıklaması
//
// # Kullanım Senaryoları
//
// - Sayfa struct'larının temel sınıfı olarak kullanma
// - Embedding yoluyla sayfa özelliklerini miras alma
// - Tüm sayfa özelliklerini merkezi bir yerden yönetme
// - Sayfa factory'lerinde kullanma
//
// # Avantajlar
//
// - Tüm sayfa işlevselliğini bir arada sunar
// - Embedding yoluyla esnek yapı
// - Nil-safe implementasyon
// - Kolay genişletme
//
// # Dezavantajlar
//
// - Tüm özellikleri içerdiği için biraz ağır olabilir
// - Bazı sayfalar tüm özelliklere ihtiyaç duymayabilir
//
// # Örnek Kullanım
//
// ```go
// type UserPage struct {
//     page.OptimizedBase
//     db *gorm.DB
// }
//
// userPage := &UserPage{db: db}
// userPage.SetSlug("users")
// userPage.SetTitle("Kullanıcılar")
// userPage.SetDescription("Kullanıcı yönetimi sayfası")
// userPage.SetIcon("users")
// userPage.SetGroup("Management")
// userPage.SetNavigationOrder(1)
// userPage.SetVisible(true)
// userPage.SetFieldResolver(&UserFieldResolver{db: db})
// userPage.SetCardResolver(&UserCardResolver{db: db})
// ```
//
// # Önemli Notlar
//
// - Embedding yoluyla kullanılması önerilir
// - Fields() ve Cards() metodları nil context ile çağrılır
// - Tüm getter metodları wrapper'dır, setter'lar mixin'lerde tanımlıdır
type OptimizedBase struct {
	Resolvable
	Navigable
	slug        string
	title       string
	description string
}

// SetSlug, sayfanın URL slug'ını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın URL'de kullanılacak slug'ını belirler. Slug, sayfaya
// erişmek için kullanılan URL yoludur (örn: "/users", "/settings").
//
// # Parametreler
//
// - `s`: Slug değeri. Boş string olabilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetSlug("users")
// base.SetSlug("user-management")
// base.SetSlug("") // Slug olmaz
// ```
//
// # Kullanım Senaryoları
//
// - Sayfa URL'lerini tanımlama
// - Navigasyon linklerini oluşturma
// - Sayfa yönlendirmesi
//
// # Önemli Notlar
//
// - Slug URL-safe olmalıdır (boşluk, özel karakterler içermemelidir)
// - Slug benzersiz olmalıdır (aynı slug'a sahip iki sayfa olmamalıdır)
// - Slug genellikle küçük harfler ve tire (-) içerir
func (b *OptimizedBase) SetSlug(s string) {
	b.slug = s
}

// Slug, sayfanın URL slug'ını döner.
//
// # Açıklama
//
// Bu metod, SetSlug ile ayarlanmış slug'ı döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Slug değeri. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetSlug("users")
// slug := base.Slug() // "users"
// ```
//
// # Önemli Notlar
//
// - Boş string dönüldüğünde slug tanımlı değildir
func (b *OptimizedBase) Slug() string {
	return b.slug
}

// SetTitle, sayfanın başlığını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın başlığını belirler. Başlık, sayfanın tarayıcı sekmesinde,
// menüde ve sayfa başında gösterilir.
//
// # Parametreler
//
// - `t`: Başlık değeri. Boş string olabilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetTitle("Kullanıcılar")
// base.SetTitle("Kullanıcı Yönetimi")
// base.SetTitle("") // Başlık olmaz
// ```
//
// # Kullanım Senaryoları
//
// - Sayfa başlığını tanımlama
// - Tarayıcı sekmesi başlığını ayarlama
// - Menü öğesi başlığını gösterme
//
// # Önemli Notlar
//
// - Başlık kullanıcı dostu olmalıdır
// - Başlık genellikle Türkçe yazılır
// - Başlık kısa ve açıklayıcı olmalıdır
func (b *OptimizedBase) SetTitle(t string) {
	b.title = t
}

// Title, sayfanın başlığını döner.
//
// # Açıklama
//
// Bu metod, SetTitle ile ayarlanmış başlığı döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Başlık değeri. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetTitle("Kullanıcılar")
// title := base.Title() // "Kullanıcılar"
// ```
//
// # Önemli Notlar
//
// - Boş string dönüldüğünde başlık tanımlı değildir
func (b *OptimizedBase) Title() string {
	return b.title
}

// SetDescription, sayfanın açıklamasını ayarlar.
//
// # Açıklama
//
// Bu metod, sayfanın açıklamasını belirler. Açıklama, sayfanın amacını ve
// işlevini açıklar. Genellikle SEO meta açıklaması olarak kullanılır.
//
// # Parametreler
//
// - `d`: Açıklama değeri. Boş string olabilir.
//
// # Dönüş Değeri
//
// Yok (void)
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetDescription("Sistem kullanıcılarını yönetin")
// base.SetDescription("Kullanıcı oluşturma, düzenleme ve silme işlemleri")
// base.SetDescription("") // Açıklama olmaz
// ```
//
// # Kullanım Senaryoları
//
// - Sayfa açıklamasını tanımlama
// - SEO meta açıklaması ayarlama
// - Sayfa hakkında bilgi sağlama
//
// # Önemli Notlar
//
// - Açıklama kısa ve açıklayıcı olmalıdır
// - Açıklama genellikle 150-160 karakter olmalıdır (SEO için)
// - Açıklama kullanıcı dostu olmalıdır
func (b *OptimizedBase) SetDescription(d string) {
	b.description = d
}

// Description, sayfanın açıklamasını döner.
//
// # Açıklama
//
// Bu metod, SetDescription ile ayarlanmış açıklamayı döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Açıklama değeri. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetDescription("Sistem kullanıcılarını yönetin")
// desc := base.Description() // "Sistem kullanıcılarını yönetin"
// ```
//
// # Önemli Notlar
//
// - Boş string dönüldüğünde açıklama tanımlı değildir
func (b *OptimizedBase) Description() string {
	return b.description
}

// Fields, sayfanın form alanlarını döner.
//
// # Açıklama
//
// Bu metod, ResolveFields'i nil context ile çağırarak sayfanın form alanlarını
// döner. Bu, basit kullanım için bir wrapper'dır.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Sayfanın form alanlarının slice'ı. Resolver nil ise boş slice.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetFieldResolver(&MyFieldResolver{})
// fields := base.Fields()
// for _, field := range fields {
//     fmt.Println(field.Name())
// }
// ```
//
// # Önemli Notlar
//
// - Context nil olarak iletilir
// - Resolver nil ise boş slice döndürülür
// - ResolveFields'in wrapper'ıdır
//
// # Uyarı
//
// Context nil olduğu için, resolver'lar context'e bağlı işlemler yapamaz.
// Daha kontrollü kullanım için ResolveFields'i doğrudan çağırın.
func (b *OptimizedBase) Fields() []fields.Element {
	return b.ResolveFields(nil)
}

// Cards, sayfanın dashboard kartlarını döner.
//
// # Açıklama
//
// Bu metod, ResolveCards'i nil context ile çağırarak sayfanın dashboard
// kartlarını döner. Bu, basit kullanım için bir wrapper'dır.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Sayfanın dashboard kartlarının slice'ı. Resolver nil ise boş slice.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetCardResolver(&MyCardResolver{})
// cards := base.Cards()
// for _, card := range cards {
//     fmt.Println(card.Title)
// }
// ```
//
// # Önemli Notlar
//
// - Context nil olarak iletilir
// - Resolver nil ise boş slice döndürülür
// - ResolveCards'in wrapper'ıdır
//
// # Uyarı
//
// Context nil olduğu için, resolver'lar context'e bağlı işlemler yapamaz.
// Daha kontrollü kullanım için ResolveCards'i doğrudan çağırın.
func (b *OptimizedBase) Cards() []widget.Card {
	return b.ResolveCards(nil)
}

// Icon, sayfanın menüde gösterilecek ikon adını döner.
//
// # Açıklama
//
// Bu metod, GetIcon'un wrapper'ıdır. Sayfanın menüde gösterilecek ikon adını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Ikon adı. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetIcon("users")
// icon := base.Icon() // "users"
// ```
//
// # Önemli Notlar
//
// - GetIcon'un wrapper'ıdır
// - Boş string dönüldüğünde ikon gösterilmez
func (b *OptimizedBase) Icon() string {
	return b.GetIcon()
}

// Group, sayfanın menüde yer alacağı grup adını döner.
//
// # Açıklama
//
// Bu metod, GetGroup'un wrapper'ıdır. Sayfanın menüde yer alacağı grup adını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Grup adı. Ayarlanmamışsa boş string.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetGroup("Management")
// group := base.Group() // "Management"
// ```
//
// # Önemli Notlar
//
// - GetGroup'un wrapper'ıdır
// - Boş string dönüldüğünde sayfa grup olmadan gösterilir
func (b *OptimizedBase) Group() string {
	return b.GetGroup()
}

// NavigationOrder, sayfanın menüde gösterilme sırasını döner.
//
// # Açıklama
//
// Bu metod, GetNavigationOrder'ın wrapper'ıdır. Sayfanın menüde gösterilme
// sırasını döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// Sıra numarası. Ayarlanmamışsa 0.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetNavigationOrder(5)
// order := base.NavigationOrder() // 5
// ```
//
// # Önemli Notlar
//
// - GetNavigationOrder'ın wrapper'ıdır
// - Küçük sayılar daha üst konumda gösterilir
// - Varsayılan değer 0'dır
func (b *OptimizedBase) NavigationOrder() int {
	return b.GetNavigationOrder()
}

// Visible, sayfanın menüde görünüp görünmeyeceğini döner.
//
// # Açıklama
//
// Bu metod, IsVisible'ın wrapper'ıdır. Sayfanın menüde görünüp görünmeyeceğini döner.
//
// # Parametreler
//
// Yok
//
// # Dönüş Değeri
//
// true ise sayfa menüde gösterilir, false ise gösterilmez.
//
// # Örnek
//
// ```go
// base := &page.OptimizedBase{}
// base.SetVisible(true)
// visible := base.Visible() // true
// ```
//
// # Önemli Notlar
//
// - IsVisible'ın wrapper'ıdır
// - Varsayılan değer false'dur
// - Görünmez sayfalar yine de erişilebilir olabilir
func (b *OptimizedBase) Visible() bool {
	return b.IsVisible()
}
