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

func AllFileInDir(dir string) map[string]os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return map[string]os.FileInfo{}
	}
	fileMap := make(map[string]os.FileInfo)
	for _, file := range files {
		filepath := path.Join(dir, file.Name())
		if file.IsDir() {
			subFiles := AllFileInDir(filepath)
			for abPath, info := range subFiles {
				fileMap[abPath] = info
			}
		} else {
			fileMap[filepath] = file
		}
	}
	return fileMap
}

func AllFIleInDirFunc(dir string, function func(info os.FileInfo) bool) map[string]os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return map[string]os.FileInfo{}
	}
	fileMap := make(map[string]os.FileInfo)
	for _, file := range files {
		if !function(file) {
			continue
		}
		filepath := path.Join(dir, file.Name())
		if file.IsDir() {
			subFiles := AllFIleInDirFunc(filepath, function)
			for abPath, info := range subFiles {
				fileMap[abPath] = info
			}
		} else {
			fileMap[filepath] = file
		}
	}
	return fileMap
}
