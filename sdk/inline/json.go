package inline

import (
	"encoding/json"
	"fmt"
)

func JsonString(o interface{}) string {
	body, _ := json.Marshal(o)
	return string(body)
}

func VString(o interface{}) string {
	return fmt.Sprintf("%+v", o)
}
