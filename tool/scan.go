package tool

import (
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/inline"
	"path"
)

type Scanner struct {
	IDLPath  string
	filename string
}

type IDLInfo struct {
	IDLName     string
	Namespace   string
	ServiceName string
	// struct info no need
	Methods []MethodTemplate
}

func (s *Scanner) Scan() (*IDLInfo, error) {
	p := s.IDLPath
	idlName := inline.FileName(p)
	contentMap := make(map[string]string)

	if !inline.Exists(p) {
		return nil, inline.NewError(ErrNilPackage, "idl path %s", p)
	}
	data, err := inline.Read(p)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "read %s err", p)
	}
	contentMap[idlName] = data

	group, err := generic.NewThriftGroup(contentMap)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "parse group %s", p)
	}
	thriftFile, ok := group.GetFile(idlName)
	if !ok {
		return nil, inline.Error("get file !ok")
	}

	serviceModel, ok := inline.UnionValue(thriftFile.ServiceMap).(*generic.ThriftServiceModel)
	if !ok {
		return nil, inline.Error("get file !ok")
	}
	idlInfo := &IDLInfo{
		IDLName:     idlName,
		Namespace:   thriftFile.Namespace,
		ServiceName: serviceModel.ServiceName,
		Methods:     nil,
	}

	for _, method := range serviceModel.MethodMap {
		idlInfo.Methods = append(idlInfo.Methods, MethodTemplate{
			MethodName: method.MethodName,
			Request:    method.Request,
			Response:   method.Response,
		})
	}
	return idlInfo, nil
}

func NewScanner(p string) *Scanner {
	filename := path.Base(p)
	return &Scanner{IDLPath: p, filename: filename}
}
