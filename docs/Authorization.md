# Yetkilendirme ve Koşullu Görünürlük (Authorization & Conditionals)

Panel SDK, kaynaklara erişimi kontrol etmek ve alanların (fields) görünürlüğünü dinamik olarak yönetmek için güçlü bir yapı sunar.

## Politikalar (Policies)

Politikalar, bir kullanıcının belirli bir kaynak üzerinde işlem yapıp yapamayacağını belirleyen mantığı kapsüller.

### Policy Arayüzü

Herhangi bir yetkilendirme mantığı `auth.Policy` arayüzünü uygulamalıdır:

```go
type Policy interface {
    ViewAny(ctx context.Context) bool             // Liste sayfasına erişim
    View(ctx context.Context, model interface{}) bool // Detay sayfasına erişim
    Create(ctx context.Context) bool              // Yeni oluşturma izni
    Update(ctx context.Context, model interface{}) bool // Güncelleme izni
    Delete(ctx context.Context, model interface{}) bool // Silme izni
}
```

### Policy Oluşturma

Örneğin, sadece adminlerin blog yazılarını silebileceği bir senaryo:

```go
type BlogPolicy struct {}

func (p *BlogPolicy) ViewAny(ctx context.Context) bool {
    return true // Herkes görebilir
}

func (p *BlogPolicy) View(ctx context.Context, model interface{}) bool {
    return true
}

func (p *BlogPolicy) Create(ctx context.Context) bool {
    user := ctx.Value("user").(*User)
    return user.IsAdmin
}

func (p *BlogPolicy) Update(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    blog := model.(*Blog)
    return user.IsAdmin || blog.UserID == user.ID
}

func (p *BlogPolicy) Delete(ctx context.Context, model interface{}) bool {
    user := ctx.Value("user").(*User)
    return user.IsAdmin
}
```

### Kaynağa Tanımlama

Politikayı kaynağınıza (Resource) `Policy()` metodu ile eklersiniz:

```go
func (r *BlogResource) Policy() auth.Policy {
    return &BlogPolicy{}
}
```

Bu tanımlandığında, Panel SDK tüm CRUD işlemlerinde (`Index`, `Show`, `Store`, `Update`, `Destroy`) otomatik olarak bu metodları kontrol eder. Erişim izni yoksa `403 Forbidden` döner.

## Koşullu Alanlar (Conditional Fields)

Bazen bir alanın sadece belirli durumlarda (örneğin sadece güncelleme formunda veya sadece adminler için) görünmesini istersiniz. Bunun için `CanSee` metodunu kullanabilirsiniz.

### Kullanım

```go
fields.Text("API Key", "api_key").
    CanSee(func(ctx context.Context) bool {
        user := ctx.Value("user").(User)
        return user.IsAdmin
    }),
```

Eğer `CanSee` false dönerse, alan API yanıtından tamamen çıkarılır ve frontend'de görüntülenmez. Ayrıca `Store` ve `Update` işlemlerinde de bu alanın değeri işlenmez.

### Hazır Kısayollar

Sık kullanılan durumlar için hazır metodlar da mevcuttur:

*   `OnList()` / `HideOnList()`
*   `OnDetail()` / `HideOnDetail()`
*   `OnForm()` / `HideOnCreate()` / `HideOnUpdate()`
*   `OnlyOnDetail()`

Örnek:

```go
fields.ID().OnlyOnDetail(), // Sadece detay sayfasında görünür
fields.Text("Password").HideOnList().HideOnUpdate(), // Listede ve güncellemede gizli
```
