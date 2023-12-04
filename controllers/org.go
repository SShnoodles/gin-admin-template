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

type OrgQuery struct {
	services.PageInfo
	Name string `form:"name"`
}

type OrgAdd struct {
	models.Org
	MenuIds []int64 `json:"menuIds"`
}

func GetOrgs(c *gin.Context) {
	var q OrgQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var orgs []models.Org
	page := services.Pagination(initializations.DB, q.PageIndex, q.PageSize, &orgs)
	c.JSON(http.StatusOK, page)
}

func GetOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var org models.Org
	err = services.FindById(&org, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, org)
}

func CreateOrg(c *gin.Context) {
	var orgAdd OrgAdd
	err := c.ShouldBindJSON(&orgAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	orgId := initializations.IdGenerate()
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		org := models.Org{
			Id:         orgId,
			Name:       orgAdd.Name,
			CreditCode: orgAdd.CreditCode,
			Address:    orgAdd.Address,
		}
		if err = initializations.DB.Create(&org).Error; err != nil {
			return err
		}
		var omr []models.OrgMenuRelation
		for _, id := range orgAdd.MenuIds {
			omr = append(omr, models.OrgMenuRelation{
				Id:     initializations.IdGenerate(),
				OrgId:  orgId,
				MenuId: id,
			})
		}
		if err = initializations.DB.Create(&omr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "新增失败")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(orgId))
}

func UpdateOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var orgAdd OrgAdd
	err = c.ShouldBindJSON(&orgAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	orgId := int64(id)
	var org models.Org
	err = services.FindById(&org, orgId)
	if err != nil {
		c.String(http.StatusBadRequest, "机构不存在")
		return
	}
	if orgAdd.Name != org.Name {
		var org models.Org
		err = services.FindByName(&org, orgAdd.Name)
		if err == nil {
			c.String(http.StatusBadRequest, "名称已存在")
			return
		}
	}
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		org.Name = orgAdd.Name
		org.CreditCode = orgAdd.CreditCode
		org.Address = orgAdd.Address
		if err = tx.Save(&org).Error; err != nil {
			return err
		}
		var oldOmr []models.OrgMenuRelation
		if err = tx.Where("org_id = ?", orgId).Find(&oldOmr).Error; err != nil {
			return err
		}
		if len(oldOmr) > 0 {
			if err = tx.Delete(oldOmr).Error; err != nil {
				return err
			}
		}
		var omr []models.OrgMenuRelation
		for _, id := range orgAdd.MenuIds {
			omr = append(omr, models.OrgMenuRelation{
				Id:     initializations.IdGenerate(),
				OrgId:  orgId,
				MenuId: id,
			})
		}
		if err = initializations.DB.Create(&omr).Error; err != nil {
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

func DeleteOrg(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var org models.Org
	err = services.FindById(&org, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
		return
	}
	err = initializations.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&org).Error; err != nil {
			return err
		}
		if err = tx.Where("org_id = ?", id).Delete(&models.OrgMenuRelation{}).Error; err != nil {
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
