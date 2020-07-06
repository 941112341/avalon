package config

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func ReadYaml(config interface{}, resource string) error {
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

func WriteYaml(config interface{}, resource string) error {
	file, err := yaml.Marshal(config)
	if err != nil {
		return errors.WithMessage(err, inline.JsonString(file))
	}
	return ioutil.WriteFile(resource, file, os.ModePerm)
}

func confList(filename string) []string {
	i := strings.LastIndex(filename, ".yaml")
	if i < 0 {
		return []string{}
	}

	defaultConfig := "config.yaml"
	env := os.Getenv("env")
	if env == "" {
		return []string{filename, defaultConfig}
	}
	first := filename[:i] + "." + env + filename[i:]
	return []string{first, filename, defaultConfig}
}

func Read(config interface{}, resource string) error {
	files := confList(resource)
	if len(files) == 0 {
		return errors.New(resource + " conf not found")
	}
	var err error
	for _, file := range files {
		err = ReadYaml(config, file)
		if err != nil {
			log.New().WithField("file", file).
				WithField("err", err.Error()).Info("write yaml err")
		}
		return nil
	}
	return errors.Wrap(err, resource)
}
