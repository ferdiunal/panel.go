# Kaynaklar (Resources)

Kaynaklar, Panel SDK'nın kalbidir. Go struct'larınızı (Modellerinizi) kullanıcı arayüzüne (UI) eşlerler.

## Bir Kaynak Tanımlama

`resource.Resource` arayüzünü uygulayan bir struct veya `resource.Base` kullanan bir yapı oluşturun:

```go
import (
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
)

type UserResource struct{
    resource.Base
}

func GetUserResource() resource.Resource {
    return &UserResource{
        Base: resource.Base{
            DataModel: &User{},
            Label:     "Kullanıcılar",
            FieldsVal: []fields.Element{
                fields.ID(),
                fields.Text("İsim", "Name"),
            },
        },
    }
}
```

### `Model`

Resource struct'ınızdaki `DataModel` alanı, SDK'ya hangi veritabanı tablosunun sorgulanacağını söyler.

```go
DataModel: &User{},
```

### `Fields`

`FieldsVal` (veya `Fields()` metodu override edilerek), görüntülenecek alanların listesini döndürür.

```go
FieldsVal: []fields.Element{
    fields.ID(),
    fields.Text("İsim", "Name"),
},
```

### `With` (Eager Loading / Ön Yükleme)

**Optimizasyon anahtardır.** N+1 sorgu sorununu önlemek için Custom Repository kullanarak ilişkileri önceden yükleyebilirsiniz.

```go
func (r *UserResource) Repository(db *gorm.DB) data.DataProvider {
    repo := data.NewGormDataProvider(db, r.DataModel)
    repo.SetWith([]string{"Company", "Profile"})
    return repo
}
```

## Yönlendirme ve Kayıt (Routing & Registration)

Kaynaklarınızı `app.go` (veya `main.go`) kurulumunda kaydedin:

```go
app.RegisterResource(GetUserResource())
```

Bu işlem aşağıdaki RESTful API uç noktalarını (endpoints) otomatik olarak oluşturur:

-   `GET /api/resource/users`
-   `GET /api/resource/users/:id`
-   `POST /api/resource/users`
-   `PUT /api/resource/users/:id`
-   `DELETE /api/resource/users/:id`

## Kart (Card) Desteği

Resource'larınıza KPI ve grafikler göstermek için kart ekleyebilirsiniz. Detaylar için **[Widgets (Cards)](Widgets)** sayfasına bakın.

```go
func (u *UserResource) Cards() []widget.Card {
    return []widget.Card{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
    }
}
```
