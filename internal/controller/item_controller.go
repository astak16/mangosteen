package controller

import (
	"log"
	"mangosteen/api"
	"mangosteen/internal/database"
	"mangosteen/sql/queries"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	ctrl.PerPage = 10
}
func (ctrl *ItemController) Create(c *gin.Context) {
	var body struct {
		Amount     int32        `json:"amount" binding:"required"`
		Kind       queries.Kind `json:"kind" binding:"required"`
		HappenedAt time.Time    `json:"happened_at" binding:"required"`
		TagIds     []int32      `json:"tag_ids" binding:"required"`
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
		if item.Kind == queries.KindInCome {
			r.Income += int32(item.Amount)
		} else {
			r.Expenses += int32(item.Amount)
		}
	}
	r.Balance = r.Income - r.Expenses
	c.JSON(http.StatusOK, r)

}
