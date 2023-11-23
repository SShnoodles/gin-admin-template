package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ResourceQuery struct {
	models.PageInfo
	Name string `form:"name"`
}

func GetResources(c *gin.Context) {
	var q ResourceQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var resource []models.Resource
	page := models.Pagination(initializations.DB, q.PageIndex, q.PageSize, &resource)
	c.JSON(http.StatusOK, page)
}

func GetResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var resource models.Resource
	err = services.FindById(&resource, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, resource)
}

func CreateResource(c *gin.Context) {
	var resource models.Resource
	err := c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	resource.Id = initializations.IdGenerate()
	err = services.Insert(resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(resource.Id))
}

func UpdateResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var resource models.Resource
	err = c.ShouldBindJSON(&resource)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	resource.Id = int64(id)
	err = services.Update(resource)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteResource(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = services.DeleteById(models.Resource{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
