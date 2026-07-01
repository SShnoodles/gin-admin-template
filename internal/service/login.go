package service

import (
	"errors"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/util"
	"time"

	"github.com/mojocn/base64Captcha"
)

var (
	ErrCaptchaInvalid = errors.New("captcha invalid")
	ErrPasswordWrong  = errors.New("password wrong")
)

type LoginResult struct {
	AccessToken  string    `json:"accessToken"`
	Expires      time.Time `json:"expires"`
	RefreshToken string    `json:"refreshToken"`
}

type CaptchaResult struct {
	CodeId string `json:"codeId"`
	Code   string `json:"code"`
}

func Login(username string, password string, codeId string, code string) (LoginResult, error) {
	if !base64Captcha.DefaultMemStore.Verify(codeId, code, true) {
		return LoginResult{}, ErrCaptchaInvalid
	}
	user, err := FindUserByUsername(username)
	if user == (domain.User{}) {
		return LoginResult{}, ErrUserNotFound
	}
	if err != nil {
		return LoginResult{}, err
	}
	if !util.VerifyPassword(password, user.Password) {
		return LoginResult{}, ErrPasswordWrong
	}
	jwt, expiresAt, err := util.GenerateToken(user.Id)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		AccessToken: jwt,
		Expires:     expiresAt,
	}, nil
}

func GenerateCaptcha() (CaptchaResult, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return CaptchaResult{}, err
	}
	return CaptchaResult{
		CodeId: id,
		Code:   b64s,
	}, nil
}
