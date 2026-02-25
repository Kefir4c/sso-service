package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
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
	Host    string        `yaml:"host" env-default:"localhost"`
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	configPath := fethConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	return MustLoadPath(configPath)
}
func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists" + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config" + err.Error())
	}

	return &cfg
}

func fethConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	if path == "" {
		path = "./config/local.yaml"
	}

	return path
}
