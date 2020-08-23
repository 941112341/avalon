package inline

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

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
		var x int64
		x, err = strconv.ParseInt(param, 10, 64)
		b = int(x)
	case reflect.Int8:
		var x int64
		b, err = strconv.ParseInt(param, 10, 8)
		b = int8(x)
	case reflect.Int16:
		var x int64
		b, err = strconv.ParseInt(param, 10, 16)
		b = int16(x)
	case reflect.Int32:
		var x int64
		b, err = strconv.ParseInt(param, 10, 32)
		b = int32(x)
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
	case reflect.Slice:
		array := strings.Split(param, ",")
		av := reflect.MakeSlice(typ, 0, 0)
		for _, p := range array {
			ele, err := Convert(typ.Elem(), p)
			if err != nil {
				continue
			}
			av = reflect.Append(av, reflect.ValueOf(ele))
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

func SetDefaultValue(o interface{}) {
	val := reflect.ValueOf(o).Elem()
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return
	}
	defer Recovery()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		if !fieldVal.IsZero() {
			continue
		}

		fieldTyp := typ.Field(i)
		defaultVal, ok := fieldTyp.Tag.Lookup("default")
		if !ok {
			continue
		}

		v, err := Convert(fieldTyp.Type, defaultVal)
		if err != nil {
			continue
		}
		if fieldVal.CanSet() {

			of := reflect.ValueOf(v)
			if of.Kind() != fieldVal.Kind() {
				WithFields("field", fieldTyp.Name).Warnln("field kind different")
				continue
			}
			fieldVal.Set(of)
		}
	}
}

func Redirect(typ reflect.Type) reflect.Type {
	if typ == nil {
		return nil
	}

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
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

func RangeSlice(o interface{}, function func(o interface{}) error) error {
	slice := reflect.ValueOf(o)
	for i := 0; i < slice.Len(); i++ {
		o := slice.Index(i).Interface()
		if err := function(o); err != nil {
			return err
		}
	}
	return nil
}
