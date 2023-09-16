package controller

import (
	"fmt"
	"log"
	"mangosteen/api"
	"mangosteen/internal/database"
	"mangosteen/sql/queries"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nav-inc/datetime"
)

type ItemController struct {
	PerPage int32
}

func (ctrl *ItemController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/items", ctrl.Create)
	v1.GET("/items", ctrl.GetPaged)
	v1.GET("/items/balance", ctrl.GetBalance)
	v1.GET("/items/summary", ctrl.GetSummary)
	ctrl.PerPage = 10
}
func (ctrl *ItemController) Create(c *gin.Context) {
	var body struct {
		Amount     int32     `json:"amount" binding:"required"`
		Kind       string    `json:"kind" binding:"required"`
		HappenedAt time.Time `json:"happened_at" binding:"required"`
		TagIds     []int32   `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	item, err := q.CreateItem(c, queries.CreateItemParams{
		UserID:     user.ID,
		Amount:     body.Amount,
		Kind:       body.Kind,
		HappenedAt: body.HappenedAt,
		TagIds:     body.TagIds,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resource": item,
	})
}
func (ctrl *ItemController) Destroy(c *gin.Context) {}
func (ctrl *ItemController) Update(c *gin.Context)  {}
func (ctrl *ItemController) Get(c *gin.Context)     {}
func (ctrl *ItemController) GetPaged(c *gin.Context) {
	var params api.GetPagedItemsRequest
	pageString, _ := c.Params.Get("page")
	if page, err := strconv.Atoi(pageString); err == nil {
		params.Page = int32(page)
	}
	if params.Page == 0 {
		params.Page = 1
	}

	happenedBefore, has := c.Params.Get("happened_before")
	if has {
		if t, err := datetime.Parse(happenedBefore, time.Local); err == nil {
			params.HappenedBefore = t
		}
	}

	happenedAfter, has := c.Params.Get("happened_after")
	if has {
		if t, err := datetime.Parse(happenedAfter, time.Local); err == nil {
			params.HappenedAfter = t
		}
	}

	q := database.NewQuery()
	items, err := q.ListItem(c, queries.ListItemParams{
		Offset: (params.Page - 1) * ctrl.PerPage,
		Limit:  ctrl.PerPage,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	count, err := q.CountItems(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, api.GetPagedItemsResponse{
		Resources: items,
		Pager: api.Pager{
			Page:    params.Page,
			PerPage: ctrl.PerPage,
			Count:   int32(count),
		},
	})
}

func (ctrl *ItemController) GetBalance(c *gin.Context) {

	// query := c.Request.URL.Query()
	// happenedBeforeString := query["happened_before"][0]
	// happenedAfterString := query["happened_after"][0]
	happenedBeforeString := c.Query("happened_before")
	happenedAfterString := c.Query("happened_after")

	happenedBefore, err := datetime.Parse(happenedBeforeString, time.Local)
	if err != nil {
		happenedBefore = time.Now().AddDate(1, 0, 0)
	}
	happenedAfter, err := datetime.Parse(happenedAfterString, time.Local)
	if err != nil {
		happenedAfter = time.Now().AddDate(-100, 0, 0)
	}

	q := database.NewQuery()
	items, err := q.ListItemsHappenedBetween(c, queries.ListItemsHappenedBetweenParams{
		HappenedBefore: happenedBefore,
		HappenedAfter:  happenedAfter,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		log.Println(err)
		return
	}

	var r api.GetBalanceResponse
	for _, item := range items {
		if item.Kind == "in_come" {
			r.Income += int32(item.Amount)
		} else {
			r.Expenses += int32(item.Amount)
		}
	}
	r.Balance = r.Income - r.Expenses
	c.JSON(http.StatusOK, r)

}

func (ctrl *ItemController) GetSummary(c *gin.Context) {
	var query api.GetSummaryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		r := api.ErrorResponse{Errors: map[string][]string{}}
		switch x := err.(type) {
		case (validator.ValidationErrors):
			for _, ve := range x {
				t := ve.Tag()
				f := ve.Field()
				if r.Errors[f] == nil {
					r.Errors[f] = []string{}
				}
				r.Errors[f] = append(r.Errors[f], t)
			}
			fmt.Println(r)
			c.JSON(http.StatusUnprocessableEntity, r)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	items, err := q.ListItemsByHappenedAtAndKind(c, queries.ListItemsByHappenedAtAndKindParams{
		UserID:         user.ID,
		Kind:           query.Kind,
		HappenedAfter:  query.HappenedAfter,
		HappenedBefore: query.HappenedBefore,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}

	if query.GroupBy == "happened_at" {
		res := api.GetSummaryHappenedAtResponse{}
		res.Total = 0
		res.Groups = []api.HappenedAtWithGroup{}

		for _, item := range items {
			k := item.HappenedAt.Format("2006-01-02")
			res.Total += item.Amount
			found := false
			for index, group := range res.Groups {
				if group.HappenedAt == k {
					found = true
					res.Groups[index].Amount += item.Amount
				}
			}
			if !found {
				res.Groups = append(res.Groups, api.HappenedAtWithGroup{
					Amount:     item.Amount,
					HappenedAt: k,
				})
			}
		}

		sort.Slice(res.Groups, func(i, j int) bool {
			return res.Groups[i].HappenedAt < res.Groups[j].HappenedAt
		})
		c.JSON(200, res)
	} else if query.GroupBy == "tag_id" {
		res := api.GetSummaryByTagIDResponse{}
		res.Total = 0
		res.Groups = []api.TagIDWithGroup{}

		for _, item := range items {
			if len(item.TagIds) == 0 {
				continue
			}
			k := item.TagIds[0]
			res.Total += item.Amount
			found := false
			for index, group := range res.Groups {
				if group.TagID == k {
					found = true
					res.Groups[index].Amount += item.Amount
				}
			}
			if !found {
				res.Groups = append(res.Groups, api.TagIDWithGroup{
					Amount: item.Amount,
					TagID:  k,
				})
			}
		}

		sort.Slice(res.Groups, func(i, j int) bool {
			return res.Groups[i].TagID < res.Groups[j].TagID
		})
		c.JSON(200, res)
	} else {
		c.JSON(422, gin.H{"message": "参数错误"})
	}
}
