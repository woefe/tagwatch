package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Registry struct {
	Auth         bool   `yaml:"auth"`
	AuthUsername string `yaml:"username"`
	AuthPassword string `yaml:"password"`
	AuthURL      string `yaml:"auth_url"`
	BaseURL      string `yaml:"base_url"`
}

type WatchConf struct {
	WatchNew bool      `yaml:"watch_new"`
	Registry *Registry `yaml:"registry"`
	Arch     string    `yaml:"arch"`
	Repo     string    `yaml:"repo"`
	Tags     []string  `yaml:"tags"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Conf struct {
	Tagwatch []*WatchConf `yaml:"tagwatch"`
	Server   *Server      `yaml:"server"`
}

func LoadConf(filename string) (*Conf, error) {
	conf := &Conf{}

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileContents, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
