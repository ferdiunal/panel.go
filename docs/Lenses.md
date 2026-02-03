# Lensler (Lenses)

Lensler, kaynaklarınızın (resources) belirli bir filtrelenmiş veya özelleştirilmiş görünümünü oluşturmanıza olanak tanır. Standart `Index` görünümünün aksine, Lensler belirli bir soruyu yanıtlamak veya özel bir rapor sunmak için tasarlanmıştır.

Örneğin, "En Çok Satış Yapan Kullanıcılar" veya "Onay Bekleyen Yorumlar" gibi durumlar için idealdir.

## Bir Lens Tanımlama

`Lens` arayüzü (interface) dört temel metodu içerir:

1.  `Name()`: UI'da görünecek isim.
2.  `Slug()`: URL'de kullanılacak benzersiz tanımlayıcı.
3.  `Query()`: Veritabanı sorgusunu özelleştirebileceğiniz yer.
4.  `Fields()`: Bu lens görünümünde hangi alanların gösterileceği.

### Örnek: En Popüler Bloglar

```go
type MostPopularBlogsLens struct{}

func (l *MostPopularBlogsLens) Name() string { 
    return "En Popüler Bloglar" 
}

func (l *MostPopularBlogsLens) Slug() string { 
    return "most-popular" 
}

func (l *MostPopularBlogsLens) Query(db *gorm.DB) *gorm.DB {
    // Sadece 1000'den fazla görüntülenen blogları filtrele
    // ve görüntülenme sayısına göre sırala
    return db.Where("views > ?", 1000).Order("views desc")
}

func (l *MostPopularBlogsLens) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Başlık", "Title"),
        fields.Number("Görüntülenme", "Views"),
        // İlişkileri de gösterebilirsiniz
        fields.Link("Yazar", "user_id"), 
    }
}
```

## Kaynağa Lens Ekleme

Oluşturduğunuz Lensi kullanmak için, ilgili Kaynağın (Resource) `Lenses` metoduna eklemeniz gerekir:

```go
func (r *BlogResource) Lenses() []resource.Lens {
    return []resource.Lens{
        &MostPopularBlogsLens{},
        &RecentBlogsLens{},
    }
}
```

## Erişim (Routing)

Kaydettiğiniz Lenslere aşağıdaki otomatik API rotası üzerinden erişilebilir:

`GET /api/resource/:resource/lens/:lens`

Örneğin, yukarıdaki "En Popüler Bloglar" lensine erişmek için:

`GET /api/resource/blogs/lens/most-popular`

Bu istek, standart `Index` yanıt formatında ancak Lens'in `Query` filtresi uygulanmış ve `Fields` tanımındaki alanları içeren bir JSON döner.
