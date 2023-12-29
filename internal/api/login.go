package api

import (
	"context"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/middleware"
	"gin-admin-template/internal/service"
	"gin-admin-template/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginInfo struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Uuid     string `json:"uuid" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

func Login(c *gin.Context) {
	var login LoginInfo
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	msg := middleware.ValidateParam(&login)
	if msg != "" {
		c.String(http.StatusBadRequest, msg)
		return
	}
	// check code
	val, err := config.RDB.Get(context.Background(), "code:"+login.Uuid).Result()
	if err != nil {
		c.String(http.StatusBadRequest, "验证码错误")
		return
	}
	defer func() {
		config.RDB.Del(context.Background(), "code:"+login.Uuid)
	}()
	if login.Code != val {
		c.String(http.StatusUnauthorized, "验证码错误")
		return
	}
	// check password
	user, err := service.FindUserByUsername(login.Username)
	if user == (domain.User{}) {
		c.String(http.StatusUnauthorized, "用户不存在")
		return
	}
	isRight := util.VerifyPassword(login.Password, user.Password)
	if !isRight {
		c.String(http.StatusUnauthorized, "密码错误")
		return
	}
	// create jwt
	jwt, err := util.GenerateToken(user.Id)
	var token = make(map[string]string)
	token["token"] = "Bearer " + jwt
	c.JSON(http.StatusOK, jwt)
}
