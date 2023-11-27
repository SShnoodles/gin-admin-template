package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MenuQuery struct {
	services.PageInfo
	Name string `form:"name"`
}

func GetMenus(c *gin.Context) {
	var q MenuQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var menus []models.Menu
	page := services.Pagination(initializations.DB, q.PageIndex, q.PageSize, &menus)
	c.JSON(http.StatusOK, page)
}

func GetMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var menu models.Menu
	err = services.FindById(&menu, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
	var menu models.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	menu.Id = initializations.IdGenerate()
	err = services.Insert(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(menu.Id))
}

func UpdateMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var menu models.Menu
	err = c.ShouldBindJSON(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	menu.Id = int64(id)
	err = services.Update(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = services.DeleteById(models.Menu{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
