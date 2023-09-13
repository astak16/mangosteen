package initialize

import (
	"log"
	"mangosteen/config"
	"mangosteen/global"

	"github.com/spf13/viper"
)

func InitViper() {
	viperConfig := struct {
		Viper config.ViperConfig
	}{}

	viper.AutomaticEnv()
	pwd := viper.GetString("ROOT_DIR")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&viperConfig); err != nil {
		log.Fatalln(err)
	}
	global.ViperConfig = viperConfig.Viper
}
