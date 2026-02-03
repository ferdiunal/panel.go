# Yetkilendirme ve Koşullu Görünürlük (Authorization & Conditionals)

Panel SDK, kaynaklara erişimi kontrol etmek ve alanların (fields) görünürlüğünü dinamik olarak yönetmek için güçlü bir yapı sunar.

## Politikalar (Policies)

Politikalar, bir kullanıcının belirli bir kaynak üzerinde işlem yapıp yapamayacağını belirleyen mantığı kapsüller.

### Policy Arayüzü

Herhangi bir yetkilendirme mantığı `auth.Policy` arayüzünü uygulamalıdır:

```go
type Policy interface {
    ViewAny(ctx *appContext.Context) bool             // Liste sayfasına erişim
    View(ctx *appContext.Context, model interface{}) bool // Detay sayfasına erişim
    Create(ctx *appContext.Context) bool              // Yeni oluşturma izni
    Update(ctx *appContext.Context, model interface{}) bool // Güncelleme izni
    Delete(ctx *appContext.Context, model interface{}) bool // Silme izni
}
```

### Rol Tabanlı Erişim (RBAC) ve İzin Dosyası

Panel, izinleri yönetmek için `permissions.toml` dosyasını kullanabilir. Bu dosya proje kök dizininde bulunur ve roller ile kaynaklar arasındaki ilişkiyi tanımlar.

Örnek `permissions.toml`:
```toml
system_roles = ["admin", "editor", "user"]

[resources]
  [resources.posts]
  label = "Blog Yazıları"
  actions = ["view_any", "view", "create", "update", "delete"]
```

### Policy ile Entegrasyon

Policy metodlarınızda `ctx.HasPermission("resource.action")` metodunu kullanarak yetki kontrolü yapabilirsiniz. `HasPermission` metodu, `admin` rolü için otomatik olarak `true` döner.

```go
type PostPolicy struct {}

func (p *PostPolicy) ViewAny(ctx *appContext.Context) bool {
    // "posts.view_any" izni var mı?
    return ctx.HasPermission("posts.view_any")
}

func (p *PostPolicy) View(ctx *appContext.Context, model interface{}) bool {
    return ctx.HasPermission("posts.view")
}

func (p *PostPolicy) Create(ctx *appContext.Context) bool {
    return ctx.HasPermission("posts.create")
}

func (p *PostPolicy) Update(ctx *appContext.Context, model interface{}) bool {
    // Karmaşık mantık: İzin VARSA VE (Admin VEYA Yazının Sahibi ise)
    post := model.(*Post)
    return ctx.HasPermission("posts.update") && (ctx.User().Role == "admin" || post.UserID == ctx.User().ID)
}

func (p *PostPolicy) Delete(ctx *appContext.Context, model interface{}) bool {
    return ctx.HasPermission("posts.delete")
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
