package tagwatch

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Registry struct {
	Auth    bool   `yaml:"auth"`
	AuthURL string `yaml:"auth_url"`
	BaseURL string `yaml:"base_url"`
}

type WatchConf struct {
	WatchNew bool      `yaml:"watch_new"`
	Registry *Registry `yaml:"registry"`
	Arch     string    `yaml:"arch"`
	Repo     string    `yaml:"repo"`
	Tags     []string  `yaml:"tags"`
}

type Conf struct {
	Tagwatch []*WatchConf `yaml:"tagwatch"`
}

func LoadConf(filename string) (*Conf, error) {
	conf := &Conf{}

	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileContents, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
