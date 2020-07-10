package config

import (
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const defaultConfig = "config.yaml"

func read(config interface{}, resource string) error {
	file, err := ioutil.ReadFile(resource)
	if err != nil {
		return errors.WithMessage(err, resource)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return errors.WithMessage(err, string(file))
	}
	return nil
}

func extendedYamlNames(filename string) (names []string) {
	i := strings.LastIndex(filename, ".yaml")

	if i < 0 {
		filename = defaultConfig
	}
	names = append(names, filename)
	env := os.Getenv("env")
	if env != "" {
		first := filename[:i] + "." + env + filename[i:]
		names = append([]string{first}, names...)
	}
	return
}

func Read(config interface{}, resource string) error {
	files := extendedYamlNames(resource)
	if len(files) == 0 {
		return errors.New(resource + " conf not found")
	}
	var err error
	for _, file := range files {
		err = read(config, file)
		if err != nil {
			log.New().WithField("file", file).
				WithField("err", err.Error()).Info("write yaml err")
		}
		return nil
	}
	return errors.Wrap(err, resource)
}
