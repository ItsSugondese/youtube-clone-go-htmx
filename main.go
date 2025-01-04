package main

import (
	"log"
	"youtube-clone/config"
	oauth2_setup "youtube-clone/config/oauth2-setup"
	socket_config "youtube-clone/config/socket-config"
	template_config "youtube-clone/config/template-config"
	"youtube-clone/internal/auth/route"
	roleModel "youtube-clone/internal/role/model"
	roleRoute "youtube-clone/internal/role/route"
	"youtube-clone/internal/temporary-attachments/model"
	tempAttachmentRoute "youtube-clone/internal/temporary-attachments/route"
	"youtube-clone/internal/user"
	baseUserModel "youtube-clone/internal/user/model"
	baseUserRoute "youtube-clone/internal/user/route"

	_ "youtube-clone/docs"
	global_gin_context "youtube-clone/global/global-gin-context"
	global_validation "youtube-clone/global/global-validation"
	"youtube-clone/pkg/common/database" // Add this line to import the database package
	"youtube-clone/pkg/common/localization"
	"youtube-clone/pkg/middleware"
	audit_middleware "youtube-clone/pkg/middleware/audit-middleware"
	cors_middleware "youtube-clone/pkg/middleware/cors-middleware"
	lang_middleware "youtube-clone/pkg/middleware/lang-middleware"
	paseto_token "youtube-clone/pkg/utils/paseto-token"

	// "cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const ( // FILL IN WITH YOURS
	bucketName = "blackpearlbucket" // FILL IN WITH YOURS
)

// init method runs before the main method so that the environment variables are loaded before the application starts
func init() {
	log.Println("Loading environment variables and database connection")
	// load .env
	config.LoadEnvVariables()

	// load database connection
	database.ConnectToDB()

	// load oAuth2 server
	oauth2_setup.SetUpOAuth2()

	// paseto setup
	setupPaseto()

	// global gin hanler setup
	global_gin_context.NewGlobalGinContext()

	// lang a
	localization.InitLocalizationManager()

	// google json locaiton
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./continual-mind-432410-g5-df1fc7f32718.json") // FILL IN WITH YOUR FILE PATH

	// Register the audit log callbacks and perform migrations
	errVal := audit_middleware.RegisterCallbacks(database.DB)
	if errVal != nil {
		panic("failed to register audit log callbacks")
	}

	database.DB.AutoMigrate(&baseUserModel.BaseUser{}, &roleModel.Role{}, &model.TemporaryAttachments{})

	database.InitializeValuesInDb()

}
func main() {
	hub := socket_config.NewHub()
	go hub.Run()

	log.Println("Starting the application")
	r := gin.Default()
	validate := validator.New()

	r.SetHTMLTemplate(template_config.LoadTemplates("./web/templates"))

	// Initialize the Google Cloud Storage client using Gin's context.
	//r.Use(func(c *gin.Context) {
	//	client, err := storage.NewClient(c.Request.Context())
	//	if err != nil {
	//		log.Fatalf("Failed to create storage client: %v", err)
	//	}
	//
	//	// Store the client in Gin's context for use in handlers.
	//	c.Set("storageClient", client)
	//
	//	utils.Uploader = &utils.ClientUploader{
	//		Cl:         client,
	//		BucketName: bucketName,
	//		//ProjectID:  projectID,
	//		UploadPath: "test-files/",
	//	}
	//	// Make sure to close the client after the request is processed.
	//	defer client.Close()
	//
	//	c.Next()
	//})

	// middlewares
	r.Use(cors_middleware.CorsMiddleware())

	r.Use(middleware.RecoveryMiddleware())
	r.Use(lang_middleware.LocalizationMiddleware(localization.InitBundle()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ws", func(c *gin.Context) {
		socket_config.ServeWs(c, hub, c.Writer, c.Request)
	})

	// payload validations
	payloadValidations()

	// Registering routes
	route.AuthRoutes(r, validate)
	tempAttachmentRoute.TempAttachmentsRoutes(r, validate)
	baseUserRoute.UserRoutes(r, validate)
	roleRoute.RoleRoutes(r, validate)

	// Run the seed function to ensure default positions and permissions are set
	//seed.SeedDefaultPositions(database.DB)
	log.Println("_____________")
	// Serve static files from the images directory
	r.Static("/images", "./images")

	r.Run()
}

func setupPaseto() {
	tokenMaker, err := paseto_token.NewPaseto("abcdefghijkl12345678901234567890")
	if err != nil {
		panic("Couldnt open tokenmaker " + err.Error())
	}

	paseto_token.TokenMaker = tokenMaker
}

func payloadValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// for user Type validation
		v.RegisterValidation("validUserType", user.ValidUserType)
		// for when both id and fileId is nil where id is uuid
		v.RegisterValidation("FileValidationIfIdNil", global_validation.RequiredIfIdNil)
		// for when both id and fileId is nil
		v.RegisterValidation("FieldValidationIfIdIsNil", global_validation.RequiredIfIdNilNotUUID)

	}
}
