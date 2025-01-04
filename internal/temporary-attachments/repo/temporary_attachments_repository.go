package repo

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"youtube-clone/internal/temporary-attachments/model"
	"youtube-clone/pkg/common/database"
)

func SaveTemporaryAttachmentsRepo(attachment model.TemporaryAttachments) (model.TemporaryAttachments, error) {
	result := database.DB.Create(&attachment)
	return attachment, result.Error
}

func FindTempAttachmentsByIdRepo(id uuid.UUID) (attachment model.TemporaryAttachments, err error) {
	if err := database.DB.
		Where("id= ?", id).
		First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, return zero value of Staff and nil error
			return model.TemporaryAttachments{}, nil
		}
		// Other errors occurred, return the error
		return model.TemporaryAttachments{}, err
	}

	// Staff found, return the Staff and nil error
	return attachment, nil
}
