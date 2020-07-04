package avalon

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOptionMerge(t *testing.T) {

	s := reflect.TypeOf(Config{})
	for i := 0; i < s.NumField(); i++ {
		fld := s.Field(i)
		fmt.Println(fld.Name)
	}
}
