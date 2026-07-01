package api

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/middleware"
	"gin-admin-template/internal/service"

	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	CodeId   string `json:"codeId" validate:"required"`
	Code     string `json:"code" validate:"required"`
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
		service.Fail(c, 400, "400", msg)
		return
	}
	result, err := service.Login(login.Username, login.Password, login.CodeId, login.Code)
	if errors.Is(err, service.ErrCaptchaInvalid) {
		service.UnauthorizedResult(c, "Error.code")
		return
	}
	if errors.Is(err, service.ErrUserNotFound) {
		service.UnauthorizedResult(c, "NotExist.user")
		return
	}
	if errors.Is(err, service.ErrPasswordWrong) {
		service.UnauthorizedResult(c, "Error.password")
		return
	}
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, result)
}

// Captcha
// @Summary captcha 验证码
// @Tags login 登录
// @Accept json
// @Produce json
// @Router /login/captcha [post]
func Captcha(c *gin.Context) {
	result, err := service.GenerateCaptcha()
	if err != nil {
		service.BadRequestResult(c, "Failed.create")
		config.Log.Error(err.Error())
		return
	}
	service.Ok(c, result)
}
