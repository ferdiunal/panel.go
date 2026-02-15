// Package fields, admin panel için alan (field) tanımlamalarını sağlar.
//
// Bu dosya, core paketinden type alias'ları ve constant'ları içe aktarır.
// Bu sayede fields paketi, core paketine doğrudan bağımlı olmadan temel tipleri kullanabilir.
package fields

import "github.com/ferdiunal/panel.go/pkg/core"

// ElementType, bir elemanın tipini belirtir.
//
// Bu tip, core.ElementType'ın bir alias'ıdır ve her field türünün (Text, Number, BelongsTo, vb.)
// hangi UI bileşenini kullanacağını belirler.
//
// Daha fazla bilgi için pkg/core/element_type.go dosyasına bakın.
type ElementType = core.ElementType

// ElementContext, bir elemanın hangi bağlamda görüntüleneceğini belirtir.
//
// Bu tip, core.ElementContext'in bir alias'ıdır ve elemanların form, liste veya detay
// sayfalarında görünürlüğünü kontrol eder.
//
// Daha fazla bilgi için pkg/core/element_context.go dosyasına bakın.
type ElementContext = core.ElementContext

// VisibilityContext, görünürlük kontrolü için bağlam bilgisini sağlar.
//
// Bu tip, core.VisibilityContext'in bir alias'ıdır ve elemanların hangi işlem
// sırasında (oluşturma, güncelleme, listeleme, vb.) görüntüleneceğini belirler.
//
// Daha fazla bilgi için pkg/core/visibility_context.go dosyasına bakın.
type VisibilityContext = core.VisibilityContext

// Resolver, alan değerlerini çözümleyen interface'dir.
//
// Bu tip, core.Resolver'ın bir alias'ıdır ve özel alan değeri çözümleme mantığı
// sağlamak için kullanılır.
//
// Daha fazla bilgi için pkg/core/resolver.go dosyasına bakın.
type Resolver = core.Resolver

const (
	TYPE_TEXT            ElementType = core.TYPE_TEXT
	TYPE_TEXTAREA        ElementType = core.TYPE_TEXTAREA
	TYPE_RICHTEXT        ElementType = core.TYPE_RICHTEXT
	TYPE_PASSWORD        ElementType = core.TYPE_PASSWORD
	TYPE_NUMBER          ElementType = core.TYPE_NUMBER
	TYPE_MONEY           ElementType = core.TYPE_MONEY
	TYPE_TEL             ElementType = core.TYPE_TEL
	TYPE_EMAIL           ElementType = core.TYPE_EMAIL
	TYPE_AUDIO           ElementType = core.TYPE_AUDIO
	TYPE_VIDEO           ElementType = core.TYPE_VIDEO
	TYPE_DATE            ElementType = core.TYPE_DATE
	TYPE_DATETIME        ElementType = core.TYPE_DATETIME
	TYPE_FILE            ElementType = core.TYPE_FILE
	TYPE_KEY_VALUE       ElementType = core.TYPE_KEY_VALUE
	TYPE_LINK            ElementType = core.TYPE_LINK
	TYPE_COLLECTION      ElementType = core.TYPE_COLLECTION
	TYPE_DETAIL          ElementType = core.TYPE_DETAIL
	TYPE_CONNECT         ElementType = core.TYPE_CONNECT
	TYPE_POLY_LINK       ElementType = core.TYPE_POLY_LINK
	TYPE_POLY_DETAIL     ElementType = core.TYPE_POLY_DETAIL
	TYPE_POLY_COLLECTION ElementType = core.TYPE_POLY_COLLECTION
	TYPE_POLY_CONNECT    ElementType = core.TYPE_POLY_CONNECT
	TYPE_BOOLEAN         ElementType = core.TYPE_BOOLEAN
	TYPE_SELECT          ElementType = core.TYPE_SELECT
	TYPE_PANEL           ElementType = core.TYPE_PANEL
	TYPE_TABS            ElementType = core.TYPE_TABS
	TYPE_STACK           ElementType = core.TYPE_STACK
	TYPE_RELATIONSHIP    ElementType = core.TYPE_RELATIONSHIP
	TYPE_BADGE           ElementType = core.TYPE_BADGE
	TYPE_CODE            ElementType = core.TYPE_CODE
	TYPE_COLOR           ElementType = core.TYPE_COLOR
	TYPE_BOOLEAN_GROUP   ElementType = core.TYPE_BOOLEAN_GROUP
)

const (
	CONTEXT_FORM   ElementContext = core.CONTEXT_FORM
	CONTEXT_DETAIL ElementContext = core.CONTEXT_DETAIL
	CONTEXT_LIST   ElementContext = core.CONTEXT_LIST

	SHOW_ON_FORM   ElementContext = core.SHOW_ON_FORM
	SHOW_ON_DETAIL ElementContext = core.SHOW_ON_DETAIL
	SHOW_ON_LIST   ElementContext = core.SHOW_ON_LIST

	HIDE_ON_LIST   ElementContext = core.HIDE_ON_LIST
	HIDE_ON_DETAIL ElementContext = core.HIDE_ON_DETAIL
	HIDE_ON_CREATE ElementContext = core.HIDE_ON_CREATE
	HIDE_ON_UPDATE ElementContext = core.HIDE_ON_UPDATE

	ONLY_ON_LIST   ElementContext = core.ONLY_ON_LIST
	ONLY_ON_DETAIL ElementContext = core.ONLY_ON_DETAIL
	ONLY_ON_CREATE ElementContext = core.ONLY_ON_CREATE
	ONLY_ON_UPDATE ElementContext = core.ONLY_ON_UPDATE
	ONLY_ON_FORM   ElementContext = core.ONLY_ON_FORM
)

const (
	ContextIndex   VisibilityContext = core.ContextIndex
	ContextDetail  VisibilityContext = core.ContextDetail
	ContextCreate  VisibilityContext = core.ContextCreate
	ContextUpdate  VisibilityContext = core.ContextUpdate
	ContextPreview VisibilityContext = core.ContextPreview
)
