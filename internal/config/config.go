package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	CertFile  string
	KeyFile   string
	HTTPPort  string
	HTTPSPort string
	IP        string
	DBfile    string
}

func LoadConfig(path string) (*EnvConfig, error) {

	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	config := &EnvConfig{
		CertFile:  getEnv("CERTFILE", "f4k3"),
		KeyFile:   getEnv("KEYFILE", "f4k3"),
		HTTPPort:  getEnv("HTTP_PORT", "80"),
		HTTPSPort: getEnv("HTTPS_PORT", "443"),
		IP:        getEnv("IP", "localhost"),
		DBfile:    getEnv("DB_NAME", "myapp.db"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
