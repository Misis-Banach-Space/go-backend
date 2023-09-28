package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LoggingLevel string

	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresDb       string

	ServerPort string
}

var Cfg *Config

func NewConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	Cfg = &Config{
		LoggingLevel:     os.Getenv("LOGGING_LEVEL"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresDb:       os.Getenv("POSTGRES_DB"),
		ServerPort:       os.Getenv("SERVER_PORT"),
	}
	return nil
}
