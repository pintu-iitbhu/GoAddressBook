package configs

import (
	"GoAddressBook/constants"
	"GoAddressBook/springcloud"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Map    map[string]interface{}
	Errors []error
}

func new() *Config {
	var errors []error
	m := make(map[string]interface{})
	return &Config{Map: m, Errors: errors}
}

func NewConfig() {
	cfg := new()
	env := os.Getenv(constants.EnvKey)
	if env == "" {
		env = constants.DevEnvironment
	}
	if env == constants.DevEnvironment {
		cfg.FromSpringFile(constants.DevConfigJsonFilePath)
	} else {
		service := constants.Service
		cfg.FromSpringCloudDomain(service, env)
	}
	cfg.SetViper()
}

func (c *Config) FromSpringCloudDomain(service, env string) {
	//TODO Implement logic to read config from cloud
}

func (c *Config) FromSpringFile(path string) *Config {
	springConfig, err := springcloud.GetFile(path)
	if err != nil {
		c.Errors = append(c.Errors, err)
		return c
	}
	springMap, err := springConfig.ToMap()
	if err != nil {
		c.Errors = append(c.Errors, err)
		return c
	}
	c.AddMap(springMap)
	c.AddEnv(springMap)
	return c
}
func (c *Config) AddMap(configurationsMap map[string]interface{}) {
	for k, v := range configurationsMap {
		c.Map[k] = v
	}
}
func (c *Config) AddEnv(configurationsMap map[string]interface{}) {
	for k, v := range configurationsMap {
		if _, ok := v.(string); ok {
			value := v.(string)
			if len(value) >= 2 && value[:2] == "${" {
				c.Map[k] = os.Getenv(value[2 : len(value)-1])
				continue
			}
		}
	}
}

func (c *Config) SetViper() {
	for k, v := range c.Map {
		viper.SetDefault(k, v)
	}
}
