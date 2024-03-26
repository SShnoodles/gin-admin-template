package api

import (
	"gin-admin-template/internal/config"
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

// Login
// @Summary login 用户登录
// @Tags login 登录
// @Accept json
// @Produce json
// @Param data body LoginInfo true "login info 信息"
// @Router /login/account [post]
func Login(c *gin.Context) {
	var login LoginInfo
	err := c.ShouldBindJSON(&login)
	if err != nil {
		service.ParamBadRequestResult(c)
		config.Log.Error(err.Error())
		return
	}
	msg := middleware.ValidateParam(&login)
	if msg != "" {
		c.String(http.StatusBadRequest, msg)
		return
	}
	// check code
	if !base64Captcha.DefaultMemStore.Verify(login.CodeId, login.Code, true) {
		service.UnauthorizedResult(c, "Error.code")
		return
	}
	// check password
	user, err := service.FindUserByUsername(login.Username)
	if user == (domain.User{}) {
		service.UnauthorizedResult(c, "NotExist.user")
		return
	}
	isRight := util.VerifyPassword(login.Password, user.Password)
	if !isRight {
		service.UnauthorizedResult(c, "Error.password")
		return
	}
	// create jwt
	jwt, expiresAt, err := util.GenerateToken(user.Id)
	result := LoginResult{
		AccessToken: jwt,
		Expires:     expiresAt,
	}
	c.JSON(http.StatusOK, result)
}

// Captcha
// @Summary captcha 验证码
// @Tags login 登录
// @Accept json
// @Produce json
// @Router /login/captcha [post]
func Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	var result = make(map[string]string)
	result["codeId"] = id
	result["code"] = b64s
	c.JSON(http.StatusOK, result)
}
