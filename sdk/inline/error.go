package inline

import (
	"errors"
	"github.com/941112341/avalon/sdk/log"
)

func Recover() {
	i := recover()
	if i != nil {
		log.New().WithField("recover", ToJsonString(i)).Errorln("panic!!")
	}
}

func RecoverErr() error {
	i := recover()
	if i == nil {
		return nil
	}
	log.New().WithField("recover", ToJsonString(i)).Errorln("panic!!")
	err, ok := i.(error)
	if !ok {
		return errors.New(ToJsonString(i))
	}
	return err
}
