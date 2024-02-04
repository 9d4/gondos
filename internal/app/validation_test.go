package app

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateFields(t *testing.T) {
	validate := validator.New()

	t.Run("Validating Fields", func(t *testing.T) {
		fields := map[string]interface{}{
			"Name":  "John",
			"Email": "john@example.com",
		}

		rules := map[string]string{
			"Name":  "required,min=3",
			"Email": "required,email",
		}

		err := validateFields(validate, fields, rules)
		assert.Nil(t, err, "Expected no validation errors, but got %v", err)
	})

	t.Run("Validating Fields with Errors", func(t *testing.T) {
		fields := map[string]interface{}{
			"Name":  "Jo",
			"Email": "john@example",
		}

		rules := map[string]string{
			"Name":  "required,min=3",
			"Email": "required,email",
		}

		err := validateFields(validate, fields, rules)
		assert.NotNil(t, err, "Expected validation errors, but got nil")
		assert.IsType(t, ValidationErrors{}, err, "Expected ValidationErrors type, but got %T", err)

		validationErrors := err.(ValidationErrors)
		assert.Equal(t, 2, len(validationErrors), "Expected 2 validation errors, but got %d", len(validationErrors))

		for _, ve := range validationErrors {
			assert.NotEmpty(t, ve.FieldName(), "Expected non-empty field name in ValidationError")
		}
	})
}
