package dto

import "mime/multipart"

type TemporaryAttachmentsDetailRequest struct {
	Attachments []*multipart.FileHeader
}
