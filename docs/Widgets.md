# Widget'lar (Widgets)

Widget'lar, resource'larınızın özet bilgilerini, grafiklerini ve metriklerini göstermenize olanak tanır. Laravel Nova Metrics sistemine benzer bir yapı sunar.

## Kullanılabilir Widget Tipleri

Şu an için iki temel widget tipi desteklenmektedir:

1.  **Value**: Tek bir sayısal değer ve (opsiyonel) değişim gösterir (örn. Toplam Kullanıcı Sayısı).
2.  **Trend**: Zaman içindeki değişimi çizgi grafik olarak gösterir (örn. Son 30 gündeki kayıtlar).

## Widget Oluşturma

### 1. Value Widget

`widget.NewCountWidget` helper'ını kullanarak hızlıca bir sayaç oluşturabilirsiniz:

```go
import "panel.go/internal/widget"

func (u *UserResource) Widgets() []widget.Widget {
    return []widget.Widget{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
        
        // Veya manuel tanımlama
        &widget.Value{
            Title: "Aktif Aboneler",
            QueryFunc: func(db *gorm.DB) (int64, error) {
                var count int64
                err := db.Model(&Subscription{}).Where("status = ?", "active").Count(&count).Error
                return count, err
            },
        },
    }
}
```

### 2. Trend Widget

Trend widget'ları, verilerin zaman içindeki dağılımını gösterir.

```go
func (u *UserResource) Widgets() []widget.Widget {
    return []widget.Widget{
        widget.NewTrendWidget("Günlük Kayıtlar", &User{}, "created_at"),
    }
}
```

## Resource'a Ekleme

Widget'ları resource'unuza eklemek için `Widgets()` metodunu implemente etmeniz yeterlidir:

```go
func (u *UserResource) Widgets() []widget.Widget {
    return []widget.Widget{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
    }
}
```
