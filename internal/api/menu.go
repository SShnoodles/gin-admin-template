package api

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
)

type MenuQuery struct {
	service.PageInfo
	Name string `form:"name"`
}
type MenuAdd struct {
	domain.Menu
	ResourceIds []string `json:"resourceIds,omitempty"`
}

// GetMenus
// @Summary List Menus 获取菜单列表
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "name 名称"
// @Router /menus [get]
func GetMenus(c *gin.Context) {
	var q MenuQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	if !service.ValidatePageInfo(q.PageInfo) {
		service.ParamBadRequestResult(c)
		return
	}
	tree, err := service.FindMenuTree()
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, tree)
}

// GetMenu
// @Summary Menu 获取菜单
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Menu ID"
// @Router /menus/{id} [get]
func GetMenu(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var menu domain.Menu
	err := service.FindById(&menu, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, menu)
}

// GetMenuResources
// @Summary Menu resources 获取菜单资源
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Menu ID"
// @Router /menus/{id}/resources [get]
func GetMenuResources(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	resourceIds, err := service.FindResourceIdsByMenuId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, resourceIds)
}

// CreateMenu
// @Summary Create menu 创建菜单
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body MenuAdd true "Menu info 菜单信息"
// @Router /menus [post]
func CreateMenu(c *gin.Context) {
	var menuAdd MenuAdd
	err := c.ShouldBindJSON(&menuAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	menuId, err := service.CreateMenu(menuAdd.Menu, menuAdd.ResourceIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrMenuPathExists) {
		service.ConflictResult(c, "Existed.path")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, domain.NewIdWrapper(menuId))
}

// UpdateMenu
// @Summary Update menu 更新菜单
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Menu ID"
// @Param data body MenuAdd true "Menu info 菜单信息"
// @Router /menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	menuId, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var menuAdd MenuAdd
	err := c.ShouldBindJSON(&menuAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	err = service.UpdateMenu(menuId, menuAdd.Menu, menuAdd.ResourceIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrMenuNotFound) {
		service.BadRequestResult(c, "NotExist.org")
		return
	}
	if errors.Is(err, service.ErrMenuPathExists) {
		service.ConflictResult(c, "Existed.path")
		return
	}
	if errors.Is(err, service.ErrInvalidMenuParent) {
		service.ParamBadRequestResult(c)
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
}

// DeleteMenu
// @Summary Delete menu 删除菜单
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Menu ID"
// @Router /menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	err := service.DeleteMenu(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.delete")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.DeleteSuccessResult())
}
