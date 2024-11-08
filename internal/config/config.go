package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

var once sync.Once
var instance *Config

func Init() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			fmt.Println("Error reading config:", err)
		}
	})
	return instance, nil
}
