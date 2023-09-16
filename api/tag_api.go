package api

import "mangosteen/sql/queries"

type GetTagRequest struct {
	ID string `json:"id"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Sign string `json:"sign" binding:"required"`
	Kind string `json:"kind" binding:"required"`
}

type GetTagResponse CreateTagResponse

type UpdateTagRequest struct {
	Name string `json:"name"`
	Sign string `json:"sign"`
	Kind string `json:"kind"`
}

type CreateTagResponse struct {
	Resource queries.Tag `json:"resource"`
}

type UpdateTagResponse CreateTagResponse

type GetPageTagsRequest struct {
	Page int32  `json:"page"`
	Kind string `json:"kind"`
}

type GetPagedTagsResponse struct {
	Resources []queries.Tag `json:"resources"`
	Pager     Pager         `json:"pager"`
}
