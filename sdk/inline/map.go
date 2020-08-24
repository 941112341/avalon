package inline

import "reflect"

func KeySet(o interface{}) []interface{} {
	set := make([]interface{}, 0)
	v := reflect.ValueOf(o).Elem()

	for i := 0; i < v.Len(); i++ {
		element := v.Index(i)
		set = append(set, element.Interface())
	}
	return set
}

func StringKeySet(o interface{}) []string {
	set := make([]string, 0)
	v := reflect.ValueOf(o).Elem()

	for i := 0; i < v.Len(); i++ {
		element := v.Index(i)
		if element.Kind() != reflect.String {
			continue
		}
		set = append(set, element.Interface().(string))
	}
	return set
}
