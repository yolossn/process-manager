package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Commands []Command `json:"commands" yaml:"commands"`
}

type Command struct {
	Command    string   `json:"command" yaml:"command"`
	Args       []string `json:"args" yaml:"args"`
	Envs       []Env    `json:"envs" yaml:"envs"`
	MaxRetries int      `json:"maxRetries" yaml:"maxRetries"`
}

type Env struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

func (c *Command) EnvStrings() []string {
	var envStrings []string
	for _, env := range c.Envs {
		envStrings = append(envStrings, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	return envStrings
}

func FromYaml(file string) (*Config, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading file")
	}

	conf := Config{}
	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshalling config")
	}
	return &conf, nil
}
