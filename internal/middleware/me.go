package middleware

import (
	"fmt"
	"mangosteen/internal/database"
	"mangosteen/internal/jwt_helper"
	"mangosteen/sql/queries"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Me(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if indexOf(whitelist, path) != -1 {
			c.Next()
			return
		}
		user, err := getMe(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("me", user)
		c.Next()
	}
}

func getMe(c *gin.Context) (queries.User, error) {
	var user queries.User
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 {
		return user, fmt.Errorf("无效的 JWT")
	}
	jwtString := auth[7:]
	if len(jwtString) == 0 {
		return user, fmt.Errorf("无效的 JWT")
	}
	token, err := jwt_helper.Parse(jwtString)
	if err != nil {
		return user, fmt.Errorf("无效的 JWT")
	}
	m, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, fmt.Errorf("无效的 JWT")
	}
	userId := int32(m["user_id"].(float64))

	q := database.NewQuery()
	user, err = q.FindUserById(c, int32(userId))
	if err != nil {
		return user, fmt.Errorf("无效的 JWT")
	}
	return user, nil
}

func indexOf(whitelist []string, path string) int {
	for i, v := range whitelist {
		if v == path {
			return i
		}
	}
	return -1
}
