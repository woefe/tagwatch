/*
 * Copyright (c) 2024. Wolfgang Popp
 *
 * This file is part of tagwatch.
 *
 * tagwatch is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tagwatch is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tagwatch.  If not, see <http://www.gnu.org/licenses/>.
 */

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
