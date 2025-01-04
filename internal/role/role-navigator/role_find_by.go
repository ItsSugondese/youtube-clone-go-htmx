package role_navigator

import (
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/internal/role/model"
	"youtube-clone/pkg/common/localization"
)

func FindRoleByIdService(name string) model.Role {
	role, getRoleError := generic_repo.FindSingleByField[model.Role]("id", name)
	if getRoleError != nil {
		panic(getRoleError)
	}

	if role == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.ROLE,
			"Second": "Id",
		}))
	}
	return *role
}

func CheckRoleExistValidationService(name string) bool {
	roleDetails, err := generic_repo.FindSingleByField[model.Role]("id", name)
	if err != nil {
		panic(err)
	}
	if roleDetails == nil {
		return false
	}
	return true
}
