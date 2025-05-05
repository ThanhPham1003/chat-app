package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"postgres"`
	Redis struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"redis"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
	Services struct {
		User struct {
			GrpcAddr string `mapstructure:"grpc_addr"`
			HttpAddr string `mapstructure:"http_addr"`
		} `mapstructure:"user"`
		Message struct {
			GrpcAddr string `mapstructure:"grpc_addr"`
			HttpAddr string `mapstructure:"http_addr"`
		} `mapstructure:"message"`
	} `mapstructure:"services"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./../../configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unmarshal error 1:", err)
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Unmarshal error:", err)
		return nil, err
	}
	return &config, nil
}
