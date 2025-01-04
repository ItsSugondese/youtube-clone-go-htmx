package pagination_utils

type PaginationRequest struct {
	Page int `json:"page"`
	Rows int `json:"rows"`
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationRequest) GetLimit() int {
	if p.Rows == 0 {
		p.Rows = 10
	}
	return p.Rows
}

func (p *PaginationRequest) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

//func (p *PaginationRequest) GetSort() string {
//	if p.Sort == "" {
//		p.Sort = "Id desc"
//	}
//	return p.Sort
//}
