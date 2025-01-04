package service

import (
	"math"
	"mime/multipart"
	"youtube-clone/enums/struct-enums/project_module"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/global/global_var"
	"youtube-clone/internal/temporary-attachments/model"
	"youtube-clone/pkg/common/database"
	"youtube-clone/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SaveTemporaryAttachmentsService(ctx *gin.Context) []model.TemporaryAttachments {

	tx := database.DB.Begin()
	tx.WithContext(ctx)
	var attachmentsList []model.TemporaryAttachments
	form, err := ctx.MultipartForm()
	if err != nil {
		panic("Error parsing form: " + err.Error())
	}

	files := form.File["attachments"] // "attachments" is the form field name for the files

	var attachments []*multipart.FileHeader

	// Loop over the files
	for _, fileHeader := range files {
		attachments = append(attachments, fileHeader)
		fileDetails := utils.SaveFile(fileHeader, project_module.ModuleNameEnums.TEMPORARY_ATTACHMENTS, global_var.ForBucket)

		attach, err := generic_repo.SaveRepo(tx, model.TemporaryAttachments{
			Name:     fileHeader.Filename,
			Location: fileDetails.FilePath,
			FileSize: math.Round(float64((fileDetails.Size/1000)*100)) / 100,
			FileType: fileDetails.FileType,
		})

		if err != nil {
			panic("Error when saving file")
		}

		attachmentsList = append(attachmentsList, attach)
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

	return attachmentsList
}
