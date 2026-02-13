package fields

import (
	"testing"
)

// TestValidationRuleApplication tests that validation rules are applied correctly
// Feature: go-field-system-phase-1, Property 2: ValidationRule Application
func TestValidationRuleApplication(t *testing.T) {
	t.Run("Required validator rejects empty values", func(t *testing.T) {
		rule := Required()
		if rule.Name != "required" {
			t.Errorf("Expected rule name 'required', got '%s'", rule.Name)
		}

		// Test with empty string
		err := ApplyValidationRule(rule, "")
		if err == nil {
			t.Error("Expected error for empty string")
		}

		// Test with nil
		err = ApplyValidationRule(rule, nil)
		if err == nil {
			t.Error("Expected error for nil")
		}

		// Test with valid value
		err = ApplyValidationRule(rule, "test")
		if err != nil {
			t.Errorf("Expected no error for valid value, got %v", err)
		}
	})

	t.Run("Email validator validates email format", func(t *testing.T) {
		rule := EmailRule()
		if rule.Name != "email" {
			t.Errorf("Expected rule name 'email', got '%s'", rule.Name)
		}

		// Test with valid email
		err := ApplyValidationRule(rule, "test@example.com")
		if err != nil {
			t.Errorf("Expected no error for valid email, got %v", err)
		}

		// Test with invalid email
		err = ApplyValidationRule(rule, "invalid-email")
		if err == nil {
			t.Error("Expected error for invalid email")
		}

		// Test with empty string (should be valid, use Required() for mandatory)
		err = ApplyValidationRule(rule, "")
		if err != nil {
			t.Errorf("Expected no error for empty email, got %v", err)
		}
	})

	t.Run("URL validator validates URL format", func(t *testing.T) {
		rule := URL()
		if rule.Name != "url" {
			t.Errorf("Expected rule name 'url', got '%s'", rule.Name)
		}

		// Test with valid URL
		err := ApplyValidationRule(rule, "https://example.com")
		if err != nil {
			t.Errorf("Expected no error for valid URL, got %v", err)
		}

		// Test with invalid URL
		err = ApplyValidationRule(rule, "not a url")
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("Min validator validates minimum value", func(t *testing.T) {
		rule := Min(10)
		if rule.Name != "min" {
			t.Errorf("Expected rule name 'min', got '%s'", rule.Name)
		}

		// Test with value above minimum
		err := ApplyValidationRule(rule, 15)
		if err != nil {
			t.Errorf("Expected no error for value above minimum, got %v", err)
		}

		// Test with value below minimum
		err = ApplyValidationRule(rule, 5)
		if err == nil {
			t.Error("Expected error for value below minimum")
		}

		// Test with value equal to minimum
		err = ApplyValidationRule(rule, 10)
		if err != nil {
			t.Errorf("Expected no error for value equal to minimum, got %v", err)
		}
	})

	t.Run("Max validator validates maximum value", func(t *testing.T) {
		rule := Max(100)
		if rule.Name != "max" {
			t.Errorf("Expected rule name 'max', got '%s'", rule.Name)
		}

		// Test with value below maximum
		err := ApplyValidationRule(rule, 50)
		if err != nil {
			t.Errorf("Expected no error for value below maximum, got %v", err)
		}

		// Test with value above maximum
		err = ApplyValidationRule(rule, 150)
		if err == nil {
			t.Error("Expected error for value above maximum")
		}

		// Test with value equal to maximum
		err = ApplyValidationRule(rule, 100)
		if err != nil {
			t.Errorf("Expected no error for value equal to maximum, got %v", err)
		}
	})

	t.Run("MinLength validator validates minimum string length", func(t *testing.T) {
		rule := MinLength(5)
		if rule.Name != "minLength" {
			t.Errorf("Expected rule name 'minLength', got '%s'", rule.Name)
		}

		// Test with string longer than minimum
		err := ApplyValidationRule(rule, "hello world")
		if err != nil {
			t.Errorf("Expected no error for string longer than minimum, got %v", err)
		}

		// Test with string shorter than minimum
		err = ApplyValidationRule(rule, "hi")
		if err == nil {
			t.Error("Expected error for string shorter than minimum")
		}

		// Test with string equal to minimum
		err = ApplyValidationRule(rule, "hello")
		if err != nil {
			t.Errorf("Expected no error for string equal to minimum, got %v", err)
		}
	})

	t.Run("MaxLength validator validates maximum string length", func(t *testing.T) {
		rule := MaxLength(10)
		if rule.Name != "maxLength" {
			t.Errorf("Expected rule name 'maxLength', got '%s'", rule.Name)
		}

		// Test with string shorter than maximum
		err := ApplyValidationRule(rule, "hello")
		if err != nil {
			t.Errorf("Expected no error for string shorter than maximum, got %v", err)
		}

		// Test with string longer than maximum
		err = ApplyValidationRule(rule, "hello world this is too long")
		if err == nil {
			t.Error("Expected error for string longer than maximum")
		}

		// Test with string equal to maximum
		err = ApplyValidationRule(rule, "1234567890")
		if err != nil {
			t.Errorf("Expected no error for string equal to maximum, got %v", err)
		}
	})

	t.Run("Pattern validator validates regex pattern", func(t *testing.T) {
		rule := Pattern("^[a-z]+$")
		if rule.Name != "pattern" {
			t.Errorf("Expected rule name 'pattern', got '%s'", rule.Name)
		}

		// Test with matching pattern
		err := ApplyValidationRule(rule, "hello")
		if err != nil {
			t.Errorf("Expected no error for matching pattern, got %v", err)
		}

		// Test with non-matching pattern
		err = ApplyValidationRule(rule, "Hello123")
		if err == nil {
			t.Error("Expected error for non-matching pattern")
		}
	})

	t.Run("Multiple validation rules are applied in order", func(t *testing.T) {
		rules := []ValidationRule{
			Required(),
			MinLength(5),
			MaxLength(20),
		}

		// Test with valid value
		value := "hello"
		for _, rule := range rules {
			err := ApplyValidationRule(rule, value)
			if err != nil {
				t.Errorf("Expected no error for valid value, got %v", err)
			}
		}

		// Test with value that fails first rule
		value = ""
		err := ApplyValidationRule(rules[0], value)
		if err == nil {
			t.Error("Expected error for empty value")
		}

		// Test with value that fails second rule
		value = "hi"
		err = ApplyValidationRule(rules[1], value)
		if err == nil {
			t.Error("Expected error for value shorter than minimum")
		}

		// Test with value that fails third rule
		value = "this is a very long string that exceeds the maximum"
		err = ApplyValidationRule(rules[2], value)
		if err == nil {
			t.Error("Expected error for value longer than maximum")
		}
	})
}

// TestBuiltInValidatorsAvailability tests that all built-in validators are available
// Feature: go-field-system-phase-1, Property 3: Built-in Validators Availability
func TestBuiltInValidatorsAvailability(t *testing.T) {
	validators := []struct {
		name string
		rule ValidationRule
	}{
		{"Required", Required()},
		{"Email", EmailRule()},
		{"URL", URL()},
		{"Min", Min(10)},
		{"Max", Max(100)},
		{"MinLength", MinLength(5)},
		{"MaxLength", MaxLength(20)},
		{"Pattern", Pattern("^[a-z]+$")},
		{"Unique", Unique("users", "email")},
		{"Exists", Exists("users", "id")},
	}

	for _, v := range validators {
		t.Run(v.name+" validator is available", func(t *testing.T) {
			if v.rule.Name == "" {
				t.Errorf("Validator %s has empty name", v.name)
			}

			if v.rule.Message == "" {
				t.Errorf("Validator %s has empty message", v.name)
			}

			// All validators should be callable
			_ = ApplyValidationRule(v.rule, "test")
		})
	}
}
