package config

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost   string `mapstructure:"service_host"`
	ServicePort   int    `mapstructure:"service_port"`
	JWTSecret     string `mapstructure:"jwt_secret"`
	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	configName := "config"
	if envCfg := os.Getenv("CONFIG_NAME"); envCfg != "" {
		configName = envCfg
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config") // на случай, если запуск из корня, а конфиг в папке config/
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("jwt_secret is not set in config.toml or .env")
	}

	log.Info("config parsed successfully")
	return &cfg, nil
}