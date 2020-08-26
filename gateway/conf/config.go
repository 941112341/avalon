package conf

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/spf13/viper"
)

func InitViper() error {
	viper.AddConfigPath(".")
	env := inline.GetEnv("env", "dev")
	fileName := fmt.Sprintf("%s.%s", "mapper", env)
	viper.SetConfigName(fileName)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.WatchConfig()
	return nil
}
