package data

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ferdiunal/panel.go/pkg/fields"
)

// eagerLoadBelongsTo, BelongsTo ilişkisini batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların BelongsTo
// ilişkilerini tek sorguda yükler.
//
// # İşlem Sırası
//
// 1. Tüm item'lardan foreign key değerlerini çıkar
// 2. GORM query builder ile ilişkili kayıtları yükle (SELECT * FROM related_table WHERE id IN (...))
// 3. İlişkili kayıtları ID'ye göre map et
// 4. Her item'a ilişkili kaydı set et
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **items**: İlişkileri yüklenecek kayıt listesi
// - **field**: BelongsTo field tanımı
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) eagerLoadBelongsTo(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
	// BelongsTo field'ından gerekli bilgileri al
	belongsToField, ok := field.(*fields.BelongsToField)
	if !ok {
		return fmt.Errorf("field is not a BelongsTo field")
	}

	foreignKey := belongsToField.GetForeignKey()
	ownerKey := belongsToField.GetOwnerKeyColumn()
	relatedTable := belongsToField.GetRelatedTableName()

	if foreignKey == "" || relatedTable == "" {
		return fmt.Errorf("invalid BelongsTo configuration: foreignKey=%s, relatedTable=%s", foreignKey, relatedTable)
	}

	// 1. Tüm item'lardan foreign key değerlerini çıkar
	foreignKeyValues := []interface{}{}
	itemsByForeignKey := map[interface{}][]interface{}{}

	for _, item := range items {
		fkValue := extractFieldValue(item, foreignKey)
		if fkValue != nil && !isZeroValue(fkValue) {
			foreignKeyValues = append(foreignKeyValues, fkValue)
			itemsByForeignKey[fkValue] = append(itemsByForeignKey[fkValue], item)
		}
	}

	if len(foreignKeyValues) == 0 {
		return nil // Hiç foreign key yok, ilişki yok
	}

	// 2. GORM query builder ile ilişkili kayıtları yükle
	safeOwnerKey := SanitizeColumnName(ownerKey)
	safeTable := SanitizeColumnName(relatedTable)

	var relatedRecords []map[string]interface{}
	err := l.db.WithContext(ctx).
		Table(safeTable).
		Where(safeOwnerKey+" IN ?", foreignKeyValues).
		Find(&relatedRecords).Error

	if err != nil {
		return fmt.Errorf("failed to load BelongsTo relationship: %w", err)
	}

	// 3. İlişkili kayıtları ID'ye göre map et
	relatedByID := map[interface{}]map[string]interface{}{}
	for _, record := range relatedRecords {
		id := record[ownerKey]
		relatedByID[id] = record
	}

	// 4. Her item'a ilişkili kaydı set et
	for fkValue, itemList := range itemsByForeignKey {
		relatedRecord := relatedByID[fkValue]
		if relatedRecord != nil {
			for _, item := range itemList {
				if err := setRelationshipData(item, field.GetRelationshipName(), relatedRecord); err != nil {
					// Log error but continue
					fmt.Printf("[WARN] Failed to set BelongsTo relationship: %v\n", err)
				}
			}
		}
	}

	return nil
}

// lazyLoadBelongsTo, tek bir kayıt için BelongsTo ilişkisini yükler.
//
// Bu metod, GORM Association API kullanarak tek bir kaydın BelongsTo
// ilişkisini yükler. Reflection kullanarak dinamik olarak struct tipini belirler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **item**: İlişkisi yüklenecek kayıt
// - **field**: BelongsTo field tanımı
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) lazyLoadBelongsTo(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// Reflection kullanarak relationship field tipini al
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	relField := itemValue.FieldByName(field.GetRelationshipName())
	if !relField.IsValid() {
		return nil, fmt.Errorf("relationship field %s not found", field.GetRelationshipName())
	}

	relType := relField.Type()

	// Yeni instance oluştur
	var relValue reflect.Value
	if relType.Kind() == reflect.Ptr {
		// Pointer tip (örn. *Author)
		relValue = reflect.New(relType.Elem())
	} else {
		// Non-pointer tip (örn. Author)
		relValue = reflect.New(relType)
	}

	// GORM Association API kullanarak ilişkiyi yükle
	err := l.db.WithContext(ctx).
		Model(item).
		Association(field.GetRelationshipName()).
		Find(relValue.Interface())

	if err != nil {
		return nil, fmt.Errorf("failed to load BelongsTo relationship: %w", err)
	}

	// İlişki verisini set et
	var actualValue interface{}
	if relType.Kind() == reflect.Ptr {
		actualValue = relValue.Interface()
	} else {
		actualValue = relValue.Elem().Interface()
	}

	if err := setRelationshipData(item, field.GetRelationshipName(), actualValue); err != nil {
		return nil, err
	}

	return actualValue, nil
}

// eagerLoadHasMany, HasMany ilişkisini batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların HasMany
// ilişkilerini tek sorguda yükler.
//
// # İşlem Sırası
//
// 1. Tüm item'lardan owner key değerlerini çıkar
// 2. GORM query builder ile ilişkili kayıtları yükle (SELECT * FROM related_table WHERE foreign_key IN (...))
// 3. İlişkili kayıtları foreign key'e göre grupla
// 4. Her item'a ilişkili kayıt listesini set et
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **items**: İlişkileri yüklenecek kayıt listesi
// - **field**: HasMany field tanımı
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) eagerLoadHasMany(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
	// HasMany field'ından gerekli bilgileri al
	hasManyField, ok := field.(*fields.HasManyField)
	if !ok {
		return fmt.Errorf("field is not a HasMany field")
	}

	foreignKey := hasManyField.GetForeignKeyColumn()
	ownerKey := hasManyField.GetOwnerKeyColumn()
	relatedTable := hasManyField.GetRelatedTableName()

	if foreignKey == "" || relatedTable == "" {
		return fmt.Errorf("invalid HasMany configuration: foreignKey=%s, relatedTable=%s", foreignKey, relatedTable)
	}

	// 1. Tüm item'lardan owner key değerlerini çıkar
	ownerKeyValues := []interface{}{}
	itemsByOwnerKey := map[interface{}]interface{}{}

	for _, item := range items {
		ownerValue := extractFieldValue(item, ownerKey)
		if ownerValue != nil && !isZeroValue(ownerValue) {
			ownerKeyValues = append(ownerKeyValues, ownerValue)
			itemsByOwnerKey[ownerValue] = item
		}
	}

	if len(ownerKeyValues) == 0 {
		return nil // Hiç owner key yok
	}

	// 2. GORM query builder ile ilişkili kayıtları yükle
	safeForeignKey := SanitizeColumnName(foreignKey)
	safeTable := SanitizeColumnName(relatedTable)

	var relatedRecords []map[string]interface{}
	err := l.db.WithContext(ctx).
		Table(safeTable).
		Where(safeForeignKey+" IN ?", ownerKeyValues).
		Find(&relatedRecords).Error

	if err != nil {
		return fmt.Errorf("failed to load HasMany relationship: %w", err)
	}

	// 3. İlişkili kayıtları foreign key'e göre grupla
	relatedByFK := map[interface{}][]map[string]interface{}{}
	for _, record := range relatedRecords {
		fkValue := record[foreignKey]
		relatedByFK[fkValue] = append(relatedByFK[fkValue], record)
	}

	// 4. Her item'a ilişkili kayıt listesini set et
	for ownerValue, item := range itemsByOwnerKey {
		relatedList := relatedByFK[ownerValue]
		if relatedList == nil {
			relatedList = []map[string]interface{}{} // Boş liste
		}
		if err := setRelationshipData(item, field.GetRelationshipName(), relatedList); err != nil {
			// Log error but continue
			fmt.Printf("[WARN] Failed to set HasMany relationship: %v\n", err)
		}
	}

	return nil
}

// lazyLoadHasMany, tek bir kayıt için HasMany ilişkisini yükler.
//
// Bu metod, GORM Association API kullanarak tek bir kaydın HasMany
// ilişkisini yükler. Reflection kullanarak dinamik olarak struct tipini belirler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **item**: İlişkisi yüklenecek kayıt
// - **field**: HasMany field tanımı
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi (slice)
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) lazyLoadHasMany(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// Reflection kullanarak relationship field tipini al
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	relField := itemValue.FieldByName(field.GetRelationshipName())
	if !relField.IsValid() {
		return nil, fmt.Errorf("relationship field %s not found", field.GetRelationshipName())
	}

	relType := relField.Type()
	if relType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("HasMany relationship field must be a slice")
	}

	// Yeni slice instance oluştur
	relValue := reflect.New(relType)

	// GORM Association API kullanarak ilişkiyi yükle
	err := l.db.WithContext(ctx).
		Model(item).
		Association(field.GetRelationshipName()).
		Find(relValue.Interface())

	if err != nil {
		return nil, fmt.Errorf("failed to load HasMany relationship: %w", err)
	}

	// İlişki verisini set et
	actualValue := relValue.Elem().Interface()
	if err := setRelationshipData(item, field.GetRelationshipName(), actualValue); err != nil {
		return nil, err
	}

	return actualValue, nil
}

// eagerLoadHasOne, HasOne ilişkisini batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların HasOne
// ilişkilerini tek sorguda yükler.
//
// # İşlem Sırası
//
// 1. Tüm item'lardan owner key değerlerini çıkar
// 2. GORM query builder ile ilişkili kayıtları yükle (SELECT * FROM related_table WHERE foreign_key IN (...))
// 3. İlişkili kayıtları foreign key'e göre map et (her foreign key için tek kayıt)
// 4. Her item'a ilişkili kaydı set et
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **items**: İlişkileri yüklenecek kayıt listesi
// - **field**: HasOne field tanımı
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) eagerLoadHasOne(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
	// HasOne field'ından gerekli bilgileri al
	hasOneField, ok := field.(*fields.HasOneField)
	if !ok {
		return fmt.Errorf("field is not a HasOne field")
	}

	foreignKey := hasOneField.GetForeignKeyColumn()
	ownerKey := hasOneField.GetOwnerKeyColumn()
	relatedTable := hasOneField.GetRelatedTableName()

	if foreignKey == "" || relatedTable == "" {
		return fmt.Errorf("invalid HasOne configuration: foreignKey=%s, relatedTable=%s", foreignKey, relatedTable)
	}

	// 1. Tüm item'lardan owner key değerlerini çıkar
	ownerKeyValues := []interface{}{}
	itemsByOwnerKey := map[interface{}]interface{}{}

	for _, item := range items {
		ownerValue := extractFieldValue(item, ownerKey)
		if ownerValue != nil && !isZeroValue(ownerValue) {
			ownerKeyValues = append(ownerKeyValues, ownerValue)
			itemsByOwnerKey[ownerValue] = item
		}
	}

	if len(ownerKeyValues) == 0 {
		return nil // Hiç owner key yok
	}

	// 2. GORM query builder ile ilişkili kayıtları yükle
	safeForeignKey := SanitizeColumnName(foreignKey)
	safeTable := SanitizeColumnName(relatedTable)

	var relatedRecords []map[string]interface{}
	err := l.db.WithContext(ctx).
		Table(safeTable).
		Where(safeForeignKey+" IN ?", ownerKeyValues).
		Find(&relatedRecords).Error

	if err != nil {
		return fmt.Errorf("failed to load HasOne relationship: %w", err)
	}

	// 3. İlişkili kayıtları foreign key'e göre map et (her foreign key için tek kayıt)
	relatedByFK := map[interface{}]map[string]interface{}{}
	for _, record := range relatedRecords {
		fkValue := record[foreignKey]
		// HasOne için sadece ilk kaydı al (birden fazla varsa ilkini kullan)
		if _, exists := relatedByFK[fkValue]; !exists {
			relatedByFK[fkValue] = record
		}
	}

	// 4. Her item'a ilişkili kaydı set et
	for ownerValue, item := range itemsByOwnerKey {
		relatedRecord := relatedByFK[ownerValue]
		if relatedRecord != nil {
			if err := setRelationshipData(item, field.GetRelationshipName(), relatedRecord); err != nil {
				// Log error but continue
				fmt.Printf("[WARN] Failed to set HasOne relationship: %v\n", err)
			}
		}
	}

	return nil
}

// lazyLoadHasOne, tek bir kayıt için HasOne ilişkisini yükler.
//
// Bu metod, GORM Association API kullanarak tek bir kaydın HasOne
// ilişkisini yükler. Reflection kullanarak dinamik olarak struct tipini belirler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **item**: İlişkisi yüklenecek kayıt
// - **field**: HasOne field tanımı
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) lazyLoadHasOne(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// Reflection kullanarak relationship field tipini al
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	relField := itemValue.FieldByName(field.GetRelationshipName())
	if !relField.IsValid() {
		return nil, fmt.Errorf("relationship field %s not found", field.GetRelationshipName())
	}

	relType := relField.Type()

	// Yeni instance oluştur
	var relValue reflect.Value
	if relType.Kind() == reflect.Ptr {
		// Pointer tip (örn. *Profile)
		relValue = reflect.New(relType.Elem())
	} else {
		// Non-pointer tip (örn. Profile)
		relValue = reflect.New(relType)
	}

	// GORM Association API kullanarak ilişkiyi yükle
	err := l.db.WithContext(ctx).
		Model(item).
		Association(field.GetRelationshipName()).
		Find(relValue.Interface())

	if err != nil {
		return nil, fmt.Errorf("failed to load HasOne relationship: %w", err)
	}

	// İlişki verisini set et
	var actualValue interface{}
	if relType.Kind() == reflect.Ptr {
		actualValue = relValue.Interface()
	} else {
		actualValue = relValue.Elem().Interface()
	}

	if err := setRelationshipData(item, field.GetRelationshipName(), actualValue); err != nil {
		return nil, err
	}

	return actualValue, nil
}

// extractFieldValue, reflection kullanarak bir item'dan field değerini çıkarır.
//
// Bu fonksiyon, hem struct hem de map tiplerini destekler.
//
// # Parametreler
//
// - **item**: Değeri çıkarılacak kayıt (struct veya map)
// - **fieldName**: Field adı (snake_case veya PascalCase)
//
// # Döndürür
//
// - interface{}: Field değeri (nil olabilir)
func extractFieldValue(item interface{}, fieldName string) interface{} {
	if item == nil {
		return nil
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// Struct field
	if v.Kind() == reflect.Struct {
		// Try PascalCase first (Go struct field naming)
		field := v.FieldByName(toPascalCase(fieldName))
		if field.IsValid() && field.CanInterface() {
			return field.Interface()
		}
		// Try exact name
		field = v.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			return field.Interface()
		}
	}

	// Map access
	if v.Kind() == reflect.Map {
		mapValue := v.Interface().(map[string]interface{})
		// Try snake_case first
		if val, ok := mapValue[fieldName]; ok {
			return val
		}
		// Try PascalCase
		if val, ok := mapValue[toPascalCase(fieldName)]; ok {
			return val
		}
	}

	return nil
}

// setRelationshipData, reflection kullanarak bir item'a ilişki verisini set eder.
//
// Bu fonksiyon, hem struct hem de map tiplerini destekler.
//
// # Parametreler
//
// - **item**: Veri set edilecek kayıt (struct veya map)
// - **relationshipName**: İlişki adı (PascalCase)
// - **data**: Set edilecek veri
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
func setRelationshipData(item interface{}, relationshipName string, data interface{}) error {
	if item == nil {
		return fmt.Errorf("item is nil")
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("item pointer is nil")
		}
		v = v.Elem()
	}

	// Struct field
	if v.Kind() == reflect.Struct {
		field := v.FieldByName(relationshipName)
		if field.IsValid() && field.CanSet() {
			dataValue := reflect.ValueOf(data)
			if dataValue.IsValid() && dataValue.Type().AssignableTo(field.Type()) {
				field.Set(dataValue)
				return nil
			}
			// Try to convert if types don't match exactly
			if dataValue.IsValid() && dataValue.Type().ConvertibleTo(field.Type()) {
				field.Set(dataValue.Convert(field.Type()))
				return nil
			}
		}
		return fmt.Errorf("cannot set field %s on struct", relationshipName)
	}

	// Map access
	if v.Kind() == reflect.Map {
		mapValue := v.Interface().(map[string]interface{})
		mapValue[relationshipName] = data
		return nil
	}

	return fmt.Errorf("unsupported item type: %s", v.Kind())
}

// isZeroValue, bir değerin zero value olup olmadığını kontrol eder.
//
// # Parametreler
//
// - **value**: Kontrol edilecek değer
//
// # Döndürür
//
// - bool: Zero value ise true
func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// toPascalCase, snake_case string'i PascalCase'e dönüştürür.
//
// # Parametreler
//
// - **s**: Dönüştürülecek string (örn. "user_id")
//
// # Döndürür
//
// - string: PascalCase string (örn. "UserID")
func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Özel durumlar
	switch s {
	case "id":
		return "ID"
	case "user_id":
		return "UserID"
	case "author_id":
		return "AuthorID"
	case "organization_id":
		return "OrganizationID"
	case "created_at":
		return "CreatedAt"
	case "updated_at":
		return "UpdatedAt"
	case "deleted_at":
		return "DeletedAt"
	}

	// Genel dönüşüm: snake_case -> PascalCase
	result := ""
	capitalize := true
	for _, ch := range s {
		if ch == '_' {
			capitalize = true
			continue
		}
		if capitalize {
			result += string(ch - 32) // Uppercase
			capitalize = false
		} else {
			result += string(ch)
		}
	}

	return result
}

// eagerLoadBelongsToMany, BelongsToMany ilişkisini batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların BelongsToMany
// ilişkilerini pivot tablo üzerinden tek sorguda yükler.
//
// # İşlem Sırası
//
// 1. Tüm item'lardan owner key değerlerini çıkar
// 2. Pivot tablodan ilişkili kayıt ID'lerini çek
// 3. İlişkili tablodan kayıtları çek
// 4. Her item'a ilişkili kayıt listesini set et
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **items**: İlişkileri yüklenecek kayıt listesi
// - **field**: BelongsToMany field tanımı
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) eagerLoadBelongsToMany(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
	// BelongsToMany field'ından gerekli bilgileri al
	belongsToManyField, ok := field.(*fields.BelongsToManyField)
	if !ok {
		return fmt.Errorf("field is not a BelongsToMany field")
	}

	pivotTable := belongsToManyField.PivotTableName
	parentColumn := belongsToManyField.ForeignKeyColumn
	relatedColumn := belongsToManyField.RelatedKeyColumn
	relatedTable := belongsToManyField.RelatedResourceSlug
	ownerKey := "id" // Ana tablonun primary key'i

	if pivotTable == "" || parentColumn == "" || relatedColumn == "" || relatedTable == "" {
		return fmt.Errorf("invalid BelongsToMany configuration")
	}

	// 1. Tüm item'lardan owner key değerlerini çıkar
	ownerKeyValues := []interface{}{}
	itemsByOwnerKey := map[interface{}]interface{}{}

	for _, item := range items {
		ownerValue := extractFieldValue(item, ownerKey)
		if ownerValue != nil && !isZeroValue(ownerValue) {
			ownerKeyValues = append(ownerKeyValues, ownerValue)
			itemsByOwnerKey[ownerValue] = item
		}
	}

	if len(ownerKeyValues) == 0 {
		return nil // Hiç owner key yok
	}

	// 2. Pivot tablodan ilişkili kayıt ID'lerini çek
	safePivotTable := SanitizeColumnName(pivotTable)
	safeParentColumn := SanitizeColumnName(parentColumn)
	safeRelatedColumn := SanitizeColumnName(relatedColumn)

	type PivotRecord struct {
		ParentID  interface{}
		RelatedID interface{}
	}

	var pivotRecords []map[string]interface{}
	err := l.db.WithContext(ctx).
		Table(safePivotTable).
		Where(safeParentColumn+" IN ?", ownerKeyValues).
		Find(&pivotRecords).Error

	if err != nil {
		return fmt.Errorf("failed to load BelongsToMany pivot records: %w", err)
	}

	// 3. İlişkili kayıt ID'lerini çıkar
	relatedIDs := []interface{}{}
	relatedIDSet := make(map[interface{}]bool)
	for _, pivot := range pivotRecords {
		relatedID := pivot[relatedColumn]
		if relatedID != nil && !relatedIDSet[relatedID] {
			relatedIDs = append(relatedIDs, relatedID)
			relatedIDSet[relatedID] = true
		}
	}

	if len(relatedIDs) == 0 {
		// Hiç ilişkili kayıt yok, boş listeler set et
		for _, item := range itemsByOwnerKey {
			if err := setRelationshipData(item, field.GetRelationshipName(), []map[string]interface{}{}); err != nil {
				fmt.Printf("[WARN] Failed to set BelongsToMany relationship: %v\n", err)
			}
		}
		return nil
	}

	// 4. İlişkili tablodan kayıtları çek
	safeRelatedTable := SanitizeColumnName(relatedTable)
	var relatedRecords []map[string]interface{}
	err = l.db.WithContext(ctx).
		Table(safeRelatedTable).
		Where("id IN ?", relatedIDs).
		Find(&relatedRecords).Error

	if err != nil {
		return fmt.Errorf("failed to load BelongsToMany related records: %w", err)
	}

	// 5. İlişkili kayıtları ID'ye göre map et
	relatedByID := map[interface{}]map[string]interface{}{}
	for _, record := range relatedRecords {
		id := record["id"]
		relatedByID[id] = record
	}

	// 6. Pivot kayıtlarını kullanarak ilişkileri grupla
	relatedByParentID := map[interface{}][]map[string]interface{}{}
	for _, pivot := range pivotRecords {
		parentID := pivot[parentColumn]
		relatedID := pivot[relatedColumn]
		if record, ok := relatedByID[relatedID]; ok {
			relatedByParentID[parentID] = append(relatedByParentID[parentID], record)
		}
	}

	// 7. Her item'a ilişkili kayıt listesini set et
	for ownerValue, item := range itemsByOwnerKey {
		relatedList := relatedByParentID[ownerValue]
		if relatedList == nil {
			relatedList = []map[string]interface{}{} // Boş liste
		}
		if err := setRelationshipData(item, field.GetRelationshipName(), relatedList); err != nil {
			// Log error but continue
			fmt.Printf("[WARN] Failed to set BelongsToMany relationship: %v\n", err)
		}
	}

	return nil
}

// lazyLoadBelongsToMany, tek bir kayıt için BelongsToMany ilişkisini yükler.
//
// Bu metod, GORM Association API kullanarak tek bir kaydın BelongsToMany
// ilişkisini yükler. Reflection kullanarak dinamik olarak struct tipini belirler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi
// - **item**: İlişkisi yüklenecek kayıt
// - **field**: BelongsToMany field tanımı
//
// # Döndürür
//
// - interface{}: Yüklenen ilişki verisi (slice)
// - error: Hata durumunda hata mesajı
func (l *GormRelationshipLoader) lazyLoadBelongsToMany(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}

	// Reflection kullanarak relationship field tipini al
	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	relField := itemValue.FieldByName(field.GetRelationshipName())
	if !relField.IsValid() {
		return nil, fmt.Errorf("relationship field %s not found", field.GetRelationshipName())
	}

	relType := relField.Type()
	if relType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("BelongsToMany relationship field must be a slice")
	}

	// Yeni slice instance oluştur
	relValue := reflect.New(relType)

	// GORM Association API kullanarak ilişkiyi yükle
	err := l.db.WithContext(ctx).
		Model(item).
		Association(field.GetRelationshipName()).
		Find(relValue.Interface())

	if err != nil {
		return nil, fmt.Errorf("failed to load BelongsToMany relationship: %w", err)
	}

	// İlişki verisini set et
	actualValue := relValue.Elem().Interface()
	if err := setRelationshipData(item, field.GetRelationshipName(), actualValue); err != nil {
		return nil, err
	}

	return actualValue, nil
}
