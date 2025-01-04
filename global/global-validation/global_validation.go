package global_validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"reflect"
)

// Generic custom validation function
func RequiredIfIdNil(fl validator.FieldLevel) bool {
	// Get the name of the field that this validation depends on
	conditionFieldName := fl.Param()

	// Get the value of the struct that is being validated
	parentStruct := fl.Top().Elem()

	// Get the field value of the conditional field (e.g. `ID`)
	conditionField := parentStruct.FieldByName(conditionFieldName)
	if !conditionField.IsValid() {
		// If the field doesn't exist, return false (invalid)
		return false
	}

	// Check if the conditional field is uuid.Nil (assuming it's UUID type)
	if conditionField.Interface() == uuid.Nil {
		// If the conditional field is uuid.Nil, the current field must not be uuid.Nil
		if fl.Field().Interface() == uuid.Nil {
			return false // Fail validation if current field is also uuid.Nil
		}
	}
	// No validation error
	return true
}

// Generic custom validation function
func RequiredIfIdNilNotUUID(fl validator.FieldLevel) bool {
	// Get the name of the field that this validation depends on
	conditionFieldName := fl.Param()

	// Get the value of the struct that is being validated
	parentStruct := fl.Top().Elem()

	// Get the field value of the conditional field (e.g. `ID`)
	conditionField := parentStruct.FieldByName(conditionFieldName)
	if !conditionField.IsValid() {
		// If the field doesn't exist, return false (invalid)
		return false
	}

	// Check if the conditional field is uuid.Nil (assuming it's UUID type)
	if conditionField.Kind() == reflect.Ptr {
		if conditionField.IsNil() {
			if fl.Field().Interface() == uuid.Nil {
				return false // Fail validation if current field is also uuid.Nil
			}
		}
	}
	// No validation error
	return true
}
