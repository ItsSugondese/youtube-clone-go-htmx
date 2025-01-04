package controller

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	response_crud_enum "youtube-clone/enums/interface-enums/response/response-crud-enum"
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_controller "youtube-clone/generics/generic-controller"
	"youtube-clone/internal/user/dto"
	"youtube-clone/internal/user/service"
	"youtube-clone/pkg/common/localization"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary register user using this api
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User details"
// @Success 200 {object} model.BaseUser
// @Router /user [post]
func RegisterUser(c *gin.Context, validate *validator.Validate) {
	var userDto dto.UserRequest

	if err := generic_controller.ControllerValidationHandler(&userDto, c, validate); err != nil {
		return
	}

	savedData := service.RegisterBaseUserService(c, userDto)
	generic_controller.GenericControllerSuccessResponseHandler(c,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": strings.ToLower(response_crud_enum.Create().String()),
		}), savedData)
}

// @Summary update user using this api
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User details"
// @Success 200 {object} model.BaseUser
// @Router /user [post]
func UpdateUser(c *gin.Context, validate *validator.Validate) {
	var userDto dto.UserRequest

	if err := generic_controller.ControllerValidationHandler(&userDto, c, validate); err != nil {
		return
	}

	savedData := service.UpdateBaseUserService(c, userDto)
	generic_controller.GenericControllerSuccessResponseHandler(c,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": strings.ToLower(response_crud_enum.Create().String()),
		}), savedData)
}

// @BasePath /api/v1

// @Summary ping example
// @Schemes
// @Description do ping
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /user/doc/:id [get]
func GetUserImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "type": "error"})
		return
	}
	service.GetUserImageService(id, c.Writer)
	// c.JSON(http.StatusOK, gin.H{"data": drivers})
}
