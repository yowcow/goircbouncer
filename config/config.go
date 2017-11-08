package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Server represents a targeted-server configuration
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

// Client represents a client-side configuration
type Client struct {
	Addr        string `yaml:"addr"`
	UsePassword bool   `yaml:"use_password"`
	Password    string `yaml:"password"`
}

// Config represents a configuration
type Config struct {
	Server *Server `yaml:"server"`
	Client *Client `yaml:"client"`
}

// Load loads config from a given YAML file
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
