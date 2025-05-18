package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type HTTPServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Config struct {
	Env          string              `yaml:"env"`
	HttpServer   *HTTPServerConfig   `yaml:"http_server"`
	CommandPaths map[string][]string `yaml:"paths"`
}

func MustLoad() *Config {
	path := MustGetPath()

	return MustLoadByPath(path)
}

func MustGetPath() string {
	path := getPath()
	if path == "" {
		log.Fatal("config path not set")
	}

	return path
}

func getPath() string {
	if path := getPathByEnv(); path != "" {
		return path
	}
	return getPathByFlag()
}

func getPathByEnv() string {
	path := os.Getenv("CONFIG_PATH")
	return path
}

func getPathByFlag() string {
	var path string

	flag.StringVar(&path, "config_path", "", "path to config file")
	flag.Parse()
	return path
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal("failed to read config: " + err.Error())
	}

	return &cfg
}
