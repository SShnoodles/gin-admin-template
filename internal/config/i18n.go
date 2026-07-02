package config

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var I18nLoc *i18n.Localizer

func InitI18n() error {
	var bundle *i18n.Bundle
	if IsDefaultLanguage() {
		bundle = i18n.NewBundle(language.English)
	} else {
		bundle = i18n.NewBundle(language.Chinese)
	}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if _, err := bundle.LoadMessageFile("locales/en.toml"); err != nil {
		return err
	}
	if _, err := bundle.LoadMessageFile("locales/zh-CN.toml"); err != nil {
		return err
	}

	if IsDefaultLanguage() {
		I18nLoc = i18n.NewLocalizer(bundle, "en")
	} else {
		I18nLoc = i18n.NewLocalizer(bundle, "zh-CN")
	}

	return nil
}
