package tool

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
)

type Scanner struct {
	IDLPath string
}

type IDLInfo struct {
	IDLName     string
	Namespace   string
	ServiceName string

	// struct info no need
	Methods []MethodTemplate
}

func (s *Scanner) Scan() (*IDLInfo, error) {
	path := s.IDLPath
	if !inline.Exists(path) {
		return nil, errors.New("file not exists")
	}

	data, err := inline.Read(path)
	if err != nil {
		return nil, errors.Wrap(err, "read err")
	}

	return s.parse(data)
}

func (s *Scanner) parse(data string) (*IDLInfo, error) {

}
