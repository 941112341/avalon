package tool

import (
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"text/template"
)

const Version = "1.0.0"

// output absolute
func build(info IDLInfo, version string) (string, error) {
	fileTemplate := FileTemplate{
		Package: info.Namespace,
		Version: version,
		IDLName: info.IDLName,
	}

	t, err := template.New("file").Parse(generateTemplate)
	if err != nil {
		return "", errors.Wrap(err, "parse template err")
	}
	var doc bytes.Buffer
	err = t.Execute(&doc, fileTemplate)
	if err != nil {
		return "", errors.Wrap(err, "execute err")
	}

	common := doc.String()

	for _, method := range info.Methods {
		t, err := template.New("method").Parse(methodTemplate)
		if err != nil {
			return "", errors.Wrap(err, "parse template err")
		}
		var doc bytes.Buffer
		err = t.Execute(&doc, method)
		if err != nil {
			return "", errors.Wrap(err, "execute err")
		}
		common = common + "\n" + doc.String()
	}
	return common, nil
}

func CreateFile(i, o string) error {
	s := NewScanner(i)
	info, err := s.Scan()
	if err != nil {
		return errors.Wrap(err, "scan")
	}
	content, err := build(*info, Version)
	if err != nil {
		return errors.Wrap(err, "build")
	}
	return ioutil.WriteFile(o, []byte(content), os.ModePerm)
}
