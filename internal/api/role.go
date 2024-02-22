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

type RoleQuery struct {
	service.PageInfo
	Name string `form:"name"`
}

type RoleAdd struct {
	domain.Role
	MenuIds []string `json:"menuIds,omitempty"`
}

type RoleOrg struct {
	domain.Role
	OrgName string `json:"orgName"`
}

func GetRoles(c *gin.Context) {
	var q RoleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Role{})
	result := service.PagedResult[RoleOrg]{
		Total: page.Total,
	}
	for _, d := range page.Data {
		var org domain.Org
		err := service.FindById(&org, d.OrgId)
		if err == nil {
			var roleOrg RoleOrg
			copier.Copy(&roleOrg, &d)
			roleOrg.OrgName = org.Name
			result.Data = append(result.Data, roleOrg)
		}
	}
	c.JSON(http.StatusOK, result)
}

func GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var role domain.Role
	err = service.FindById(&role, id)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, role)
}

func GetRoleMenus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	menusIds, err := service.FindMenuIdsByRoleId(id)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, menusIds)
}

func GetOrgRoles(c *gin.Context) {
	orgId, err := strconv.ParseInt(c.Param("orgId"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	roles, err := service.FindRolesByOrgId(orgId)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func CreateRole(c *gin.Context) {
	var roleAdd RoleAdd
	err := c.ShouldBindJSON(&roleAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	roleId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		role := domain.Role{
			Id:    roleId,
			Name:  roleAdd.Name,
			Code:  roleAdd.Code,
			OrgId: roleAdd.OrgId,
		}
		if err = tx.Create(&role).Error; err != nil {
			return err
		}
		if len(roleAdd.MenuIds) > 0 {
			var rmr []domain.RoleMenuRelation
			for _, id := range roleAdd.MenuIds {
				menuId, _ := strconv.ParseInt(id, 10, 64)
				rmr = append(rmr, domain.RoleMenuRelation{
					Id:     config.IdGenerate(),
					RoleId: roleId,
					MenuId: menuId,
				})
			}
			if err = tx.Create(&rmr).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "新增失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(roleId))
}

func UpdateRole(c *gin.Context) {
	roleId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var roleAdd RoleAdd
	err = c.ShouldBindJSON(&roleAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	var role domain.Role
	err = service.FindById(&role, roleId)
	if err != nil {
		c.String(http.StatusBadRequest, "角色不存在")
		config.Log.Error(err.Error())
		return
	}
	if roleAdd.Code != role.Code {
		var role domain.Role
		err = service.FindByCode(&role, roleAdd.Code)
		if err == nil {
			c.String(http.StatusBadRequest, "编码已存在")
			return
		}
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		role.Name = roleAdd.Name
		role.Code = roleAdd.Code
		role.OrgId = roleAdd.OrgId
		if err = tx.Save(&role).Error; err != nil {
			return err
		}
		var oldRmr []domain.RoleMenuRelation
		if err = tx.Where("role_id = ?", roleId).Find(&oldRmr).Error; err != nil {
			return err
		}
		if len(oldRmr) > 0 {
			if err = tx.Delete(oldRmr).Error; err != nil {
				return err
			}
		}
		if len(roleAdd.MenuIds) > 0 {
			var rmr []domain.RoleMenuRelation
			for _, id := range roleAdd.MenuIds {
				menuId, _ := strconv.ParseInt(id, 10, 64)
				rmr = append(rmr, domain.RoleMenuRelation{
					Id:     config.IdGenerate(),
					RoleId: roleId,
					MenuId: menuId,
				})
			}
			if err = tx.Create(&rmr).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var role domain.Role
	err = service.FindById(&role, id)
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
		config.Log.Error(err.Error())
		return
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&role).Error; err != nil {
			return err
		}
		if err = tx.Where("role_id = ?", id).Delete(&domain.RoleMenuRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
