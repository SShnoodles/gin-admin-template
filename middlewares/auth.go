package middlewares

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, models.NewMessageWrapper("无权操作"))
			c.Abort()
			return
		}

		claims, err := utils.ValidateUserToken(token, initializations.AppConfig.Jwt.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.NewMessageWrapper("无权操作"))
			c.Abort()
			return
		}
		id, _ := claims.GetSubject()
		c.Set("UserId", id)
	}
}
