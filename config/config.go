package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret_key"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// DashScopeConfig 存储 DashScope 平台配置
type DashScopeConfig struct {
	ApiKey string `mapstructure:"api_key"`
}


type Config struct {
	Database  DBConfig
	JWT       JWTConfig
	DashScope DashScopeConfig
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}
}
