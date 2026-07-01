package api

import (
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

// GetVersion
// @Summary Version 版本
// @Tags project 项目
// @Accept json
// @Produce json
// @Router /project/version [get]
func GetVersion(c *gin.Context) {
	service.Ok(c, Version)
}
