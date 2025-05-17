package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type (
	Database struct {
		Host     string `yaml:"host" envDefault:"mongo"`
		Port     string `yaml:"port" envDefault:"27017"`
		Username string `yaml:"username" envDefault:"root"`
		Password string `yaml:"password" envDefault:"example"`
		Name     string `yaml:"name" envDefault:"doto"`
	}

	Serving struct {
		Host string `yaml:"host" envDefault:"doto"`
		Port string `yaml:"port" envDefault:"8080"`
	}
)

type Config struct {
	Database Database `yaml:"database"`
	Serving  Serving  `yaml:"serving"`
}

func New(name string) (*Config, error) {
	cfg, err := NewConfigFromFile(name)
	if err != nil {
		return nil, err
	}

	if err := NewConfigFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewConfigFromFile(name string) (*Config, error) {
	cfg := &Config{}

	v := viper.New()

	v.SetConfigType("yaml")

	v.SetConfigFile(name)

	if err := v.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := v.MergeInConfig(); err != nil {
		if errors.Is(err, &viper.ConfigParseError{}) {
			return nil, fmt.Errorf("merge config: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func NewConfigFromEnv(cfg *Config) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}

func (d *Database) ToDSN() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

var (
	//go:embed default-config.yaml
	configBytes []byte
)
