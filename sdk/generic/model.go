package generic

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"reflect"
	"strconv"
	"strings"
)

type lazyFieldProperties func() fieldProperties

// 对应 request/response
type structProperty struct {
	structName string
	thriftName structType

	properties lazyFieldProperties
}

type fieldProperty struct {
	fieldName  string
	thriftName string
	jsonName   string

	index int16

	kind      reflect.Kind
	elemKind  reflect.Kind
	keyKind   reflect.Kind
	valueKind reflect.Kind

	childProps lazyFieldProperties
}

func (fp *fieldProperty) convert() *CommonTStruct {
	if fp == nil {
		return nil
	}
	cts := &CommonTStruct{}
	cts.FieldName = fp.thriftName
	cts.StructName = fp.fieldName
	cts.Type = convert(fp.kind)
	cts.ID = fp.index
	cts.JSONPath = fp.jsonName

	switch fp.kind {
	case reflect.Array, reflect.Slice:
		acts := &CommonTStruct{Type: convert(fp.elemKind)}
		cts.ArrayStruct = acts
	case reflect.Map:
		cts.MapKeyStruct = &CommonTStruct{Type: convert(fp.keyKind)}
		cts.MapValueStruct = &CommonTStruct{Type: convert(fp.valueKind)}
	case reflect.Struct:
		cts.FieldMap = LazyField{
			lazy: func() []*CommonTStruct {
				return fp.childProps().convert()
			},
		}
	}
	return cts
}

type fieldProperties []*fieldProperty

func (f fieldProperties) convert() []*CommonTStruct {
	cts := make([]*CommonTStruct, 0)
	for _, property := range f {
		cts = append(cts, property.convert())
	}
	return cts
}

type structType string

var StructType = struct {
	Request  structType
	Response structType
}{
	Request:  "request",
	Response: "response",
}

type ThriftParser struct {
	MethodName string
}

func (s *ThriftParser) parseField(itype reflect.Type) lazyFieldProperties {

	return func() fieldProperties {
		if itype.Kind() == reflect.Ptr {
			itype = itype.Elem()
		}
		fieldProperties := make(fieldProperties, 0)
		for i := 0; i < itype.NumField(); i++ {
			fld := itype.Field(i)
			fname := fld.Name

			tstring, ok := fld.Tag.Lookup("thrift")
			if !ok {
				continue
			}
			jstring, ok := fld.Tag.Lookup("json")
			if !ok {
				jstring = tstring
			}

			if tss := strings.Split(tstring, ","); len(tss) < 2 {
				continue
			} else {
				tname, sCnt := tss[0], tss[1]
				cnt, err := strconv.ParseInt(sCnt, 10, 16)
				if err != nil {
					inline.WithFields("cnt", sCnt).Debugln("thrift parse cnt err")
					continue
				}

				jss := strings.Split(jstring, ",")

				fp := &fieldProperty{
					fieldName:  fname,
					thriftName: tname,
					jsonName:   jss[0],
					index:      int16(cnt),
				}
				ftype := fld.Type
				if ftype.Kind() == reflect.Ptr {
					ftype = ftype.Elem()
				}
				switch ftype.Kind() {
				case reflect.Slice:
					fp.kind = reflect.Slice
					fp.elemKind = ftype.Elem().Kind()
				case reflect.Array:
					fp.kind = reflect.Array
					fp.elemKind = ftype.Elem().Kind()
				case reflect.Map:
					fp.kind = reflect.Map
					fp.valueKind = ftype.Elem().Kind()
					fp.keyKind = ftype.Key().Kind()
				case reflect.Struct:
					fp.kind = reflect.Struct

					fp.childProps = s.parseField(ftype)
				default:
					fp.kind = ftype.Kind()
				}

				fieldProperties = append(fieldProperties, fp)
			}

		}
		return fieldProperties
	}
}

func (s *ThriftParser) ParseModel(thriftName structType, iface interface{}) (*CommonTStruct, error) {
	itype := reflect.TypeOf(iface)
	if itype.Kind() == reflect.Ptr {
		itype = itype.Elem()
	}

	properties := s.parseField(itype)

	sp := &structProperty{
		structName: itype.Name(),
		thriftName: thriftName,
		properties: properties,
	}
	return s.doParseModel(sp)
}

func (s *ThriftParser) doParseModel(sp *structProperty) (*CommonTStruct, error) {

	if sp.thriftName == StructType.Request {
		return s.doParseRequest(sp)
	} else if sp.thriftName == StructType.Response {
		return s.doParseResponse(sp)
	}
	return nil, errors.New("invalid thrift name")
}

func (s *ThriftParser) doParseRequest(sp *structProperty) (*CommonTStruct, error) {

	request := &CommonTStruct{
		ID:         1,
		StructName: sp.structName,
		FieldName:  string(sp.thriftName),
		JSONPath:   string(sp.thriftName),
		Type:       thrift.STRUCT,
		FieldMap: LazyField{
			lazy: func() []*CommonTStruct {
				return sp.properties().convert()
			},
		},
	}
	args := &CommonTStruct{
		StructName: fmt.Sprintf("%s_args", s.MethodName),
		Type:       thrift.STRUCT,
		FieldMap: LazyField{
			lazy: func() []*CommonTStruct {
				return []*CommonTStruct{
					request,
				}
			},
		},
	}
	return args, nil
}

func (s *ThriftParser) doParseResponse(sp *structProperty) (*CommonTStruct, error) {

	response := &CommonTStruct{
		ID:         0,
		StructName: sp.structName,
		FieldName:  string(sp.thriftName),
		JSONPath:   string(sp.thriftName),
		Type:       thrift.STRUCT,
		FieldMap: LazyField{
			lazy: func() []*CommonTStruct {
				return sp.properties().convert()
			},
		},
	}

	result := &CommonTStruct{
		StructName: fmt.Sprintf("%s_result", s.MethodName),
		Type:       thrift.STRUCT,
		FieldMap: LazyField{
			lazy: func() []*CommonTStruct {
				return []*CommonTStruct{
					response,
				}
			},
		},
	}
	return result, nil
}

func convert(kind reflect.Kind) (ttype thrift.TType) {
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
