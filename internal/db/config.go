package db

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"database"`
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
	var config Config
	viper.SetConfigFile("./../../configs/config.yaml")
	if err := viper.ReadInConfig(); err == nil {
		if err := viper.Unmarshal(&config); err != nil {
			return nil, err
		}
		return &config, nil
	}
	// Override with environment variables
	if v := os.Getenv("DB_HOST"); v != "" {
		config.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := parseInt(v); err == nil {
			config.Database.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		config.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		config.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		config.Database.DBName = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		config.Redis.Addr = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		config.JWT.Secret = v
	}
	if v := os.Getenv("USER_GRPC_ADDR"); v != "" {
		config.Services.User.GrpcAddr = v
	}
	if v := os.Getenv("USER_HTTP_ADDR"); v != "" {
		config.Services.User.HttpAddr = v
	}
	if v := os.Getenv("MESSAGE_GRPC_ADDR"); v != "" {
		config.Services.Message.GrpcAddr = v
	}
	if v := os.Getenv("MESSAGE_HTTP_ADDR"); v != "" {
		config.Services.Message.HttpAddr = v
	}

	return &config, nil
}

func parseInt(s string) (int, error) {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	return v, err
}
