package server

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/spf13/viper"
	"testing"
)

func TestViper(t *testing.T) {
	zk := Zookeeper{}
	viper.AddConfigPath("../")
	viper.SetConfigName("avalon")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.UnmarshalKey("Zookeeper", &zk)
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(zk))
}
