package config

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	env := inline.GetEnv("env", "dev")
	first := filename[:i] + "." + env + filename[i:]
	names = append([]string{first}, names...)
	return
}

func Read(config interface{}, resource string) (err error) {
	files := extendedYamlNames(resource)
	for _, file := range files {
		err = read(config, file)
		if err != nil {
			inline.Warningln("read file err", inline.NewPair("resource", file), inline.NewPair("err", inline.VString(err)))
		}
		return nil
	}
	return errors.Wrap(err, resource)
}
