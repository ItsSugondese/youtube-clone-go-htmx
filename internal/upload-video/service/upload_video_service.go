package service

import (
	generic_repo "youtube-clone/generics/generic-repo"
	upload_video_navigator "youtube-clone/internal/upload-video/upload-video-navigator"
	"youtube-clone/internal/upload-video/dto"
	"youtube-clone/internal/upload-video/model"
	"youtube-clone/internal/upload-video/repo"
	"youtube-clone/pkg/common/database"
	pagination_utils "youtube-clone/pkg/utils/pagination-utils"
	"encoding/json"
	dto_utils "youtube-clone/pkg/utils/dto-utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadVideoService interface {
	SaveUploadVideoService(ctx *gin.Context, uploadVideodto dto.UploadVideoRequest) (dto.UploadVideoResponse)
    GetAllUploadVideo(ctx *gin.Context, ) []model.UploadVideo
    FindAllUploadVideosPaginatedService(ctx *gin.Context, request dto.UploadVideoPaginationRequest) *pagination_utils.PaginationResponse
    GetUploadVideoDetailsByIdService(ctx *gin.Context, id uuid.UUID) *dto.UploadVideoResponse
    DeleteUploadVideoByIdService(ctx *gin.Context, uploadVideoId uuid.UUID)
}

type uploadVideoService struct {
	repo repo.UploadVideoRepo
}

func NewUploadVideoService(repo repo.UploadVideoRepo) UploadVideoService {
	return &uploadVideoService{
		repo: repo,
	}
}

func (s *uploadVideoService) SaveUploadVideoService(ctx *gin.Context, uploadVideodto dto.UploadVideoRequest) dto.UploadVideoResponse {
	tx := database.DB.Begin()
	tx = tx.WithContext(ctx)


	var uploadVideoDetails model.UploadVideo

	if uploadVideodto.ID != uuid.Nil {
		uploadVideoDetails = upload_video_navigator.FindUploadVideoByIdService(uploadVideodto.ID)
	}

	dto_utils.DtoConvertErrorHandled(uploadVideodto, &uploadVideoDetails)

	var savedUploadVideo model.UploadVideo
	var saveUploadVideoError error

	if uploadVideodto.ID == uuid.Nil {
		savedUploadVideo, saveUploadVideoError = generic_repo.SaveRepo(tx, uploadVideoDetails)
	} else {
		savedUploadVideo, saveUploadVideoError = generic_repo.UpdateRepo(tx, uploadVideoDetails)
	}

	if saveUploadVideoError != nil {
		tx.Rollback()
		panic(saveUploadVideoError)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	var response dto.UploadVideoResponse
    dto_utils.DtoConvertErrorHandled(savedUploadVideo, &response)

    return response // Successfully return the saved Upload video
}

func (s *uploadVideoService) GetAllUploadVideo(ctx *gin.Context, ) []model.UploadVideo {
	response, err := generic_repo.FindAll[model.UploadVideo]()

	if err != nil {
		panic(err)
	}
	return response
}

func (s *uploadVideoService) FindAllUploadVideosPaginatedService(ctx *gin.Context, request dto.UploadVideoPaginationRequest) *pagination_utils.PaginationResponse {
	return s.repo.FindAllUploadVideoPaginatedRepo(request, pagination_utils.PaginationResponse{})
}

func (s *uploadVideoService) GetUploadVideoDetailsByIdService(ctx *gin.Context, id uuid.UUID) *dto.UploadVideoResponse {
	var uploadVideoResponse dto.UploadVideoResponse
    details := upload_video_navigator.FindUploadVideoByIdService(id)

	jsonData, _ := json.Marshal(details)
	jsonUnmarshalError := json.Unmarshal(jsonData, &uploadVideoResponse)
	if jsonUnmarshalError != nil {
		panic(jsonUnmarshalError)
	}

	return &uploadVideoResponse
}

func (s *uploadVideoService) DeleteUploadVideoByIdService(ctx *gin.Context, uploadVideoId uuid.UUID) {
	tx := database.DB.Begin()
	tx = tx.WithContext(ctx)

	uploadVideoDetails := upload_video_navigator.FindUploadVideoByIdService(uploadVideoId)


	// Delete the user location
	err := generic_repo.DeleteByStructRepo(tx, uploadVideoDetails)
	if err != nil {
		panic(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}
