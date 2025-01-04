package globaldto

import (
	"youtube-clone/enums/interface-enums/response/response-status-enum"
)

type ApiResponse struct {
	Status  response_status_enum.ResponseStatusEnum `json:"status"`
	Message string                                  `json:"message"`
	Data    interface{}                             `json:"data"`
}
