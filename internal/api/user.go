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

type UserQuery struct {
	service.PageInfo
	Username string `form:"username"`
	Mobile   string `form:"mobile"`
}

type UserAdd struct {
	domain.User
	RoleIds []int64 `json:"roleIds"`
}

func GetUsers(c *gin.Context) {
	var q UserQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var users []domain.User
	page := service.Pagination(config.DB, q.PageIndex, q.PageSize, &users)
	c.JSON(http.StatusOK, page)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var user domain.User
	err = service.FindById(&user, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var userAdd UserAdd
	err := c.ShouldBindJSON(&userAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	userId := config.IdGenerate()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		user := domain.User{
			Id:       userId,
			Username: userAdd.Username,
			RealName: userAdd.RealName,
			WorkNo:   userAdd.WorkNo,
			Password: userAdd.Password,
			OrgId:    userAdd.OrgId,
		}
		if err = tx.Create(&user).Error; err != nil {
			return err
		}
		var urr []domain.UserRoleRelation
		for _, id := range userAdd.RoleIds {
			urr = append(urr, domain.UserRoleRelation{
				Id:     config.IdGenerate(),
				UserId: userId,
				RoleId: id,
				OrgId:  userAdd.OrgId,
			})
		}
		if err = tx.Create(&urr).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, domain.NewIdWrapper(userId))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var userAdd UserAdd
	err = c.ShouldBindJSON(&userAdd)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	userId := int64(id)
	var user domain.User
	err = service.FindById(&user, userId)
	if err != nil {
		c.String(http.StatusBadRequest, "用户不存在")
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
		if err = tx.Where("user_id = ?", user).Find(&oldUrr).Error; err != nil {
			return err
		}
		if len(oldUrr) > 0 {
			if err = tx.Delete(oldUrr).Error; err != nil {
				return err
			}
		}

		var urr []domain.UserRoleRelation
		for _, id := range userAdd.RoleIds {
			urr = append(urr, domain.UserRoleRelation{
				Id:     config.IdGenerate(),
				UserId: userId,
				RoleId: id,
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
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("更新成功"))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	var user domain.User
	err = service.FindById(&user, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "数据不存在")
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
		return
	}
	c.JSON(http.StatusOK, domain.NewMessageWrapper("删除成功"))
}
