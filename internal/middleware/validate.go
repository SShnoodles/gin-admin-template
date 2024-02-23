package middleware

import (
	"gin-admin-template/internal/config"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var Validate *validator.Validate
var Trans ut.Translator

func init() {
	Validate = validator.New()

	if config.IsDefaultLanguage() {
		uni := ut.New(en.New())
		trans, _ := uni.GetTranslator("en")
		err := entranslations.RegisterDefaultTranslations(Validate, trans)
		if err != nil {
			config.Log.Fatal(err)
		}
		Trans = trans
	} else {
		uni := ut.New(zh.New())
		trans, _ := uni.GetTranslator("zh")
		err := zhtranslations.RegisterDefaultTranslations(Validate, trans)
		if err != nil {
			config.Log.Fatal(err)
		}
		Trans = trans
	}
}

func ValidateParam(param interface{}) string {
	if err := Validate.Struct(param); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err.Translate(Trans)
		}
	}
	return ""
}
