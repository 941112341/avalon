package inline

import "encoding/json"

func JsonString(o interface{}) string {
	body, _ := json.Marshal(o)
	return string(body)
}
