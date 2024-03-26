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

type OrgQuery struct {
	service.PageInfo
	Name string `form:"name"`
}

type OrgAdd struct {
	domain.Org
	MenuIds []string `json:"menuIds,omitempty"`
}

// GetOrgs
// @Summary List orgs 获取机构列表
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "name 名称"
// @Router /orgs [get]
func GetOrgs(c *gin.Context) {
	var q OrgQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	if q.PageSize == 0 {
		var orgs []domain.Org
		err := service.FindAll(&orgs)
		if err != nil {
			service.BadRequestResult(c, "Failed.query")
			config.Log.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, orgs)
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Org{})
	c.JSON(http.StatusOK, page)
}

// GetOrg
// @Summary Org 获取机构
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Org ID"
// @Router /orgs/{id} [get]
func GetOrg(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var org domain.Org
	err = service.FindById(&org, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, org)
}

// GetOrgMenus
// @Summary Org menus 获取机构菜单
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Org ID"
// @Router /orgs/{id}/menus [get]
func GetOrgMenus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	menusIds, err := service.FindMenuIdsByOrgId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, menusIds)
}

// CreateOrg
// @Summary Create org 创建机构
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body OrgAdd true "Org info 机构信息"
// @Router /orgs [post]
func CreateOrg(c *gin.Context) {
	var orgAdd OrgAdd
	err := c.ShouldBindJSON(&orgAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var org domain.Org
	err = service.FindByName(&org, orgAdd.Name)
	if err == nil {
		service.ConflictResult(c, "Existed.name")
		return
	}

	orgId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		org := domain.Org{
			Id:         orgId,
			Name:       orgAdd.Name,
			CreditCode: orgAdd.CreditCode,
			Address:    orgAdd.Address,
		}
		if err = tx.Create(&org).Error; err != nil {
			return err
		}
		if orgAdd.MenuIds != nil {
			var omr []domain.OrgMenuRelation
			for _, id := range orgAdd.MenuIds {
				menuId, _ := strconv.ParseInt(id, 10, 64)
				omr = append(omr, domain.OrgMenuRelation{
					Id:     config.IdGenerate(),
					OrgId:  orgId,
					MenuId: menuId,
				})
			}
			if err = tx.Create(&omr).Error; err != nil {
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
	c.JSON(http.StatusOK, domain.NewIdWrapper(orgId))
}

// UpdateOrg
// @Summary Update org 更新机构
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Org ID"
// @Param data body OrgAdd true "Org info 机构信息"
// @Router /orgs/{id} [put]
func UpdateOrg(c *gin.Context) {
	orgId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var orgAdd OrgAdd
	err = c.ShouldBindJSON(&orgAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var org domain.Org
	err = service.FindById(&org, orgId)
	if err != nil {
		service.BadRequestResult(c, "NotExist.org")
		config.Log.Error(err.Error())
		return
	}
	if orgAdd.Name != org.Name {
		var org domain.Org
		err = service.FindByName(&org, orgAdd.Name)
		if err == nil {
			service.ConflictResult(c, "Existed.name")
			return
		}
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		org.Name = orgAdd.Name
		org.CreditCode = orgAdd.CreditCode
		org.Address = orgAdd.Address
		if err = tx.Save(&org).Error; err != nil {
			return err
		}
		var oldOmr []domain.OrgMenuRelation
		if err = tx.Where("org_id = ?", orgId).Find(&oldOmr).Error; err != nil {
			return err
		}
		if len(oldOmr) > 0 {
			if err = tx.Delete(oldOmr).Error; err != nil {
				return err
			}
		}
		if orgAdd.MenuIds != nil {
			var omr []domain.OrgMenuRelation
			for _, id := range orgAdd.MenuIds {
				menuId, _ := strconv.ParseInt(id, 10, 64)
				omr = append(omr, domain.OrgMenuRelation{
					Id:     config.IdGenerate(),
					OrgId:  orgId,
					MenuId: menuId,
				})
			}
			if err = tx.Create(&omr).Error; err != nil {
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

// DeleteOrg
// @Summary Delete org 删除机构
// @Tags orgs 机构
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Org ID"
// @Router /orgs/{id} [delete]
func DeleteOrg(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var org domain.Org
	err = service.FindById(&org, id)
	if err != nil {
		service.BadRequestResult(c, "NotExist.org")
		config.Log.Error(err.Error())
		return
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&org).Error; err != nil {
			return err
		}
		if err = tx.Where("org_id = ?", id).Delete(&domain.OrgMenuRelation{}).Error; err != nil {
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
