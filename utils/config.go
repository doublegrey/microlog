package utils

import "github.com/BurntSushi/toml"

var Config config

type config struct {
	Host  string
	Port  int
	Auth  bool
	Token string
	Kafka kafka
	UI    ui
}

type kafka struct {
	Brokers      []string
	Topic        string
	CustomTopics bool   `toml:"allow_custom_topic"`
	TopicPrefix  string `toml:"custom_topic_prefix"`
}

type ui struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
}

func (c *config) Parse() error {
	_, err := toml.DecodeFile("./config.toml", c)
	return err
}
