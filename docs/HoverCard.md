# Hover Card - İlişki Field'ları için Dinamik Veri Gösterimi

Hover card özelliği, ilişki field'larının (BelongsTo, HasOne, MorphTo) index ve detail sayfalarında fare ile üzerine gelindiğinde dinamik olarak veri gösterilmesini sağlar. Bu özellik, field resolver handler yaklaşımı ile çalışır ve tam kontrol sağlar.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Temel Kullanım](#temel-kullanım)
- [API Endpoint'leri](#api-endpointleri)
- [Özellikler](#özellikler)
- [Örnekler](#örnekler)
- [Güvenlik](#güvenlik)
- [Performans](#performans)

## Genel Bakış

Hover card özelliği şu bileşenlerden oluşur:

1. **HoverCardConfig**: Hover card konfigürasyonu
2. **HoverCardResolver**: Hover card verilerini çözen callback fonksiyonu
3. **Field Metodları**: `HoverCard()` ve `ResolveHoverCard()` metodları
4. **API Endpoint**: `/api/resource/{resource}/resolver/{field}` endpoint'i

## Temel Kullanım

### 1. Hover Card Struct'ını Tanımlama

Önce hover card'da gösterilecek verileri içeren bir struct tanımlayın:

```go
type AuthorHoverCard struct {
    Avatar string `json:"avatar"`
    Name   string `json:"name"`
    Email  string `json:"email"`
    Phone  string `json:"phone"`
    Role   string `json:"role"`
}
```

### 2. Field'a Hover Card Ekleme

Field'a hover card eklemek için `HoverCard()` ve `ResolveHoverCard()` metodlarını kullanın:

```go
field := fields.BelongsTo("Author", "author_id", "authors").
    DisplayUsing("name").
    HoverCard(&AuthorHoverCard{}).
    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
        // İlişkili kaydı veritabanından al
        author := &Author{}
        if err := db.First(author, relatedID).Error; err != nil {
            return nil, err
        }

        // Hover card verisini döndür
        return &AuthorHoverCard{
            Avatar: author.Avatar,
            Name:   author.Name,
            Email:  author.Email,
            Phone:  author.Phone,
            Role:   author.Role,
        }, nil
    })
```

## API Endpoint'leri

Hover card verileri şu endpoint'ler üzerinden alınır:

### GET İsteği

```
GET /api/resource/{resource}/resolver/{field}?id={related_id}
```

**Parametreler:**
- `resource`: Kaynak adı (örn: "posts", "users")
- `field`: Field adı (örn: "author_id", "profile")
- `id`: İlişkili kaydın ID'si (query parameter)
- `record_id`: Ana kaydın ID'si (opsiyonel, query parameter)

**Örnek:**
```bash
curl -X GET "http://localhost:8080/api/resource/posts/resolver/author_id?id=5"
```

### POST İsteği

```
POST /api/resource/{resource}/resolver/{field}
Content-Type: application/json

{
  "id": 5,
  "record_id": 123
}
```

**Örnek:**
```bash
curl -X POST "http://localhost:8080/api/resource/posts/resolver/author_id" \
  -H "Content-Type: application/json" \
  -d '{"id": 5, "record_id": 123}'
```

### Response Format

Başarılı durumda:
```json
{
  "data": {
    "avatar": "https://example.com/avatar.jpg",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1 234 567 8900",
    "role": "Admin"
  }
}
```

Hata durumunda:
```json
{
  "error": "Field not found"
}
```

## Özellikler

### 1. HTTP Method Desteği

Hover card endpoint'i şu HTTP method'larını destekler:
- **GET**: Veri okuma
- **POST**: Veri gönderme (body ile)
- **PATCH**: Kısmi güncelleme
- **DELETE**: Silme işlemi

### 2. Esnek Parametre Desteği

Parametreler hem query string hem de request body ile gönderilebilir:

```go
// Query string
GET /api/resource/posts/resolver/author_id?id=5&record_id=123

// Request body
POST /api/resource/posts/resolver/author_id
{
  "id": 5,
  "record_id": 123
}
```

### 3. Context Desteği

Resolver fonksiyonu, ana kaydı (record) parametre olarak alır. Bu sayede hover card verisi, ana kayda göre özelleştirilebilir:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    // Ana kayıt üzerinden işlem yapabilirsiniz
    post := record.(*Post)

    // İlişkili kaydı al
    author := &Author{}
    db.First(author, relatedID)

    // Ana kayda göre özelleştir
    if post.Status == "published" {
        return &AuthorHoverCard{
            Avatar: author.Avatar,
            Name:   author.Name,
            Email:  author.Email,
        }, nil
    }

    // Draft için daha az bilgi göster
    return &AuthorHoverCard{
        Name: author.Name,
    }, nil
})
```

## Örnekler

### BelongsTo Field

```go
type AuthorHoverCard struct {
    Avatar   string `json:"avatar"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Bio      string `json:"bio"`
    PostCount int   `json:"post_count"`
}

field := fields.BelongsTo("Author", "author_id", "authors").
    DisplayUsing("name").
    HoverCard(&AuthorHoverCard{}).
    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
        author := &Author{}
        if err := db.First(author, relatedID).Error; err != nil {
            return nil, err
        }

        // Yazar'ın post sayısını hesapla
        var postCount int64
        db.Model(&Post{}).Where("author_id = ?", relatedID).Count(&postCount)

        return &AuthorHoverCard{
            Avatar:    author.Avatar,
            Name:      author.Name,
            Email:     author.Email,
            Bio:       author.Bio,
            PostCount: int(postCount),
        }, nil
    })
```

### HasOne Field

```go
type ProfileHoverCard struct {
    Avatar   string `json:"avatar"`
    Bio      string `json:"bio"`
    Location string `json:"location"`
    Website  string `json:"website"`
}

field := fields.HasOne("Profile", "profile", "profiles").
    ForeignKey("user_id").
    HoverCard(&ProfileHoverCard{}).
    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
        profile := &Profile{}
        if err := db.First(profile, relatedID).Error; err != nil {
            return nil, err
        }

        return &ProfileHoverCard{
            Avatar:   profile.Avatar,
            Bio:      profile.Bio,
            Location: profile.Location,
            Website:  profile.Website,
        }, nil
    })
```

### MorphTo Field

```go
type CommentableHoverCard struct {
    Thumbnail string `json:"thumbnail"`
    Title     string `json:"title"`
    Type      string `json:"type"`
    Author    string `json:"author"`
}

field := fields.NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    }).
    HoverCard(&CommentableHoverCard{}).
    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
        // MorphTo için tip bilgisini al
        comment := record.(*Comment)
        morphType := comment.CommentableType

        // Tip'e göre ilişkili kaydı al
        switch morphType {
        case "post":
            post := &Post{}
            if err := db.First(post, relatedID).Error; err != nil {
                return nil, err
            }
            return &CommentableHoverCard{
                Thumbnail: post.FeaturedImage,
                Title:     post.Title,
                Type:      "Post",
                Author:    post.Author.Name,
            }, nil

        case "video":
            video := &Video{}
            if err := db.First(video, relatedID).Error; err != nil {
                return nil, err
            }
            return &CommentableHoverCard{
                Thumbnail: video.Thumbnail,
                Title:     video.Title,
                Type:      "Video",
                Author:    video.Creator.Name,
            }, nil
        }

        return nil, fmt.Errorf("unknown morph type: %s", morphType)
    })
```

### Koşullu Veri Gösterimi

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    // Context'ten kullanıcı bilgisini al
    user := ctx.Value("user").(*User)

    author := &Author{}
    db.First(author, relatedID)

    // Kullanıcı rolüne göre farklı veri göster
    if user.Role == "admin" {
        return &AuthorHoverCard{
            Avatar: author.Avatar,
            Name:   author.Name,
            Email:  author.Email,
            Phone:  author.Phone,
            Role:   author.Role,
        }, nil
    }

    // Normal kullanıcılar için sadece temel bilgiler
    return &AuthorHoverCard{
        Avatar: author.Avatar,
        Name:   author.Name,
    }, nil
})
```

## Güvenlik

### 1. Authorization Kontrolü

Resolver içinde authorization kontrolü yapın:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    // Kullanıcı kontrolü
    user := ctx.Value("user").(*User)
    if user == nil {
        return nil, fmt.Errorf("unauthorized")
    }

    // İzin kontrolü
    if !user.Can("view_author_details") {
        return nil, fmt.Errorf("permission denied")
    }

    // Veriyi döndür
    author := &Author{}
    db.First(author, relatedID)
    return &AuthorHoverCard{...}, nil
})
```

### 2. Hassas Veri Filtreleme

Hassas verileri döndürmeden önce filtreleyin:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    author := &Author{}
    db.First(author, relatedID)

    // Hassas verileri filtrele
    return &AuthorHoverCard{
        Avatar: author.Avatar,
        Name:   author.Name,
        Email:  maskEmail(author.Email), // Email'i maskele
        // Phone alanını döndürme
    }, nil
})

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }
    return parts[0][:2] + "***@" + parts[1]
}
```

### 3. Rate Limiting

API endpoint'ine rate limiting uygulayın (middleware ile):

```go
// Rate limiting middleware
app.Use("/api/resource/*/resolver/*", rateLimiter.New(rateLimiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))
```

## Performans

### 1. Caching

Hover card verilerini cache'leyin:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    // Cache key oluştur
    cacheKey := fmt.Sprintf("author_hover:%v", relatedID)

    // Cache'den kontrol et
    if cached, found := cache.Get(cacheKey); found {
        return cached, nil
    }

    // Veritabanından al
    author := &Author{}
    db.First(author, relatedID)

    hoverCard := &AuthorHoverCard{
        Avatar: author.Avatar,
        Name:   author.Name,
        Email:  author.Email,
    }

    // Cache'e kaydet (5 dakika)
    cache.Set(cacheKey, hoverCard, 5*time.Minute)

    return hoverCard, nil
})
```

### 2. N+1 Sorgu Problemini Önleme

Eager loading kullanın:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    author := &Author{}

    // İlişkili verileri eager load et
    if err := db.Preload("Posts").Preload("Comments").First(author, relatedID).Error; err != nil {
        return nil, err
    }

    return &AuthorHoverCard{
        Avatar:       author.Avatar,
        Name:         author.Name,
        Email:        author.Email,
        PostCount:    len(author.Posts),
        CommentCount: len(author.Comments),
    }, nil
})
```

### 3. Gereksiz Veri Döndürmeme

Sadece gerekli verileri döndürün:

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    author := &Author{}

    // Sadece gerekli alanları seç
    if err := db.Select("id", "avatar", "name", "email").First(author, relatedID).Error; err != nil {
        return nil, err
    }

    return &AuthorHoverCard{
        Avatar: author.Avatar,
        Name:   author.Name,
        Email:  author.Email,
    }, nil
})
```

## Frontend Kullanımı

### TypeScript/JavaScript

```typescript
// GET isteği
const fetchHoverCardData = async (resource: string, field: string, id: number) => {
  try {
    const response = await axios.get(
      `/api/resource/${resource}/resolver/${field}`,
      { params: { id } }
    );
    return response.data.data;
  } catch (error) {
    console.error('Hover card data fetch failed:', error);
    return null;
  }
};

// Kullanım
const authorData = await fetchHoverCardData('posts', 'author_id', 5);
console.log(authorData);
// { avatar: "...", name: "...", email: "...", phone: "..." }
```

### React Hook

```typescript
import { useState, useEffect } from 'react';
import axios from 'axios';

const useHoverCardData = (resource: string, field: string, id: number | null) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!id) return;

    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await axios.get(
          `/api/resource/${resource}/resolver/${field}`,
          { params: { id } }
        );
        setData(response.data.data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [resource, field, id]);

  return { data, loading, error };
};

// Kullanım
const AuthorHoverCard = ({ authorId }) => {
  const { data, loading, error } = useHoverCardData('posts', 'author_id', authorId);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error loading data</div>;
  if (!data) return null;

  return (
    <div>
      <img src={data.avatar} alt={data.name} />
      <h3>{data.name}</h3>
      <p>{data.email}</p>
    </div>
  );
};
```

## Hata Yönetimi

### Backend

```go
ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.RelationshipField) (interface{}, error) {
    author := &Author{}

    // Veritabanı hatası
    if err := db.First(author, relatedID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("author not found")
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    // Validation hatası
    if author.Status == "deleted" {
        return nil, fmt.Errorf("author has been deleted")
    }

    return &AuthorHoverCard{
        Avatar: author.Avatar,
        Name:   author.Name,
        Email:  author.Email,
    }, nil
})
```

### Frontend

```typescript
const fetchHoverCardData = async (resource: string, field: string, id: number) => {
  try {
    const response = await axios.get(
      `/api/resource/${resource}/resolver/${field}`,
      { params: { id } }
    );
    return { data: response.data.data, error: null };
  } catch (error) {
    if (error.response?.status === 404) {
      return { data: null, error: 'Record not found' };
    }
    if (error.response?.status === 403) {
      return { data: null, error: 'Permission denied' };
    }
    return { data: null, error: 'Failed to load data' };
  }
};
```

## Best Practices

1. **Minimal Veri Döndürme**: Sadece hover card'da gösterilecek verileri döndürün
2. **Cache Kullanımı**: Sık erişilen verileri cache'leyin
3. **Authorization**: Her zaman authorization kontrolü yapın
4. **Error Handling**: Kapsamlı hata yönetimi uygulayın
5. **Performance**: N+1 sorgu problemine dikkat edin
6. **Security**: Hassas verileri filtreleyin
7. **Rate Limiting**: API endpoint'ine rate limiting uygulayın
8. **Logging**: Hataları ve önemli olayları loglayın

## Sorun Giderme

### Hover Card Verisi Gelmiyor

1. Field'ın hover card konfigürasyonunu kontrol edin
2. Resolver callback'inin doğru tanımlandığını kontrol edin
3. API endpoint'inin doğru çağrıldığını kontrol edin
4. Network tab'ında request/response'u inceleyin

### Authorization Hatası

1. Kullanıcının oturum açtığını kontrol edin
2. Kullanıcının gerekli izinlere sahip olduğunu kontrol edin
3. Resolver içinde authorization kontrolünü kontrol edin

### Performance Problemi

1. N+1 sorgu problemini kontrol edin
2. Cache kullanımını kontrol edin
3. Gereksiz veri döndürülmediğini kontrol edin
4. Database query'lerini optimize edin

## İlgili Dokümantasyon

- [Relationships](Relationships.md) - İlişki field'ları hakkında genel bilgi
- [Field Resolver](FieldResolver.md) - Field resolver hakkında detaylı bilgi
- [API Reference](API.md) - API endpoint'leri hakkında detaylı bilgi
