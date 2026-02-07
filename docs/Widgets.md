# Widget'lar (Cards)

Widget'lar (veya Cards), resource'larınızın özet bilgilerini, grafiklerini ve metriklerini göstermenize olanak tanır.

## Kullanılabilir Kart Tipleri

Şu an için iki temel kart tipi desteklenmektedir:

1.  **Value**: Tek bir sayısal değer ve (opsiyonel) değişim gösterir (örn. Toplam Kullanıcı Sayısı).
2.  **Trend**: Zaman içindeki değişimi çizgi grafik olarak gösterir (örn. Son 30 gündeki kayıtlar).

## Kart Oluşturma

### 1. Value Card

`widget.NewCountWidget` helper'ını kullanarak hızlıca bir sayaç oluşturabilirsiniz:

```go
import "github.com/ferdiunal/panel.go/pkg/widget"

func (u *UserResource) Cards() []widget.Card {
    return []widget.Card{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
        
        // Veya manuel tanımlama (Custom Card)
        widget.NewCard("Aktif Aboneler", "value-metric").
            SetContent(calculateSubscribers()),
    }
}
```

### 2. Trend Card

Trend widget'ları, verilerin zaman içindeki dağılımını gösterir.

```go
func (u *UserResource) Cards() []widget.Card {
    return []widget.Card{
        widget.NewTrendWidget("Günlük Kayıtlar", &User{}, "created_at"),
    }
}
```

## Resource'a Ekleme

Kartları resource'unuza eklemek için `Cards()` metodunu implemente etmeniz yeterlidir:

```go
func (u *UserResource) Cards() []widget.Card {
    return []widget.Card{
        widget.NewCountWidget("Toplam Kullanıcı", &User{}),
    }
}
```

## Sayfalara Ekleme (Dashboard)

Kartlar sadece resource'larda değil, `Page` (Sayfa) yapılarında da kullanılabilir. Örneğin Dashboard sayfasında:

```go
func (d *Dashboard) Cards() []widget.Card {
    return []widget.Card{
        widget.NewCountWidget("Toplam Kullanıcı", &user.User{}),
    }
}
```
