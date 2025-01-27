package dto

import (
	"github.com/google/uuid"
)

type UploadVideoResponse struct {
	ID             uuid.UUID `json:"id"`
}
