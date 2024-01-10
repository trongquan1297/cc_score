package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// DatabaseConfig struct represents the database configuration.
type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

// Config struct represents the structure of the configuration file (config.yml).
type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig loads the configuration from the specified file.
func LoadConfig(filename string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return config, nil
}
