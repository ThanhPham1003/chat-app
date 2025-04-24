package db

import (
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Redis struct {
		Addr string
	}
	JWT struct {
		Secret string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./../../configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
