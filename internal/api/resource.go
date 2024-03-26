package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResourceQuery struct {
	service.PageInfo
	Name string `form:"name"`
	Path string `form:"path"`
}

// GetResources
// @Summary List resources 获取资源列表
// @Tags resources 资源
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "name 名称"
// @Param path query string false "path 路径"
// @Router /resources [get]
func GetResources(c *gin.Context) {
	var q ResourceQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	if q.PageSize == 0 {
		var resources []domain.Resource
		err = service.FindAll(&resources)
		if err != nil {
			service.BadRequestResult(c, "Failed.query")
			config.Log.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, resources)
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Resource{})
	c.JSON(http.StatusOK, page)
}
