package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `env:"DATABASE_URL" env-required:"true"`
	Port        string `env:"PORT" env-required:"true"`
}

func MustLoad() *Config {

	errorFromLoadingEnv := godotenv.Load()
	if errorFromLoadingEnv != nil {
		log.Fatal("Error loading .env file")
	}

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	errFromReadingYAML := cleanenv.ReadConfig(configPath, &cfg)
	if errFromReadingYAML != nil {
		log.Fatalf("can not read config file: %s", errFromReadingYAML.Error())
	}
	errFromReadingENV := cleanenv.ReadEnv(&cfg)
	if errFromReadingENV != nil {
		log.Fatalf("can not read env variables: %s", errFromReadingENV.Error())
	}

	return &cfg
}
