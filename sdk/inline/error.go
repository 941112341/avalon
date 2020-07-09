package inline

import "github.com/941112341/avalon/sdk/log"

func Recover() {
	i := recover()
	if i != nil {
		log.New().WithField("recover", ToJsonString(i)).Errorln("panic!!")
	}
}
