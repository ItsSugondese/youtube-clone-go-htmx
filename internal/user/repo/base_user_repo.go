package repo

import (
	"errors"
	"github.com/google/uuid"
	"youtube-clone/internal/user/model"
	"youtube-clone/pkg/common/database"

	"gorm.io/gorm"
)

func SaveBaseUser(tx *gorm.DB, user model.BaseUser) (*model.BaseUser, error) {
	result := tx.Create(&user)
	return &user, result.Error
}

func UpdateBaseUser(tx *gorm.DB, user model.BaseUser) (model.BaseUser, error) {
	result := tx.Model(&user).
		Updates(&user)
	return user, result.Error
}

func DoesUserExists(userType string, id uuid.UUID) (bool, error) {
	var exists bool
	query := `
		SELECT COALESCE((SELECT true FROM ? WHERE user_id = ?), false);`

	// Perform the raw query with GORM
	err := database.DB.Raw(query, gorm.Expr(userType+"s"), id).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func FindUserByPhoneNumberRepo(phoneNumber string) (user model.BaseUser, err error) {
	if err := database.DB.
		Where("phone_number = ?", phoneNumber).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, return zero value of Staff and nil error
			return model.BaseUser{}, nil
		}
		// Other errors occurred, return the error
		return model.BaseUser{}, err
	}

	// Staff found, return the Staff and nil error
	return user, nil
}

func FindUserByColumnRepo(value any, column string) (user *model.BaseUser, err error) {
	if err := database.DB.
		Where(column+"= ?", value).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, return zero value of Staff and nil error
			return nil, nil
		}
		// Other errors occurred, return the error
		return nil, err
	}

	// Staff found, return the Staff and nil error
	return user, nil
}
