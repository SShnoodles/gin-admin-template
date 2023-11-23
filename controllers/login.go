package controllers

import (
	"context"
	"gin-admin-template/initializations"
	"gin-admin-template/middlewares"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"gin-admin-template/utils"
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
	msg := middlewares.ValidateParam(&login)
	if msg != "" {
		c.String(http.StatusBadRequest, msg)
		return
	}
	// check code
	val, err := initializations.RDB.Get(context.Background(), "code:"+login.Uuid).Result()
	if err != nil {
		c.String(http.StatusBadRequest, "验证码错误")
		return
	}
	defer func() {
		initializations.RDB.Del(context.Background(), "code:"+login.Uuid)
	}()
	if login.Code != val {
		c.String(http.StatusUnauthorized, "验证码错误")
		return
	}
	// check password
	user, err := services.FindUserByUsername(login.Username)
	if user == (models.User{}) {
		c.String(http.StatusUnauthorized, "用户不存在")
		return
	}
	isRight := utils.VerifyPassword(login.Password, user.Password)
	if !isRight {
		c.String(http.StatusUnauthorized, "密码错误")
		return
	}
	// create jwt
	jwt, err := utils.GenerateToken(user.Id)
	var token = make(map[string]string)
	token["token"] = "Bearer " + jwt
	c.JSON(http.StatusOK, jwt)
}
