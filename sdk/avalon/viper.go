package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/spf13/viper"
	"reflect"
	"sync"
)

// implements set config
var once sync.Once

type ViperBean interface {
	Bean
	Key() string
}

type Viper struct {
	Bean
	Config func() error
}

func (v *Viper) Key() string {
	vb, ok := v.Bean.(ViperBean)
	if ok {
		return vb.Key()
	}

	return reflect.TypeOf(v.Bean).Elem().Name()
}

func (v *Viper) Initial() (err error) {
	once.Do(func() {
		err = v.Config()
	})
	if err := viper.UnmarshalKey(v.Key(), &v.Bean); err != nil {
		return err
	}

	if err := InitialBySub(v.Bean); err != nil {
		return err
	}

	InitialByDefault(v.Bean)
	return v.Bean.Initial()
}

// 实现不同组成部分从不同key读取
func InitialBySub(bean Bean) error {
	typ := reflect.TypeOf(bean).Elem()
	defer inline.Recovery()
	val := reflect.ValueOf(bean).Elem()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		viperKey, ok := field.Tag.Lookup("viper")
		if !ok {
			continue
		}

		o := viper.Get(viperKey)
		// danger
		fldVal := val.Field(i)
		if fldVal.CanSet() && o != nil {
			of := reflect.ValueOf(o)
			if fldVal.Kind() != of.Kind() {
				inline.WithFields("field", field.Name).Warnln("kind different")
				continue
			}
			fldVal.Set(of)
		}
	}
	return nil
}

func InitialByDefault(bean Bean) {
	inline.SetDefaultValue(bean)
}

func NewBean(bean Bean) Bean {
	return &Viper{Bean: bean, Config: func() error {
		viper.SetConfigName("avalon")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$GOPATH/src/github.com/941112341/avalon/sdk/avalon")
		return viper.ReadInConfig()
	}}
}

func NewBeanFunc(bean Bean, function func() error) Bean {
	return &Viper{Bean: bean, Config: function}
}
