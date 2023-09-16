package api

type Pager struct {
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	Count   int32 `json:"count"`
}

type ErrorResponse struct {
	Errors map[string][]string `json:"errors"`
}
