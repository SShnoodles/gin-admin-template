package middleware

import (
	"context"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
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
		timesStr, _ := config.RDB.Get(config.CTX, ip).Result()
		if timesStr == "" {
			config.RDB.Set(context.Background(), ip, 1, time.Second)
			return
		}
		times, _ := strconv.Atoi(timesStr)
		if times+1 > timesMax {
			c.JSON(http.StatusBadRequest, domain.NewMessageWrapper("限制操作"))
			c.Abort()
			return
		}
		config.RDB.Set(context.Background(), ip, times+1, time.Second)
	}
}
