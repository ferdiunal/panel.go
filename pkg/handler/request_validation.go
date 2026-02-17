package handler

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	panelcontext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	paneli18n "github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	validationErrorCode      = "VALIDATION_ERROR"
	validationMessagesProp   = "validation_messages"
	validationMessagesAlias  = "messages"
	validationDefaultMessage = "Invalid value"
)

var (
	requestPayloadValidator = newRequestPayloadValidator()
	safeIdentifierRegexp    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.]*$`)
)

type requestValidationErrors struct {
	fieldErrors map[string][]string
}

func newRequestValidationErrors() *requestValidationErrors {
	return &requestValidationErrors{
		fieldErrors: make(map[string][]string),
	}
}

func (e *requestValidationErrors) add(field string, message string) {
	field = strings.TrimSpace(field)
	message = strings.TrimSpace(message)
	if field == "" || message == "" {
		return
	}
	e.fieldErrors[field] = append(e.fieldErrors[field], message)
}

func (e *requestValidationErrors) hasAny() bool {
	return len(e.fieldErrors) > 0
}

func (e *requestValidationErrors) response(c *fiber.Ctx) fiber.Map {
	message := paneli18n.TransWithFallback(c, "error.validationError", "Validation error")
	return fiber.Map{
		"error":   message,
		"code":    validationErrorCode,
		"errors":  e.fieldErrors,
		"details": e.fieldErrors,
	}
}

func newRequestPayloadValidator() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("panel_regex", validatePanelRegex)
	return v
}

func validatePanelRegex(fl validator.FieldLevel) bool {
	pattern := strings.TrimSpace(fl.Param())
	if pattern == "" {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(pattern)
	if err != nil {
		return false
	}

	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return true
	}

	matched, err := regexp.MatchString(string(decoded), value)
	return err == nil && matched
}

func (h *FieldHandler) validateCreatePayload(c *panelcontext.Context, payload map[string]interface{}) *requestValidationErrors {
	return h.validateRequestPayload(c, payload, fields.ContextCreate, "")
}

func (h *FieldHandler) validateUpdatePayload(c *panelcontext.Context, id string, payload map[string]interface{}) *requestValidationErrors {
	return h.validateRequestPayload(c, payload, fields.ContextUpdate, id)
}

func (h *FieldHandler) validateRequestPayload(
	c *panelcontext.Context,
	payload map[string]interface{},
	visibilityCtx fields.VisibilityContext,
	recordID string,
) *requestValidationErrors {
	if c == nil {
		return nil
	}

	resourceCtx := c.Resource()
	if resourceCtx != nil {
		resourceCtx.VisibilityCtx = visibilityCtx
	}

	elements := h.resolveValidationElements(c, resourceCtx)
	if len(elements) == 0 {
		return nil
	}

	validationErrors := newRequestValidationErrors()
	db := h.resolveValidationDB(c)

	for _, element := range elements {
		if element == nil {
			continue
		}
		if resourceCtx != nil && !element.IsVisible(resourceCtx) {
			continue
		}
		if setter, ok := element.(interface{ SetContextForI18n(*fiber.Ctx) }); ok {
			setter.SetContextForI18n(c.Ctx)
		}

		key := strings.TrimSpace(element.GetKey())
		if key == "" {
			continue
		}

		serialized := element.JsonSerialize()
		if readOnly, ok := serialized["read_only"].(bool); ok && readOnly {
			continue
		}
		if disabled, ok := serialized["disabled"].(bool); ok && disabled {
			continue
		}

		value, hasValue := payload[key]
		rules := collectFieldValidationRules(element, serialized)
		customValidators := collectFieldCustomValidators(element)
		if len(rules) == 0 && len(customValidators) == 0 {
			continue
		}

		fieldLabel := resolveValidationFieldLabel(element, serialized, key)
		messageOverrides := resolveValidationMessageOverrides(serialized["props"])

		for _, rule := range rules {
			if shouldSkipValidationRule(rule, hasValue, value, visibilityCtx) {
				continue
			}

			message := h.runValidationRule(c, db, element, rule, fieldLabel, key, value, recordID, messageOverrides)
			if message != "" {
				validationErrors.add(key, message)
			}
		}

		if !hasValue || isEmptyValidationValue(value) {
			continue
		}

		for _, customValidator := range customValidators {
			if err := customValidator(value, c); err != nil {
				validationErrors.add(key, err.Error())
			}
		}
	}

	if !validationErrors.hasAny() {
		return nil
	}

	return validationErrors
}

func (h *FieldHandler) resolveValidationElements(c *panelcontext.Context, resourceCtx *core.ResourceContext) []fields.Element {
	if resourceCtx != nil && len(resourceCtx.Elements) > 0 {
		elements := make([]fields.Element, 0, len(resourceCtx.Elements))
		for _, element := range resourceCtx.Elements {
			if element == nil {
				continue
			}
			elements = append(elements, element)
		}
		if len(elements) > 0 {
			return elements
		}
	}

	return h.getElements(c)
}

func (h *FieldHandler) runValidationRule(
	c *panelcontext.Context,
	db *gorm.DB,
	element fields.Element,
	rule fields.ValidationRule,
	fieldLabel string,
	fieldKey string,
	value interface{},
	recordID string,
	messageOverrides map[string]string,
) string {
	normalizedRule := normalizeValidationRuleName(rule.Name)
	message := h.resolveValidationMessage(c, rule, fieldLabel, fieldKey, messageOverrides)

	switch normalizedRule {
	case "required":
		if isEmptyValidationValue(value) {
			return message
		}
		return ""
	case "unique":
		if h.validateUniqueRule(db, rule, value, recordID) {
			return ""
		}
		return message
	case "exists":
		if h.validateExistsRule(db, rule, value) {
			return ""
		}
		return message
	}

	tag := buildValidatorTag(rule)
	if tag == "" {
		return ""
	}

	normalizedValue := normalizeValidationValue(element, value)
	if err := requestPayloadValidator.Var(normalizedValue, tag); err != nil {
		return message
	}

	return ""
}

func (h *FieldHandler) resolveValidationMessage(
	c *panelcontext.Context,
	rule fields.ValidationRule,
	fieldLabel string,
	fieldKey string,
	messageOverrides map[string]string,
) string {
	templateData := buildValidationTemplateData(fieldLabel, fieldKey, rule)

	if override := lookupValidationMessageOverride(messageOverrides, rule.Name); override != "" {
		return localizeOrFallbackMessage(c.Ctx, override, templateData)
	}

	if custom := resolveRuleCustomMessage(rule); custom != "" {
		return localizeOrFallbackMessage(c.Ctx, custom, templateData)
	}

	translationKey := "validation." + normalizeValidationTranslationKey(rule.Name)
	localized := paneli18n.Trans(c.Ctx, translationKey, templateData)
	if localized != translationKey {
		return localized
	}

	if fallback := strings.TrimSpace(rule.Message); fallback != "" {
		return fallback
	}

	return paneli18n.TransWithFallback(c.Ctx, "validation.invalid", validationDefaultMessage, templateData)
}

func localizeOrFallbackMessage(c *fiber.Ctx, message string, templateData map[string]interface{}) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}

	// If a translation key is provided as custom message, resolve it first.
	if !strings.ContainsAny(message, " \t\n") {
		localized := paneli18n.Trans(c, message, templateData)
		if localized != message {
			return localized
		}
	}

	return interpolateValidationMessage(message, templateData)
}

func collectFieldValidationRules(element fields.Element, serialized map[string]interface{}) []fields.ValidationRule {
	rawRules := element.GetValidationRules()
	rules := make([]fields.ValidationRule, 0, len(rawRules)+1)

	hasRequiredRule := false
	for _, rawRule := range rawRules {
		rule, ok := rawRule.(fields.ValidationRule)
		if !ok {
			continue
		}
		if normalizeValidationRuleName(rule.Name) == "required" {
			hasRequiredRule = true
		}
		rules = append(rules, rule)
	}

	if required, ok := serialized["required"].(bool); ok && required && !hasRequiredRule {
		rules = append([]fields.ValidationRule{{Name: "required"}}, rules...)
	}

	return rules
}

func collectFieldCustomValidators(element fields.Element) []fields.ValidatorFunc {
	rawValidators := element.GetCustomValidators()
	validators := make([]fields.ValidatorFunc, 0, len(rawValidators))
	for _, rawValidator := range rawValidators {
		validatorFunc, ok := rawValidator.(fields.ValidatorFunc)
		if !ok {
			continue
		}
		validators = append(validators, validatorFunc)
	}
	return validators
}

func resolveValidationFieldLabel(element fields.Element, serialized map[string]interface{}, fallback string) string {
	if label, ok := serialized["label"].(string); ok && strings.TrimSpace(label) != "" {
		return strings.TrimSpace(label)
	}
	if name, ok := serialized["name"].(string); ok && strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	}
	if name := strings.TrimSpace(element.GetName()); name != "" {
		return name
	}
	return fallback
}

func resolveValidationMessageOverrides(rawProps interface{}) map[string]string {
	props, ok := rawProps.(map[string]interface{})
	if !ok {
		return nil
	}

	overrides := make(map[string]string)
	merge := func(raw interface{}) {
		switch typed := raw.(type) {
		case map[string]interface{}:
			for key, value := range typed {
				message, ok := value.(string)
				if !ok || strings.TrimSpace(message) == "" {
					continue
				}
				overrides[normalizeValidationRuleName(key)] = strings.TrimSpace(message)
			}
		case map[string]string:
			for key, value := range typed {
				if strings.TrimSpace(value) == "" {
					continue
				}
				overrides[normalizeValidationRuleName(key)] = strings.TrimSpace(value)
			}
		}
	}

	merge(props[validationMessagesProp])
	merge(props[validationMessagesAlias])

	if len(overrides) == 0 {
		return nil
	}

	return overrides
}

func lookupValidationMessageOverride(overrides map[string]string, ruleName string) string {
	if len(overrides) == 0 {
		return ""
	}

	normalizedRule := normalizeValidationRuleName(ruleName)
	if message, ok := overrides[normalizedRule]; ok {
		return message
	}

	return ""
}

func shouldSkipValidationRule(
	rule fields.ValidationRule,
	hasValue bool,
	value interface{},
	visibilityCtx fields.VisibilityContext,
) bool {
	normalizedRule := normalizeValidationRuleName(rule.Name)

	if normalizedRule == "required" {
		if !hasValue && visibilityCtx == fields.ContextUpdate {
			return true
		}
		return false
	}

	if !hasValue {
		return true
	}

	return isEmptyValidationValue(value)
}

func normalizeValidationRuleName(ruleName string) string {
	normalized := strings.ToLower(strings.TrimSpace(ruleName))
	normalized = strings.ReplaceAll(normalized, "-", "_")

	switch normalized {
	case "min_length", "minlength":
		return "minlength"
	case "max_length", "maxlength":
		return "maxlength"
	default:
		return normalized
	}
}

func normalizeValidationTranslationKey(ruleName string) string {
	switch normalizeValidationRuleName(ruleName) {
	case "minlength":
		return "minLength"
	case "maxlength":
		return "maxLength"
	default:
		return normalizeValidationRuleName(ruleName)
	}
}

func resolveRuleCustomMessage(rule fields.ValidationRule) string {
	message := strings.TrimSpace(rule.Message)
	if message == "" {
		return ""
	}

	defaultMessage := defaultRuleMessage(rule)
	if defaultMessage != "" && defaultMessage == message {
		return ""
	}

	return message
}

func defaultRuleMessage(rule fields.ValidationRule) string {
	switch normalizeValidationRuleName(rule.Name) {
	case "required":
		return "This field is required"
	case "email":
		return "This field must be a valid email address"
	case "url":
		return "This field must be a valid URL"
	case "min":
		if len(rule.Parameters) > 0 {
			return fmt.Sprintf("This field must be at least %v", rule.Parameters[0])
		}
	case "max":
		if len(rule.Parameters) > 0 {
			return fmt.Sprintf("This field must be at most %v", rule.Parameters[0])
		}
	case "minlength":
		if len(rule.Parameters) > 0 {
			return fmt.Sprintf("This field must be at least %v characters", rule.Parameters[0])
		}
	case "maxlength":
		if len(rule.Parameters) > 0 {
			return fmt.Sprintf("This field must be at most %v characters", rule.Parameters[0])
		}
	case "pattern":
		return "This field format is invalid"
	case "unique":
		return "This value already exists"
	case "exists":
		return "This value does not exist"
	}
	return ""
}

func buildValidationTemplateData(fieldLabel string, fieldKey string, rule fields.ValidationRule) map[string]interface{} {
	firstParam := ""
	if len(rule.Parameters) > 0 {
		firstParam = fmt.Sprintf("%v", rule.Parameters[0])
	}

	templateData := map[string]interface{}{
		"Field":   fieldLabel,
		"Key":     fieldKey,
		"Param":   firstParam,
		"Value":   firstParam,
		"Min":     firstParam,
		"Max":     firstParam,
		"Length":  firstParam,
		"Pattern": firstParam,
	}

	if len(rule.Parameters) > 1 {
		templateData["Param2"] = fmt.Sprintf("%v", rule.Parameters[1])
	}

	return templateData
}

func buildValidatorTag(rule fields.ValidationRule) string {
	normalizedRule := normalizeValidationRuleName(rule.Name)

	switch normalizedRule {
	case "email":
		return "email"
	case "url":
		return "url"
	case "min":
		if len(rule.Parameters) == 0 {
			return ""
		}
		return "min=" + fmt.Sprintf("%v", rule.Parameters[0])
	case "max":
		if len(rule.Parameters) == 0 {
			return ""
		}
		return "max=" + fmt.Sprintf("%v", rule.Parameters[0])
	case "minlength":
		if len(rule.Parameters) == 0 {
			return ""
		}
		return "min=" + fmt.Sprintf("%v", rule.Parameters[0])
	case "maxlength":
		if len(rule.Parameters) == 0 {
			return ""
		}
		return "max=" + fmt.Sprintf("%v", rule.Parameters[0])
	case "pattern":
		if len(rule.Parameters) == 0 {
			return ""
		}
		pattern := fmt.Sprintf("%v", rule.Parameters[0])
		return "panel_regex=" + base64.StdEncoding.EncodeToString([]byte(pattern))
	default:
		return ""
	}
}

func normalizeValidationValue(element fields.Element, value interface{}) interface{} {
	switch element.GetType() {
	case fields.TYPE_NUMBER, fields.TYPE_MONEY:
		switch typed := value.(type) {
		case string:
			if parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64); err == nil {
				return parsed
			}
		}
	}

	return value
}

func isEmptyValidationValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed) == ""
	case []string:
		return len(typed) == 0
	case []interface{}:
		return len(typed) == 0
	}

	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Ptr, reflect.Interface:
		return reflectValue.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return reflectValue.Len() == 0
	}

	return false
}

func interpolateValidationMessage(message string, templateData map[string]interface{}) string {
	if len(templateData) == 0 {
		return message
	}

	replacerArgs := make([]string, 0, len(templateData)*6)
	for key, value := range templateData {
		valueStr := fmt.Sprintf("%v", value)
		keyLower := strings.ToLower(key)
		replacerArgs = append(
			replacerArgs,
			"{{."+key+"}}", valueStr,
			"{"+key+"}", valueStr,
			"{"+keyLower+"}", valueStr,
			":"+keyLower, valueStr,
		)
	}

	replacerArgs = append(replacerArgs, ":attribute", fmt.Sprintf("%v", templateData["Field"]))
	return strings.NewReplacer(replacerArgs...).Replace(message)
}

func (h *FieldHandler) resolveValidationDB(c *panelcontext.Context) *gorm.DB {
	if c != nil {
		if db := c.DB(); db != nil {
			return db
		}
	}

	if provider, ok := h.Provider.(*data.GormDataProvider); ok {
		return provider.DB
	}

	return nil
}

func (h *FieldHandler) validateUniqueRule(
	db *gorm.DB,
	rule fields.ValidationRule,
	value interface{},
	recordID string,
) bool {
	if db == nil {
		return true
	}

	table, column, ok := parseRuleTableAndColumn(rule.Parameters)
	if !ok {
		return true
	}

	values := collectValidationValues(value)
	if len(values) == 0 {
		return true
	}
	if len(values) > 1 {
		// Multi-value unique checks are not supported in this mode.
		return true
	}

	query := db.Table(table).Where(clause.Eq{
		Column: clause.Column{Name: column},
		Value:  values[0],
	})

	modelTable, primaryColumn, hasModelSchema := h.resolveProviderTableSchema(db)
	if hasModelSchema && recordID != "" && strings.EqualFold(modelTable, table) && primaryColumn != "" {
		query = query.Where(clause.Neq{
			Column: clause.Column{Name: primaryColumn},
			Value:  recordID,
		})
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return true
	}

	return count == 0
}

func (h *FieldHandler) validateExistsRule(db *gorm.DB, rule fields.ValidationRule, value interface{}) bool {
	if db == nil {
		return true
	}

	table, column, ok := parseRuleTableAndColumn(rule.Parameters)
	if !ok {
		return true
	}

	values := collectValidationValues(value)
	if len(values) == 0 {
		return true
	}

	query := db.Table(table).Distinct(column)
	if len(values) == 1 {
		query = query.Where(clause.Eq{
			Column: clause.Column{Name: column},
			Value:  values[0],
		})
	} else {
		query = query.Where(clause.IN{
			Column: clause.Column{Name: column},
			Values: values,
		})
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return true
	}

	if len(values) == 1 {
		return count > 0
	}

	return count == int64(len(uniqueValidationValues(values)))
}

func (h *FieldHandler) resolveProviderTableSchema(db *gorm.DB) (string, string, bool) {
	provider, ok := h.Provider.(*data.GormDataProvider)
	if !ok || provider == nil || provider.Model == nil || db == nil {
		return "", "", false
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(provider.Model); err != nil || stmt.Schema == nil {
		return "", "", false
	}
	if stmt.Schema.PrioritizedPrimaryField == nil {
		return stmt.Schema.Table, "", true
	}

	primaryColumn := stmt.Schema.PrioritizedPrimaryField.DBName
	if primaryColumn == "" {
		primaryColumn = stmt.Schema.PrioritizedPrimaryField.Name
	}

	return stmt.Schema.Table, primaryColumn, true
}

func parseRuleTableAndColumn(parameters []interface{}) (string, string, bool) {
	if len(parameters) < 2 {
		return "", "", false
	}

	table := strings.TrimSpace(fmt.Sprintf("%v", parameters[0]))
	column := strings.TrimSpace(fmt.Sprintf("%v", parameters[1]))
	if table == "" || column == "" {
		return "", "", false
	}

	if !safeIdentifierRegexp.MatchString(table) || !safeIdentifierRegexp.MatchString(column) {
		return "", "", false
	}

	return table, column, true
}

func collectValidationValues(value interface{}) []interface{} {
	if isEmptyValidationValue(value) {
		return nil
	}

	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
		return []interface{}{value}
	}

	values := make([]interface{}, 0, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		item := reflectValue.Index(i).Interface()
		if isEmptyValidationValue(item) {
			continue
		}
		values = append(values, item)
	}

	return values
}

func uniqueValidationValues(values []interface{}) []interface{} {
	unique := make([]interface{}, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, value := range values {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, value)
	}

	return unique
}
