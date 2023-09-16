package controller

import (
	"mangosteen/api"
	"mangosteen/internal/database"
	"mangosteen/sql/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	PerPage int32
}

func (ctrl *TagController) RegisterRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1.POST("/tags", ctrl.Create)
	v1.GET("/tags", ctrl.GetPaged)
	v1.GET("/tags/:id", ctrl.Get)
	v1.PATCH("/tags/:id", ctrl.Update)
	v1.DELETE("/tags/:id", ctrl.Destroy)
	ctrl.PerPage = 10
}

func (ctrl *TagController) Create(c *gin.Context) {
	var body api.CreateTagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}
	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.CreateTag(c, queries.CreateTagParams{
		UserID: user.ID,
		Kind:   body.Kind,
		Name:   body.Name,
		Sign:   body.Sign,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, api.CreateTagResponse{
		Resource: tag,
	})
}

func (ctrl *TagController) Destroy(c *gin.Context) {
	idString, has := c.Params.Get("id")
	if !has {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}
	q := database.NewQuery()
	err = q.DeleteTag(c, int32(id))
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (ctrl *TagController) Update(c *gin.Context) {
	var body api.UpdateTagRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}
	idString, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}

	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.UpdateTag(c, queries.UpdateTagParams{
		ID:     int32(id),
		UserID: user.ID,
		Kind:   body.Kind,
		Name:   body.Name,
		Sign:   body.Sign,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, api.UpdateTagResponse{
		Resource: tag,
	})
}

func (ctrl *TagController) Get(c *gin.Context) {
	idString, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(422, gin.H{"message ": "请求参数有误"})
		return
	}

	me, _ := c.Get("me")
	user, _ := me.(queries.User)
	q := database.NewQuery()
	tag, err := q.GetTag(c, queries.GetTagParams{
		ID:     int32(id),
		UserID: user.ID,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, api.GetTagResponse{
		Resource: tag,
	})
}

func (ctrl *TagController) GetPaged(c *gin.Context) {
	me, _ := c.Get("me")
	var params api.GetPageTagsRequest
	pageString, _ := c.Params.Get("page")
	if page, err := strconv.Atoi(pageString); err == nil {
		params.Page = int32(page)
	}
	if params.Page == 0 {
		params.Page = 1
	}

	kind, _ := c.Params.Get("kind")
	if kind == "" {
		kind = "expenses"
	}

	q := database.NewQuery()
	tags, err := q.ListTags(c, queries.ListTagsParams{
		Offset: (params.Page - 1) * ctrl.PerPage,
		Limit:  ctrl.PerPage,
		Kind:   kind,
		UserID: me.(queries.User).ID,
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
	c.JSON(http.StatusOK, api.GetPagedTagsResponse{
		Resources: tags,
		Pager: api.Pager{
			Page:    params.Page,
			PerPage: ctrl.PerPage,
			Count:   int32(count),
		},
	})
}
