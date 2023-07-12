package utils

import (
	"vngitSub/model"
	"github.com/spf13/viper"
)

//LoadConfig - load default configs
func LoadConfig(path string) (config model.Default, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}