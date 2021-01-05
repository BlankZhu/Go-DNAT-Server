package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var conf ProxyConfig

// MySQLConfig config of mysql connection
type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	IP       string `yaml:"ip"`
	Port     uint16 `yaml:"port"`
	Schema   string `yaml:"schema"`
}

// ProxyConfig config of the proxy
type ProxyConfig struct {
	Port      uint16      `yaml:"port"`
	PrivateIP string      `yaml:"private"`
	PublicIP  string      `yaml:"public"`
	MySQL     MySQLConfig `yaml:"mysql"`
}

// LoadYAMLConfigFromFile load yaml config from given file path
func LoadYAMLConfigFromFile(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return err
	}
	return nil
}

// Get get the config
func Get() *ProxyConfig {
	return &conf
}
