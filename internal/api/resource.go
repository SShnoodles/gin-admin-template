package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResourceQuery struct {
	service.PageInfo
	Name string `form:"name"`
}

func GetResources(c *gin.Context) {
	var q ResourceQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var resource []domain.Resource
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, &resource)
	c.JSON(http.StatusOK, page)
}

func GetResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var resource domain.Resource
	err = service.FindById(&resource, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, resource)
}

func CreateResource(c *gin.Context) {
	var resource domain.Resource
	err := c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	resource.Id = config.IdGenerate()
	err = service.Insert(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(resource.Id))
}

func UpdateResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var resource domain.Resource
	err = c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	resource.Id = int64(id)
	err = service.Update(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = service.DeleteById(domain.Resource{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
