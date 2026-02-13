package products

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cargo.go/entity"
	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// GetProductActions, Products resource için özel action'ları döndürür
//
// Bu fonksiyon, OptimizedBase'den gelen varsayılan action'lara (Seçilenleri Sil)
// ek olarak Products resource'una özel action'ları ekler.
//
// Döndürür: resource.Action listesi
//
// Varsayılan Action'lar (OptimizedBase'den):
// 1. Seçilenleri Sil - Checkbox ile seçilen kayıtları siler (tüm resource'larda var)
//
// Özel Action'lar (Products için):
// 2. Tümünü Sil - TÜM ürünleri siler (çok destructive, admin only)
// 3. Seçilenleri Dışa Aktar - CSV formatında export
//
// Kullanım:
//   func (r *ProductResource) GetActions() []resource.Action {
//       return GetProductActions(r)
//   }
//
// API Endpoint'leri:
//   GET  /api/products/actions              # Action listesini al
//   POST /api/products/actions/delete-selected  # Seçilenleri sil (varsayılan)
//   POST /api/products/actions/delete-all       # Tümünü sil (özel)
//   POST /api/products/actions/export-selected  # CSV export (özel)
func GetProductActions(r *resource.OptimizedBase) []resource.Action {
	// Varsayılan action'ları al (Seçilenleri Sil)
	actions := r.GetDefaultActions()

	// ============================================================================
	// Özel Action 1: Tümünü Sil
	// ============================================================================
	// Veritabanındaki TÜM ürünleri siler. Çok tehlikeli!
	// Sadece test ortamlarında veya acil durumlarda kullanılmalıdır.
	deleteAll := action.New("Tümünü Sil").
		SetIcon("trash").
		SetSlug("delete-all").
		Destructive().
		Confirm("⚠️ DİKKAT: TÜM ürünler silinecek! Bu işlem geri alınamaz ve tüm ürün verileriniz kaybolacak. Devam etmek istediğinizden EMİN misiniz?").
		ConfirmButton("Evet, Tümünü Sil").
		CancelButton("İptal").
		ShowOnlyOnIndex().
		Handle(func(ctx *action.ActionContext) error {
			// Transaction içinde çalış (hata durumunda rollback)
			return ctx.DB.Transaction(func(tx *gorm.DB) error {
				// Tüm ürünleri sil
				result := tx.Where("1 = 1").Delete(&entity.Product{})
				if result.Error != nil {
					return fmt.Errorf("ürünler silinirken hata oluştu: %w", result.Error)
				}

				// Silinen kayıt sayısını logla
				fmt.Printf("Toplam %d ürün silindi\n", result.RowsAffected)
				return nil
			})
		}).
		AuthorizeUsing(func(ctx *action.ActionContext) bool {
			// TODO: Gerçek yetkilendirme kontrolü ekleyin
			// Örnek: user := ctx.User.(*User); return user.IsAdmin
			return true // Şimdilik herkese izin ver (geliştirme ortamı)
		})

	// ============================================================================
	// Özel Action 2: Seçilenleri Dışa Aktar
	// ============================================================================
	// Seçili ürünleri CSV dosyası olarak dışa aktarır.
	// Dosya exports/ klasörüne kaydedilir.
	exportSelected := action.New("Seçilenleri Dışa Aktar").
		SetIcon("download").
		SetSlug("export-selected").
		ShowOnlyOnIndex().
		Handle(func(ctx *action.ActionContext) error {
			// exports klasörünü oluştur
			exportsDir := "exports"
			if err := os.MkdirAll(exportsDir, 0755); err != nil {
				return fmt.Errorf("exports klasörü oluşturulamadı: %w", err)
			}

			// CSV dosyası oluştur (timestamp ile)
			timestamp := time.Now().Format("20060102_150405")
			filename := filepath.Join(exportsDir, fmt.Sprintf("products_%s.csv", timestamp))
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("CSV dosyası oluşturulamadı: %w", err)
			}
			defer file.Close()

			// CSV writer oluştur
			writer := csv.NewWriter(file)
			defer writer.Flush()

			// Header satırını yaz
			header := []string{"ID", "Name", "Organization", "Created At", "Updated At"}
			if err := writer.Write(header); err != nil {
				return fmt.Errorf("CSV header yazılamadı: %w", err)
			}

			// Her ürün için satır yaz
			for _, item := range ctx.Models {
				product, ok := item.(*entity.Product)
				if !ok {
					return fmt.Errorf("geçersiz kayıt tipi")
				}

				// Organization adını al (eager load edilmişse)
				orgName := ""
				if product.Organization != nil {
					orgName = product.Organization.Name
				}

				// CSV satırını oluştur
				row := []string{
					fmt.Sprintf("%d", product.ID),
					product.Name,
					orgName,
					product.CreatedAt.Format("2006-01-02 15:04:05"),
					product.UpdatedAt.Format("2006-01-02 15:04:05"),
				}

				if err := writer.Write(row); err != nil {
					return fmt.Errorf("CSV satırı yazılamadı: %w", err)
				}
			}

			fmt.Printf("CSV dosyası oluşturuldu: %s\n", filename)
			return nil
		})

	// Özel action'ları ekle
	actions = append(actions, deleteAll, exportSelected)

	return actions
}

