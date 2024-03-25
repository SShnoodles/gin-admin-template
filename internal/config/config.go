package config

import (
	"github.com/spf13/viper"
	"os"
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

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(wd)
	v.AddConfigPath(wd + "/config")
	v.SetConfigType("yml")
	err = v.ReadInConfig()
	if err := v.ReadInConfig(); err != nil {
		Log.Fatal("read conf failed ", err)
	}
	if err := v.Unmarshal(&AppConfig); err != nil {
		Log.Fatal("unable to decode into struct ", err)
	}
}

func IsDefaultLanguage() bool {
	return AppConfig.Language == "en"
}
