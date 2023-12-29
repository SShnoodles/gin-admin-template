package middleware

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, domain.NewMessageWrapper("无权操作"))
			c.Abort()
			return
		}

		claims, err := util.ValidateUserToken(token, config.AppConfig.Jwt.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.NewMessageWrapper("无权操作"))
			c.Abort()
			return
		}
		id, _ := claims.GetSubject()
		c.Set("UserId", id)
	}
}
