package model

import (
	generic_models "youtube-clone/generics/generic-models"
)

type UploadVideo struct {
	generic_models.AuditModel
}

func (b *UploadVideo) HasAuditModel() bool {
	return true
}
