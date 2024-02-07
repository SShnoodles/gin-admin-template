package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
)

type MenuQuery struct {
	service.PageInfo
	Name string `form:"name"`
}
type MenuTree struct {
	domain.Menu
	Children []*MenuTree `json:"children"`
}

func GetMenus(c *gin.Context) {
	var q MenuQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var menus []domain.Menu
	err = service.FindAll(&menus)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	var menuTree []*MenuTree
	for _, menu := range menus {
		var mt MenuTree
		copier.Copy(&mt, menu)
		menuTree = append(menuTree, &mt)
	}
	tree := buildTree(menuTree, 0)
	c.JSON(http.StatusOK, tree)
}

func buildTree(menuTree []*MenuTree, pid int64) []*MenuTree {
	var children []*MenuTree
	for _, menu := range menuTree {
		if menu.Pid == pid {
			menu.Children = buildTree(menuTree, menu.Id)
			children = append(children, menu)
		}
	}
	return children
}

func GetMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	var menu domain.Menu
	err = service.FindById(&menu, id)
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
	_, err = service.FindMenuByPath(menu.Path)
	if err == nil {
		c.String(http.StatusBadRequest, "路径已存在")
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	var menu domain.Menu
	err = c.ShouldBindJSON(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	menu.Id = id
	var oldMenu domain.Menu
	err = service.FindById(&oldMenu, menu.Id)
	if err != nil {
		c.String(http.StatusBadRequest, "机构不存在")
		return
	}
	if menu.Path != oldMenu.Path {
		_, err = service.FindMenuByPath(menu.Path)
		if err == nil {
			c.String(http.StatusBadRequest, "路径已存在")
			return
		}
	}

	err = service.Update(&menu)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	err = service.DeleteById(domain.Menu{}, id)
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	err = service.DeleteMenusByPid(id)
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
