package api

import (
	"mangosteen/sql/queries"
	"time"
)

type GetSummaryRequest struct {
	HappenedBefore time.Time `form:"happened_before" binding:"required"`
	HappenedAfter  time.Time `form:"happened_after" binding:"required"`
	Kind           string    `form:"kind" binding:"required,oneof=expenses in_come"`
	GroupBy        string    `form:"group_by" binding:"required,oneof=tag_id happened_at"`
}

type HappenedAtWithGroup struct {
	Amount     int32  `json:"amount"`
	HappenedAt string `json:"happened_at"`
}

type GetSummaryHappenedAtResponse struct {
	Groups []HappenedAtWithGroup `json:"groups"`
	Total  int32                 `json:"total"`
}

type TagIDWithGroup struct {
	TagID  int32       `json:"tag_id"`
	Amount int32       `json:"amount"`
	Tag    queries.Tag `json:"tag"`
}

type GetSummaryByTagIDResponse struct {
	Groups []TagIDWithGroup `json:"groups"`
	Total  int32            `json:"total"`
}

type CreateItemRequest struct {
	Amount     int32     `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIds     []int32   `json:"tag_ids" binding:"required"`
}

type CreateItemResponse struct {
	Resource queries.Item
}

type GetPagedItemsRequest struct {
	Page           int32     `json:"page"`
	HappenedAfter  time.Time `json:"happened_after"`
	HappenedBefore time.Time `json:"happened_before"`
}

type GetPagedItemsResponse struct {
	Resources []queries.Item `json:"resources"`
	Pager     Pager          `json:"pager"`
}

type GetBalanceResponse struct {
	Income   int32 `json:"income"`
	Expenses int32 `json:"expenses"`
	Balance  int32 `json:"balance"`
}
