package controller

import (
	response_crud_enum "youtube-clone/enums/interface-enums/response/response-crud-enum"
	localization_enums "youtube-clone/enums/struct-enums/localization-enums"
	"youtube-clone/enums/struct-enums/project_module"
	generic_controller "youtube-clone/generics/generic-controller"
	globaldto "youtube-clone/global/global_dto"
	"youtube-clone/internal/upload-video/dto"
	"youtube-clone/internal/upload-video/service"
	"youtube-clone/pkg/common/localization"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type UploadVideoController interface {
	SaveUploadVideo(ctx *gin.Context, validate *validator.Validate)
    GetAllUploadVideoDetails(ctx *gin.Context)
    GetAllUploadVideoDetailsPaginated(ctxtx *gin.Context, validate *validator.Validate)
    GetUploadVideoDetailsById(ctx *gin.Context)
    DeleteUploadVideoById(ctx *gin.Context)
}

type uploadVideoController struct {
	service service.UploadVideoService
}

func NewUploadVideoController(service service.UploadVideoService) UploadVideoController {
	return &uploadVideoController{
		service: service,
	}
}

// @Summary register Upload video using this api
// @Schemes
// @Description
// @Tags Upload video
// @Accept json
// @Produce json
// @Param upload_video body dto.UploadVideoRequest true "Upload video details"
// @Success 200 {object} dto.UploadVideoResponse
// @Router /upload-video [post]
func (c *uploadVideoController) SaveUploadVideo(ctx *gin.Context, validate *validator.Validate) {

	var uploadVideoDto dto.UploadVideoRequest

	if err := generic_controller.ControllerValidationHandler(&uploadVideoDto, ctx, validate); err != nil {
		return
	}

	savedData := c.service.SaveUploadVideoService(ctx, uploadVideoDto)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO,
			"Second": strings.ToLower(response_crud_enum.Create().String()),
		}), savedData)
}

// @Summary get all Upload video
// @Schemes
// @Description
// @Tags Upload video
// @Accept json
// @Produce json
// @Success 200 {object} model.UploadVideo
// @Router /faq [get]
func (c *uploadVideoController) GetAllUploadVideoDetails(ctx *gin.Context) {

	getData := c.service.GetAllUploadVideo(ctx)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO + " Details",
			"Second": strings.ToLower(response_crud_enum.Get().String()),
		}), getData)
}

// @Summary get all Upload video details paginated
// @Schemes
// @Description
// @Tags Upload video
// @Accept json
// @Produce json
// @Param upload_video body dto.UploadVideoPaginationRequest true "Upload video details"
// @Success 200 {object} pagination_utils.PaginationResponse{Data=[]dto.UploadVideoResponse}
// @Router /upload-video/paginated [post]
func (c *uploadVideoController) GetAllUploadVideoDetailsPaginated(ctx *gin.Context, validate *validator.Validate) {
    var paginatedRequest dto.UploadVideoPaginationRequest

	if err := generic_controller.ControllerValidationHandler(&paginatedRequest, ctx, validate); err != nil {
		return
	}

	getData := c.service.FindAllUploadVideosPaginatedService(ctx, paginatedRequest)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO + " Details",
			"Second": strings.ToLower(response_crud_enum.Get().String()),
		}), getData)

}

// @Summary get Upload video details by passing id on url
// @Schemes
// @Description
// @Tags Upload video
// @Accept json
// @Produce json
// @Param id path int true "Upload video ID"
// @Success 200 {object} dto.UploadVideoResponse
// @Router /upload-video/:id [get]
func (c *uploadVideoController) GetUploadVideoDetailsById(ctx *gin.Context) {
	id, parseError := uuid.Parse(ctx.Param("id"))

	if parseError != nil {
		panic(parseError)
	}

	getData := c.service.GetUploadVideoDetailsByIdService(ctx, id)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO + " Details",
			"Second": strings.ToLower(response_crud_enum.Get().String()),
		}), getData)

}

// @Summary Delete Upload video
// @Description This API deletes an existing Upload video details by its ID.
// @Tags Upload video
// @Param id path string true "Upload video ID"
// @Success 200 {string} string "Upload video deleted successfully"
// @Failure 404 {object} bool "Upload video not found"
// @Failure 500 {object} bool "Internal server error"
// @Router /upload-video/{id} [delete]
func (c *uploadVideoController) DeleteUploadVideoById(ctx *gin.Context) {
	// Parse the location ID from the URL
	uploadVideoId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		panic(&globaldto.PanicObject{
			Err:        err,
			StatusCode: 400,
		})
	}

	// Call the service to delete the user location
	c.service.DeleteUploadVideoByIdService(ctx, uploadVideoId)

	generic_controller.GenericControllerSuccessResponseHandler(ctx,
		localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.API_OPERATION, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.UPLOAD_VIDEO,
			"Second": strings.ToLower(response_crud_enum.Delete().String()),
		}), nil)
}
