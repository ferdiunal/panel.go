# Widget'lar (Cards) - Legacy Teknik Akış

Widget'lar (veya Cards), resource'larınızın özet bilgilerini, grafiklerini ve metriklerini göstermenize olanak tanır.

## Bu Doküman Ne Zaman Okunmalı?

Önerilen sıra:
1. [Başlarken](Getting-Started)
2. [Kaynaklar (Resource)](Resources)
3. [Lensler (Lenses)](Lenses)
4. Bu doküman (`Widgets`)

## Hızlı Widget Akışı

1. İhtiyaca göre kart tipi seç (`Value` veya `Trend`).
2. Resource veya Page üzerinde `Cards()` metodunda kartları döndür.
3. Kart verisini mümkün olduğunca hafif sorgularla üret.
4. Büyük veri setlerinde widget hesaplarını cache stratejisiyle destekle.

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

## Sık Hata Kontrolü (Widget)

- Kart görünmüyor: `Cards()` metodunun ilgili resource/page üzerinde gerçekten implement edildiğini kontrol edin.
- Değer yanlış: widget içinde kullanılan model/sorgu alan eşleşmelerini doğrulayın.
- Trend boş: zaman alanı (`created_at` vb.) yanlış veya null olabilir.
- Performans düşüyor: ağır metrik hesaplarını istek anında değil cache/ön-hesaplama ile çalıştırın.

## Sonraki Adım

- Lens ile birlikte kullanım için: [Lensler (Lenses)](Lenses)
- Toplu aksiyon + metrik kombinasyonu için: [Action'lar](Actions)
- Sayfa entegrasyonu için: [Sayfalar (Pages)](Pages)
