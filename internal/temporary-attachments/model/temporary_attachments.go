package model

import (
	"github.com/google/uuid"
	"youtube-clone/constants/file_type_constants"
)

type TemporaryAttachments struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`

	Name     string                       `json:"name"`
	Location string                       `json:"location"`
	FileSize float64                      `json:"file_size"`
	FileType file_type_constants.FileType `json:"file_type"`
}

func (b *TemporaryAttachments) HasAuditModel() bool {
	return false
}
