package inline

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

func SetField(any interface{}, fieldName string, to interface{}) (err error) {
	i, err := GetFieldAddress(any, fieldName)
	if err != nil {
		return err
	}
	str, _ := jsoniter.Marshal(to)
	return jsoniter.Unmarshal(str, i)
}

func GetField(any interface{}, fieldName string) (i interface{}, err error) {
	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("any is not ptr")
	}
	v = v.Elem()
	fld := v.FieldByName(fieldName)
	if !fld.IsValid() {
		return nil, fmt.Errorf("%s field not found", fieldName)
	}
	defer func() {
		err = RecoverErr()
	}()
	i = fld.Interface()
	return
}

func GetFieldAddress(any interface{}, fieldName string) (i interface{}, err error) {
	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("any is not ptr")
	}
	v = v.Elem()
	fld := v.FieldByName(fieldName)
	if !fld.IsValid() {
		return nil, fmt.Errorf("%s field not found", fieldName)
	}
	defer func() {
		err = RecoverErr()
	}()
	i = fld.Addr().Interface()
	return
}

func Set(fld reflect.Value, param string) (err error) {
	defer func() {
		err = RecoverErr()
	}()
	b, err := Convert(fld.Type(), param)
	if err != nil {
		return err
	}
	fld.Set(reflect.ValueOf(b))
	return
}

func Convert(typ reflect.Type, param string) (b interface{}, err error) {
	switch typ.Kind() {
	case reflect.Bool:
		b, err = strconv.ParseBool(param)
	case reflect.Float32:
		b, err = strconv.ParseFloat(param, 32)
	case reflect.Float64:
		b, err = strconv.ParseFloat(param, 64)
	case reflect.Int:
		b, err = strconv.ParseInt(param, 10, 64)
	case reflect.Int8:
		b, err = strconv.ParseInt(param, 10, 8)
	case reflect.Int16:
		b, err = strconv.ParseInt(param, 10, 16)
	case reflect.Int32:
		b, err = strconv.ParseInt(param, 10, 32)
	case reflect.Int64:
		b, err = strconv.ParseInt(param, 10, 64)
	case reflect.Uint:
		b, err = strconv.ParseUint(param, 10, 64)
	case reflect.Uint8:
		b, err = strconv.ParseUint(param, 10, 8)
	case reflect.Uint16:
		b, err = strconv.ParseUint(param, 10, 16)
	case reflect.Uint32:
		b, err = strconv.ParseUint(param, 10, 32)
	case reflect.Uint64:
		b, err = strconv.ParseUint(param, 10, 64)
	case reflect.String:
		b = param
	case reflect.Array:
		av := reflect.MakeSlice(typ, typ.Len(), typ.Len())
		array := strings.Split(param, ",")
		for _, p := range array {
			ele, err := Convert(typ.Elem(), p)
			if err != nil {
				continue
			}
			av = reflect.AppendSlice(av, reflect.ValueOf(ele))
		}
		b = av.Interface()
	case reflect.Ptr:
		x, serr := Convert(typ.Elem(), param)
		b, err = &x, serr
	default:
		err = errors.New("don't support")
	}
	return
}

func Redirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}
