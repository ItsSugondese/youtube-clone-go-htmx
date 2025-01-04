package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"youtube-clone/internal/auth/controller"
)

func AuthRoutes(r *gin.Engine, validate *validator.Validate) {
	auths := r.Group("/auth/")
	{
		auths.POST("login", func(c *gin.Context) {
			controller.LoginUser(c, validate)
		})
		auths.POST("register-client/oauth", func(c *gin.Context) {
        			controller.RegisterOAuth2Client(c, validate)
        })
		//auths.POST("verify", VerifyToken)

	}
}
