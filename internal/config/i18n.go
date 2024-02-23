package config

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var I18nLoc *i18n.Localizer

func init() {
	var bundle *i18n.Bundle
	if IsDefaultLanguage() {
		bundle = i18n.NewBundle(language.English)
	} else {
		bundle = i18n.NewBundle(language.Chinese)
	}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("locales/en.toml")
	bundle.MustLoadMessageFile("locales/zh-CN.toml")

	if IsDefaultLanguage() {
		I18nLoc = i18n.NewLocalizer(bundle, "en")
	} else {
		I18nLoc = i18n.NewLocalizer(bundle, "zh-CN")
	}

}
