package model

import (
	generic_models "youtube-clone/generics/generic-models"
	"youtube-clone/internal/role/model"
)

type BaseUser struct {
	generic_models.AuditModel

	FullName    *string      `json:"fullName"`
	Email       string       `json:"email" gorm:"unique;not null"`
	Password    string       `json:"password"`
	UserType    string       `json:"userType"`
	Roles       []model.Role `json:"role" gorm:"many2many:user_role;association_autoupdate:false;association_autocreate:false"`
	ProfilePath *string      `json:"profilePath"`
}

func (b *BaseUser) HasAuditModel() bool {
	return true
}
