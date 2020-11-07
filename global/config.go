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
	} `yaml:"email"`
	Bcrypt struct {
		Cost int `yaml:"cost"`
	} `yaml:"brypt"`
	Database struct {
		Host        string `yaml:"host"`
		Port        uint   `yaml:"port"`
		Name        string `yaml:"name"`
		Connections struct {
			MaxIdle     int `yaml:"max-idle"`
			MaxOpen     int `yaml:"max-open"`
			MaxIdleTime int `yaml:"max-idle-time"`
			MaxLifetime int `yaml:"max-lifetime"`
		} `yaml:"connections"`
		Credentials struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
	} `yaml:"database"`
}

func LoadConfig(configFileName string) (*Config, error) {
	config := new(Config)
	configFileBytes, err := ioutil.ReadFile(configFileName)
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
