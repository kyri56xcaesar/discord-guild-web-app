package server

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	CertFile       string
	KeyFile        string
	HTTPPort       string
	HTTPSPort      string
	IP             string
	DBfile         string
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
}

func loadConfig(path string) (*EnvConfig, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	config := &EnvConfig{
		CertFile:       getEnv("CERTFILE", "f4k3"),
		KeyFile:        getEnv("KEYFILE", "f4k3"),
		HTTPPort:       getEnv("HTTP_PORT", "80"),
		HTTPSPort:      getEnv("HTTPS_PORT", "443"),
		IP:             getEnv("IP", "localhost"),
		DBfile:         getEnv("DB_NAME", "dads.db"),
		AllowedOrigins: parseValues("ALLOWED_ORIGINS", []string{"None"}),
		AllowedHeaders: parseValues("ALLOWED_HEADERS", nil),
		AllowedMethods: parseValues("ALLOWED_METHODS", nil),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func parseValues(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		values := strings.SplitAfter(value, ",")
		return values
	}

	return fallback
}
