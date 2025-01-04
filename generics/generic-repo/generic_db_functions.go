package generic_repo

import (
	"reflect"
    "youtube-clone/pkg/common/database"
	"errors"
	"gorm.io/gorm"
)

type OnlyStructs interface {
	~struct{}
}

func SaveRepo[T any](tx *gorm.DB, model T) (T, error) {
	result := tx.Create(&model)
	return model, result.Error
}

func UpdateRepo[T any](tx *gorm.DB, model T) (T, error) {
	result := tx.Model(&model).
		Updates(&model)

	return model, result.Error
}

func FindSingleByField[T any](field string, value interface{}) (*T, error) {
	var record T
	if err := database.DB.
		Where(field+" = ?", value).
		First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, return nil without an error
			return nil, nil
		}
		// Other errors occurred, return the error
		return nil, err
	}
	return &record, nil
}

func FindAll[T any]() ([]T, error) {
	var records []T
	if err := database.DB.Find(&records).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found, return an empty slice and nil error
			return []T{}, nil
		}
		// Other errors occurred, return the error
		return nil, err
	}
	return records, nil
}

func DeleteByStructRepo[T any](tx *gorm.DB, model T) error {
	// Use reflection to set the 'IsDeleting' field to true, if it exists
	val := reflect.ValueOf(&model)

	// If model is a pointer, get the element it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if 'IsDeleting' field exists and is a boolean
	if field := val.FieldByName("IsDeleting"); field.IsValid() && field.Kind() == reflect.Bool {
		field.SetBool(true)
	}

	result := tx.Model(&model).Updates(&model)
	return result.Error
}
