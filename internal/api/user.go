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

func GetUsers(c *gin.Context) {
	var q UserQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
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

func GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var user domain.User
	err = service.FindById(&user, id)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

func GetUserRoles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	roles, err := service.FindRoleIdsByUserId(id)
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func CreateUser(c *gin.Context) {
	var userAdd UserAdd
	err := c.ShouldBindJSON(&userAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	user, _ := service.FindUserByUsername(userAdd.Username)
	if user != (domain.User{}) {
		c.String(http.StatusBadRequest, "用户已存在")
		config.Log.Error(err.Error())
		return
	}

	userId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		password, _ := util.EncryptedPassword(userAdd.Password)
		user := domain.User{
			Id:       userId,
			Username: userAdd.Username,
			RealName: userAdd.RealName,
			WorkNo:   userAdd.WorkNo,
			Password: password,
			OrgId:    userAdd.OrgId,
		}
		if err = tx.Create(&user).Error; err != nil {
			return err
		}
		if userAdd.RoleIds != nil {
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
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(userId))
}

func UpdateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}
	var userAdd UserAdd
	err = c.ShouldBindJSON(&userAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		config.Log.Error(err.Error())
		return
	}
	var user domain.User
	err = service.FindById(&user, userId)
	if err != nil {
		c.String(http.StatusBadRequest, "用户不存在")
		config.Log.Error(err.Error())
		return
	}
	if user.Username != userAdd.Username {
		_, err := service.FindUserByUsername(userAdd.Username)
		if err == nil {
			c.String(http.StatusBadRequest, "用户名已存在")
			return
		}
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		user.Username = userAdd.Username
		user.Password = userAdd.Password
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
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "参数不正确")
		config.Log.Error(err.Error())
		return
	}

	var user domain.User
	err = service.FindById(&user, id)
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
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
		c.String(http.StatusBadRequest, "删除失败")
		config.Log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
