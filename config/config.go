package config

import (
	"os"

	"github.com/dmitrykharchenko95/otus_user/internal/database"
	"github.com/dmitrykharchenko95/otus_user/internal/server"
)

type Config struct {
	Server server.Config   `json:"server"`
	DB     database.Config `json:"db"`
}

func NewFromENVs() Config {
	return Config{
		Server: server.Config{
			Address: os.Getenv("SVC_ADDRESS"),
		},
		DB: database.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Database: os.Getenv("DB_NAME"),
		},
	}
}
