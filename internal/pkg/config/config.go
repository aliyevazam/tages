package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	LogLevel    string
	ServiceHost string
	ServicePort string
	FilePath    string
}

func Load() Config {
	godotenv.Load("./sample.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		Environment: conf.GetString("ENVIRONMENT"),
		LogLevel:    conf.GetString("LOG_LEVEL"),
		ServiceHost: conf.GetString("SERVICE_HOST"),
		ServicePort: conf.GetString("SERVICE_PORT"),
		FilePath:    conf.GetString("FILE_PATH"),
	}
	return cfg
}
