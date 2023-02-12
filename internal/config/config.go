package config

import (
	"ChatGo/pkg/logging"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Mongo struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		DB   string `yaml:"db"`
		URI  string
	} `yaml:"mongo"`
	SigningKey string `yaml:"signingkey"`
	Salt       string `yaml:"salt"`
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Get config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("./app.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
		instance.Mongo.URI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority", instance.Mongo.User, instance.Mongo.Pass, instance.Mongo.Host, instance.Mongo.Port)
		logger.Debugf("config %v", instance)
	})
	return instance
}
