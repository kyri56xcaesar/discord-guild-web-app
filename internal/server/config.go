package server

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ConfigPath     string
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

	split := strings.Split(path, "/")

	config := &EnvConfig{
		ConfigPath:     split[len(split)-1],
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

// CertFile string, KeyFile string, HTTPPort string, HTTPSPort string, IP string, DBfile string, AllowedOrigins []string, AllowedHeaders []string
// AllowedMethods []string
func (cfg *EnvConfig) toString() string {
	var strBuilder strings.Builder

	reflectedValues := reflect.ValueOf(cfg).Elem()
	reflectedTypes := reflect.TypeOf(cfg).Elem()

	strBuilder.WriteString(fmt.Sprintf("CONFIG: %s\n", cfg.ConfigPath))

	for i := 0; i < reflectedValues.NumField(); i++ {
		fieldName := reflectedTypes.Field(i).Name
		fieldValue := reflectedValues.Field(i).Interface()
		if i < 9 {
			strBuilder.WriteString(fmt.Sprintf("%d.  ", i+1))
		} else {
			strBuilder.WriteString(fmt.Sprintf("%d. ", i+1))
		}
		if len(fieldName) < 5 {
			strBuilder.WriteString(fmt.Sprintf("%v\t\t\t-> %v\n", fieldName, fieldValue))
		} else if len(fieldName) < 14 {
			strBuilder.WriteString(fmt.Sprintf("%v\t\t-> %v\n", fieldName, fieldValue))
		} else {
			strBuilder.WriteString(fmt.Sprintf("%v\t-> %v\n", fieldName, fieldValue))
		}
	}

	strBuilder.WriteString("\n")

	return strBuilder.String()
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
