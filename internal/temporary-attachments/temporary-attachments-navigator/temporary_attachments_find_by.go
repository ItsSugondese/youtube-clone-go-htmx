package temporary_attachments_navigator

import (
	"github.com/google/uuid"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/internal/temporary-attachments/model"
)

func FindByIdService(id uuid.UUID) model.TemporaryAttachments {
	attachment, err := generic_repo.FindSingleByField[model.TemporaryAttachments]("id", id)
	if err != nil {
		panic("Didn't find attachment with that id")
	}

	return *attachment
}
