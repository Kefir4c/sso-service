package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	MigrationsPath string        `yaml:"migrations-path" env-default:"./migrations"`
	Storage        StorageConfig `yaml:"storage"`
	Cache          CacheConfig   `yaml:"cache"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	TokenTTL       time.Duration `yaml:"tokenTTL" env-default:"1h"`
}

type StorageConfig struct {
	Type     string `yaml:"type" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-default:"sso_db"`
}

type CacheConfig struct {
	Driver string `yaml:"driver" env-default:"redis"`
	Host   string `yaml:"host" env-default:"localhost"`
	Port   int    `yaml:"port" env-default:"6379"`
	DB     int    `yaml:"db" env-default:"0"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func mustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file does not exist: %s\n\n"+
			"💡 Quick start:\n"+
			"   cp config/local.yaml.example config/local.yaml\n"+
			"   go run main.go\n\n"+
			"📌 Or set custom path:\n"+
			"   go run main.go -config=./config/prod.yaml\n"+
			"   export CONFIG_PATH=./config/prod.yaml",
			configPath))
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(fmt.Sprintf("Failed to read config: %v", err))
	}

	return &config
}
