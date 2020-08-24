package conf

import (
	"github.com/spf13/viper"
)

func InitViper() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("mapper")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.WatchConfig()
	return nil
}
