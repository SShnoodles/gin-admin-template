package middlewares

import (
	"context"
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Limit 计数限流
func Limit(timesMax int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		timesStr, _ := initializations.RDB.Get(initializations.CTX, ip).Result()
		if timesStr == "" {
			initializations.RDB.Set(context.Background(), ip, 1, time.Minute)
			return
		}
		times, _ := strconv.Atoi(timesStr)
		if times+1 > timesMax {
			c.JSON(http.StatusBadRequest, models.NewMessageWrapper("限制操作"))
			c.Abort()
			return
		}
		initializations.RDB.Set(context.Background(), ip, times+1, time.Minute)
	}
}
