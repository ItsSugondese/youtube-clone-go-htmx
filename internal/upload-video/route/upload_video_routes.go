package route

import (
	"youtube-clone/internal/upload-video/controller"
	repo2 "youtube-clone/internal/upload-video/repo"
    "youtube-clone/internal/upload-video/service"
	authentication_middleware "youtube-clone/pkg/middleware/authentication-middleware"
	paseto_token "youtube-clone/pkg/utils/paseto-token"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func UploadVideoRoutes(r *gin.Engine, validate *validator.Validate, db *gorm.DB) {
    repo := repo2.NewUploadVideoRepo(db)
	services := service.NewUploadVideoService(repo)
	controllers := controller.NewUploadVideoController(services)
	uploadVideoRouting := r.Group("/upload-video", authentication_middleware.PasetoAuthMiddleware(*paseto_token.TokenMaker))
	{
		uploadVideoRouting.POST("", func(c *gin.Context) {
			controllers.SaveUploadVideo(c, validate)
		})
		uploadVideoRouting.POST("/paginated", func(c *gin.Context) {
            controllers.GetAllUploadVideoDetailsPaginated(c, validate)
        })
		uploadVideoRouting.GET("", controllers.GetAllUploadVideoDetails)
		uploadVideoRouting.GET("/:id", controllers.GetUploadVideoDetailsById)
		uploadVideoRouting.DELETE("/:id", controllers.DeleteUploadVideoById)
	}
}
