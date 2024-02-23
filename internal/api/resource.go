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
		service.ParamBadRequestResult(c, err)
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Resource{})
	c.JSON(http.StatusOK, page)
}

func GetResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	var resource domain.Resource
	err = service.FindById(&resource, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query", err)
		return
	}
	c.JSON(http.StatusOK, resource)
}

func CreateResource(c *gin.Context) {
	var resource domain.Resource
	err := c.ShouldBindJSON(&resource)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	resource.Id = config.IdGenerate()
	err = service.Insert(&resource)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(resource.Id))
}

func UpdateResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	var resource domain.Resource
	err = c.ShouldBindJSON(&resource)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	resource.Id = id
	err = service.Update(&resource)
	if err != nil {
		service.BadRequestResult(c, "Failed.update", err)
		return
	}
	c.JSON(http.StatusOK, service.UpdateSuccessResult())
}

func DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	err = service.DeleteById(domain.Resource{}, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.delete", err)
		return
	}
	c.JSON(http.StatusOK, service.DeleteSuccessResult())
}
