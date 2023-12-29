package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MenuQuery struct {
	service.PageInfo
	Name string `form:"name"`
}

func GetMenus(c *gin.Context) {
	var q MenuQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var menus []domain.Menu
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, &menus)
	c.JSON(http.StatusOK, page)
}

func GetMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var menu domain.Menu
	err = service.FindById(&menu, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
	var menu domain.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	menu.Id = config.IdGenerate()
	err = service.Insert(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(menu.Id))
}

func UpdateMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var menu domain.Menu
	err = c.ShouldBindJSON(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	menu.Id = int64(id)
	err = service.Update(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = service.DeleteById(domain.Menu{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
