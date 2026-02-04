package fields

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// ValidationRule represents a validation rule for a field
type ValidationRule struct {
	Name       string        // Rule name (required, email, min, etc.)
	Parameters []interface{} // Rule parameters
	Message    string        // Error message
}

// ValidatorFunc is a function that validates a value
type ValidatorFunc func(value interface{}, context interface{}) error

// ValidationResult represents the result of validation
type ValidationResult struct {
	IsValid bool
	Errors  map[string][]string
}

// Built-in validators

// Required returns a validation rule for required fields
func Required() ValidationRule {
	return ValidationRule{
		Name:    "required",
		Message: "This field is required",
	}
}

// EmailRule returns a validation rule for email format
func EmailRule() ValidationRule {
	return ValidationRule{
		Name:    "email",
		Message: "This field must be a valid email address",
	}
}

// URL returns a validation rule for URL format
func URL() ValidationRule {
	return ValidationRule{
		Name:    "url",
		Message: "This field must be a valid URL",
	}
}

// Min returns a validation rule for minimum value
func Min(min interface{}) ValidationRule {
	return ValidationRule{
		Name:       "min",
		Parameters: []interface{}{min},
		Message:    fmt.Sprintf("This field must be at least %v", min),
	}
}

// Max returns a validation rule for maximum value
func Max(max interface{}) ValidationRule {
	return ValidationRule{
		Name:       "max",
		Parameters: []interface{}{max},
		Message:    fmt.Sprintf("This field must be at most %v", max),
	}
}

// MinLength returns a validation rule for minimum string length
func MinLength(length int) ValidationRule {
	return ValidationRule{
		Name:       "minLength",
		Parameters: []interface{}{length},
		Message:    fmt.Sprintf("This field must be at least %d characters", length),
	}
}

// MaxLength returns a validation rule for maximum string length
func MaxLength(length int) ValidationRule {
	return ValidationRule{
		Name:       "maxLength",
		Parameters: []interface{}{length},
		Message:    fmt.Sprintf("This field must be at most %d characters", length),
	}
}

// Pattern returns a validation rule for regex pattern matching
func Pattern(pattern string) ValidationRule {
	return ValidationRule{
		Name:       "pattern",
		Parameters: []interface{}{pattern},
		Message:    "This field format is invalid",
	}
}

// Unique returns a validation rule for database uniqueness
func Unique(table, column string) ValidationRule {
	return ValidationRule{
		Name:       "unique",
		Parameters: []interface{}{table, column},
		Message:    "This value already exists",
	}
}

// Exists returns a validation rule for database existence
func Exists(table, column string) ValidationRule {
	return ValidationRule{
		Name:       "exists",
		Parameters: []interface{}{table, column},
		Message:    "This value does not exist",
	}
}

// ValidateRequired validates that a value is not empty
func ValidateRequired(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is required")
	}

	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("value is required")
		}
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("value is required")
		}
	case []string:
		if len(v) == 0 {
			return fmt.Errorf("value is required")
		}
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	_, err := mail.ParseAddress(str)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateURL validates URL format
func ValidateURL(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("url must be a string")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	_, err := url.ParseRequestURI(str)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// ValidateMin validates minimum value
func ValidateMin(value interface{}, min interface{}) error {
	switch v := value.(type) {
	case int:
		minVal, ok := min.(int)
		if !ok {
			return fmt.Errorf("min parameter must be int")
		}
		if v < minVal {
			return fmt.Errorf("value must be at least %d", minVal)
		}
	case float64:
		minVal, ok := min.(float64)
		if !ok {
			return fmt.Errorf("min parameter must be float64")
		}
		if v < minVal {
			return fmt.Errorf("value must be at least %f", minVal)
		}
	case string:
		minVal, ok := min.(int)
		if !ok {
			return fmt.Errorf("min parameter must be int for string length")
		}
		if len(v) < minVal {
			return fmt.Errorf("value must be at least %d characters", minVal)
		}
	}

	return nil
}

// ValidateMax validates maximum value
func ValidateMax(value interface{}, max interface{}) error {
	switch v := value.(type) {
	case int:
		maxVal, ok := max.(int)
		if !ok {
			return fmt.Errorf("max parameter must be int")
		}
		if v > maxVal {
			return fmt.Errorf("value must be at most %d", maxVal)
		}
	case float64:
		maxVal, ok := max.(float64)
		if !ok {
			return fmt.Errorf("max parameter must be float64")
		}
		if v > maxVal {
			return fmt.Errorf("value must be at most %f", maxVal)
		}
	case string:
		maxVal, ok := max.(int)
		if !ok {
			return fmt.Errorf("max parameter must be int for string length")
		}
		if len(v) > maxVal {
			return fmt.Errorf("value must be at most %d characters", maxVal)
		}
	}

	return nil
}

// ValidateMinLength validates minimum string length
func ValidateMinLength(value interface{}, length int) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("minLength validation requires string value")
	}

	if len(str) < length {
		return fmt.Errorf("value must be at least %d characters", length)
	}

	return nil
}

// ValidateMaxLength validates maximum string length
func ValidateMaxLength(value interface{}, length int) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("maxLength validation requires string value")
	}

	if len(str) > length {
		return fmt.Errorf("value must be at most %d characters", length)
	}

	return nil
}

// ValidatePattern validates regex pattern matching
func ValidatePattern(value interface{}, pattern string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("pattern validation requires string value")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %v", err)
	}

	if !matched {
		return fmt.Errorf("value does not match required pattern")
	}

	return nil
}

// ApplyValidationRule applies a validation rule to a value
func ApplyValidationRule(rule ValidationRule, value interface{}) error {
	switch rule.Name {
	case "required":
		return ValidateRequired(value)
	case "email":
		return ValidateEmail(value)
	case "url":
		return ValidateURL(value)
	case "min":
		if len(rule.Parameters) > 0 {
			return ValidateMin(value, rule.Parameters[0])
		}
	case "max":
		if len(rule.Parameters) > 0 {
			return ValidateMax(value, rule.Parameters[0])
		}
	case "minLength":
		if len(rule.Parameters) > 0 {
			if length, ok := rule.Parameters[0].(int); ok {
				return ValidateMinLength(value, length)
			}
		}
	case "maxLength":
		if len(rule.Parameters) > 0 {
			if length, ok := rule.Parameters[0].(int); ok {
				return ValidateMaxLength(value, length)
			}
		}
	case "pattern":
		if len(rule.Parameters) > 0 {
			if pattern, ok := rule.Parameters[0].(string); ok {
				return ValidatePattern(value, pattern)
			}
		}
	case "unique":
		// Unique validation requires database access, implement in handler
		return nil
	case "exists":
		// Exists validation requires database access, implement in handler
		return nil
	}

	return nil
}

// ConditionalValidator represents a validator that applies conditionally
type ConditionalValidator struct {
	Condition func(context interface{}) bool
	Validator ValidatorFunc
}

// ApplyConditionalValidation applies a validator only if condition is met
func ApplyConditionalValidation(validator ConditionalValidator, value interface{}, context interface{}) error {
	if validator.Condition(context) {
		return validator.Validator(value, context)
	}
	return nil
}

// CustomValidatorRegistry stores custom validators
type CustomValidatorRegistry struct {
	validators map[string]ValidatorFunc
}

// NewCustomValidatorRegistry creates a new custom validator registry
func NewCustomValidatorRegistry() *CustomValidatorRegistry {
	return &CustomValidatorRegistry{
		validators: make(map[string]ValidatorFunc),
	}
}

// Register registers a custom validator
func (r *CustomValidatorRegistry) Register(name string, validator ValidatorFunc) {
	r.validators[name] = validator
}

// Get retrieves a custom validator
func (r *CustomValidatorRegistry) Get(name string) (ValidatorFunc, bool) {
	validator, ok := r.validators[name]
	return validator, ok
}

// Apply applies a custom validator
func (r *CustomValidatorRegistry) Apply(name string, value interface{}, context interface{}) error {
	validator, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("custom validator '%s' not found", name)
	}
	return validator(value, context)
}

// Global custom validator registry
var globalValidatorRegistry = NewCustomValidatorRegistry()

// RegisterCustomValidator registers a global custom validator
func RegisterCustomValidator(name string, validator ValidatorFunc) {
	globalValidatorRegistry.Register(name, validator)
}

// GetCustomValidator retrieves a global custom validator
func GetCustomValidator(name string) (ValidatorFunc, bool) {
	return globalValidatorRegistry.Get(name)
}

// ApplyCustomValidator applies a global custom validator
func ApplyCustomValidator(name string, value interface{}, context interface{}) error {
	return globalValidatorRegistry.Apply(name, value, context)
}
