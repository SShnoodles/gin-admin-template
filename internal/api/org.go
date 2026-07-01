package api

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
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
		service.Ok(c, orgs)
		return
	}
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.Org{})
	service.Ok(c, page)
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
	service.Ok(c, org)
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
	service.Ok(c, menusIds)
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
	orgId, err := service.CreateOrg(orgAdd.Org, orgAdd.MenuIds)
	if errors.Is(err, service.ErrOrgNameExists) {
		service.ConflictResult(c, "Existed.name")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, domain.NewIdWrapper(orgId))
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
	err = service.UpdateOrg(orgId, orgAdd.Org, orgAdd.MenuIds)
	if errors.Is(err, service.ErrOrgNotFound) {
		service.BadRequestResult(c, "NotExist.org")
		return
	}
	if errors.Is(err, service.ErrOrgNameExists) {
		service.ConflictResult(c, "Existed.name")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
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
	err = service.DeleteOrg(id)
	if errors.Is(err, service.ErrOrgNotFound) {
		service.BadRequestResult(c, "NotExist.org")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.delete")
		config.Log.Error(err.Error())
		return
	}

	service.Ok(c, service.DeleteSuccessResult())
}
