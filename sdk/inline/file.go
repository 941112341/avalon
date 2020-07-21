package inline

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Read(path string) (string, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FileName(filepath string) string {
	file := path.Base(filepath)
	ext := path.Ext(filepath)
	return strings.TrimSuffix(file, ext)
}
