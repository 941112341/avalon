package inline

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// o need a list
func RandomList(o interface{}) interface{} {
	v := reflect.ValueOf(o)
	l := v.Len()
	i := rand.Intn(l)
	return v.Index(i).Interface()
}
