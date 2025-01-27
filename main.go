package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"youtube-clone/config"
	oauth2_setup "youtube-clone/config/oauth2-setup"
	socket_config "youtube-clone/config/socket-config"
	"youtube-clone/internal/auth/route"
	roleModel "youtube-clone/internal/role/model"
	roleRoute "youtube-clone/internal/role/route"
	"youtube-clone/internal/temporary-attachments/model"
	tempAttachmentRoute "youtube-clone/internal/temporary-attachments/route"
	uploadVideoRoute "youtube-clone/internal/upload-video/route"
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

const (
	defaultChunkSize = 1024 * 1024 // 1MB chunks
	maxRetries       = 3
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

	//r.SetHTMLTemplate(template_config.LoadTemplates("./web/templates"))
	r.LoadHTMLGlob("./web/templates/*.html")
	// Serve static files like CSS, JS, images
	r.Static("/web/static", "./web/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "HTMX with Gin"})
	})
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

	// Example of a simple file upload handler in Go using Gin
	r.POST("/upload", func(c *gin.Context) {
		// Create the uploads directory if it doesn't exist
		outputDir := "/home/lazybot/Desktop/uploads"
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create upload directory"})
			return
		}

		// Set the destination file path
		destFilePath := filepath.Join(outputDir, "uploaded_chunk.mp4")

		// Open the file in append mode, create it if it doesn't exist
		destFile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to open file"})
			return
		}
		defer destFile.Close()

		// Append the raw request body to the file
		_, err = io.Copy(destFile, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save file"})
			return
		}

		// Respond with success
		c.JSON(200, gin.H{"message": "File uploaded and appended successfully"})
	})

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

	chunkSize := defaultChunkSize
	if size, ok := os.LookupEnv("CHUNK_SIZE"); ok {
		fmt.Sscanf(size, "%d", &chunkSize)
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <file_path>")
	}
	filePath := os.Args[1]

	config := Config{ChunkSize: chunkSize, ServerURL: "http://localhost:3000/upload"}

	chunker := &DefaultFileChunker{chunkSize: config.ChunkSize}
	uploader := &DefaultUploader{serverURL: config.ServerURL}

	baseName := filepath.Base(filePath)
	// Remove the extension
	fileNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	generatedFileName := FileNameGenerator(fileNameWithoutExt)
	//chunks, err := chunker.ChunkFile(filePath)
	chunks, err := chunker.ChunkVideo(filePath, "/home/lazybot/Desktop/sample", 5, generatedFileName)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	err = synchronizeChunks(chunks, uploader, &wg, generatedFileName)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	//changeChan := make(chan bool)
	//go watchFile(filePath, changeChan)
	//
	//for {
	//	select {
	//	case <-changeChan:
	//		log.Println("File changed, re-chunking and synchronizing...")
	//		chunks, err = chunker.ChunkFile(filePath)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		err = synchronizeChunks(chunks, metadata, uploader, &wg, &mu)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		wg.Wait()
	//
	//		err = metadataManager.SaveMetadata(fmt.Sprintf("%s.metadata.json", filePath), metadata)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	case <-time.After(10 * time.Second):
	//		log.Println("No changes detected, checking again...")
	//	}
	//}

	r.POST("/upload", func(c *gin.Context) {
		// Parse metadata from the request
		fileNameFromServer := c.Query("file_name")
		chunkNameFromServer := c.Query("chunk_name")
		chunkIndex, _ := strconv.Atoi(c.Query("chunk_index"))
		totalChunks, _ := strconv.Atoi(c.Query("total_chunks"))

		if chunkNameFromServer == "" || chunkIndex < 0 || totalChunks <= 0 {
			c.JSON(400, gin.H{"error": "Invalid chunk metadata"})
			return
		}

		// Set the destination file path (e.g., /uploads/chunkNameFromServer.mp4)
		outputDir := "/home/lazybot/Desktop/uploads"
		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create upload directory"})
			return
		}

		chunkName := filepath.Base(chunkNameFromServer)
		destFilePath := filepath.Join(outputDir, chunkName)

		// Open the file in append mode, create it if it doesn't exist
		destFile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to open destination file"})
			return
		}
		defer destFile.Close()

		// Append the request body to the destination file
		_, err = io.Copy(destFile, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to append chunk"})
			return
		}

		// write into text file from here

		destTextFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.txt", fileNameFromServer))

		var file *os.File

		fileNamesSlices := ReadFromFileAndReturnFileEntrySlice(destTextFilePath)
		if len(fileNamesSlices) != totalChunks-1 {
			file, err = os.OpenFile(destTextFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			// Append or write data
			text := fmt.Sprintf("%d file '%s'\n", chunkIndex, chunkName)
			_, err = file.WriteString(text)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			file, err = os.OpenFile(destTextFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		}
		if err != nil {
			fmt.Println("Error opening or creating file:", err)
			return
		}
		defer file.Close()

		fileNamesSlices = append(fileNamesSlices, FileEntry{
			Index: chunkIndex,
			File:  chunkName,
		})

		fmt.Println("Data written successfully.")
		// Check if the upload is complete
		if len(fileNamesSlices) == totalChunks {
			//tempFile := CreateTempFile(fileNamesSlices)

			sort.Slice(fileNamesSlices, func(i, j int) bool {
				return fileNamesSlices[i].Index < fileNamesSlices[j].Index
			})
			// Write slice data to the temporary file
			var text string
			for _, val := range fileNamesSlices {
				text += fmt.Sprintf("file '%s'\n", val.File)
			}
			_, err = file.WriteString(text)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

			outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.mp4", fileNameFromServer))
			cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", destTextFilePath, "-c", "copy", outputFile)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Errorf("ffmpeg command failed: %w", err)
			}

			log.Println("Video merge successful:", outputFile)
			err := os.Remove(destTextFilePath)
			if err != nil {
				return
			}

			for _, val := range fileNamesSlices {
				err := os.Remove(filepath.Join(outputDir, val.File))
				if err != nil {
					return
				}
			}
			c.JSON(200, gin.H{"message": "All chunks uploaded successfully"})
		} else {
			c.JSON(200, gin.H{"message": "Chunk uploaded successfully"})

		}
	})
	// Registering routes
	route.AuthRoutes(r, validate)
	tempAttachmentRoute.TempAttachmentsRoutes(r, validate)
	baseUserRoute.UserRoutes(r, validate)
	roleRoute.RoleRoutes(r, validate)
	uploadVideoRoute.UploadVideoRoutes(r, validate, database.DB)

	// Run the seed function to ensure default positions and permissions are set
	//seed.SeedDefaultPositions(database.DB)
	log.Println("_____________")
	// Serve static files from the images directory
	r.Static("/images", "./images")

	//r.Run()
}

func FileNameGenerator(name string) string {
	rand.Seed(time.Now().UnixNano())

	// Generate a random 4-digit number
	randomDigits := rand.Intn(9000) + 1000 // Ensures it's a 4-digit number

	// Get the current date and time in the desired format
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102150405.000") // Format: yyyyMMddHHmmss.mmm

	// Combine to form the string
	filename := fmt.Sprintf("%s%d%s", name, randomDigits, formattedTime)

	return filename
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
