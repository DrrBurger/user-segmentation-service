package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Env string `yaml:"local"`
		HTTP
		Log
		PG
		Hasher
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		URL      string `env:"PG_URL"`
		URLLocal string `env:"PG_URL_LOCAL"`
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DBName   string `env:"POSTGRES_DB"`
	}

	Hasher struct {
		Salt string `env:"HASHER_SALT"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := &Config{}

	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
