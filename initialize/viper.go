package initialize

import (
	"log"
	"mangosteen/config"
	"mangosteen/global"

	"github.com/spf13/viper"
)

func InitViper() {
	config := struct {
		Viper config.ViperConfig
		Jwt   struct{ Path string }
	}{}

	viper.AutomaticEnv()
	pwd := viper.GetString("ROOT_DIR")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	global.ViperConfig = config.Viper
	global.RootPath = pwd
	global.JwtPath = pwd + "/" + config.Jwt.Path
}
