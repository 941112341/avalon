package inline

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

func SetFieldJSON(any interface{}, fieldName string, to interface{}) (err error) {
	defer func() {
		if serr, ok := recover().(error); ok {
			err = serr
		}
	}()
	fld, err := getField(any, fieldName)
	if err != nil {
		return PrependErrorFmt(err, "set field json %s", fieldName)
	}
	typ := fld.Type()
	isPtr := typ.Kind() == reflect.Ptr
	if isPtr {
		typ = typ.Elem()
	}
	unptr := reflect.New(typ).Interface()
	data, err := jsoniter.Marshal(to)
	if err != nil {
		return errors.WithMessage(err, "marshal to "+ToJsonString(to))
	}
	err = jsoniter.Unmarshal(data, &unptr)
	if err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	if !fld.CanSet() {
		return errors.New(fieldName + " fld cannot set")
	}
	fld.Set(reflect.ValueOf(unptr))
	return nil
}

func SetField(any interface{}, fieldName string, to interface{}) (err error) {
	defer func() {
		if serr, ok := recover().(error); ok {
			err = serr
		}
	}()
	fld, err := getField(any, fieldName)
	if err != nil {
		return PrependErrorFmt(err, "SetField get field %s", fieldName)
	}
	if !fld.CanSet() {
		return errors.New(fieldName + " fld cannot set")
	}
	fld.Set(reflect.ValueOf(to))
	return
}

func GetField(any interface{}, fieldName string) (i interface{}, err error) {
	fld, err := getField(any, fieldName)
	if err != nil {
		return nil, PrependErrorFmt(err, "GetField %s", fieldName)
	}
	i = fld.Interface()
	return
}

func GetFieldAddress(any interface{}, fieldName string) (i interface{}, err error) {

	fld, err := getField(any, fieldName)
	if err != nil {
		return nil, errors.Wrap(err, "GetFieldAddress")
	}

	i = fld.Addr().Interface()
	return
}

func Set(fld reflect.Value, param string) (err error) {
	defer func() {
		if serr, ok := recover().(error); ok {
			err = serr
		}
	}()
	b, err := Convert(fld.Type(), param)
	if err != nil {
		return err
	}
	fld.Set(reflect.ValueOf(b))
	return
}

func getField(any interface{}, fieldName string) (fld reflect.Value, err error) {
	defer func() {
		if serr, ok := recover().(error); ok {
			err = serr
		}
	}()
	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("any is not ptr")
	}
	v = v.Elem()
	fld = v.FieldByName(fieldName)
	if !fld.IsValid() {
		return reflect.Value{}, fmt.Errorf("%s field not found", fieldName)
	}
	return fld, nil
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

func UnionValue(o interface{}) interface{} {
	typ := reflect.ValueOf(o)
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		if typ.Len() > 0 {
			return typ.Index(0).Interface()
		}
	case reflect.Map:
		if typ.Len() > 0 {
			key := typ.MapKeys()[0]
			return typ.MapIndex(key).Interface()
		}
	}
	return nil
}

func IsType(o interface{}, kind reflect.Kind) bool {
	return reflect.TypeOf(o).Kind() == kind
}
