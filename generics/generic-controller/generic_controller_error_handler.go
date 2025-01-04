package generic_controller

import (
	globaldto "youtube-clone/global/global_dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ControllerValidationHandler(T any, c *gin.Context, validate *validator.Validate) error {
	if err := c.ShouldBindJSON(&T); err != nil {
		panic(&globaldto.PanicObject{
			Err:        err,
			StatusCode: http.StatusBadRequest,
		})
	}

	if err := validate.Struct(T); err != nil {
		panic(&globaldto.PanicObject{
			Err:        err,
			StatusCode: http.StatusBadRequest,
		})
	}
	return nil
}
