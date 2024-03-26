package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
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
	var menus []domain.Menu
	err = service.FindAll(&menus)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
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

// GetMenu
// @Summary Menu 获取菜单
// @Tags menus 菜单
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Menu ID"
// @Router /menus/{id} [get]
func GetMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var menu domain.Menu
	err = service.FindById(&menu, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, menu)
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	resourceIds, err := service.FindResourceIdsByMenuId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, resourceIds)
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
	_, err = service.FindMenuByPath(menuAdd.Path)
	if err == nil {
		service.ConflictResult(c, "Existed.path")
		return
	}
	menuId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		menu := domain.Menu{
			Id: menuId,
		}
		copier.Copy(&menu, &menuAdd)
		if err = tx.Create(&menu).Error; err != nil {
			return err
		}
		if len(menuAdd.ResourceIds) > 0 {
			var mrr []domain.MenuResourceRelation
			for _, id := range menuAdd.ResourceIds {
				resourceId, _ := strconv.ParseInt(id, 10, 64)
				mrr = append(mrr, domain.MenuResourceRelation{
					Id:         config.IdGenerate(),
					ResourceId: resourceId,
					MenuId:     menuId,
				})
			}
			if err = tx.Create(&mrr).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(menuId))
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
	menuId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var menuAdd MenuAdd
	err = c.ShouldBindJSON(&menuAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var menu domain.Menu
	err = service.FindById(&menu, menuAdd.Id)
	if err != nil {
		service.BadRequestResult(c, "NotExist.org")
		config.Log.Error(err.Error())
		return
	}
	if menuAdd.Path != menu.Path {
		_, err = service.FindMenuByPath(menuAdd.Path)
		if err == nil {
			service.ConflictResult(c, "Existed.path")
			return
		}
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		copier.Copy(&menu, menuAdd)
		if err = tx.Save(&menu).Error; err != nil {
			return err
		}
		var oldMrr []domain.MenuResourceRelation
		if err = tx.Where("menu_id = ?", menuId).Find(&oldMrr).Error; err != nil {
			return err
		}
		if len(oldMrr) > 0 {
			if err = tx.Delete(oldMrr).Error; err != nil {
				return err
			}
		}
		if len(menuAdd.ResourceIds) > 0 {
			var mrr []domain.MenuResourceRelation
			for _, id := range menuAdd.ResourceIds {
				resourceId, _ := strconv.ParseInt(id, 10, 64)
				mrr = append(mrr, domain.MenuResourceRelation{
					Id:         config.IdGenerate(),
					ResourceId: resourceId,
					MenuId:     menuId,
				})
			}
			if err = tx.Create(&mrr).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, service.UpdateSuccessResult())
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&domain.Menu{}, id).Error; err != nil {
			return err
		}
		if err = tx.Where("pid = ?", id).Delete(&domain.Menu{}).Error; err != nil {
			return err
		}
		if err = tx.Where("menu_id = ?", id).Delete(&domain.MenuResourceRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		service.BadRequestResult(c, "Failed.delete")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, service.DeleteSuccessResult())
}
