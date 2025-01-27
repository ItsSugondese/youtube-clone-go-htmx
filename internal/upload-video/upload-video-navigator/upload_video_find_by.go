package upload_video_navigator

import (
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	"youtube-clone/internal/upload-video/model"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/pkg/common/localization"
	"github.com/google/uuid"
)

func FindUploadVideoByIdService(id uuid.UUID) model.UploadVideo {
	uploadVideo, err := generic_repo.FindSingleByField[model.UploadVideo]("id", id)

	if err != nil {
		panic(err)
	}
	if uploadVideo == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO,
			"Second": "Id",
		}))
	}
	return *uploadVideo
}
