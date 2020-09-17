package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config is a collection of commands.
type Config struct {
	Commands []Command `json:"commands" yaml:"commands"`
}

// Command defines the config for a command.
type Command struct {
	Command    string   `json:"command" yaml:"command"`
	Args       []string `json:"args" yaml:"args"`
	Envs       []env    `json:"envs" yaml:"envs"`
	MaxRetries int      `json:"maxRetries" yaml:"maxRetries"`
}

type env struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// EnvStrings converts env key value to array of strings "key=value".
func (c *Command) EnvStrings() []string {
	var envStrings []string
	for _, env := range c.Envs {
		envStrings = append(envStrings, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	return envStrings
}

// FromYaml reads config from the yaml file.
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
