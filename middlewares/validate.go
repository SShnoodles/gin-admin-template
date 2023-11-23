package middlewares

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
)

var Validate *validator.Validate
var Trans ut.Translator

func init() {
	Validate = validator.New()

	// 中文翻译器
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	err := zhtranslations.RegisterDefaultTranslations(Validate, trans)
	if err != nil {
		log.Fatal(err)
	}
	Trans = trans
}

func ValidateParam(param interface{}) string {
	if err := Validate.Struct(param); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err.Translate(Trans)
		}
	}
	return ""
}
