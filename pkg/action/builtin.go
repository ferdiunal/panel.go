// Bu paket, panel uygulamasında kullanılan yerleşik (built-in) action'ları tanımlar.
// Action'lar, seçilen modeller üzerinde toplu işlemler gerçekleştirmek için kullanılır.
// Örneğin: CSV export, toplu silme, onaylama vb.
package action

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// Bu fonksiyon, seçilen modelleri CSV dosyasına aktarmak için bir action oluşturur.
//
// Kullanım Senaryosu:
// - Kullanıcılar, panel arayüzünden seçtikleri kayıtları CSV formatında dışa aktarabilirler
// - Raporlama ve veri analizi için kullanılır
// - Toplu veri transferi için idealdir
//
// Parametreler:
// - filename: Dışa aktarılan dosyanın adı. Boş string ise otomatik olarak "export_TIMESTAMP.csv" adıyla oluşturulur
//
// Dönüş Değeri:
// - *BaseAction: Yapılandırılmış action pointer'ı, method chaining için kullanılabilir
//
// Kullanım Örneği:
//   action := ExportCSV("users.csv")
//   // Sonuç: storage/exports/users_20240207_150405.csv dosyası oluşturulur
//
// Önemli Notlar:
// - Dışa aktarılan dosyalar "storage/exports" dizinine kaydedilir
// - Dosya adına otomatik olarak timestamp eklenir (çakışmaları önlemek için)
// - Sadece public (dışa aktarılan) struct alanları CSV'ye yazılır
// - Eğer model listesi boşsa hata döndürülür
func ExportCSV(filename string) *BaseAction {
	return New("Export as CSV").
		SetIcon("download").
		Handle(func(ctx *ActionContext) error {
			return exportToCSV(ctx.Models, filename)
		})
}

// Bu fonksiyon, seçilen modelleri veritabanından silmek için bir action oluşturur.
//
// Kullanım Senaryosu:
// - Kullanıcılar, panel arayüzünden toplu olarak kayıtları silebilirler
// - Veri temizleme ve yönetimi için kullanılır
// - Yanlışlıkla silmeyi önlemek için onay mekanizması vardır
//
// Dönüş Değeri:
// - *BaseAction: Yapılandırılmış action pointer'ı, method chaining için kullanılabilir
//
// Kullanım Örneği:
//   action := Delete()
//   // Kullanıcı onay verirse, seçilen tüm modeller silinir
//
// Önemli Notlar:
// - Bu action "Destructive" olarak işaretlenmiştir (kırmızı uyarı gösterir)
// - Silme işleminden önce kullanıcıdan onay istenir
// - Eğer herhangi bir model silinirken hata oluşursa, işlem durdurulur ve hata döndürülür
// - Veritabanı seviyesinde cascade delete kuralları uygulanır
// - Silinen veriler geri alınamaz, bu nedenle dikkatli kullanılmalıdır
func Delete() *BaseAction {
	return New("Delete").
		SetIcon("trash").
		Destructive().
		Confirm("Are you sure you want to delete these items?").
		ConfirmButton("Delete").
		Handle(func(ctx *ActionContext) error {
			for _, model := range ctx.Models {
				if err := ctx.DB.Delete(model).Error; err != nil {
					return err
				}
			}
			return nil
		})
}

// Bu fonksiyon, seçilen modelleri onaylamak için bir action oluşturur.
//
// Kullanım Senaryosu:
// - Kullanıcılar, panel arayüzünden toplu olarak kayıtları onaylayabilirler
// - İş akışı yönetimi ve onay süreçleri için kullanılır
// - Örneğin: başvuruları onaylama, siparişleri onaylama vb.
//
// Dönüş Değeri:
// - *BaseAction: Yapılandırılmış action pointer'ı, method chaining için kullanılabilir
//
// Kullanım Örneği:
//   action := Approve()
//   // Seçilen modellerin Status alanı "approved" olarak ayarlanır
//   // veya Approved alanı true olarak ayarlanır
//
// Önemli Notlar:
// - Fonksiyon, modellerde "Status" veya "Approved" alanı arar
// - Önce "Status" alanını kontrol eder ve "approved" string değeri ayarlar
// - Eğer "Status" alanı yoksa, "Approved" boolean alanını true olarak ayarlar
// - Eğer her iki alan da yoksa, model değiştirilmez
// - Reflection kullanarak dinamik olarak alanları bulur ve ayarlar
// - Onay işleminden önce kullanıcıdan onay istenir
// - Her model başarıyla onaylandıktan sonra veritabanına kaydedilir
func Approve() *BaseAction {
	return New("Approve").
		SetIcon("check").
		Confirm("Are you sure you want to approve these items?").
		Handle(func(ctx *ActionContext) error {
			// Varsayılan onaylama mantığı - özel action'larla geçersiz kılınabilir
			for _, model := range ctx.Models {
				// Modelin reflection değerini al
				v := reflect.ValueOf(model)
				// Eğer pointer ise, işaret ettiği değeri al
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				// Modelin struct olup olmadığını kontrol et
				if v.Kind() == reflect.Struct {
					// Status alanını bulmaya ve ayarlamaya çalış
					statusField := v.FieldByName("Status")
					if statusField.IsValid() && statusField.CanSet() {
						statusField.SetString("approved")
						ctx.DB.Save(model)
						continue
					}

					// Approved alanını bulmaya ve ayarlamaya çalış
					approvedField := v.FieldByName("Approved")
					if approvedField.IsValid() && approvedField.CanSet() {
						approvedField.SetBool(true)
						ctx.DB.Save(model)
						continue
					}
				}
			}
			return nil
		})
}

// Bu fonksiyon, verilen modelleri CSV dosyasına aktarır.
//
// Kullanım Senaryosu:
// - ExportCSV action'ı tarafından dahili olarak çağrılır
// - Modelleri CSV formatına dönüştürür ve dosyaya yazar
// - Toplu veri dışa aktarma işlemini gerçekleştirir
//
// Parametreler:
// - models: Dışa aktarılacak model nesnelerinin slice'ı ([]interface{})
// - filename: Oluşturulacak CSV dosyasının adı
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   models := []interface{}{user1, user2, user3}
//   err := exportToCSV(models, "users.csv")
//   if err != nil {
//       log.Fatal(err)
//   }
//
// Önemli Notlar:
// - Dışa aktarılan dosyalar "storage/exports" dizinine kaydedilir
// - Dizin yoksa otomatik olarak oluşturulur (0755 izinleriyle)
// - Dosya adına otomatik olarak timestamp eklenir (YYYYMMDD_HHMMSS formatında)
// - Sadece public (dışa aktarılan) struct alanları CSV'ye yazılır
// - Unexported (küçük harfle başlayan) alanlar atlanır
// - İlk modelden başlık satırı (header) oluşturulur
// - Tüm modeller aynı struct tipinde olmalıdır
// - Eğer model listesi boşsa hata döndürülür
// - CSV yazma işlemi sırasında hata oluşursa, işlem durdurulur
func exportToCSV(models []interface{}, filename string) error {
	// Modellerin boş olup olmadığını kontrol et
	if len(models) == 0 {
		return fmt.Errorf("no models to export")
	}

	// Dışa aktarma dizinini oluştur (yoksa)
	// Dizin yapısı: storage/exports/
	exportsDir := filepath.Join("storage", "exports")
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create exports directory: %w", err)
	}

	// Dosya adına timestamp ekle (çakışmaları önlemek için)
	// Format: YYYYMMDD_HHMMSS (örn: 20240207_150405)
	timestamp := time.Now().Format("20060102_150405")
	if filename == "" {
		// Eğer dosya adı belirtilmemişse, varsayılan ad kullan
		filename = fmt.Sprintf("export_%s.csv", timestamp)
	} else {
		// Dosya adına timestamp ekle (uzantıyı koru)
		// Örn: users.csv -> users_20240207_150405.csv
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		filename = fmt.Sprintf("%s_%s%s", name, timestamp, ext)
	}

	// Tam dosya yolunu oluştur
	filePath := filepath.Join(exportsDir, filename)

	// CSV dosyasını oluştur
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// CSV yazıcısını oluştur
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// İlk modelden alan adlarını al
	firstModel := models[0]
	v := reflect.ValueOf(firstModel)

	// Eğer pointer ise, işaret ettiği değeri al
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Modelin struct olup olmadığını kontrol et
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct")
	}

	// Struct tipini al
	t := v.Type()
	var headers []string
	var fieldIndices []int

	// Tüm public alanları topla
	// Unexported (küçük harfle başlayan) alanlar atlanır
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Unexported alanları atla (PkgPath boş değilse unexported demektir)
		if field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
		fieldIndices = append(fieldIndices, i)
	}

	// CSV başlık satırını yaz
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Tüm modellerin verilerini CSV'ye yaz
	for _, model := range models {
		v := reflect.ValueOf(model)

		// Eğer pointer ise, işaret ettiği değeri al
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// Bu model için bir satır oluştur
		var row []string
		for _, idx := range fieldIndices {
			fieldValue := v.Field(idx)
			// Alan değerini string'e dönüştür
			row = append(row, fmt.Sprintf("%v", fieldValue.Interface()))
		}

		// Satırı CSV'ye yaz
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
