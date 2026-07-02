package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server       Server
	Logging      Logging
	Datasource   Datasource
	Redis        Redis
	Jwt          Jwt
	Language     string
	Verification Verification
}

type Server struct {
	Port int
}

type Logging struct {
	File  File
	Level string
}

type File struct {
	Name string
	Path string
}

type Datasource struct {
	Driver   string
	Url      string
	Username string
	Password string
}

type Jwt struct {
	Secret string
	Expire int
}

type Verification struct {
	ResourceEnabled bool
}

type Redis struct {
	Addr     string
	Password string
	Db       int
}

var AppConfig = &Config{
	Language: "en",
}

func Init() error {
	if err := LoadConfig(); err != nil {
		return err
	}
	InitLogger()
	if err := InitI18n(); err != nil {
		return err
	}
	if err := InitIDGenerator(); err != nil {
		return err
	}
	InitRedis()
	if err := InitDB(); err != nil {
		return err
	}
	return nil
}

func LoadConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(wd)
	v.AddConfigPath(wd + "/config")
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read conf failed: %w", err)
	}
	if err := v.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}

func IsDefaultLanguage() bool {
	return AppConfig.Language == "en"
}
