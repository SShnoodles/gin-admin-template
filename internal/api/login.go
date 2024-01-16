package api

import (
	"context"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/middleware"
	"gin-admin-template/internal/service"
	"gin-admin-template/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"time"
)

type LoginInfo struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	CodeId   string `json:"codeId" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type LoginResult struct {
	AccessToken  string    `json:"accessToken"`
	Expires      time.Time `json:"expires"`
	RefreshToken string    `json:"refreshToken"`
}

var ctx = context.Background()

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
	if !base64Captcha.DefaultMemStore.Verify(login.CodeId, login.Code, true) {
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
	jwt, expiresAt, err := util.GenerateToken(user.Id)
	result := LoginResult{
		AccessToken: "Bearer " + jwt,
		Expires:     expiresAt,
	}
	c.JSON(http.StatusOK, result)
}

func Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		c.String(http.StatusUnauthorized, "创建验证码失败")
		return
	}
	var result = make(map[string]string)
	result["codeId"] = id
	result["code"] = b64s
	c.JSON(http.StatusOK, result)
}
