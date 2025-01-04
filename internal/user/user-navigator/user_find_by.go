package user_navigator

import (
	"github.com/google/uuid"
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/internal/user/model"
	"youtube-clone/internal/user/repo"
	"youtube-clone/pkg/common/localization"
)

func FindUserByIdService(id uuid.UUID) model.BaseUser {
	userDetails, err := generic_repo.FindSingleByField[model.BaseUser]("id", id)

	if err != nil {
		panic(err)
	}

	if userDetails == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": "Id",
		}))
	}
	return *userDetails
}

func FindUserByEmailService(email string) *model.BaseUser {
	userDetails, err := repo.FindUserByColumnRepo(email, "email")
	if err != nil {
		panic(err)
	}
	if userDetails == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": "email",
		}))
	}
	return userDetails
}

func CheckUserByEmailExistValidationService(email string) bool {
	userDetails, err := repo.FindUserByColumnRepo(email, "email")
	if err != nil {
		panic(err)
	}
	if userDetails == nil {
		return false
	}
	return true
}
