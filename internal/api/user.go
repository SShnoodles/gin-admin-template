package api

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserQuery struct {
	service.PageInfo
	Username string `form:"username"`
	Mobile   string `form:"mobile"`
}

type UserAdd struct {
	domain.User
	RoleIds []string `json:"roleIds,omitempty"`
}

type UserOrg struct {
	domain.User
	OrgName string `json:"orgName"`
}

type UserPassword struct {
	OldPassword string `form:"oldPassword"`
	NewPassword string `form:"newPassword"`
}

// GetUsers
// @Summary List users 用户列表
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "name 名称"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var q UserQuery
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
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, []domain.User{})
	result := service.PagedResult[UserOrg]{
		Total: page.Total,
	}
	for _, d := range page.Data {
		var org domain.Org
		err := service.FindById(&org, d.OrgId)
		if err == nil {
			var userOrg UserOrg
			copier.Copy(&userOrg, &d)
			userOrg.OrgName = org.Name
			userOrg.Password = ""
			result.Data = append(result.Data, userOrg)
		}
	}
	service.Ok(c, result)
}

// GetUser
// @Summary User 获取用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var user domain.User
	err := service.FindById(&user, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	user.Password = ""

	service.Ok(c, user)
}

// GetUserRoles
// @Summary User roles 获取用户角色
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Router /users/{id}/roles [get]
func GetUserRoles(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	roles, err := service.FindRoleIdsByUserId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, roles)
}

// CreateUser
// @Summary Create user 创建用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body UserAdd true "User info 用户信息"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var userAdd UserAdd
	err := c.ShouldBindJSON(&userAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	userId, err := service.CreateUser(userAdd.User, userAdd.RoleIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrUserExists) {
		service.BadRequestResult(c, "Existed.user")
		return
	}
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, domain.NewIdWrapper(userId))
}

// UpdateUser
// @Summary Update users 更新用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param data body UserAdd true "User info 用户信息"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	userId, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	var userAdd UserAdd
	err := c.ShouldBindJSON(&userAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	err = service.UpdateUserInfo(userId, userAdd.User, userAdd.RoleIds)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		service.BadRequestResult(c, "NotExist.user")
		return
	}
	if errors.Is(err, service.ErrUserExists) {
		service.ConflictResult(c, "Existed.user")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
}

// EnabledUser
// @Summary Enabled user 启用/禁用用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Router /users/{id}/enabled [put]
func EnabledUser(c *gin.Context) {
	userId, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}
	err := service.ToggleUserEnabled(userId)
	if errors.Is(err, service.ErrUserNotFound) {
		service.BadRequestResult(c, "NotExist.user")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
}

// ChangeUserPassword
// @Summary Change user password 修改用户密码
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param data body UserPassword true "User password 用户密码"
// @Router /users/change-password [put]
func ChangeUserPassword(c *gin.Context) {
	userId := c.GetInt64("UserId")

	var userPassword UserPassword
	err := c.ShouldBindJSON(&userPassword)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}

	err = service.ChangeUserPassword(userId, userPassword.OldPassword, userPassword.NewPassword)
	if errors.Is(err, service.ErrInvalidParam) {
		service.ParamBadRequestResult(c)
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		service.BadRequestResult(c, "NotExist.user")
		return
	}
	if errors.Is(err, service.ErrUserPasswordWrong) {
		service.BadRequestResult(c, "Error.password")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.UpdateSuccessResult())
}

// DeleteUser
// @Summary Delete user 删除用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id, ok := service.ParseIdParam(c, "id")
	if !ok {
		return
	}

	err := service.DeleteUser(id)
	if errors.Is(err, service.ErrUserNotFound) {
		service.BadRequestResult(c, "NotExist.user")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.delete")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, service.DeleteSuccessResult())
}
