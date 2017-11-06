package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Host        string   `yaml:"host"`
	Addr        string   `yaml:"addr"`
	Secure      bool     `yaml:"secure"`
	UsePassword bool     `yaml:"use_password"`
	Password    string   `yaml:"password"`
	User        string   `yaml:"user"`
	Nick        string   `yaml:"nick"`
	Channels    []string `yaml:"channels"`
}

type Client struct {
	Addr        string `yaml:"addr"`
	UsePassword bool   `yaml:"use_password"`
	Password    string `yaml:"password"`
}

type Config struct {
	Server *Server `yaml:"server"`
	Client *Client `yaml:"client"`
}

func Load(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
