package models

type Paginator struct {
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
}

type Response struct {
	Data       any `json:"data"`
	Pagination Paginator
}
