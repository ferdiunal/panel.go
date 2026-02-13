# DialogField - Modal/Dialog Field Sistemi

## Genel Bakış

DialogField, Panel.go projesine eklenen modal/dialog tabanlı form ve wizard sistemidir. Kullanıcıdan modal içinde veri toplamak, multi-step wizard formları oluşturmak ve sayfa geçişlerinde kullanıcıyı bilgilendirmek için kullanılır.

### Temel Özellikler

- **İki Mod**: Basit form veya multi-step wizard
- **Esnek Tetikleme**: Varsayılan açık veya buton ile tetiklenebilir
- **Özelleştirilebilir**: Dialog boyutu, başlık, açıklama
- **Progress Indicator**: Wizard mode için adım göstergesi
- **Skip Functionality**: Wizard adımları atlanabilir
- **UniversalResourceForm Entegrasyonu**: Mevcut form sistemi ile tam uyumlu
- **Field Registry Entegrasyonu**: Diğer field'lar gibi otomatik render

---

## Mimari

### Backend (Go)

```
pkg/fields/
└── dialog.go              # DialogField struct ve builder methods
```

**Ana Bileşenler:**
- `DialogField`: Ana field struct'ı (Schema'yı embed eder)
- `DialogContentType`: Content tipi enum (form, wizard)
- `DialogStep`: Wizard adımı struct'ı
- `Dialog()`: Helper fonksiyon (field oluşturma)

### Frontend (React/TypeScript)

```
web/src/
├── types/
│   └── dialog.ts                    # Type definitions
├── components/
│   └── fields/
│       ├── DialogField.tsx          # Ana component
│       ├── DialogContent.tsx        # Form content
│       └── DialogWizard.tsx         # Wizard content
└── components/forms/fields/
    └── index.ts                     # Field registry (güncellendi)
```

**Ana Bileşenler:**
- `DialogField`: Ana component (modal yönetimi)
- `DialogContent`: Basit form content wrapper
- `DialogWizard`: Multi-step wizard component
- `DialogFieldProps`: TypeScript type definitions

---

## Backend Kullanımı

### 1. Basit Form Dialog

Tek adımlı form için kullanılır. Varsayılan açık veya buton ile tetiklenebilir.

#### Varsayılan Açık Mode

```go
func (r *UserResource) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.ID("ID"),
        fields.Text("Ad", "name"),
        fields.Email("Email", "email"),

        // Profil tamamlama dialog'u - sayfa açıldığında otomatik açılır
        fields.Dialog("Profil Tamamla", "profile_completion").
            DefaultOpen(true).                                    // Otomatik açılır
            DialogTitle("Profilinizi Tamamlayın").
            DialogDesc("Lütfen eksik bilgilerinizi doldurun").
            DialogSize("md").                                     // Orta boyut
            Content([]core.Element{
                fields.Text("Telefon", "phone").Required(),
                fields.Text("Adres", "address").Required(),
                fields.Date("Doğum Tarihi", "birth_date"),
            }).
            OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
                // Form tamamlandığında çağrılır
                userID := ctx.Locals("user_id")
                // Veriyi kaydet
                return db.UpdateUserProfile(userID, data)
            }),
    }
}
```

#### Buton ile Tetikleme

```go
fields.Dialog("Ayarlar", "settings_dialog").
    TriggerButton("Gelişmiş Ayarlar").                    // Buton metni
    TriggerIcon("⚙️").                                    // Buton ikonu
    DialogTitle("Gelişmiş Ayarlar").
    Content([]core.Element{
        fields.Switch("Debug Mode", "debug_mode"),
        fields.Text("API Key", "api_key"),
        fields.Number("Timeout (sn)", "timeout").Min(1).Max(300),
    }).
    OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
        return saveSettings(data)
    })
```

### 2. Wizard Dialog (Multi-Step)

Çok adımlı form süreçleri için kullanılır. Her adım kendi field'larına sahiptir.

```go
fields.Dialog("Onboarding", "onboarding_wizard").
    TriggerButton("Başlangıç Rehberini Başlat").
    TriggerIcon("🚀").
    DialogTitle("Hoş Geldiniz!").
    DialogDesc("Hızlı bir kurulum ile başlayalım").
    DialogSize("lg").                                     // Büyük boyut
    Wizard([]fields.DialogStep{
        // Adım 1: Kişisel Bilgiler
        {
            Index:       0,
            Title:       "Kişisel Bilgiler",
            Description: "Önce sizi tanıyalım",
            Fields: []core.Element{
                fields.Text("Ad Soyad", "full_name").Required(),
                fields.Email("Email", "email").Required(),
                fields.Tel("Telefon", "phone"),
            },
            CanSkip: false,                               // Atlanamaz
        },
        // Adım 2: Şirket Bilgileri
        {
            Index:       1,
            Title:       "Şirket Bilgileri",
            Description: "Şirketiniz hakkında bilgi verin",
            Fields: []core.Element{
                fields.Text("Şirket Adı", "company_name").Required(),
                fields.Select("Sektör", "industry").Options(map[string]string{
                    "tech":    "Teknoloji",
                    "finance": "Finans",
                    "health":  "Sağlık",
                }),
                fields.Number("Çalışan Sayısı", "employee_count"),
            },
            CanSkip: true,                                // Atlanabilir
        },
        // Adım 3: Tercihler
        {
            Index:       2,
            Title:       "Tercihler",
            Description: "Son olarak tercihlerinizi belirleyin",
            Fields: []core.Element{
                fields.Switch("Email Bildirimleri", "email_notifications"),
                fields.Switch("SMS Bildirimleri", "sms_notifications"),
                fields.Select("Dil", "language").Options(map[string]string{
                    "tr": "Türkçe",
                    "en": "English",
                }),
            },
            CanSkip: true,
        },
    }).
    OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
        // Tüm wizard tamamlandığında çağrılır
        // data: tüm adımlardan toplanan veri
        userID := ctx.Locals("user_id")
        return completeOnboarding(userID, data)
    }).
    OnSkip(func(ctx *fiber.Ctx) error {
        // Wizard atlandığında çağrılır
        userID := ctx.Locals("user_id")
        return markOnboardingSkipped(userID)
    })
```

### 3. Conditional Visibility (Policy ile)

DialogField'ı sadece belirli kullanıcılara göstermek için:

```go
fields.Dialog("Admin Ayarları", "admin_settings").
    TriggerButton("Admin Paneli").
    CanSee(func(ctx *core.ResourceContext) bool {
        // Sadece admin kullanıcılar görebilir
        user := ctx.User
        return user.IsAdmin()
    }).
    Content([]core.Element{
        fields.Switch("Maintenance Mode", "maintenance_mode"),
        fields.Text("System Message", "system_message"),
    })
```

### 4. Dialog Boyutları

```go
// Küçük dialog (form için)
DialogSize("sm")    // max-w-sm (384px)

// Orta dialog (varsayılan)
DialogSize("md")    // max-w-md (448px)

// Büyük dialog (wizard için)
DialogSize("lg")    // max-w-lg (512px)

// Çok büyük dialog
DialogSize("xl")    // max-w-xl (576px)

// Tam ekran
DialogSize("full")  // max-w-full
```

---

## Frontend Kullanımı

DialogField, field registry'ye kayıtlı olduğu için otomatik olarak render edilir. Manuel kullanım gerekmez.

### Manuel Kullanım (Gerekirse)

```tsx
import { DialogField } from '@/components/fields/DialogField';

function MyComponent() {
  return (
    <DialogField
      name="profile_completion"
      label="Profil Tamamla"
      defaultOpen={true}
      dialogTitle="Profilinizi Tamamlayın"
      dialogDesc="Lütfen eksik bilgilerinizi doldurun"
      contentType="form"
      fields={[
        { key: 'phone', name: 'Telefon', view: 'text-field', required: true },
        { key: 'address', name: 'Adres', view: 'text-field', required: true },
      ]}
      onChange={(data) => console.log('Data:', data)}
    />
  );
}
```

### Wizard Mode

```tsx
<DialogField
  name="onboarding"
  label="Onboarding"
  triggerButton="Başlangıç Rehberini Başlat"
  dialogTitle="Hoş Geldiniz!"
  contentType="wizard"
  steps={[
    {
      index: 0,
      title: 'Kişisel Bilgiler',
      description: 'Önce sizi tanıyalım',
      fields: [
        { key: 'name', name: 'Ad', view: 'text-field', required: true },
        { key: 'email', name: 'Email', view: 'email-field', required: true },
      ],
      can_skip: false,
    },
    {
      index: 1,
      title: 'Tercihler',
      description: 'Tercihlerinizi belirleyin',
      fields: [
        { key: 'notifications', name: 'Bildirimler', view: 'switch-field' },
      ],
      can_skip: true,
    },
  ]}
  onChange={(data) => console.log('Wizard completed:', data)}
/>
```

---

## API Referansı

### Backend (Go)

#### DialogField Methods

```go
// Dialog oluşturma
Dialog(name string, key ...string) *DialogField

// Trigger ayarları
DefaultOpen(open bool) *DialogField                      // Varsayılan açık
TriggerButton(text string) *DialogField                  // Buton metni
TriggerIcon(icon string) *DialogField                    // Buton ikonu

// Content ayarları
Content(fields []core.Element) *DialogField              // Basit form
Wizard(steps []DialogStep) *DialogField                  // Multi-step wizard

// Dialog ayarları
DialogTitle(title string) *DialogField                   // Dialog başlığı
DialogDesc(desc string) *DialogField                     // Dialog açıklaması
DialogSize(size string) *DialogField                     // Dialog boyutu

// Callbacks
OnComplete(fn func(*fiber.Ctx, map[string]any) error) *DialogField
OnSkip(fn func(*fiber.Ctx) error) *DialogField

// Visibility (Schema'dan miras)
CanSee(fn VisibilityFunc) *DialogField
OnlyOnForm() *DialogField
OnlyOnList() *DialogField
OnlyOnDetail() *DialogField
```

#### DialogStep Struct

```go
type DialogStep struct {
    Index       int            // Adım sırası (0'dan başlar)
    Title       string         // Adım başlığı
    Description string         // Adım açıklaması
    Fields      []core.Element // Adımdaki field'lar
    CanSkip     bool           // Adım atlanabilir mi?
}
```

#### DialogContentType

```go
const (
    DialogContentForm   DialogContentType = "form"    // Basit form
    DialogContentWizard DialogContentType = "wizard"  // Multi-step wizard
)
```

### Frontend (TypeScript)

#### DialogFieldProps

```typescript
interface DialogFieldProps {
  // Field props (FieldRenderer'dan gelir)
  name: string;
  label: string;
  value?: Record<string, any>;
  onChange?: (value: Record<string, any>) => void;
  error?: string;
  disabled?: boolean;
  required?: boolean;
  helpText?: string;
  className?: string;

  // DialogField özel props
  defaultOpen?: boolean;
  triggerButton?: string;
  triggerIcon?: string;
  contentType: 'form' | 'wizard';
  fields?: FieldDefinition[];
  steps?: WizardStep[];
  dialogTitle?: string;
  dialogDesc?: string;
  dialogSize?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
}
```

#### WizardStep

```typescript
interface WizardStep {
  index: number;
  title: string;
  description?: string;
  fields: FieldDefinition[];
  can_skip: boolean;
}
```

---

## Kullanım Senaryoları

### 1. Profil Tamamlama

Kullanıcı eksik bilgileri varsa sayfa açıldığında otomatik açılan dialog:

```go
fields.Dialog("Profil Tamamla", "profile_completion").
    DefaultOpen(true).
    DialogTitle("Profilinizi Tamamlayın").
    DialogDesc("Hesabınızı kullanmaya devam etmek için lütfen bilgilerinizi tamamlayın").
    Content([]core.Element{
        fields.Text("Telefon", "phone").Required(),
        fields.Text("Adres", "address").Required(),
    }).
    CanSee(func(ctx *core.ResourceContext) bool {
        // Sadece profili eksik kullanıcılara göster
        user := ctx.User
        return user.Phone == "" || user.Address == ""
    })
```

### 2. Onboarding Wizard

Yeni kullanıcılar için adım adım kurulum:

```go
fields.Dialog("Onboarding", "onboarding").
    TriggerButton("Kurulumu Başlat").
    DialogTitle("Hoş Geldiniz!").
    Wizard([]fields.DialogStep{
        {
            Index:   0,
            Title:   "Hesap Bilgileri",
            Fields:  []core.Element{/* ... */},
            CanSkip: false,
        },
        {
            Index:   1,
            Title:   "Tercihler",
            Fields:  []core.Element{/* ... */},
            CanSkip: true,
        },
    }).
    CanSee(func(ctx *core.ResourceContext) bool {
        // Sadece onboarding tamamlanmamış kullanıcılara göster
        user := ctx.User
        return !user.OnboardingCompleted
    })
```

### 3. Hızlı Eylem Dialog'u

Buton ile tetiklenen hızlı işlem formu:

```go
fields.Dialog("Hızlı Not", "quick_note").
    TriggerButton("Not Ekle").
    TriggerIcon("📝").
    DialogTitle("Hızlı Not Ekle").
    DialogSize("sm").
    Content([]core.Element{
        fields.Text("Başlık", "title").Required(),
        fields.Textarea("İçerik", "content").Required(),
    }).
    OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
        return createNote(data)
    })
```

### 4. Multi-Step Form (Karmaşık Veri Girişi)

Karmaşık veri girişi için wizard:

```go
fields.Dialog("Ürün Ekle", "add_product").
    TriggerButton("Yeni Ürün").
    DialogTitle("Ürün Ekle").
    DialogSize("lg").
    Wizard([]fields.DialogStep{
        {
            Index:   0,
            Title:   "Temel Bilgiler",
            Fields: []core.Element{
                fields.Text("Ürün Adı", "name").Required(),
                fields.Textarea("Açıklama", "description"),
                fields.Number("Fiyat", "price").Required(),
            },
            CanSkip: false,
        },
        {
            Index:   1,
            Title:   "Stok ve Kategori",
            Fields: []core.Element{
                fields.Number("Stok", "stock").Required(),
                fields.Select("Kategori", "category_id")./* ... */,
            },
            CanSkip: false,
        },
        {
            Index:   2,
            Title:   "Görseller",
            Fields: []core.Element{
                fields.Image("Ana Görsel", "main_image"),
                // fields.Images("Galeri", "gallery"),
            },
            CanSkip: true,
        },
    })
```

---

## Best Practices

### 1. Dialog Boyutu Seçimi

- **sm**: Basit formlar (2-3 field)
- **md**: Orta formlar (4-6 field) - varsayılan
- **lg**: Wizard veya karmaşık formlar
- **xl**: Çok karmaşık formlar
- **full**: Tam ekran gerekli formlar

### 2. Wizard Adım Sayısı

- **Optimal**: 2-4 adım
- **Maksimum**: 5-6 adım
- Çok fazla adım kullanıcı deneyimini olumsuz etkiler

### 3. CanSkip Kullanımı

- İlk adım genellikle atlanamaz (`CanSkip: false`)
- Opsiyonel bilgiler için `CanSkip: true`
- Son adım genellikle atlanamaz

### 4. DefaultOpen Kullanımı

- Sadece kritik durumlarda kullanın
- Kullanıcı deneyimini bozabilir
- CanSee ile birlikte kullanarak sadece gerekli kullanıcılara gösterin

### 5. Callback Kullanımı

```go
// ✅ İyi: Hata kontrolü
OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
    if err := validateData(data); err != nil {
        return err
    }
    return saveData(data)
})

// ❌ Kötü: Hata kontrolü yok
OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
    saveData(data)
    return nil
})
```

---

## Troubleshooting

### Dialog Açılmıyor

**Sorun**: Dialog butonu görünüyor ama tıklandığında açılmıyor.

**Çözüm**:
1. Browser console'da hata var mı kontrol edin
2. DialogField'ın field registry'ye kayıtlı olduğundan emin olun
3. Frontend build'i yeniden yapın: `cd web && npm run build`

### Wizard Adımları Geçmiyor

**Sorun**: Wizard'da "İleri" butonuna tıklandığında sonraki adıma geçmiyor.

**Çözüm**:
1. Form validation hatası olabilir - required field'ları kontrol edin
2. Browser console'da hata var mı kontrol edin
3. UniversalResourceForm'un onSubmit callback'i çağrılıyor mu kontrol edin

### Data Kaydedilmiyor

**Sorun**: Dialog tamamlandığında data kaydedilmiyor.

**Çözüm**:
1. OnComplete callback'inin tanımlı olduğundan emin olun
2. Callback içinde hata dönüyor mu kontrol edin
3. Backend log'larını kontrol edin

### TypeScript Hataları

**Sorun**: Frontend build'de TypeScript hataları alıyorum.

**Çözüm**:
```bash
cd web
npm run build
```

Hatalar devam ediyorsa:
1. `web/src/types/dialog.ts` dosyasının var olduğundan emin olun
2. DialogField component'lerinin import edildiğinden emin olun
3. `web/src/components/forms/fields/index.ts` dosyasında DialogField'ın kayıtlı olduğundan emin olun

---

## Örnekler

### Tam Örnek: User Resource

```go
package resources

import (
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/core"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/gofiber/fiber/v2"
)

type UserResource struct {
    *resource.OptimizedBase
}

func (r *UserResource) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.ID("ID"),
        fields.Text("Ad", "name").Required(),
        fields.Email("Email", "email").Required(),
        fields.Tel("Telefon", "phone"),
        fields.Text("Adres", "address"),

        // Profil tamamlama dialog'u
        fields.Dialog("Profil Tamamla", "profile_completion").
            DefaultOpen(true).
            DialogTitle("Profilinizi Tamamlayın").
            DialogDesc("Lütfen eksik bilgilerinizi doldurun").
            Content([]core.Element{
                fields.Tel("Telefon", "phone").Required(),
                fields.Text("Adres", "address").Required(),
                fields.Date("Doğum Tarihi", "birth_date"),
            }).
            CanSee(func(ctx *core.ResourceContext) bool {
                user := ctx.User
                return user.Phone == "" || user.Address == ""
            }).
            OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
                userID := ctx.Locals("user_id").(uint)
                return db.Model(&User{}).
                    Where("id = ?", userID).
                    Updates(data).Error
            }),

        // Onboarding wizard
        fields.Dialog("Onboarding", "onboarding_wizard").
            TriggerButton("Başlangıç Rehberini Başlat").
            TriggerIcon("🚀").
            DialogTitle("Hoş Geldiniz!").
            DialogSize("lg").
            Wizard([]fields.DialogStep{
                {
                    Index:   0,
                    Title:   "Kişisel Bilgiler",
                    Fields: []core.Element{
                        fields.Text("Ad Soyad", "full_name").Required(),
                        fields.Email("Email", "email").Required(),
                    },
                    CanSkip: false,
                },
                {
                    Index:   1,
                    Title:   "Tercihler",
                    Fields: []core.Element{
                        fields.Switch("Email Bildirimleri", "email_notifications"),
                        fields.Switch("SMS Bildirimleri", "sms_notifications"),
                    },
                    CanSkip: true,
                },
            }).
            CanSee(func(ctx *core.ResourceContext) bool {
                user := ctx.User
                return !user.OnboardingCompleted
            }).
            OnComplete(func(ctx *fiber.Ctx, data map[string]any) error {
                userID := ctx.Locals("user_id").(uint)
                data["onboarding_completed"] = true
                return db.Model(&User{}).
                    Where("id = ?", userID).
                    Updates(data).Error
            }),
    }
}
```

---

## Changelog

### v1.0.0 (2026-02-08)

**Eklenen Özellikler:**
- ✅ DialogField backend implementasyonu
- ✅ DialogField frontend component'leri
- ✅ Basit form mode
- ✅ Multi-step wizard mode
- ✅ Progress indicator
- ✅ Skip functionality
- ✅ Özelleştirilebilir dialog boyutu
- ✅ UniversalResourceForm entegrasyonu
- ✅ Field registry entegrasyonu
- ✅ TypeScript type definitions
- ✅ Comprehensive documentation

---

## Lisans

Bu özellik Panel.go projesinin bir parçasıdır ve aynı lisans altında dağıtılır.

---

## Destek

Sorularınız veya sorunlarınız için:
- GitHub Issues: https://github.com/ferdiunal/panel.go/issues
- Documentation: https://panel.go/docs
