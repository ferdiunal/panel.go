# Tabs Field ve UI FileSystem Güncellemeleri

Bu dokümantasyon, Panel.go'ya eklenen yeni özellikler ve güncellemeleri açıklar:
1. **Tabs Component Enhancement**: Tab pozisyonlandırma ve variant desteği
2. **ResponsiveModal Enhancement**: Tam ekran modu
3. **TabsField Implementation**: Backend ve frontend tabs field sistemi
4. **UI FileSystem Update**: Plugin UI öncelik sistemi

## İçindekiler

- [Tabs Component Enhancement](#tabs-component-enhancement)
- [ResponsiveModal Enhancement](#responsivemodal-enhancement)
- [TabsField Implementation](#tabsfield-implementation)
- [UI FileSystem Update](#ui-filesystem-update)

---

## Tabs Component Enhancement

### Genel Bakış

Tabs component'ine `side` prop'u eklendi. Bu, tab'ların pozisyonunu kontrol etmeyi sağlar.

### Yeni Özellikler

**Side Prop:**
- `top`: Tab'lar üstte (varsayılan)
- `bottom`: Tab'lar altta
- `left`: Tab'lar solda
- `right`: Tab'lar sağda

**Backward Compatibility:**
- Mevcut `orientation` prop hala çalışır
- `orientation="horizontal"` → `side="top"` (varsayılan)
- `orientation="vertical"` → `side="left"` (varsayılan)

### Kullanım Örnekleri

```tsx
// Top (varsayılan)
<Tabs defaultValue="tab1">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
    <TabsTrigger value="tab2">Tab 2</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
  <TabsContent value="tab2">Content 2</TabsContent>
</Tabs>

// Bottom
<Tabs defaultValue="tab1" side="bottom">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
</Tabs>

// Left
<Tabs defaultValue="tab1" side="left">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
</Tabs>

// Right
<Tabs defaultValue="tab1" side="right">
  <TabsList>
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
</Tabs>
```

### Variant Desteği

TabsList component'i `variant` prop'unu destekler:
- `default`: Arka plan ve padding ile standart görünüm
- `line`: Sadece alt çizgi ile minimal görünüm

```tsx
<Tabs defaultValue="tab1">
  <TabsList variant="line">
    <TabsTrigger value="tab1">Tab 1</TabsTrigger>
  </TabsList>
  <TabsContent value="tab1">Content 1</TabsContent>
</Tabs>
```

### Teknik Detaylar

**Dosya:** `/web/src/components/ui/tabs.tsx`

**Props Interface:**
```typescript
interface TabsProps extends React.ComponentProps<typeof TabsPrimitive.Root> {
  side?: "top" | "bottom" | "left" | "right"
  orientation?: "horizontal" | "vertical" // deprecated
}
```

**CSS Logic:**
- `side="top"`: `flex-col` (TabsList üstte, content altta)
- `side="bottom"`: `flex-col-reverse` (content üstte, TabsList altta)
- `side="left"`: `flex-row` (TabsList solda, content sağda)
- `side="right"`: `flex-row-reverse` (content solda, TabsList sağda)

**Line Variant Indicator:**
- `side="top"`: Alt çizgi (bottom-[-5px])
- `side="bottom"`: Üst çizgi (top-[-5px])
- `side="left"`: Sağ çizgi (-right-1)
- `side="right"`: Sol çizgi (-left-1)

---

## ResponsiveModal Enhancement

### Genel Bakış

ResponsiveModal component'ine tam ekran modu eklendi. Kullanıcılar modal'ı tam ekran yapabilir.

### Yeni Özellikler

**Fullscreen Props:**
- `defaultFullscreen`: Uncontrolled mode için varsayılan tam ekran durumu
- `fullscreen`: Controlled mode için tam ekran durumu
- `onFullscreenChange`: Tam ekran durumu değiştiğinde callback
- `showFullscreenButton`: Tam ekran butonunu göster/gizle (varsayılan: true)

**Fullscreen Button:**
- Header'da close button'un solunda konumlanır
- Maximize2/Minimize2 icon ile toggle
- Desktop'ta gösterilir, mobile'da (drawer) gösterilmez

### Kullanım Örnekleri

```tsx
// Uncontrolled mode
<ResponsiveModal
  title="Modal Başlığı"
  defaultFullscreen={false}
  onFullscreenChange={(fullscreen) => console.log(fullscreen)}
>
  <div>Modal içeriği</div>
</ResponsiveModal>

// Controlled mode
const [isFullscreen, setIsFullscreen] = useState(false)
<ResponsiveModal
  title="Modal Başlığı"
  fullscreen={isFullscreen}
  onFullscreenChange={setIsFullscreen}
>
  <div>Modal içeriği</div>
</ResponsiveModal>

// Fullscreen button'u gizle
<ResponsiveModal
  title="Modal Başlığı"
  showFullscreenButton={false}
>
  <div>Modal içeriği</div>
</ResponsiveModal>
```

### Teknik Detaylar

**Dosya:** `/web/src/components/ui/responsive-modal.tsx`

**Props Interface:**
```typescript
interface ResponsiveModalProps {
  // ... mevcut props
  defaultFullscreen?: boolean
  fullscreen?: boolean
  onFullscreenChange?: (fullscreen: boolean) => void
  showFullscreenButton?: boolean
}
```

**State Management:**
```typescript
const [internalFullscreen, setInternalFullscreen] = React.useState(defaultFullscreen)
const isFullscreen = fullscreen !== undefined ? fullscreen : internalFullscreen

const toggleFullscreen = () => {
  const newValue = !isFullscreen
  if (fullscreen === undefined) {
    setInternalFullscreen(newValue)
  }
  onFullscreenChange?.(newValue)
}
```

**CSS Classes:**
- **Dialog Fullscreen**: `w-screen h-screen max-w-none rounded-none`
- **Dialog Normal**: `max-w-[calc(100%-2rem)] sm:max-w-md`
- **Sheet Fullscreen**: `w-screen h-screen max-w-none`
- **Sheet Normal**: `data-[side=right]:w-3/4 data-[side=right]:sm:max-w-sm`
- **Transition**: `transition-all duration-200 ease-in-out`

---

## TabsField Implementation

### Genel Bakış

Panel.go'ya yeni bir field tipi eklendi: **TabsField**. Bu field, alanları tab'lara ayırmak için kullanılır.

### Kullanım Senaryoları

- **Çoklu Dil Desteği**: Türkçe, İngilizce, vb. tab'ları ile çeviri alanları
- **Kategorize Edilmiş Formlar**: Genel Bilgiler, Adres, İletişim tab'ları
- **Karmaşık Form Organizasyonu**: Uzun formları mantıksal bölümlere ayırma
- **İlgili Alan Grupları**: Benzer alanları bir arada gösterme

### Backend Implementation (Go)

#### 1. Type Constants

**Dosya:** `/pkg/core/types.go`

```go
TYPE_TABS ElementType = "tabs"
```

**Dosya:** `/pkg/fields/enum.go`

```go
TYPE_TABS ElementType = core.TYPE_TABS
```

#### 2. TabsField Struct

**Dosya:** `/pkg/fields/tabs.go`

```go
// Tab, bir tab'ın yapısını temsil eder
type Tab struct {
    Value  string         `json:"value"`  // Benzersiz tanımlayıcı
    Label  string         `json:"label"`  // Görünen ad
    Fields []core.Element `json:"fields"` // Tab içindeki field'lar
}

// TabsField, alanları tab'lara ayırmak için bir konteyner temsil eder
type TabsField struct {
    *Schema
    Tabs []Tab
}
```

#### 3. Constructor ve Methods

```go
// Tabs, yeni bir tabs konteyner oluşturur
func Tabs(title string) *TabsField {
    schema := NewField(title)
    schema.View = "tabs-field"
    schema.Type = TYPE_TABS

    return &TabsField{
        Schema: schema,
        Tabs:   []Tab{},
    }
}

// AddTab, tabs konteyner'a yeni bir tab ekler
func (t *TabsField) AddTab(value, label string, fields ...core.Element) *TabsField {
    t.Tabs = append(t.Tabs, Tab{
        Value:  value,
        Label:  label,
        Fields: fields,
    })
    return t
}

// WithSide, tab'ların pozisyonunu ayarlar
func (t *TabsField) WithSide(side string) *TabsField {
    t.Props["side"] = side
    return t
}

// WithVariant, tab'ların görünüm stilini ayarlar
func (t *TabsField) WithVariant(variant string) *TabsField {
    t.Props["variant"] = variant
    return t
}

// WithDefaultTab, varsayılan aktif tab'ı ayarlar
func (t *TabsField) WithDefaultTab(value string) *TabsField {
    t.Props["defaultTab"] = value
    return t
}

// GetFields, tüm tab'lardaki alanları döndürür
func (t *TabsField) GetFields() []core.Element {
    var fields []core.Element
    for _, tab := range t.Tabs {
        fields = append(fields, tab.Fields...)
    }
    return fields
}
```

#### 4. Kullanım Örneği

```go
package main

import (
    "github.com/ferdiunal/panel.go/pkg/fields"
)

func ProductFieldResolver() []core.Element {
    return []core.Element{
        fields.ID(),

        // Çoklu dil desteği için tabs
        fields.Tabs("Ürün Bilgileri").
            AddTab("tr", "Türkçe",
                fields.Text("Başlık", "title_tr").Required(),
                fields.Textarea("Açıklama", "description_tr"),
                fields.RichText("İçerik", "content_tr"),
            ).
            AddTab("en", "English",
                fields.Text("Title", "title_en").Required(),
                fields.Textarea("Description", "description_en"),
                fields.RichText("Content", "content_en"),
            ).
            WithSide("top").
            WithVariant("line").
            WithDefaultTab("tr"),

        // Kategorize edilmiş form alanları
        fields.Tabs("Ürün Detayları").
            AddTab("general", "Genel Bilgiler",
                fields.Text("SKU", "sku"),
                fields.Number("Fiyat", "price"),
                fields.Number("Stok", "stock"),
            ).
            AddTab("seo", "SEO",
                fields.Text("Meta Başlık", "meta_title"),
                fields.Textarea("Meta Açıklama", "meta_description"),
                fields.Text("Slug", "slug"),
            ).
            AddTab("images", "Görseller",
                fields.Image("Ana Görsel", "main_image"),
                fields.Image("Galeri", "gallery"),
            ).
            WithSide("left").
            WithDefaultTab("general"),
    }
}
```

### Frontend Implementation (TypeScript/React)

#### 1. Field Components

**Form Variant:** `/web/src/components/fields/form/TabsField.tsx`

```tsx
export const TabsFormField: React.FC<FormFieldProps> = ({
  field,
  name,
  label,
  error,
  disabled = false,
  required = false,
  helpText,
}) => {
  const tabs = (field.props?.tabs || []) as Array<{
    value: string;
    label: string;
    fields?: any[];
  }>;
  const side = (field.props?.side || 'top') as "top" | "bottom" | "left" | "right";
  const defaultTab = (field.props?.defaultTab || tabs[0]?.value || '') as string;

  return (
    <FieldLayout
      name={name}
      label={label}
      error={error}
      required={required}
      helpText={helpText}
      disabled={disabled}
    >
      <Tabs defaultValue={defaultTab} side={side}>
        <TabsList>
          {tabs.map((tab) => (
            <TabsTrigger key={tab.value} value={tab.value}>
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>
        {tabs.map((tab) => (
          <TabsContent key={tab.value} value={tab.value}>
            <div className="space-y-4">
              {/* Tab içindeki field'lar render edilir */}
            </div>
          </TabsContent>
        ))}
      </Tabs>
    </FieldLayout>
  );
};
```

**Detail Variant:** `/web/src/components/fields/detail/TabsField.tsx`

```tsx
export const TabsDetailField: React.FC<DetailFieldProps> = ({ field }) => {
  const tabs = field.props?.tabs || [];
  const side = field.props?.side || 'top';
  const defaultTab = field.props?.defaultTab || tabs[0]?.value;

  return (
    <FieldLayout
      name={field.key}
      label={field.name || field.label}
      helpText={field.help_text}
    >
      <Tabs defaultValue={defaultTab} side={side}>
        <TabsList>
          {tabs.map((tab: any) => (
            <TabsTrigger key={tab.value} value={tab.value}>
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>
        {tabs.map((tab: any) => (
          <TabsContent key={tab.value} value={tab.value}>
            <div className="space-y-4">
              {/* Tab içindeki field'lar render edilir */}
            </div>
          </TabsContent>
        ))}
      </Tabs>
    </FieldLayout>
  );
};
```

**Index Variant:** `/web/src/components/fields/index/TabsField.tsx`

```tsx
export const TabsIndexField: React.FC<IndexFieldProps> = ({ field }) => {
  const tabs = field.props?.tabs || [];

  return (
    <FieldLayout
      name={field.key}
      label={field.name || field.label}
      hideLabel={true}
    >
      <div className="text-sm text-muted-foreground">
        {tabs.length > 0 ? `${tabs.length} tab` : '—'}
      </div>
    </FieldLayout>
  );
};
```

#### 2. Field Registry

**Dosya:** `/web/src/components/forms/fields/index.ts`

```typescript
// Import
import { TabsFormField } from '@/components/fields/form/TabsField';
import { TabsDetailField } from '@/components/fields/detail/TabsField';
import { TabsIndexField } from '@/components/fields/index/TabsField';

// Memoize
export const MemoizedTabsField = React.memo(TabsFormField);

// Register
export function registerAllFields() {
  // ... diğer field'lar

  // Tabs field
  fieldRegistry.register('tabs', MemoizedTabsField as any);
  fieldRegistry.register('tabs-field', MemoizedTabsField as any);
  fieldRegistry.register('tabs-field-form', TabsFormField as any);
  fieldRegistry.register('tabs-field-detail', TabsDetailField as any);
  fieldRegistry.register('tabs-field-index', TabsIndexField as any);
}
```

### Özellikler

**Backend:**
- ✅ Type-safe field definition
- ✅ Method chaining API
- ✅ Tab pozisyonlandırma (top, bottom, left, right)
- ✅ Tab variant (default, line)
- ✅ Varsayılan tab seçimi
- ✅ GetFields() ile tüm field'ları toplama

**Frontend:**
- ✅ Form, Detail, Index variant'ları
- ✅ FieldLayout pattern ile tutarlı UI
- ✅ Tabs component entegrasyonu
- ✅ Type-safe props interface
- ✅ Field registry entegrasyonu

---

## UI FileSystem Update

### Genel Bakış

Panel.go'nun UI dosya sistemi güncellendi. Artık plugin UI'ları otomatik olarak algılanır ve kullanılır.

### Sorun

**Önceki Durum:**
- UI dosyaları binary'ye gömülü (embedded)
- Plugin sistemi UI build'i `assets/ui/` dizinine kopyalıyor
- Ama `GetFileSystem()` fonksiyonu sadece embedded UI'ı kullanıyordu
- Plugin UI'ları kullanılamıyordu

### Çözüm

**Yeni Durum:**
- `GetFileSystem()` fonksiyonu önce `assets/ui/` dizinini kontrol eder
- Dizin varsa plugin UI'ı kullanır
- Yoksa embedded UI'a fallback yapar
- Kullanıcı hiçbir config değişikliği yapmadan plugin UI'ı kullanabilir

### Dosya Sistemi Öncelik Sırası

**Dosya:** `/pkg/panel/assets.go`

```go
func GetFileSystem(useEmbed bool) (fs.FS, error) {
    // 1. Önce assets/ui/ dizinini kontrol et (plugin UI)
    if _, err := os.Stat("assets/ui"); err == nil {
        return os.DirFS("assets/ui"), nil
    }

    // 2. Embedded UI kullan (varsayılan Panel.go UI)
    if useEmbed {
        return fs.Sub(embedFS, "ui")
    }

    // 3. Fallback: Geliştirme ortamı
    return nil, nil
}
```

### Öncelik Sırası

1. **Plugin UI** (`assets/ui/`):
   - Plugin sistemi kullanıldığında UI build bu dizine kopyalanır
   - Custom field'lar ve plugin UI'ları içerir
   - Öncelikli olarak kontrol edilir

2. **Embedded UI** (`pkg/panel/ui/`):
   - Binary'ye gömülü varsayılan Panel.go UI
   - Plugin UI yoksa otomatik olarak kullanılır
   - Harici dosya bağımlılığı yok

3. **Geliştirme Modu** (`useEmbed=false`):
   - nil döndürür
   - Çağıran taraf `os.DirFS()` kullanarak disk'ten yükler
   - Sıcak yenileme (hot reload) için uygun

### Kullanım Senaryoları

#### Senaryo 1: Plugin Kullanımı (Üretim)

```go
// Plugin oluşturulduğunda UI build assets/ui/ dizinine kopyalanır
fs, err := panel.GetFileSystem(true)
// fs -> assets/ui/ (plugin UI)
```

#### Senaryo 2: Varsayılan Panel (Üretim)

```go
// assets/ui/ dizini yoksa, embedded UI kullanılır
fs, err := panel.GetFileSystem(true)
// fs -> pkg/panel/ui/ (embedded)
```

#### Senaryo 3: Geliştirme Ortamı

```go
fs, err := panel.GetFileSystem(false)
if err != nil || fs == nil {
    fs = os.DirFS("./pkg/panel/ui")
}
```

### Plugin Workflow

**1. Plugin Oluşturma:**
```bash
cd examples/cargo.go
panel plugin create importer
```

**2. UI Build:**
```bash
# Plugin CLI otomatik olarak UI build alır
# Build output: assets/ui/
```

**3. Uygulama Başlatma:**
```go
// GetFileSystem() otomatik olarak assets/ui/ dizinini algılar
fs, err := panel.GetFileSystem(true)
// fs -> assets/ui/ (plugin UI ile custom field'lar)
```

### Avantajlar

- ✅ **Otomatik Algılama**: Plugin UI otomatik olarak algılanır
- ✅ **Zero Config**: Kullanıcı hiçbir config değişikliği yapmaz
- ✅ **Fallback Desteği**: Embedded UI her zaman fallback olarak çalışır
- ✅ **Plugin Desteği**: Custom field'lar ve plugin UI'ları sorunsuz çalışır
- ✅ **Backward Compatible**: Mevcut projeler etkilenmez

### Teknik Detaylar

**Import Güncellemesi:**
```go
import (
    "embed"
    "io/fs"
    "os"  // Yeni eklendi
)
```

**Dosya Kontrolü:**
```go
if _, err := os.Stat("assets/ui"); err == nil {
    // Dizin var, plugin UI kullan
    return os.DirFS("assets/ui"), nil
}
```

**Performans:**
- Plugin UI (assets/ui/): O(n) erişim süresi (disk I/O'ya bağlı)
- Embedded UI: O(1) erişim süresi (bellekten)
- Dizin kontrolü: O(1) (os.Stat çok hızlı)

---

## Özet

### Yapılan Değişiklikler

**1. Tabs Component:**
- ✅ `side` prop eklendi (top, bottom, left, right)
- ✅ Backward compatibility korundu
- ✅ Line variant indicator pozisyonları güncellendi

**2. ResponsiveModal:**
- ✅ Fullscreen mode eklendi
- ✅ Controlled ve uncontrolled mode desteği
- ✅ Fullscreen button eklendi
- ✅ Smooth animation

**3. TabsField:**
- ✅ Backend implementation (Go)
- ✅ Frontend implementation (TypeScript/React)
- ✅ Field registry entegrasyonu
- ✅ Form, Detail, Index variant'ları

**4. UI FileSystem:**
- ✅ Plugin UI öncelik sistemi
- ✅ Otomatik algılama
- ✅ Embedded UI fallback
- ✅ Zero config

### Değiştirilen Dosyalar

**Backend:**
- `/pkg/core/types.go`
- `/pkg/fields/enum.go`
- `/pkg/fields/tabs.go` (yeni)
- `/pkg/panel/assets.go`

**Frontend:**
- `/web/src/components/ui/tabs.tsx`
- `/web/src/components/ui/responsive-modal.tsx`
- `/web/src/components/fields/form/TabsField.tsx` (yeni)
- `/web/src/components/fields/detail/TabsField.tsx` (yeni)
- `/web/src/components/fields/index/TabsField.tsx` (yeni)
- `/web/src/components/forms/fields/index.ts`

### Kullanıcı Etkisi

**Kullanıcılar için:**
- ✅ Yeni TabsField ile çoklu dil desteği kolaylaştı
- ✅ Plugin UI'ları otomatik olarak çalışır
- ✅ Config değişikliği gerekmez
- ✅ Mevcut projeler etkilenmez

**Geliştiriciler için:**
- ✅ Type-safe field definition
- ✅ Method chaining API
- ✅ Tutarlı field pattern
- ✅ Plugin sistemi ile entegrasyon

---

## İlgili Dokümantasyon

- [Plugin Sistemi](PLUGIN_SYSTEM.md)
- [Field'lar](Fields.md)
- [UI Component'leri](../web/README.md)

---

**Son Güncelleme:** 2026-02-13
**Versiyon:** 1.0.0
