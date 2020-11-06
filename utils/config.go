package utils

import "github.com/BurntSushi/toml"

var Config config

type config struct {
	Host  string
	Port  int
	Auth  bool
	Token string
}

func (c *config) Parse() error {
	_, err := toml.DecodeFile("./config.toml", c)
	return err
}
