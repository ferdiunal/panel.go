# Kaynaklar (Resources)

Kaynaklar, Panel SDK'nın kalbidir. Go struct'larınızı (Modellerinizi) kullanıcı arayüzüne (UI) eşlerler.

## Bir Kaynak Tanımlama

`Resource` arayüzünü (interface) uygulayın:

```go
type Resource interface {
    Model() interface{}
    Fields() []fields.Element
    With() []string
}
```

### `Model()`

GORM modelinizin bir örneğini döndürür. Bu, SDK'ya hangi veritabanı tablosunun sorgulanacağını söyler.

```go
func (u *UserResource) Model() interface{} {
    return &User{}
}
```

### `Fields()`

Görüntülenecek alanların listesini (slice) döndürür.

```go
func (u *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("İsim", "Name"),
    }
}
```

### `With()` (Eager Loading / Ön Yükleme)

**Optimizasyon anahtardır.** `With` metodu, verilerinizi görüntülerken N+1 sorgu sorununa yol açmamak için hangi ilişkilerin önceden yükleneceğini (eager-loaded) belirtmenizi sağlar.

Eğer bir `Link`, `Collection` veya `Detail` alanı kullanıyorsanız, ilgili ilişkinin ismini burada listelemelisiniz.

```go
func (u *UserResource) With() []string {
    // "Company" ve "Profile" ilişkilerini önceden yükle
    return []string{"Company", "Profile"}
}
```

SDK, bu sorguları güçlü bir şekilde tiplemek (strongly type) için generic reflection kullanır, bu sayede GORM'un `Preload` işlevselliği generic arayüzün arkasında bile mükemmel çalışır.

## Yönlendirme ve Kayıt (Routing & Registration)

Kaynaklarınızı `app.go` kurulumunda kaydedin:

```go
panel.Register("users", &UserResource{})
panel.Register("posts", &PostResource{})
```

Bu işlem aşağıdaki RESTful API uç noktalarını (endpoints) otomatik olarak oluşturur:

-   `GET /api/resource/users`
-   `GET /api/resource/users/:id`
-   `POST /api/resource/users`
-   `PUT /api/resource/users/:id`
-   `DELETE /api/resource/users/:id`

## Widget Desteği

Resource'larınıza KPI ve grafikler göstermek için widget ekleyebilirsiniz. Detaylar için **[Widgets](Widgets)** sayfasına bakın.

```go
func (u *UserResource) Widgets() []widget.Widget {
    return []widget.Widget{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
    }
}
```
