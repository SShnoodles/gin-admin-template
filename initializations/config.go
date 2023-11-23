package initializations

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Server     Server
	Datasource Datasource
	Redis      Redis
	Jwt        Jwt
}

type Server struct {
	Port int
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

type Redis struct {
	Addr     string
	Password string
	Db       int
}

var AppConfig = &Config{}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(wd)
	v.AddConfigPath(wd + "/initializations")
	v.SetConfigType("yml")
	err = v.ReadInConfig()
	if err := v.ReadInConfig(); err != nil {
		log.Fatal("read conf failed ", err)
	}
	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatal("unable to decode into struct ", err)
	}
}
