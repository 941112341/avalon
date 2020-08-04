package generic

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/apache/thrift/lib/go/thrift"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrRegexMatch inline.AvalonErrorCode = 100 + iota
	ErrNilStruct
)

type ThriftContext interface {
	GetStruct(base, structName string) (*ThriftStructModel, bool)
	GetService(base, serviceName string) (*ThriftServiceModel, bool)
	GetFile(base string) (*ThriftFileModel, bool)
	GetMethod(base, service, method string) (*ThriftMethodModel, bool)
	Ptr(base, structName string) (*ThriftStructModel, error)
}

type ThriftGroup struct {
	ModelMap map[string]*ThriftFileModel
}

func (t *ThriftGroup) GetMethod(base, service, method string) (*ThriftMethodModel, bool) {
	serviceModel, ok := t.GetService(base, service)
	if !ok {
		for _, model := range t.ModelMap {
			for _, serviceModel := range model.ServiceMap {
				if methodModel, ok := serviceModel.MethodMap[method]; ok {
					return methodModel, true
				}
			}
		}
		return nil, false
	}
	methodModel, ok := serviceModel.MethodMap[method]
	return methodModel, ok
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
	base, structName = strings.Trim(base, " "), strings.Trim(structName, " ")
	b, ok := t.GetFile(base)
	if !ok {
		for _, model := range t.ModelMap {
			if structModel, ok := model.StructMap[structName]; ok {
				return structModel, ok
			}
		}
		return nil, false
	}
	m, ok := b.StructMap[structName]
	return m, ok
}

func (t *ThriftGroup) GetService(base, serviceName string) (*ThriftServiceModel, bool) {
	b, ok := t.GetFile(base)
	if !ok {
		for _, model := range t.ModelMap {
			if serviceModel, ok := model.ServiceMap[serviceName]; ok {
				return serviceModel, ok
			}
		}
		return nil, false
	}
	m, ok := b.ServiceMap[serviceName]
	return m, ok
}

func NewThriftGroup(maps map[string]string) (*ThriftGroup, error) {
	ctx := &ThriftGroup{
		ModelMap: map[string]*ThriftFileModel{},
	}
	for base, content := range maps {
		if inline.IsEmpty(content) {
			continue
		}
		fileModel := NewThriftFileModel(base)
		if err := fileModel.Parse(content); err != nil {
			return nil, inline.PrependErrorFmt(err, "parse group err")
		}
		ctx.ModelMap[base] = fileModel
	}
	return ctx, nil
}

func NewThriftGroupBase(base []string) (*ThriftGroup, error) {

	maps := make(map[string]string)
	for _, s := range base {
		file := s
		suffix := ".thrift"
		if !strings.HasSuffix(s, suffix) {
			file = s + suffix
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, inline.PrependErrorFmt(err, "file %s", file)
		}
		content := string(data)
		maps[inline.FileName(file)] = content
	}

	grp, err := NewThriftGroup(maps)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "new group")
	}
	return grp, nil
}

type ThriftFileModel struct {
	Base      string
	Namespace string
	Language  string
	Include   []string

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

	pattern = regexp.MustCompile(`include[ \t]+"(\w+).thrift"`)
	sss := pattern.FindAllStringSubmatch(content, -1)
	for _, ss := range sss {
		t.Include = append(t.Include, ss[1])
	}

	pattern = regexp.MustCompile(`struct[ \t]+\w+[ \t]+[{][^}]*[}]`)
	sss = pattern.FindAllStringSubmatch(content, -1)
	for _, ss := range sss {
		structModel := NewThriftStructModel(t.Base)
		if err := structModel.Parse(ss[0]); err != nil {
			return inline.PrependErrorFmt(err, "parse struct")
		}
		t.StructMap[structModel.StructName] = structModel
	}

	pattern = regexp.MustCompile(`service[ \t]+\w+[ \t]+[{][^}]*[}]`)
	sss = pattern.FindAllStringSubmatch(content, -1)
	for _, ss := range sss {

		serviceModel := NewThriftServiceModel(t.Base)
		if err := serviceModel.Parse(ss[0]); err != nil {
			return inline.PrependErrorFmt(err, "parse service")
		}
		t.ServiceMap[serviceModel.ServiceName] = serviceModel
	}

	return nil
}

func NewThriftFileModel(base string) *ThriftFileModel {
	return &ThriftFileModel{ServiceMap: map[string]*ThriftServiceModel{}, StructMap: map[string]*ThriftStructModel{}, Base: base}
}

type ThriftStructModel struct {
	Base       string
	StructName string

	FieldMap map[int16]*ThriftFieldModel
}

func NewThriftStructModel(base string) *ThriftStructModel {
	return &ThriftStructModel{FieldMap: map[int16]*ThriftFieldModel{}, Base: base}
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
		fieldModel := NewThriftFieldModel(t.Base)
		if err := fieldModel.Parse(line); err != nil {
			return inline.PrependErrorFmt(err, "parse line %s", line)
		}
		t.FieldMap[fieldModel.Idx] = fieldModel
	}
	return nil
}

func NewThriftFieldModel(base string) *ThriftFieldModel {
	return &ThriftFieldModel{Base: base}
}

type ThriftFieldModel struct {
	Base      string
	FieldName string
	Idx       int16
	//Tag string
	TType       thrift.TType
	OptionalVar bool

	// thrift
	StructTypeName string
}

func (t *ThriftFieldModel) TypeName() string {
	return t.StructTypeName
}

func (t *ThriftFieldModel) ID() int16 {
	return t.Idx
}

func (t *ThriftFieldModel) JsonPath() string {
	return t.FieldName
}

func (t *ThriftFieldModel) ThriftName() string {
	return t.FieldName
}

func (t *ThriftFieldModel) Type() thrift.TType {
	return t.TType
}

func (t *ThriftFieldModel) Optional() bool {
	return t.OptionalVar
}

type regexpTemplate struct {
	Index    string
	Optional string
	TypeName string
	Name     string
}

// 命名变量组会更好
func (t *ThriftFieldModel) Parse(content string) error {
	content = strings.Trim(content, " ")
	pattern := regexp.MustCompile(`(?P<Index>\d+)[ \t]*:[ \t]*(?P<optional>optional|required)?[ \t]*(?P<TypeName>[\w.]+|list<.+>|map<.+>)[\t ]+(?P<Name>\w+)`)
	template := regexpTemplate{}
	err := inline.SubNameMatchStruct(pattern, content, &template)
	if err != nil {
		return inline.PrependErrorFmt(err, "match err %s", content)
	}
	id, err := strconv.ParseInt(template.Index, 10, 16)
	if err != nil {
		return inline.PrependErrorFmt(err, "parse int err")
	}
	t.Idx = int16(id)
	if template.Optional == "optional" {
		t.OptionalVar = true
	}
	t.StructTypeName = template.TypeName
	t.FieldName = template.Name
	t.TType = String2Thrift(t.StructTypeName)
	return nil
}

func (t *ThriftFieldModel) Elem() (m *ThriftFieldModel) {
	if t.TType != thrift.LIST {
		return
	}
	types := inline.Unwrap(`<(.*)>`, t.StructTypeName)
	ttype := String2Thrift(types)
	return &ThriftFieldModel{
		Base:           t.Base,
		TType:          ttype,
		OptionalVar:    false,
		StructTypeName: types,
	}
}

func (t *ThriftFieldModel) KVElem() (m, n *ThriftFieldModel) {
	if t.TType != thrift.MAP {
		return
	}
	types := inline.Unwraps(`<(.*),(.*)>`, t.StructTypeName)
	ktypeString, vtypeString := types[0], types[1]
	ktype, vtype := String2Thrift(ktypeString), String2Thrift(vtypeString)
	return &ThriftFieldModel{
			Base:           t.Base,
			TType:          ktype,
			OptionalVar:    false,
			StructTypeName: ktypeString,
		}, &ThriftFieldModel{
			Base:           t.Base,
			TType:          vtype,
			OptionalVar:    false,
			StructTypeName: vtypeString,
		}
}

func String2Thrift(types string) thrift.TType {
	types = strings.Trim(types, " ")
	var tt thrift.TType
	if types == "i8" {
		tt = thrift.I08
	} else if types == "i16" {
		tt = thrift.I16
	} else if types == "i32" {
		tt = thrift.I32
	} else if types == "i64" {
		tt = thrift.I64
	} else if types == "bool" {
		tt = thrift.BOOL
	} else if strings.HasPrefix(types, "map") {
		tt = thrift.MAP
	} else if types == "double" {
		tt = thrift.DOUBLE
	} else if strings.HasPrefix(types, "list") {
		tt = thrift.LIST
	} else if types == "string" {
		tt = thrift.STRING
	} else {
		// is struct
		tt = thrift.STRUCT
	}
	return tt
}

type ThriftServiceModel struct {
	Base        string
	ServiceName string

	MethodMap map[string]*ThriftMethodModel
}

func NewThriftServiceModel(base string) *ThriftServiceModel {
	return &ThriftServiceModel{MethodMap: map[string]*ThriftMethodModel{}, Base: base}
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
