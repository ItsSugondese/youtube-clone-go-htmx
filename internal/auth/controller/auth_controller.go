package controller

import (
	"strings"
	response_crud_enum "youtube-clone/enums/interface-enums/response/response-crud-enum"
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_controller "youtube-clone/generics/generic-controller"
	"youtube-clone/internal/auth/dto"
	"youtube-clone/internal/auth/service"
	"youtube-clone/pkg/common/localization"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary login to the site using this api
// @Schemes
// @Description
// @Tags Auth
// @Accept json
// @Produce json
// @Param driver body dto.AuthRequest true "Driver details"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/login [post]
func LoginUser(c *gin.Context, validate *validator.Validate) {
	var authRequest dto.AuthRequest
	if err := generic_controller.ControllerValidationHandler(&authRequest, c, validate); err != nil {
		return
	}

	authResponse := service.LoginService(authRequest)

	generic_controller.GenericControllerSuccessResponseHandler(c,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": strings.ToLower(response_crud_enum.Create().String()),
		}), authResponse)

}

// @Summary register OAuth Client for your site
// @Schemes
// @Description
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body dto.AuthRequest true "auth details"
// @Success 200 {object} dto.OAuth2ClientResponse
// @Router /auth/register-client/oauth [post]
func RegisterOAuth2Client(ctx *gin.Context, validate *validator.Validate) {
	var clientRequest dto.OAuth2ClientRequest
	if err := generic_controller.ControllerValidationHandler(&clientRequest, ctx, validate); err != nil {
		return
	}

	authResponse := service.RegisterOAuth2ClientService(ctx, clientRequest)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
			"Client Registration success", authResponse)
}

// VerifyToken is a handler function that verifies the provided token
//func VerifyToken(c *gin.Context) {
//	var authVerify AuthVerify
//	if err := c.ShouldBindJSON(&authVerify); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "type": "error", "status": 400})
//		return
//	}
//
//	verified, err := VerifyTokenService(authVerify)
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "type": "error", "status": 401})
//		return
//	}
//
//	if verified {
//
//		c.JSON(http.StatusOK, gin.H{"message": "verify successful", "status": 200, "type": "success"})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "verify faild", "status": 400, "type": "error"})
//}
