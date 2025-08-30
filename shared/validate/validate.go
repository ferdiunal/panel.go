package validate

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type ValidationError map[string]map[string]string

func ValidateStruct(s interface{}) ValidationError {
	errors := make(map[string]map[string]string)
	err := Validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[strings.ToLower(err.Field())] = map[string]string{
				"message": err.Tag(),
			}
		}
	}

	return errors
}
