package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath  string `env:"DATABASE_URL" env-required:"true"`
	Port         string `env:"PORT" env-required:"true"`
	GCPProjectID string `env:"GCP_PROJECT_ID" env-required:"true"`

	PickupRequestSubscriptionID       string `yaml:"pickup_request_subscription_id" env:"PICKUP_REQUEST_SUBSCRIPTION_ID" env-required:"true"`
	DriverLocationSubscriptionID      string `yaml:"driver_location_subscription_id" env:"DRIVER_LOCATION_SUBSCRIPTION_ID" env-required:"true"`
	AcceptPickupRequestSubscriptionID string `yaml:"accept_pickup_request_subscription_id" env:"ACCEPT_PICKUP_REQUEST_SUBSCRIPTION_ID" env-required:"true"`
	RejectPickupRequestSubscriptionID string `yaml:"reject_pickup_request_subscription_id" env:"REJECT_PICKUP_REQUEST_SUBSCRIPTION_ID" env-required:"true"`
	StartDeliverySubscriptionID       string `yaml:"start_delivery_subscription_id" env:"START_DELIVERY_SUBSCRIPTION_ID" env-required:"true"`
	EndDeliverySubscriptionID         string `yaml:"end_delivery_subscription_id" env:"END_DELIVERY_SUBSCRIPTION_ID" env-required:"true"`
	AssignDriverSubscriptionID        string `yaml:"assign_driver_subscription_id" env:"ASSIGN_DRIVER_SUBSCRIPTION_ID" env-required:"true"`
	UnassignDriverSubscriptionID      string `yaml:"unassign_driver_subscription_id" env:"UNASSIGN_DRIVER_SUBSCRIPTION_ID" env-required:"true"`
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
