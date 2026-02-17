package user

import (
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
)

// ============================================================================
// Bu yapı, kullanıcı yönetimi için yetkilendirme (authorization) kurallarını
// belirler ve uygulamak için kullanılır.
//
// **Amaç:**
// - Kullanıcı listesi görüntüleme, oluşturma, güncelleme ve silme işlemleri
//   için yetki kontrolü sağlamak
// - Rol tabanlı erişim kontrolü (RBAC) uygulamak
// - Kendini silme gibi tehlikeli işlemleri engellemek
//
// **Yetki Kaynağı:**
// - Tüm yetkilendirme kuralları `examples/simple/permissions.toml` dosyasından
//   yönetilir
// - Yetkilendirme kontrolleri `appContext.Context.HasPermission()` metodu
//   aracılığıyla yapılır
//
// **Kullanım Senaryoları:**
// 1. Admin panelinde kullanıcı listesini görüntüleme
// 2. Yeni kullanıcı oluşturma işlemi
// 3. Mevcut kullanıcı bilgilerini güncelleme
// 4. Kullanıcı hesabını silme (kendini silme hariç)
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	ctx := appContext.NewContext(user, permissions)
//
//	// Listeyi görüntüleme yetkisi kontrolü
//	if policy.ViewAny(ctx) {
//		// Kullanıcı listesini göster
//	}
//
//	// Belirli kullanıcıyı silme yetkisi kontrolü
//	if policy.Delete(ctx, targetUser) {
//		// Kullanıcıyı sil
//	}
//
// **Önemli Notlar:**
// - Bu yapı, resource bazlı bir yetkilendirme sistemi uygulamaktadır
// - Tüm metotlar receiver olarak `UserPolicy` değerini alır (pointer değil)
// - Delete metodu, kendini silme işlemini engelleme mantığı içerir
// ============================================================================
type UserPolicy struct{}

// ============================================================================
// Bu metod, kimliği doğrulanmış kullanıcının tüm kullanıcı listesini
// görüntüleme yetkisine sahip olup olmadığını kontrol eder.
//
// **Parametreler:**
// - ctx (*appContext.Context): Kimliği doğrulanmış kullanıcının bağlamı ve
//   yetkileri içeren context nesnesi
//
// **Dönüş Değeri:**
// - bool: Kullanıcı "users.view_any" yetkisine sahipse true, aksi takdirde false
//
// **Kullanım Senaryoları:**
// - Kullanıcı yönetim panelinin erişim kontrolü
// - Tüm kullanıcıların listelendiği sayfaya erişim izni
// - Raporlama ve analitik sayfalarında kullanıcı verilerine erişim
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	if policy.ViewAny(ctx) {
//		users := getAllUsers()
//		displayUserList(users)
//	} else {
//		showAccessDeniedError()
//	}
//
// **Önemli Notlar:**
// - Bu metod, belirli bir kullanıcı kaydı için değil, genel liste erişimi için
//   kontrol yapar
// - "users.view_any" izni, genellikle admin veya yönetici rolüne verilir
// - Context nil ise panic oluşabilir, bu nedenle context'in geçerliliğini
//   önceden kontrol edin
// ============================================================================
func (p UserPolicy) ViewAny(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.view_any")
}

// ============================================================================
// Bu metod, kimliği doğrulanmış kullanıcının belirli bir kullanıcı kaydını
// görüntüleme yetkisine sahip olup olmadığını kontrol eder.
//
// **Parametreler:**
// - ctx (*appContext.Context): Kimliği doğrulanmış kullanıcının bağlamı ve
//   yetkileri içeren context nesnesi
// - model (any): Görüntülenecek kullanıcı kaydı (genellikle *domainUser.User)
//
// **Dönüş Değeri:**
// - bool: Kullanıcı "users.view" yetkisine sahipse true, aksi takdirde false
//
// **Kullanım Senaryoları:**
// - Belirli bir kullanıcının profil sayfasına erişim kontrolü
// - Kullanıcı detay bilgilerinin görüntülenmesi
// - Kullanıcı bilgilerinin API aracılığıyla alınması
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	targetUser := getUserByID(userID)
//	if policy.View(ctx, targetUser) {
//		displayUserDetails(targetUser)
//	} else {
//		showAccessDeniedError()
//	}
//
// **Önemli Notlar:**
// - model parametresi genellikle *domainUser.User türünde olmalıdır
// - Bu metod, model parametresini kontrol etmez, sadece genel "users.view"
//   yetkisini kontrol eder
// - Belirli kullanıcıya özel erişim kuralları için bu metodu genişletebilirsiniz
// ============================================================================
func (p UserPolicy) View(ctx *appContext.Context, model any) bool {
	return ctx.HasPermission("users.view")
}

// ============================================================================
// Bu metod, kimliği doğrulanmış kullanıcının yeni bir kullanıcı oluşturma
// yetkisine sahip olup olmadığını kontrol eder.
//
// **Parametreler:**
// - ctx (*appContext.Context): Kimliği doğrulanmış kullanıcının bağlamı ve
//   yetkileri içeren context nesnesi
//
// **Dönüş Değeri:**
// - bool: Kullanıcı "users.create" yetkisine sahipse true, aksi takdirde false
//
// **Kullanım Senaryoları:**
// - Yeni kullanıcı oluşturma formunun görüntülenmesi
// - Yeni kullanıcı kaydı oluşturma işleminin gerçekleştirilmesi
// - Toplu kullanıcı ekleme işlemleri
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	if policy.Create(ctx) {
//		// Yeni kullanıcı oluşturma formunu göster
//		showCreateUserForm()
//	} else {
//		showAccessDeniedError()
//	}
//
// **Önemli Notlar:**
// - Bu metod, yeni kullanıcı oluşturma işleminin başında kontrol edilmelidir
// - "users.create" izni, genellikle admin veya insan kaynakları yöneticisine
//   verilir
// - Yeni kullanıcı oluşturulurken ek doğrulamalar (email benzersizliği vb.)
//   yapılmalıdır
// ============================================================================
func (p UserPolicy) Create(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.create")
}

// ============================================================================
// Bu metod, kimliği doğrulanmış kullanıcının mevcut bir kullanıcı kaydını
// güncelleme yetkisine sahip olup olmadığını kontrol eder.
//
// **Parametreler:**
// - ctx (*appContext.Context): Kimliği doğrulanmış kullanıcının bağlamı ve
//   yetkileri içeren context nesnesi
// - model (any): Güncellenecek kullanıcı kaydı (genellikle *domainUser.User)
//
// **Dönüş Değeri:**
// - bool: Kullanıcı "users.update" yetkisine sahipse true, aksi takdirde false
//
// **Kullanım Senaryoları:**
// - Kullanıcı profil bilgilerinin güncellenmesi
// - Kullanıcı rolü veya izinlerinin değiştirilmesi
// - Kullanıcı hesabı durumunun (aktif/pasif) değiştirilmesi
// - Kullanıcı iletişim bilgilerinin güncellenmesi
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	targetUser := getUserByID(userID)
//	if policy.Update(ctx, targetUser) {
//		// Kullanıcı güncelleme formunu göster
//		showUpdateUserForm(targetUser)
//	} else {
//		showAccessDeniedError()
//	}
//
// **Önemli Notlar:**
// - model parametresi genellikle *domainUser.User türünde olmalıdır
// - Bu metod, model parametresini kontrol etmez, sadece genel "users.update"
//   yetkisini kontrol eder
// - Belirli kullanıcıya özel güncelleme kuralları için bu metodu genişletebilirsiniz
// - Kendi profilini güncellemek için ayrı bir yetki tanımlanabilir
// ============================================================================
func (p UserPolicy) Update(ctx *appContext.Context, model any) bool {
	return ctx.HasPermission("users.update")
}

// ============================================================================
// Bu metod, kimliği doğrulanmış kullanıcının bir kullanıcı kaydını silme
// yetkisine sahip olup olmadığını kontrol eder. Özellikle, kendini silme
// işlemini engeller.
//
// **Parametreler:**
// - ctx (*appContext.Context): Kimliği doğrulanmış kullanıcının bağlamı ve
//   yetkileri içeren context nesnesi
// - model (any): Silinecek kullanıcı kaydı (genellikle *domainUser.User)
//
// **Dönüş Değeri:**
// - bool: Kullanıcı silme yetkisine sahipse ve kendini silmiyorsa true,
//   aksi takdirde false
//
// **Kullanım Senaryoları:**
// - Kullanıcı hesabının silinmesi
// - Eski veya deaktif kullanıcıların sistemden kaldırılması
// - Veri temizleme ve bakım işlemleri
//
// **Örnek Kullanım:**
//
//	policy := UserPolicy{}
//	targetUser := getUserByID(userID)
//	if policy.Delete(ctx, targetUser) {
//		// Kullanıcıyı sil
//		deleteUser(targetUser)
//		showSuccessMessage("Kullanıcı başarıyla silindi")
//	} else {
//		showAccessDeniedError()
//	}
//
// **Kontrol Akışı:**
// 1. model nil ise true döner (genel yetki kontrolü)
// 2. model *domainUser.User türüne dönüştürülür
// 3. Dönüştürme başarısız ise false döner
// 4. ctx nil ise false döner
// 5. ctx.User() nil ise false döner
// 6. Silinecek kullanıcı ID'si, kimliği doğrulanmış kullanıcı ID'sine eşitse
//    false döner (kendini silme engellenir)
// 7. Tüm kontroller geçerse true döner
//
// **Önemli Notlar:**
// - **Kendini Silme Engelleme**: Bu metod, bir kullanıcının kendi hesabını
//   silmesini engeller. Bu, yanlışlıkla hesap kaybını önlemek için önemlidir.
// - **Nil Kontrolleri**: Metod, context ve kullanıcı bilgilerinin nil olup
//   olmadığını kontrol eder.
// - **Tür Dönüştürme**: model parametresi *domainUser.User türüne başarıyla
//   dönüştürülmelidir.
// - **Uyarı**: model nil ise, metod true döner. Bu, genel silme yetkisini
//   kontrol etmek için kullanılabilir, ancak belirli bir kullanıcıyı silmek
//   için model parametresi gereklidir.
// - **Veri Tabanı Bütünlüğü**: Silme işleminden önce, ilişkili kayıtların
//   (örneğin, kullanıcının oluşturduğu veriler) nasıl işleneceğini belirleyin.
// ============================================================================
func (p UserPolicy) Delete(ctx *appContext.Context, model any) bool {
	// Genel yetki kontrolü: model nil ise, sadece yetki kontrolü yapılır
	if model == nil {
		return true
	}

	// model parametresini *domainUser.User türüne dönüştür
	userModel, ok := model.(*domainUser.User)
	if !ok {
		// Dönüştürme başarısız ise false döner
		return false
	}

	// Context nil ise false döner (kimliği doğrulanmış kullanıcı bilgisi yok)
	if ctx == nil {
		return false
	}

	// Kimliği doğrulanmış kullanıcı bilgisini al
	authUser := ctx.User()
	if authUser == nil {
		// Kimliği doğrulanmış kullanıcı bilgisi yok ise false döner
		return false
	}

	// Kendini silmeyi engelle: Silinecek kullanıcı ID'si, kimliği doğrulanmış
	// kullanıcı ID'sine eşitse false döner
	if userModel.ID == authUser.ID {
		return false
	}

	// Tüm kontroller geçerse true döner (silme işlemine izin ver)
	return true
}
