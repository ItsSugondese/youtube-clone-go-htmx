package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	response_crud_enum "youtube-clone/enums/interface-enums/response/response-crud-enum"
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_controller "youtube-clone/generics/generic-controller"
	"youtube-clone/internal/temporary-attachments/dto"
	"youtube-clone/internal/temporary-attachments/service"
	"youtube-clone/pkg/common/localization"
)

// @Summary Get list of attachment id
// @Schemes
// @Description do ping
// @Tags Temporary Attachments
// @Accept multipart/form-data
// @Produce json
// @Success 200 {array} int
// @Router /temporary-attachments [post]
// post /temporary-attachments
func CreateTemporaryAttachments(c *gin.Context) {
	var attachmentsDetailRequest dto.TemporaryAttachmentsDetailRequest
	if err := c.ShouldBind(&attachmentsDetailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the uploaded file

	// calling service

	savedData := service.SaveTemporaryAttachmentsService(c)

	generic_controller.GenericControllerSuccessResponseHandler(c,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": strings.ToLower(response_crud_enum.Create().String()),
		}), savedData)
}
