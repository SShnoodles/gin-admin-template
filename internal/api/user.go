package api

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"gin-admin-template/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	c.JSON(http.StatusOK, result)
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var user domain.User
	err = service.FindById(&user, id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	user.Password = ""

	c.JSON(http.StatusOK, user)
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	roles, err := service.FindRoleIdsByUserId(id)
	if err != nil {
		service.BadRequestResult(c, "Failed.query")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
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
	user, _ := service.FindUserByUsername(userAdd.Username)
	if user != (domain.User{}) {
		service.BadRequestResult(c, "Existed.user")
		config.Log.Error(err.Error())
		return
	}

	userId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		password, _ := util.EncryptedPassword(util.DefaultPassword)
		user := domain.User{
			Id:       userId,
			Username: userAdd.Username,
			RealName: userAdd.RealName,
			WorkNo:   userAdd.WorkNo,
			Password: password,
			OrgId:    userAdd.OrgId,
			Enabled:  true,
		}
		if err = tx.Create(&user).Error; err != nil {
			return err
		}
		if len(userAdd.RoleIds) > 0 {
			var urr []domain.UserRoleRelation
			for _, id := range userAdd.RoleIds {
				roleId, _ := strconv.ParseInt(id, 10, 64)
				urr = append(urr, domain.UserRoleRelation{
					Id:     config.IdGenerate(),
					UserId: userId,
					RoleId: roleId,
					OrgId:  userAdd.OrgId,
				})
			}
			if err = tx.Create(&urr).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(userId))
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
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var userAdd UserAdd
	err = c.ShouldBindJSON(&userAdd)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var user domain.User
	err = service.FindById(&user, userId)
	if err != nil {
		service.BadRequestResult(c, "NotExist.user")
		config.Log.Error(err.Error())
		return
	}
	if user.Username != userAdd.Username {
		_, err := service.FindUserByUsername(userAdd.Username)
		if err == nil {
			service.ConflictResult(c, "Existed.user")
			return
		}
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		user.Username = userAdd.Username
		user.RealName = userAdd.RealName
		user.WorkNo = userAdd.WorkNo
		user.OrgId = userAdd.OrgId
		if err = tx.Save(&user).Error; err != nil {
			return err
		}
		var oldUrr []domain.UserRoleRelation
		if err = tx.Where("user_id = ?", userId).Find(&oldUrr).Error; err != nil {
			return err
		}
		if len(oldUrr) > 0 {
			if err = tx.Delete(oldUrr).Error; err != nil {
				return err
			}
		}

		if len(userAdd.RoleIds) > 0 {
			var urr []domain.UserRoleRelation
			for _, id := range userAdd.RoleIds {
				roleId, _ := strconv.ParseInt(id, 10, 64)
				urr = append(urr, domain.UserRoleRelation{
					Id:     config.IdGenerate(),
					UserId: userId,
					RoleId: roleId,
					OrgId:  userAdd.OrgId,
				})
			}
			if err = tx.Create(&urr).Error; err != nil {
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

// EnabledUser
// @Summary Enabled user 启用/禁用用户
// @Tags users 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Router /users/{id}/enabled [put]
func EnabledUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	var user domain.User
	err = service.FindById(&user, userId)
	if err != nil {
		service.BadRequestResult(c, "NotExist.user")
		config.Log.Error(err.Error())
		return
	}
	err = config.DB.Model(&user).UpdateColumn("enabled", !user.Enabled).Error
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, service.UpdateSuccessResult())
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

	var user domain.User
	err = service.FindById(&user, userId)
	if err != nil {
		service.BadRequestResult(c, "NotExist.user")
		config.Log.Error(err.Error())
		return
	}
	isRight := util.VerifyPassword(userPassword.OldPassword, user.Password)
	if !isRight {
		service.BadRequestResult(c, "Error.password")
		return
	}
	password, _ := util.EncryptedPassword(userPassword.NewPassword)
	err = config.DB.Model(&user).UpdateColumn("password", password).Error
	if err != nil {
		service.BadRequestResult(c, "Failed.update")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, service.UpdateSuccessResult())
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}

	var user domain.User
	err = service.FindById(&user, id)
	if err != nil {
		service.BadRequestResult(c, "NotExist.user")
		config.Log.Error(err.Error())
		return
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Delete(&user).Error; err != nil {
			return err
		}
		if err = tx.Where("user_id = ?", id).Delete(&domain.UserRoleRelation{}).Error; err != nil {
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
