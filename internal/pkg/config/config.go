package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment      string
	LogLevel         string
	ServiceHost      string
	ServicePort      string
	FilePath         string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
}

func Load() Config {
	godotenv.Load("./.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		Environment:      conf.GetString("ENVIRONMENT"),
		LogLevel:         conf.GetString("LOG_LEVEL"),
		ServiceHost:      conf.GetString("SERVICE_HOST"),
		ServicePort:      conf.GetString("SERVICE_PORT"),
		FilePath:         conf.GetString("FILE_PATH"),
		PostgresHost:     conf.GetString("POSTGRES_HOST"),
		PostgresPort:     conf.GetInt("POSTGRES_PORT"),
		PostgresUser:     conf.GetString("POSTGRES_USER"),
		PostgresPassword: conf.GetString("POSTGRES_PASSWORD"),
		PostgresDatabase: conf.GetString("POSTGRES_DB"),
	}
	return cfg
}
