package config

import (
	"os"
)

const (
	Port = "PORT"

	MongoUri = "MONGO_URI"
)

type Config struct {
	MongoURI string
	Port     string
}

func Load() Config {
	var cfg Config

	cfg.Port = valueFromEnvOrDefault(Port, "8080")
	cfg.MongoURI = valueFromEnvOrDefault(MongoUri, "mongodb://user:password@localhost:27017")

	return cfg
}

func valueFromEnvOrDefault(tag string, defaultValue string) string {
	d, ok := os.LookupEnv(tag)
	if !ok {
		return defaultValue
	}

	return d
}
