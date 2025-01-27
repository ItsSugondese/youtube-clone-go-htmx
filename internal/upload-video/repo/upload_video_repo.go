package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	generic_repo "youtube-clone/generics/generic-repo"
	"youtube-clone/internal/upload-video/dto"
	dto_utils "youtube-clone/pkg/utils/dto-utils"
	pagination_utils "youtube-clone/pkg/utils/pagination-utils"
)

type UploadVideoRepo interface {
	FindAllUploadVideoPaginatedRepo(pagination dto.UploadVideoPaginationRequest, response pagination_utils.PaginationResponse) *pagination_utils.PaginationResponse
	FindUploadVideoDetailsByIdRepo(id uuid.UUID) (response *dto.UploadVideoResponse, err error)
}

type uploadVideoRepo struct {
	db *gorm.DB
}

func NewUploadVideoRepo(db *gorm.DB) UploadVideoRepo {
	return &uploadVideoRepo{db: db}
}

func (r *uploadVideoRepo) FindAllUploadVideoPaginatedRepo(pagination dto.UploadVideoPaginationRequest, response pagination_utils.PaginationResponse) *pagination_utils.PaginationResponse {

	query := `
    `

	// Store the results in a map
	var resultMap []map[string]interface{}
	r.db.Scopes(generic_repo.RawQueryPaginate(&pagination.PaginationRequest, &response, r.db, query)).Find(&resultMap)
	tempDtos := dto_utils.ConvertSlice[map[string]interface{}, dto.UploadVideoResponse](resultMap)

	response.Data = tempDtos
	return &response
}

func (r *uploadVideoRepo) FindUploadVideoDetailsByIdRepo(id uuid.UUID) (response *dto.UploadVideoResponse, err error) {
	query := `
    `

	var resultMap map[string]interface{}
	err = r.db.Raw(query, id).Scan(&resultMap).Error

	if err != nil {
		return nil, err
	}
	err = dto_utils.DtoConvertErrorHandledReturnError(resultMap, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}
