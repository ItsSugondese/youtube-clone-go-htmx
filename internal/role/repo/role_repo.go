package repo

import (
	"gorm.io/gorm"
	"youtube-clone/internal/role/model"
)

func SaveRoleRepo(tx *gorm.DB, role *model.Role) (*model.Role, error) {
	result := tx.Create(&role)
	return role, result.Error
}
