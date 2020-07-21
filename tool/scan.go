package tool

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"path"
	"regexp"
	"strings"
)

var (
	Namespace  = regexp.MustCompile(`namespace[ \t]+go[ \t]+(\w+)`)
	IDLService = regexp.MustCompile(`service[ \t]+(.*)Service[ \t]+[{]`)
	Method     = regexp.MustCompile(`(\w+)Response[ \t]+(\w+)\(1:[ \t]*(\w+)Request[ \t]+[\w]+\)`)
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
	if !inline.Exists(p) {
		return nil, errors.New("file not exists")
	}

	data, err := inline.Read(p)
	if err != nil {
		return nil, errors.Wrap(err, "read err")
	}

	return s.parse(data)
}

func (s *Scanner) parse(data string) (*IDLInfo, error) {

	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return nil, errors.New("data is nil")
	}
	namespaceLine := ""
	index := 0
	for idx, line := range lines {
		if strings.Trim(line, " ") == "" {
			continue
		}
		namespaceLine = line
		index = idx
		break
	}
	// read first 第一行应该为namespace
	subMatches := Namespace.FindStringSubmatch(namespaceLine)
	if len(subMatches) < 2 {
		return nil, fmt.Errorf("namespace expect in 1st fact %s", lines[0])
	}
	idlInfo := IDLInfo{Namespace: subMatches[1], IDLName: inline.Ucfirst(inline.FileName(s.filename))}
	for idx := index; idx < len(lines); idx++ {
		line := lines[idx]
		if strings.Trim(line, " ") == "" {
			continue
		}
		if idx < 2 {
			continue
		}
		subStr := IDLService.FindStringSubmatch(line)
		if len(subStr) < 2 {
			continue
		}

		idlInfo.ServiceName = subStr[1]
		// 开始解析method
		for i := idx + 1; i < len(lines); i++ {
			arr := Method.FindStringSubmatch(lines[i])
			if len(arr) < 4 {
				continue
			}
			m := MethodTemplate{
				MethodName: arr[2],
				Request:    arr[3] + "Request",
				Response:   arr[1] + "Response",
			}
			idlInfo.Methods = append(idlInfo.Methods, m)
		}
		break
	}
	return &idlInfo, nil
}

func NewScanner(p string) *Scanner {
	filename := path.Base(p)
	return &Scanner{IDLPath: p, filename: filename}
}
