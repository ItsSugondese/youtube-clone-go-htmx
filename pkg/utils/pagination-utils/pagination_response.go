package pagination_utils

type PaginationResponse struct {
	TotalPages       int         `json:"totalPages"`
	TotalElements    int64       `json:"totalElements"`
	NoOfElements     int         `json:"noOfElements"`
	CurrentPageIndex int         `json:"currentPageIndex"`
	Data             interface{} `json:"data"`
}
