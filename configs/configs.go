package configs

import (
	"fmt"
	"os"

	"github.com/ai-readone/go-url-shortner/logger"

	"gopkg.in/yaml.v2"
)

type Database struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DbName   string `yaml:"db_name"`
}

type Config struct {
	Server         string   `yaml:"server"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
	Database       Database `yaml:"database"`
}

var conf Config

// Initializes the configuration variables
// using the configs/configs.yaml file
func LoadConfig(path string) *Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to read config file, error: %v ", err))
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to process config, error: %v", err))
	}

	return &conf
}
