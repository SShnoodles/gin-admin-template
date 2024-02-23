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
		service.ParamBadRequestResult(c, err)
		return
	}
	var menus []domain.Menu
	err = service.FindAll(&menus)
	if err != nil {
		service.BadRequestResult(c, "Failed.query", err)
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
		service.ParamBadRequestResult(c, err)
		return
	}
	var menu domain.Menu
	err = service.FindById(&menu, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query", err)
		return
	}
	c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
	var menu domain.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	_, err = service.FindMenuByPath(menu.Path)
	if err == nil {
		service.ConflictResult(c, "Existed.path")
		return
	}
	menu.Id = config.IdGenerate()
	err = service.Insert(&menu)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(menu.Id))
}

func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	var menu domain.Menu
	err = c.ShouldBindJSON(&menu)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	menu.Id = id
	var oldMenu domain.Menu
	err = service.FindById(&oldMenu, menu.Id)
	if err != nil {
		service.BadRequestResult(c, "NotExist.org", err)
		return
	}
	if menu.Path != oldMenu.Path {
		_, err = service.FindMenuByPath(menu.Path)
		if err == nil {
			service.ConflictResult(c, "Existed.path")
			return
		}
	}

	err = service.Update(&menu)
	if err != nil {
		service.BadRequestResult(c, "Failed.update", err)
		return
	}
	c.JSON(http.StatusOK, service.UpdateSuccessResult())
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c, err)
		return
	}
	err = service.DeleteById(domain.Menu{}, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.delete", err)
		return
	}
	err = service.DeleteMenusByPid(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.delete", err)
		return
	}
	c.JSON(http.StatusOK, service.DeleteSuccessResult())
}
