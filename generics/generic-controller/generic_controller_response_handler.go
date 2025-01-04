package generic_controller

import (
	"youtube-clone/enums/interface-enums/response/response-status-enum"
	globaldto "youtube-clone/global/global_dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GenericControllerSuccessResponseHandler(c *gin.Context, message string, data interface{}) {
	response := globaldto.ApiResponse{
		Status:  response_status_enum.Success(),
		Message: message,
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}
