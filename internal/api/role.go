package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
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
	MenuIds []int64 `json:"menuIds"`
}

func GetRoles(c *gin.Context) {
	var q RoleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Role{})
	c.JSON(http.StatusOK, page)
}

func GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var role domain.Role
	err = service.FindById(&role, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, role)
}

func CreateRole(c *gin.Context) {
	var roleAdd RoleAdd
	err := c.ShouldBindJSON(&roleAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
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
		var rmr []domain.RoleMenuRelation
		for _, id := range roleAdd.MenuIds {
			rmr = append(rmr, domain.RoleMenuRelation{
				Id:     config.IdGenerate(),
				RoleId: roleId,
				MenuId: id,
			})
		}
		if err = tx.Create(&rmr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "新增失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(roleId))
}

func UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var roleAdd RoleAdd
	err = c.ShouldBindJSON(&roleAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	roleId := int64(id)
	var role domain.Role
	err = service.FindById(&role, roleId)
	if err != nil {
		c.String(http.StatusBadRequest, "角色不存在")
		return
	}
	if roleAdd.Code != role.Code {
		var role domain.Role
		err = service.FindByCode(&role, roleAdd.Name)
		if err == nil {
			c.String(http.StatusBadRequest, "名称已存在")
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
		var rmr []domain.RoleMenuRelation
		for _, id := range roleAdd.MenuIds {
			rmr = append(rmr, domain.RoleMenuRelation{
				Id:     config.IdGenerate(),
				RoleId: roleId,
				MenuId: id,
			})
		}
		if err = config.DB.Create(&rmr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var role domain.Role
	err = service.FindById(&role, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
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
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
