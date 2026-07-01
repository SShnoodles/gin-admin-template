package api

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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

// GetRoles
// @Summary List roles 获取机构列表
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "name 名称"
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	var q RoleQuery
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
	service.Ok(c, result)
}

// GetRole
// @Summary Role 获取角色
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Org ID"
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var role domain.Role
	err := service.FindById(&role, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, role)
}

// GetRoleMenus
// @Summary Role menus 获取角色菜单
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Router /roles/{id}/menus [get]
func GetRoleMenus(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	menusIds, err := service.FindMenuIdsByRoleId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, menusIds)
}

// GetOrgRoles
// @Summary Org roles 获取机构角色
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param orgId path string true "Org ID"
// @Router /roles/orgs/{orgId} [get]
func GetOrgRoles(c *gin.Context) {
	orgId, ok := service.ParseIdParam(c, "orgId")
	if !ok {
		return
	}
	roles, err := service.FindRolesByOrgId(orgId)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, roles)
}

// CreateRole
// @Summary Create role 创建角色
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body RoleAdd true "Role info 角色信息"
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var roleAdd RoleAdd
	err := c.ShouldBindJSON(&roleAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	roleId, err := service.CreateRole(roleAdd.Role, roleAdd.MenuIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, domain.NewIdWrapper(roleId))
}

// UpdateRole
// @Summary Update role 更新角色
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param data body RoleAdd true "Role info 角色信息"
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
	roleId, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var roleAdd RoleAdd
	err := c.ShouldBindJSON(&roleAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	err = service.UpdateRole(roleId, roleAdd.Role, roleAdd.MenuIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrRoleNotFound) {
		service.BadRequestResult(c, "NotExist.role")
		return
	}
	if errors.Is(err, service.ErrRoleCodeExists) {
		service.ConflictResult(c, "Existed.code")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
}

// DeleteRole
// @Summary Delete role 删除角色
// @Tags roles 角色
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	err := service.DeleteRole(id)
	if errors.Is(err, service.ErrRoleNotFound) {
		service.BadRequestResult(c, "NotExist.role")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.delete")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.DeleteSuccessResult())
}
