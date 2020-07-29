package generic

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrRegexMatch inline.AvalonErrorCode = 100 + iota
	ErrNilPackage
	ErrNilStruct
	ErrUnknownType
)

type Parser interface {
	Parse(content string) error
}

type ThriftContext interface {
	GetStruct(base, structName string) (*ThriftStructModel, bool)
	GetService(base, serviceName string) (*ThriftServiceModel, bool)
	GetFile(base string) (*ThriftFileModel, bool)
	Ptr(base, structName string) (*ThriftStructModel, error)
}

type ThriftGroup struct {
	ContentMap map[string]string
	ModelMap   map[string]*ThriftFileModel
}

func (t *ThriftGroup) GetFile(base string) (*ThriftFileModel, bool) {
	m, ok := t.ModelMap[base]
	return m, ok
}

func (t *ThriftGroup) Ptr(base, structName string) (*ThriftStructModel, error) {
	pkg := strings.Split(structName, ".")
	var sModel *ThriftStructModel
	if len(pkg) == 2 {

		s, ok := t.GetStruct(pkg[0], pkg[1])
		if !ok {
			return nil, inline.NewError(ErrNilStruct, "struct %s in %s is nil", pkg[1], pkg[0])
		}
		sModel = s
	} else {
		s, ok := t.GetStruct(base, pkg[0])
		if !ok {
			return nil, inline.NewError(ErrNilStruct, "struct %s  is nil", pkg[0])
		}
		sModel = s
	}
	return sModel, nil
}

func (t *ThriftGroup) GetStruct(base, structName string) (*ThriftStructModel, bool) {

	b, ok := t.GetFile(base)
	if !ok {
		return nil, false
	}
	m, ok := b.StructMap[structName]
	return m, ok
}

func (t *ThriftGroup) GetService(base, serviceName string) (*ThriftServiceModel, bool) {
	b, ok := t.GetFile(base)
	if !ok {
		return nil, false
	}
	m, ok := b.ServiceMap[serviceName]
	return m, ok
}

func NewThriftGroup(maps map[string]string) (*ThriftGroup, error) {
	ctx := &ThriftGroup{
		ContentMap: maps,
		ModelMap:   map[string]*ThriftFileModel{},
	}
	for base, content := range maps {
		if inline.IsEmpty(content) {
			continue
		}
		fileModel := NewThriftFileModel()
		if err := fileModel.Parse(content); err != nil {
			return nil, inline.PrependErrorFmt(err, "parse group err")
		}
		ctx.ModelMap[base] = fileModel
	}
	return ctx, nil
}

type ThriftFileModel struct {
	group ThriftGroup

	Namespace string
	Language  string
	Include   string

	StructMap  map[string]*ThriftStructModel
	ServiceMap map[string]*ThriftServiceModel
}

func (t *ThriftFileModel) Parse(content string) error {
	pattern := regexp.MustCompile(`namespace[ \t]+(\w+)[ \t]+(\w+)`)
	ss := pattern.FindStringSubmatch(content)
	if len(ss) < 3 {
		return inline.NewError(ErrRegexMatch, "parse file %s", content)
	}
	t.Namespace = ss[2]
	t.Language = ss[1]

	pattern = regexp.MustCompile(`struct[ \t]+\w+[ \t]+[{][^}]*[}]`)
	ss = pattern.FindStringSubmatch(content)
	for _, s := range ss {
		if inline.IsEmpty(s) {
			continue
		}
		structModel := NewThriftStructModel()
		if err := structModel.Parse(s); err != nil {
			return inline.PrependErrorFmt(err, "parse struct")
		}
		t.StructMap[structModel.StructName] = structModel
	}

	pattern = regexp.MustCompile(`service[ \t]+\w+[ \t]+[{][^}]*[}]`)
	ss = pattern.FindStringSubmatch(content)
	for _, s := range ss {
		if inline.IsEmpty(s) {
			continue
		}
		serviceModel := NewThriftServiceModel()
		if err := serviceModel.Parse(s); err != nil {
			return inline.PrependErrorFmt(err, "parse service")
		}
		t.ServiceMap[serviceModel.ServiceName] = serviceModel
	}

	return nil
}

func NewThriftFileModel() *ThriftFileModel {
	return &ThriftFileModel{ServiceMap: map[string]*ThriftServiceModel{}, StructMap: map[string]*ThriftStructModel{}}
}

type ThriftStructModel struct {
	StructName string

	FieldMap map[int16]*ThriftFieldModel
}

func NewThriftStructModel() *ThriftStructModel {
	return &ThriftStructModel{FieldMap: map[int16]*ThriftFieldModel{}}
}

func (t *ThriftStructModel) Parse(content string) error {
	content = strings.Trim(content, " ")
	pattern := regexp.MustCompile(`struct[ \t]+(\w+)[ \t]+[{][^}]*[}]`)
	ss := pattern.FindStringSubmatch(content)
	if len(ss) < 2 {
		return inline.NewError(ErrRegexMatch, "parse struct %s", content)
	}
	t.StructName = ss[1]
	splits := strings.Split(content, "\n")
	for i := 1; i < len(splits)-1; i++ {
		line := splits[i]
		if inline.IsEmpty(line) {
			continue
		}
		fieldModel := NewThriftFieldModel()
		if err := fieldModel.Parse(line); err != nil {
			return inline.PrependErrorFmt(err, "parse line %s", line)
		}
		t.FieldMap[fieldModel.Idx] = fieldModel
	}
	return nil
}

func NewThriftFieldModel() *ThriftFieldModel {
	return &ThriftFieldModel{}
}

type ThriftFieldModel struct {
	FieldName string
	Idx       int16
	//Tag string
	Type thrift.TType

	structTypeName string
}

func (t *ThriftFieldModel) Parse(content string) error {
	content = strings.Trim(content, " ")
	pattern := regexp.MustCompile(`(\d+)[ \t]*:[ \t]*(?:optional)?[ \t]*([a-zA-Z.0-9]+)(?:<.*>)?[\t ]+(\w+)`)
	ss := pattern.FindStringSubmatch(content)
	if len(ss) < 4 {
		return inline.NewError(ErrRegexMatch, "parse field %s", content)
	}
	iString := ss[1]
	idx, err := strconv.ParseInt(iString, 10, 16)
	if err != nil {
		return inline.PrependErrorFmt(err, "idx: %s", iString)
	}
	t.Idx = int16(idx)

	types := ss[2]
	if types == "i8" {
		t.Type = thrift.I08
	} else if types == "i16" {
		t.Type = thrift.I16
	} else if types == "i32" {
		t.Type = thrift.I32
	} else if types == "i64" {
		t.Type = thrift.I64
	} else if types == "bool" {
		t.Type = thrift.BOOL
	} else if strings.Contains(types, ".") {
		t.Type = thrift.STRUCT
	} else if strings.HasPrefix(types, "map") {
		t.Type = thrift.MAP
	} else if strings.HasPrefix(types, "list") {
		t.Type = thrift.LIST
	} else if strings.HasPrefix(types, "double") {
		t.Type = thrift.DOUBLE
	} else if types == "string" {
		t.Type = thrift.STRING
	} else {
		// is struct
		t.Type = thrift.STRUCT
	}

	if t.Type == thrift.STRUCT {
		t.structTypeName = ss[2]
	}

	t.FieldName = ss[3]
	return nil
}

type ThriftServiceModel struct {
	ServiceName string

	MethodMap map[string]*ThriftMethodModel
}

func NewThriftServiceModel() *ThriftServiceModel {
	return &ThriftServiceModel{MethodMap: map[string]*ThriftMethodModel{}}
}

func (t *ThriftServiceModel) Parse(content string) error {
	content = strings.Trim(content, " ")
	pattern := regexp.MustCompile(`service[ \t]+(\w+)[ \t]+[{][^}]*[}]`)
	ss := pattern.FindStringSubmatch(content)
	if len(ss) < 2 {
		return inline.NewError(ErrRegexMatch, "parse service %s", content)
	}
	t.ServiceName = ss[1]
	splits := strings.Split(content, "\n")
	for i := 1; i < len(splits)-1; i++ {
		line := splits[i]
		if inline.IsEmpty(line) {
			continue
		}
		serviceModel := NewThriftMethodModel()
		if err := serviceModel.Parse(line); err != nil {
			return inline.PrependErrorFmt(err, "parse line %s", line)
		}
		t.MethodMap[serviceModel.MethodName] = serviceModel
	}
	return nil
}

type ThriftMethodModel struct {
	RequestName string
	Request     string // type name
	Response    string // type name
	MethodName  string
}

func NewThriftMethodModel() *ThriftMethodModel {
	return &ThriftMethodModel{}
}

func (t *ThriftMethodModel) Parse(content string) error {
	content = strings.Trim(content, " ")
	pattern := regexp.MustCompile(`([a-zA-Z.0-9]+)[ \t]+(\w+)[(].*:[ \t]*([a-zA-Z.]+)[ \t]+(\w+)[)]`)
	groups := pattern.FindStringSubmatch(content)
	if len(groups) != 5 {
		return inline.NewError(ErrRegexMatch, "parse service %s", content)
	}

	t.Response = groups[1]
	t.MethodName = groups[2]
	t.Request = groups[3]
	t.RequestName = groups[4]
	return nil
}
