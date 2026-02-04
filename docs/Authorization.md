# Yetkilendirme (Authorization) Rehberi

Yetkilendirme, kullanıcıların hangi işlemleri yapabileceğini kontrol eder. Panel.go'da bu, Policy (Politika) sistemi aracılığıyla yapılır.

## Temel Kavramlar

### Policy Nedir?

Policy, bir resource (kaynak) üzerinde yapılabilecek işlemleri kontrol eden bir yapıdır. Her resource'un kendi politikası vardır.

### Desteklenen İşlemler

```go
type Policy interface {
	// Tüm kaynakları görme izni
	ViewAny(ctx *context.Context) bool

	// Belirli bir kaynağı görme izni
	View(ctx *context.Context, model any) bool

	// Kaynak oluşturma izni
	Create(ctx *context.Context) bool

	// Kaynak güncelleme izni
	Update(ctx *context.Context, model any) bool

	// Kaynak silme izni
	Delete(ctx *context.Context, model any) bool

	// Kaynak geri yükleme izni
	Restore(ctx *context.Context, model any) bool

	// Kaynak kalıcı silme izni
	ForceDelete(ctx *context.Context, model any) bool
}
```

## Basit Policy Örneği

Herkesin tüm işlemleri yapabileceği bir policy:

```go
type PostPolicy struct{}

func (p *PostPolicy) ViewAny(ctx *context.Context) bool {
	return true // Herkes yazıları görebilir
}

func (p *PostPolicy) View(ctx *context.Context, model any) bool {
	return true // Herkes yazıları görebilir
}

func (p *PostPolicy) Create(ctx *context.Context) bool {
	return true // Herkes yazı oluşturabilir
}

func (p *PostPolicy) Update(ctx *context.Context, model any) bool {
	return true // Herkes yazı güncelleyebilir
}

func (p *PostPolicy) Delete(ctx *context.Context, model any) bool {
	return true // Herkes yazı silebilir
}

func (p *PostPolicy) Restore(ctx *context.Context, model any) bool {
	return false // Geri yükleme desteklenmiyor
}

func (p *PostPolicy) ForceDelete(ctx *context.Context, model any) bool {
	return false // Kalıcı silme desteklenmiyor
}
```

## Gelişmiş Policy Örneği

Rol tabanlı yetkilendirme:

```go
type PostPolicy struct{}

// Yardımcı fonksiyonlar
func isAdmin(ctx *context.Context) bool {
	// Context'ten kullanıcı bilgisini al
	user := ctx.User()
	return user != nil && user.Role == "admin"
}

func isEditor(ctx *context.Context) bool {
	user := ctx.User()
	return user != nil && user.Role == "editor"
}

func isAuthor(ctx *context.Context, post *Post) bool {
	user := ctx.User()
	return user != nil && post.AuthorID == user.ID
}

// Policy metodları
func (p *PostPolicy) ViewAny(ctx *context.Context) bool {
	// Herkes yazıları görebilir
	return true
}

func (p *PostPolicy) View(ctx *context.Context, model any) bool {
	post := model.(*Post)
	// Yayınlanmış yazılar herkes görebilir
	if post.Status == "published" {
		return true
	}
	// Taslak yazıları sadece yazar ve admin görebilir
	return isAuthor(ctx, post) || isAdmin(ctx)
}

func (p *PostPolicy) Create(ctx *context.Context) bool {
	// Sadece editor ve admin yazı oluşturabilir
	return isEditor(ctx) || isAdmin(ctx)
}

func (p *PostPolicy) Update(ctx *context.Context, model any) bool {
	post := model.(*Post)
	// Yazar kendi yazısını güncelleyebilir
	if isAuthor(ctx, post) {
		return true
	}
	// Admin herhangi bir yazıyı güncelleyebilir
	return isAdmin(ctx)
}

func (p *PostPolicy) Delete(ctx *context.Context, model any) bool {
	post := model.(*Post)
	// Yazar kendi yazısını silebilir
	if isAuthor(ctx, post) {
		return true
	}
	// Admin herhangi bir yazıyı silebilir
	return isAdmin(ctx)
}

func (p *PostPolicy) Restore(ctx *context.Context, model any) bool {
	// Sadece admin geri yükleyebilir
	return isAdmin(ctx)
}

func (p *PostPolicy) ForceDelete(ctx *context.Context, model any) bool {
	// Sadece admin kalıcı silebilir
	return isAdmin(ctx)
}
```

## Kullanıcı Bilgisine Erişim

Context'ten kullanıcı bilgisine erişebilirsiniz:

```go
func (p *PostPolicy) Update(ctx *context.Context, model any) bool {
	// Kullanıcı bilgisini al
	user := ctx.User()
	
	if user == nil {
		return false // Giriş yapmamış
	}

	// Kullanıcı bilgilerini kullan
	post := model.(*Post)
	
	// Kullanıcı ID'si
	if post.AuthorID == user.ID {
		return true
	}

	// Kullanıcı rolü
	if user.Role == "admin" {
		return true
	}

	return false
}
```

## Örnek: Yorum Yönetimi

```go
type Comment struct {
	ID        string
	PostID    string
	AuthorID  string
	Content   string
	Status    string // pending, approved, rejected
	CreatedAt time.Time
}

type CommentPolicy struct{}

func (p *CommentPolicy) ViewAny(ctx *context.Context) bool {
	// Herkes onaylı yorumları görebilir
	return true
}

func (p *CommentPolicy) View(ctx *context.Context, model any) bool {
	comment := model.(*Comment)
	user := ctx.User()

	// Onaylı yorumlar herkes görebilir
	if comment.Status == "approved" {
		return true
	}

	// Yorum yazarı kendi yorumunu görebilir
	if user != nil && comment.AuthorID == user.ID {
		return true
	}

	// Admin tüm yorumları görebilir
	if user != nil && user.Role == "admin" {
		return true
	}

	return false
}

func (p *CommentPolicy) Create(ctx *context.Context) bool {
	// Giriş yapmış kullanıcılar yorum yapabilir
	return ctx.User() != nil
}

func (p *CommentPolicy) Update(ctx *context.Context, model any) bool {
	comment := model.(*Comment)
	user := ctx.User()

	// Yorum yazarı kendi yorumunu güncelleyebilir
	if user != nil && comment.AuthorID == user.ID {
		return true
	}

	// Admin herhangi bir yorumu güncelleyebilir
	if user != nil && user.Role == "admin" {
		return true
	}

	return false
}

func (p *CommentPolicy) Delete(ctx *context.Context, model any) bool {
	comment := model.(*Comment)
	user := ctx.User()

	// Yorum yazarı kendi yorumunu silebilir
	if user != nil && comment.AuthorID == user.ID {
		return true
	}

	// Admin herhangi bir yorumu silebilir
	if user != nil && user.Role == "admin" {
		return true
	}

	return false
}

func (p *CommentPolicy) Restore(ctx *context.Context, model any) bool {
	// Sadece admin geri yükleyebilir
	user := ctx.User()
	return user != nil && user.Role == "admin"
}

func (p *CommentPolicy) ForceDelete(ctx *context.Context, model any) bool {
	// Sadece admin kalıcı silebilir
	user := ctx.User()
	return user != nil && user.Role == "admin"
}
```

## İpuçları

1. **Güvenlik**: Her zaman context'i kontrol edin, nil olabilir
2. **Performans**: Policy metodlarında ağır veritabanı sorguları yapmayın
3. **Tutarlılık**: Aynı kuralları API ve UI'da uygulayın
4. **Loglama**: Yetkilendirme başarısızlıklarını loglayın
5. **Test**: Policy'leri kapsamlı şekilde test edin

## Örnek: Departman Tabanlı Yetkilendirme

```go
type EmployeePolicy struct{}

func (p *EmployeePolicy) ViewAny(ctx *context.Context) bool {
	user := ctx.User()
	// HR ve admin tüm çalışanları görebilir
	return user != nil && (user.Department == "HR" || user.Role == "admin")
}

func (p *EmployeePolicy) View(ctx *context.Context, model any) bool {
	employee := model.(*Employee)
	user := ctx.User()

	if user == nil {
		return false
	}

	// Kendi bilgilerini görebilir
	if employee.ID == user.ID {
		return true
	}

	// Aynı departmandaki yönetici görebilir
	if user.Role == "manager" && employee.Department == user.Department {
		return true
	}

	// HR ve admin görebilir
	if user.Department == "HR" || user.Role == "admin" {
		return true
	}

	return false
}

func (p *EmployeePolicy) Create(ctx *context.Context) bool {
	user := ctx.User()
	// Sadece HR ve admin oluşturabilir
	return user != nil && (user.Department == "HR" || user.Role == "admin")
}

func (p *EmployeePolicy) Update(ctx *context.Context, model any) bool {
	employee := model.(*Employee)
	user := ctx.User()

	if user == nil {
		return false
	}

	// Kendi bilgilerini güncelleyebilir
	if employee.ID == user.ID {
		return true
	}

	// HR ve admin güncelleyebilir
	return user.Department == "HR" || user.Role == "admin"
}

func (p *EmployeePolicy) Delete(ctx *context.Context, model any) bool {
	user := ctx.User()
	// Sadece admin silebilir
	return user != nil && user.Role == "admin"
}

func (p *EmployeePolicy) Restore(ctx *context.Context, model any) bool {
	user := ctx.User()
	return user != nil && user.Role == "admin"
}

func (p *EmployeePolicy) ForceDelete(ctx *context.Context, model any) bool {
	user := ctx.User()
	return user != nil && user.Role == "admin"
}
```

## Sonraki Adımlar

- [Başlangıç Rehberi](./Getting-Started.md) - Temel kurulum
- [Alanlar Rehberi](./Fields.md) - Alan tanımı
- [API Referansı](./API-Reference.md) - Tüm metodlar
