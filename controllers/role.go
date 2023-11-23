package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RoleQuery struct {
	models.PageInfo
	Name string `form:"name"`
}

func GetRoles(c *gin.Context) {
	var q RoleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var roles []models.Role
	page := models.Pagination(initializations.DB, q.PageIndex, q.PageSize, &roles)
	c.JSON(http.StatusOK, page)
}

func GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var role models.Role
	err = services.FindById(&role, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, role)
}

func CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	role.Id = initializations.IdGenerate()
	err = services.Insert(role)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(role.Id))
}

func UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var role models.Role
	err = c.ShouldBindJSON(&role)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	role.Id = int64(id)
	err = services.Update(role)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = services.DeleteById(models.Role{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
