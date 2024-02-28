package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const Version = "1.0.0"

// GetVersion
// @Summary Version 版本
// @Tags project 项目
// @Accept json
// @Produce json
// @Router /project/version [get]
func GetVersion(c *gin.Context) {
	c.String(http.StatusOK, Version)
}
