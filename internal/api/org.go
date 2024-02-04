package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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

func GetOrgs(c *gin.Context) {
	var q OrgQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Org{})
	c.JSON(http.StatusOK, page)
}

func GetOrg(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	var org domain.Org
	err = service.FindById(&org, id)
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
		log.Println(err)
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var org domain.Org
	err = service.FindByName(&org, orgAdd.Name)
	if err == nil {
		c.String(http.StatusBadRequest, "名称已存在")
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
		c.String(http.StatusBadRequest, "新增失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(orgId))
}

func UpdateOrg(c *gin.Context) {
	orgId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	var orgAdd OrgAdd
	err = c.ShouldBindJSON(&orgAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var org domain.Org
	err = service.FindById(&org, orgId)
	if err != nil {
		c.String(http.StatusBadRequest, "机构不存在")
		return
	}
	if orgAdd.Name != org.Name {
		var org domain.Org
		err = service.FindByName(&org, orgAdd.Name)
		if err == nil {
			c.String(http.StatusBadRequest, "名称已存在")
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
			if err = config.DB.Create(&omr).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteOrg(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		return
	}
	var org domain.Org
	err = service.FindById(&org, id)
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
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
		c.String(http.StatusBadRequest, "删除失败")
		return
	}

	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
