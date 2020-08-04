package generic

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
)

func Type2thrift(typ reflect.Type) thrift.TType {
	switch kind := typ.Kind(); kind {
	case reflect.Ptr:
		return Kind2thrift(typ.Elem().Kind())
	default:
		return Kind2thrift(kind)
	}
}

func Kind2thrift(kind reflect.Kind) (ttype thrift.TType) {
	switch kind {
	case reflect.String:
		ttype = thrift.STRING
	case reflect.Map:
		ttype = thrift.MAP
	case reflect.Int64, reflect.Int:
		ttype = thrift.I64
	case reflect.Int32:
		ttype = thrift.I32
	case reflect.Int16:
		ttype = thrift.I16
	case reflect.Int8:
		ttype = thrift.I08
	case reflect.Float64, reflect.Float32:
		ttype = thrift.DOUBLE
	case reflect.Array, reflect.Slice:
		ttype = thrift.LIST
	case reflect.Struct:
		ttype = thrift.STRUCT
	case reflect.Bool:
		ttype = thrift.BOOL
	default:
		inline.WithFields("kind", kind).Debugln("unsupport kind")
	}
	return
}
