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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&viperConfig); err != nil {
		log.Fatalln(err)
	}
	global.ViperConfig = viperConfig.Viper
}
