package global

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Email struct {
		ApiKey string `yaml:"api-key"`
		Name   string `yaml:"name"`
		From   string `yaml:"from"`
	}
}

func LoadConfig(configFileName string) (*Config, error) {
	config := new(Config)
	configFileBytes, err := ioutil.ReadFile(configFileName)
	log.Printf("config-file: %s", string(configFileBytes))
	if err != nil {
		log.Printf("Failed to read config file: config-file-location=%s", configFileName)
		return nil, err
	}
	err = yaml.Unmarshal(configFileBytes, config)
	if err != nil {
		log.Printf("Failed to parse config file yaml: err=%+v", err)
		return nil, err
	}
	return config, nil
}
