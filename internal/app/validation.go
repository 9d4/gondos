package app

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	return "validation errors"
}

// ValidationError wraps validator.ValidationErrors with the addition of FieldName.
// Because using validate.Var we can't get the field name of what field is not validating.
type ValidationError struct {
	validator.ValidationErrors
	fieldName string
}

func (ve ValidationError) FieldName() string {
	return ve.fieldName
}

func newValidationError(fieldName string, err error) error {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	customVe := ValidationError{
		fieldName:        fieldName,
		ValidationErrors: ve,
	}
	return customVe
}

func validateField(v *validator.Validate, fieldName string, field interface{}, tag string) error {
	err := v.Var(field, tag)
	if err != nil {
		return newValidationError(fieldName, err)
	}
	return nil
}

// validateFields runs validateField in parallel.
// The returned error type is ValidationErrors
func validateFields(v *validator.Validate, fields map[string]interface{}, rules map[string]string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(fields))

	for fieldName, value := range fields {
		wg.Add(1)
		go func(v *validator.Validate, fn string, val interface{}, tag string) {
			if err := validateField(v, fn, val, tag); err != nil {
				errCh <- err
			}
			wg.Done()
		}(v, fieldName, value, rules[fieldName])
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errors ValidationErrors
	for err := range errCh {
		if ve, ok := err.(ValidationError); ok {
			errors = append(errors, ve)
		}
	}

	if len(errors) != 0 {
		return errors
	}
	return nil
}
