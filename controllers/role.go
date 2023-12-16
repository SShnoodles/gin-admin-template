package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RoleQuery struct {
	services.PageInfo
	Name string `form:"name"`
}

type RoleAdd struct {
	models.Role
	MenuIds []int64 `json:"menuIds"`
}

func GetRoles(c *gin.Context) {
	var q RoleQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var roles []models.Role
	page := services.Pagination(initializations.DB, q.PageIndex, q.PageSize, &roles)
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
	var roleAdd RoleAdd
	err := c.ShouldBindJSON(&roleAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	roleId := initializations.IdGenerate()
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		role := models.Role{
			Id:    roleId,
			Name:  roleAdd.Name,
			Code:  roleAdd.Code,
			OrgId: roleAdd.OrgId,
		}
		if err = tx.Create(&role).Error; err != nil {
			return err
		}
		var rmr []models.RoleMenuRelation
		for _, id := range roleAdd.MenuIds {
			rmr = append(rmr, models.RoleMenuRelation{
				Id:     initializations.IdGenerate(),
				RoleId: roleId,
				MenuId: id,
			})
		}
		if err = initializations.DB.Create(&rmr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "新增失败")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(roleId))
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
	var role models.Role
	err = services.FindById(&role, roleId)
	if err != nil {
		c.String(http.StatusBadRequest, "角色不存在")
		return
	}
	if roleAdd.Code != role.Code {
		var role models.Role
		err = services.FindByCode(&role, roleAdd.Name)
		if err == nil {
			c.String(http.StatusBadRequest, "名称已存在")
			return
		}
	}
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		role.Name = roleAdd.Name
		role.Code = roleAdd.Code
		role.OrgId = roleAdd.OrgId
		if err = tx.Save(&role).Error; err != nil {
			return err
		}
		var oldRmr []models.RoleMenuRelation
		if err = tx.Where("role_id = ?", roleId).Find(&oldRmr).Error; err != nil {
			return err
		}
		if len(oldRmr) > 0 {
			if err = tx.Delete(oldRmr).Error; err != nil {
				return err
			}
		}
		var rmr []models.RoleMenuRelation
		for _, id := range roleAdd.MenuIds {
			rmr = append(rmr, models.RoleMenuRelation{
				Id:     initializations.IdGenerate(),
				RoleId: roleId,
				MenuId: id,
			})
		}
		if err = initializations.DB.Create(&rmr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var role models.Role
	err = services.FindById(&role, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
		return
	}
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&role).Error; err != nil {
			return err
		}
		if err = tx.Where("role_id = ?", id).Delete(&models.RoleMenuRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
