package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Storage StorageConfig `mapstructure:"storage"`
	Auth    AuthConfig    `mapstructure:"auth"`
}

type ServerConfig struct {
	ListenAddress string `mapstructure:"listen_address"`
	Environment   string `mapstructure:"environment"`
}

type StorageConfig struct {
	DSN string `mapstructure:"dsn"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
