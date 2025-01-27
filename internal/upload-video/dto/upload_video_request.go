package dto

import (
	"github.com/google/uuid"
)

type UploadVideoRequest struct {
	ID         uuid.UUID `json:"id"`
//	FileId     uuid.UUID `json:"fileId" binding:"FileValidationIfIdNil=ID"`
}
