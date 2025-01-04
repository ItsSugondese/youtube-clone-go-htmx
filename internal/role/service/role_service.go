package service

import (
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/internal/role/dto"
	"youtube-clone/internal/role/model"
	role_navigator "youtube-clone/internal/role/role-navigator"
	"youtube-clone/pkg/common/database"
)

func CreateRoleService(dto *dto.RoleRequest) *model.Role {
	tx := database.DB.Begin()

	exists := role_navigator.CheckRoleExistValidationService(dto.Name)

	if exists {
		panic("Role already exists")
	}

	savedRole, saveRoleError := generic_repo.SaveRepo(tx, &model.Role{ID: dto.Name})
	if saveRoleError != nil {
		panic(saveRoleError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
	return savedRole
}
