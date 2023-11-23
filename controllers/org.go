package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrgQuery struct {
	models.PageInfo
	Name string `form:"name"`
}

func GetOrgs(c *gin.Context) {
	var q OrgQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var orgs []models.Org
	page := models.Pagination(initializations.DB, q.PageIndex, q.PageSize, &orgs)
	c.JSON(http.StatusOK, page)
}

func GetOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var org models.Org
	err = services.FindById(&org, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, org)
}

func CreateOrg(c *gin.Context) {
	var org models.Org
	err := c.ShouldBindJSON(&org)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	org.Id = initializations.IdGenerate()
	err = services.Insert(org)
	if err != nil {
		c.String(http.StatusBadRequest, "新增失败")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(org.Id))
}

func UpdateOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var org models.Org
	err = c.ShouldBindJSON(&org)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	org.Id = int64(id)
	err = services.Update(org)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = services.DeleteById(models.Org{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
