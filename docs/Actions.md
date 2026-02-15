# Action'lar Rehberi (Legacy Teknik Akış)

Action, resource kayıtları üzerinde toplu işlem çalıştırmanızı sağlar.

Örnek kullanım:
- Toplu onaylama
- Toplu durum değiştirme
- Seçili kayıtları dışa aktarma
- Sistem seviyesinde (ID seçmeden) bakım işlemi

## Bu Doküman Ne Zaman Okunmalı?

Önerilen sıra:
1. [Başlarken](Getting-Started)
2. [Kaynaklar (Resource)](Resources)
3. [Yetkilendirme](Authorization)
4. Bu doküman (`Actions`)

## Kısa Özet

Action akışı:
1. Resource `GetActions()` ile action'ları döner.
2. Frontend `GET /api/resource/:resource/actions` ile listeyi alır.
3. Kullanıcı action seçer ve `POST /api/resource/:resource/actions/:action` çağrılır.
4. Backend policy + doğrulama + `CanRun` kontrolünden sonra `Execute` çalıştırır.

## Hızlı Uygulama Akışı

1. Action'ı `action.New(...)` ile tanımla ve benzersiz `slug` ver.
2. Gerekli kullanıcı girdilerini `WithFields(...)` ile ekle.
3. İş mantığını `Handle(...)` içinde yaz.
4. Çalışma yetkisini `AuthorizeUsing(...)` ve policy zinciriyle sınırla.
5. İhtiyaca göre `Standalone()` veya `Sole()` davranışını seç.

## Action Tanımlama (Fluent API)

`pkg/action` içinde en pratik yol `action.New(...)` ile tanımlamaktır.

```go
func (r *UserResource) GetActions() []resource.Action {
	return []resource.Action{
		action.New("Kullanıcıyı Pasifleştir").
			SetSlug("deactivate-user").
			SetIcon("user-x").
			Confirm("Seçili kullanıcıları pasifleştirmek istediğinize emin misiniz?").
			ConfirmButton("Pasifleştir").
			CancelButton("İptal").
			Destructive().
			WithFields(
				fields.Text("Neden", "reason").Required(),
			).
			Handle(func(ctx *action.ActionContext) error {
				reason, _ := ctx.Fields["reason"].(string)

				for _, item := range ctx.Models {
					user, ok := item.(*entity.User)
					if !ok {
						continue
					}
					user.IsActive = false
					user.DeactivateReason = reason
					if err := ctx.DB.Save(user).Error; err != nil {
						return err
					}
				}
				return nil
			}).
			AuthorizeUsing(func(ctx *action.ActionContext) bool {
				// basit örnek, burada rol kontrolü yapılabilir
				return true
			}),
	}
}
```

## Standalone ve Sole

### Standalone action

ID seçmeden çalışabilir.

```go
action.New("Cache Temizle").
	SetSlug("clear-cache").
	Standalone().
	Handle(func(ctx *action.ActionContext) error {
		return nil
	})
```

Sunucu kuralı:
- `standalone` değilse boş `ids` kabul edilmez.

### Sole action

Sadece tek kayıt seçildiğinde çalışır.

```go
action.New("MFA Sıfırla").
	SetSlug("reset-mfa").
	Sole().
	Handle(func(ctx *action.ActionContext) error {
		// ctx.Models burada tek kayıt olmalı
		return nil
	})
```

Sunucu kuralı:
- `sole` action için birden fazla `id` gönderilirse istek `400` döner.

## Zorunlu Action Field Doğrulaması

Action field'larında `Required()` işaretlenen alanlar, backend tarafında da doğrulanır.

Eksik alan gönderilirse `400` döner:

```json
{
  "error": "Reason is required"
}
```

## Endpoint'ler

### Action listesi

`GET /api/resource/:resource/actions`

Örnek yanıt:

```json
{
  "actions": [
    {
      "name": "Kullanıcıyı Pasifleştir",
      "slug": "deactivate-user",
      "icon": "user-x",
      "confirmText": "Seçili kullanıcıları pasifleştirmek istediğinize emin misiniz?",
      "confirmButtonText": "Pasifleştir",
      "cancelButtonText": "İptal",
      "destructive": true,
      "onlyOnIndex": false,
      "onlyOnDetail": false,
      "showInline": false,
      "standalone": false,
      "sole": false,
      "fields": []
    }
  ]
}
```

### Action çalıştırma

`POST /api/resource/:resource/actions/:action`

Body:

```json
{
  "ids": ["1", "2"],
  "fields": {
    "reason": "Toplu pasifleştirme"
  }
}
```

Başarılı yanıt:

```json
{
  "message": "Action executed successfully on 2 item(s)",
  "count": 2
}
```

## Lens Üzerinden Action

Lens görünümünde action endpoint'leri ayrı olarak da kullanılabilir:

- `GET /api/resource/:resource/lens/:lens/actions`
- `POST /api/resource/:resource/lens/:lens/actions/:action`

Bu endpoint'ler, Nova parity için lens route'u altında action erişimi sağlar.

## Yetki Akışı

Sunucu tarafında sırayla:
1. `Policy.Update(...)`
2. Action bulundu mu kontrolü
3. `standalone/sole` doğrulaması
4. Required field doğrulaması
5. `CanRun(...)`
6. `Execute(...)`

## Sık Hata

### "No items selected"

- Action `Standalone()` değil.
- `ids` boş gönderildi.

### "This action can only run on a single item"

- Action `Sole()` ama birden fazla `id` gönderildi.

### "Action cannot be executed in this context"

- `CanRun(...)` false döndü.

### "Action not found"

- URL'deki action slug ile `GetActions()` içindeki slug eşleşmiyor.

## Sonraki Adım

- Lens aksiyonları için: [Lensler (Lenses)](Lenses)
- Kart/özet görünümü için: [Widget'lar (Cards)](Widgets)
- Uçtan uca API davranışı için: [API Referansı](API-Reference)
