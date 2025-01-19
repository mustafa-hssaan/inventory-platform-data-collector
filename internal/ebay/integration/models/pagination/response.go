package pagination

type PaginationResponse struct {
	Total int `json:"total"`
	Size  int `json:"size"`
	Limit int `json:"limit"`
}
