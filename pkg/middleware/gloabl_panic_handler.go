package middleware

import (
	response_status_enum "youtube-clone/enums/interface-enums/response/response-status-enum"
	globaldto "youtube-clone/global/global_dto"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				switch err := err.(type) {

				case *globaldto.PanicObject:
					// Handle the custom panic
					response := &globaldto.ApiResponse{
						Status:  response_status_enum.Fail(),
						Message: fmt.Sprintf("%v", err.Err),
						Data:    []string{err.Err.Error()},
					}
					c.JSON(err.StatusCode, response)
				default:
					// Handle generic panic
					// Log the error
					fmt.Fprintf(os.Stderr, "Panic recovered: %s\n", err)

					errors := recoverAndConvertToErrors(err)

					response := &globaldto.ApiResponse{
						Status:  response_status_enum.Fail(),
						Message: fmt.Sprintf("%v", err),
						Data:    errors,
					}
					c.JSON(http.StatusInternalServerError, response)
				}

				//c.AbortWithStatus(http.StatusInternalServerError)
				c.Abort()
			}
		}()

		// Process the next handler in the chain
		c.Next()
	}
}

func recoverAndConvertToErrors(T any) []string {
	// Check if T is already a slice of errors
	if errList, ok := T.([]error); ok {
		strErrors := make([]string, len(errList))
		for i, err := range errList {
			strErrors[i] = err.Error()
		}
		return strErrors
	}

	// Check if T is a single error
	if singleErr, ok := T.(error); ok {
		return []string{singleErr.Error()}
	}

	// Handle other types by converting to a single error in a slice
	return []string{fmt.Sprintf("%v", T)}
}
