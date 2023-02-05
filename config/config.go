package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	SigningKey string `yaml:"signingkey"`
	Salt       string `yaml:"salt"`
	URI        string
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile("./config/app.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
		config.URI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority", config.User, config.Pass, config.Host, config.Port)
	})
	return &config
}
