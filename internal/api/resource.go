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
		config.Log.Error(err.Error())
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Resource{})
	c.JSON(http.StatusOK, page)
}

func GetResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var resource domain.Resource
	err = service.FindById(&resource, id)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, resource)
}

func CreateResource(c *gin.Context) {
	var resource domain.Resource
	err := c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	resource.Id = config.IdGenerate()
	err = service.Insert(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(resource.Id))
}

func UpdateResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var resource domain.Resource
	err = c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	resource.Id = id
	err = service.Update(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	err = service.DeleteById(domain.Resource{}, id)
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
