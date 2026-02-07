// Package fields, admin panel için alan (field) tanımlamalarını sağlar.
//
// Bu dosya, core paketinden type alias'ları içe aktarır ve fields paketinde kullanılabilir hale getirir.
// Bu sayede fields paketi, core paketine doğrudan bağımlı olmadan temel tipleri kullanabilir.
package fields

import "github.com/ferdiunal/panel.go/pkg/core"

// Element, admin panel'deki bir UI elemanını temsil eder.
//
// Bu tip, core.Element'in bir alias'ıdır ve tüm field türleri (Text, Number, BelongsTo, vb.)
// bu interface'i implement eder.
//
// Daha fazla bilgi için pkg/core/element.go dosyasına bakın.
type Element = core.Element

// VisibilityFunc, bir elemanın görünürlüğünü kontrol eden callback fonksiyonudur.
//
// Bu tip, core.VisibilityFunc'ın bir alias'ıdır ve elemanların dinamik olarak
// gösterilip gizlenmesini sağlar.
//
// # Kullanım Örneği
//
//	field.ShowIf(func(ctx VisibilityContext) bool {
//	    return ctx == ContextCreate
//	})
//
// Daha fazla bilgi için pkg/core/visibility.go dosyasına bakın.
type VisibilityFunc = core.VisibilityFunc

// StorageCallbackFunc, dosya depolama işlemlerini özelleştiren callback fonksiyonudur.
//
// Bu tip, core.StorageCallbackFunc'ın bir alias'ıdır ve dosya yükleme alanlarında
// (Image, File, Video, Audio) depolama davranışını özelleştirmek için kullanılır.
//
// # Kullanım Örneği
//
//	field.Storage(func(file interface{}) (string, error) {
//	    // Özel depolama mantığı
//	    return "path/to/file", nil
//	})
//
// Daha fazla bilgi için pkg/core/storage.go dosyasına bakın.
type StorageCallbackFunc = core.StorageCallbackFunc
