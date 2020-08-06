package inline

import "reflect"

func MustStringIfaceMap(o interface{}) map[string]interface{} {
	val := reflect.ValueOf(o)
	typ := val.Type()
	if typ.Kind() != reflect.Map {
		return nil
	}

	m := make(map[string]interface{})
	values := val.MapKeys()

	for _, key := range values {
		value := val.MapIndex(key)
		k := MustString(key.Interface())
		m[k] = value.Interface()
	}
	return m
}

func MustIfaceList(o interface{}) []interface{} {
	val := reflect.ValueOf(o)
	typ := val.Type()
	if typ.Kind() != reflect.Slice || typ.Kind() != reflect.Array {
		return nil
	}

	a := make([]interface{}, 0)
	l := typ.Len()
	for i := 0; i < l; i++ {
		v := val.Index(i)
		a = append(a, v.Interface())
	}
	return a
}
