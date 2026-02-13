package auth

import "github.com/ferdiunal/panel.go/pkg/context"

// Bu interface, uygulamada yetkilendirme (authorization) politikalarını tanımlar.
//
// Policy interface'i, kullanıcıların belirli kaynaklar (resources) üzerinde
// gerçekleştirebilecekleri işlemleri kontrol etmek için kullanılır. RBAC (Role-Based
// Access Control) veya ABAC (Attribute-Based Access Control) gibi yetkilendirme
// modellerinin uygulanmasında temel rol oynar.
//
// # Kullanım Senaryoları
//
// 1. **Kaynak Erişim Kontrolü**: Kullanıcıların belirli kaynakları görüp göremeyeceğini
//    kontrol etmek için kullanılır.
//
// 2. **İşlem Yetkilendirmesi**: Oluşturma, güncelleme, silme gibi işlemlerin
//    gerçekleştirilip gerçekleştirilemeyeceğini belirler.
//
// 3. **Rol Tabanlı Erişim**: Farklı roller (admin, editor, viewer) için farklı
//    yetkilendirme kuralları tanımlanabilir.
//
// # Örnek Kullanım
//
//	type UserPolicy struct{}
//
//	func (p *UserPolicy) ViewAny(ctx *context.Context) bool {
//	    return ctx.User.IsAuthenticated()
//	}
//
//	func (p *UserPolicy) View(ctx *context.Context, model interface{}) bool {
//	    user := ctx.User
//	    targetUser := model.(*User)
//	    return user.ID == targetUser.ID || user.IsAdmin()
//	}
//
//	func (p *UserPolicy) Create(ctx *context.Context) bool {
//	    return ctx.User.HasPermission("users.create")
//	}
//
//	func (p *UserPolicy) Update(ctx *context.Context, model interface{}) bool {
//	    user := ctx.User
//	    targetUser := model.(*User)
//	    return user.ID == targetUser.ID || user.IsAdmin()
//	}
//
//	func (p *UserPolicy) Delete(ctx *context.Context, model interface{}) bool {
//	    return ctx.User.IsAdmin()
//	}
//
// # Önemli Notlar
//
// - Her metod, yetkilendirme kararını hızlı bir şekilde vermeli (performans önemli)
// - Context'ten kullanıcı bilgisi ve diğer gerekli veriler alınmalı
// - Metotlar thread-safe olmalı
// - Yetkilendirme kararları loglama yapılmalı (audit trail)
// - Veritabanı sorguları minimize edilmeli (cache kullanılabilir)
type Policy interface {
	// Bu metod, belirli bir kaynak türünün herhangi bir örneğini görüntüleme
	// yetkisini kontrol eder. Genellikle liste sayfalarında kullanılır.
	//
	// # Parametreler
	//
	// - ctx: İstek bağlamı, kullanıcı bilgisi ve diğer veriler içerir
	//
	// # Dönüş Değeri
	//
	// - true: Kullanıcı bu kaynak türünün herhangi bir örneğini görebilir
	// - false: Kullanıcı bu kaynak türünü görüntülemek için yetkilendirilmemiş
	//
	// # Kullanım Senaryosu
	//
	// Kullanıcı, "Kullanıcılar" listesini açmaya çalıştığında, ViewAny() çağrılır.
	// Eğer false dönerse, liste sayfasına erişim reddedilir.
	//
	// # Örnek
	//
	//	if !policy.ViewAny(ctx) {
	//	    return errors.New("bu listeyi görüntülemek için yetkiniz yok")
	//	}
	ViewAny(ctx *context.Context) bool

	// Bu metod, belirli bir kaynak örneğini (model) görüntüleme yetkisini kontrol eder.
	// Genellikle detay sayfalarında kullanılır.
	//
	// # Parametreler
	//
	// - ctx: İstek bağlamı, kullanıcı bilgisi ve diğer veriler içerir
	// - model: Görüntülenmek istenen kaynak örneği (interface{} olarak geçilir)
	//
	// # Dönüş Değeri
	//
	// - true: Kullanıcı bu spesifik kaynağı görebilir
	// - false: Kullanıcı bu kaynağı görüntülemek için yetkilendirilmemiş
	//
	// # Kullanım Senaryosu
	//
	// Kullanıcı, ID=5 olan bir kullanıcının detaylarını görmek istediğinde,
	// View() çağrılır. Eğer kullanıcı kendi profilini görüyorsa veya admin ise
	// true döner, aksi takdirde false döner.
	//
	// # Örnek
	//
	//	user := &User{ID: 5, Name: "John"}
	//	if !policy.View(ctx, user) {
	//	    return errors.New("bu kaynağı görüntülemek için yetkiniz yok")
	//	}
	View(ctx *context.Context, model interface{}) bool

	// Bu metod, yeni bir kaynak oluşturma yetkisini kontrol eder.
	// Genellikle "Yeni Oluştur" sayfasında ve form gönderildiğinde kullanılır.
	//
	// # Parametreler
	//
	// - ctx: İstek bağlamı, kullanıcı bilgisi ve diğer veriler içerir
	//
	// # Dönüş Değeri
	//
	// - true: Kullanıcı bu kaynak türünün yeni bir örneğini oluşturabilir
	// - false: Kullanıcı yeni kaynak oluşturmak için yetkilendirilmemiş
	//
	// # Kullanım Senaryosu
	//
	// Kullanıcı, yeni bir ürün oluşturmak istediğinde, Create() çağrılır.
	// Eğer kullanıcı "ürün.oluştur" izni varsa true döner.
	//
	// # Örnek
	//
	//	if !policy.Create(ctx) {
	//	    return errors.New("yeni kaynak oluşturmak için yetkiniz yok")
	//	}
	//	newProduct := &Product{Name: "Yeni Ürün"}
	//	db.Create(newProduct)
	Create(ctx *context.Context) bool

	// Bu metod, mevcut bir kaynağı güncelleme yetkisini kontrol eder.
	// Genellikle düzenleme sayfasında ve form gönderildiğinde kullanılır.
	//
	// # Parametreler
	//
	// - ctx: İstek bağlamı, kullanıcı bilgisi ve diğer veriler içerir
	// - model: Güncellenecek kaynak örneği (interface{} olarak geçilir)
	//
	// # Dönüş Değeri
	//
	// - true: Kullanıcı bu kaynağı güncelleyebilir
	// - false: Kullanıcı bu kaynağı güncellemek için yetkilendirilmemiş
	//
	// # Kullanım Senaryosu
	//
	// Kullanıcı, ID=10 olan bir ürünü düzenlemek istediğinde, Update() çağrılır.
	// Eğer kullanıcı bu ürünün sahibi veya admin ise true döner.
	//
	// # Örnek
	//
	//	product := &Product{ID: 10, Name: "Ürün"}
	//	if !policy.Update(ctx, product) {
	//	    return errors.New("bu kaynağı güncellemek için yetkiniz yok")
	//	}
	//	product.Name = "Güncellenmiş Ürün"
	//	db.Save(product)
	Update(ctx *context.Context, model interface{}) bool

	// Bu metod, bir kaynağı silme yetkisini kontrol eder.
	// Genellikle silme işlemi gerçekleştirilmeden önce çağrılır.
	//
	// # Parametreler
	//
	// - ctx: İstek bağlamı, kullanıcı bilgisi ve diğer veriler içerir
	// - model: Silinecek kaynak örneği (interface{} olarak geçilir)
	//
	// # Dönüş Değeri
	//
	// - true: Kullanıcı bu kaynağı silebilir
	// - false: Kullanıcı bu kaynağı silmek için yetkilendirilmemiş
	//
	// # Kullanım Senaryosu
	//
	// Kullanıcı, bir ürünü silmek istediğinde, Delete() çağrılır.
	// Genellikle sadece admin veya kaynağın sahibi silme işlemini gerçekleştirebilir.
	//
	// # Örnek
	//
	//	product := &Product{ID: 10}
	//	if !policy.Delete(ctx, product) {
	//	    return errors.New("bu kaynağı silmek için yetkiniz yok")
	//	}
	//	db.Delete(product)
	//
	// # Uyarı
	//
	// Silme işlemi geri alınamaz olduğu için, Delete() metodu diğer metotlardan
	// daha katı yetkilendirme kuralları uygulamalıdır.
	Delete(ctx *context.Context, model interface{}) bool
}
